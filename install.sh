#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_PATH="${BASH_SOURCE[0]:-$0}"
ROOT_DIR="$(cd "$(dirname "$SCRIPT_PATH")" && pwd)"
APP_DIR="${APP_DIR:-/opt/chega_junto_concurseiro}"
REPOSITORY_URL="https://github.com/reinaldorocha/cjc.git"
ACTION="${1:-install}"
DEPLOY_DIR="$ROOT_DIR/deploy"
ENV_FILE="$DEPLOY_DIR/.env.production"
COMPOSE_FILE="$DEPLOY_DIR/compose.production.yml"

if [[ "${EUID}" -ne 0 ]]; then
  SUDO="sudo"
else
  SUDO=""
fi

install_git() {
  if command -v git >/dev/null 2>&1; then
    return
  fi

  if ! command -v apt-get >/dev/null 2>&1; then
    echo "Install Git and run this script again." >&2
    exit 1
  fi

  $SUDO apt-get update
  $SUDO apt-get install -y git
}

prepare_application_directory() {
  if [[ "$ROOT_DIR" == "$APP_DIR" ]]; then
    if [[ "$ACTION" == "--update" ]]; then
      install_git
      $SUDO git -C "$APP_DIR" pull --ff-only
    fi
    return
  fi

  install_git
  if [[ -d "$APP_DIR/.git" ]]; then
    $SUDO git -C "$APP_DIR" pull --ff-only
  elif [[ -e "$APP_DIR" ]]; then
    echo "The application directory already exists and is not a Git repository: $APP_DIR" >&2
    exit 1
  else
    $SUDO mkdir -p "$(dirname "$APP_DIR")"
    $SUDO git clone "$REPOSITORY_URL" "$APP_DIR"
  fi

  exec $SUDO "$APP_DIR/install.sh"
}

install_docker() {
  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    return
  fi

  if ! command -v apt-get >/dev/null 2>&1; then
    echo "Install Docker Engine and Docker Compose Plugin, then run this script again." >&2
    exit 1
  fi

  $SUDO apt-get update
  . /etc/os-release
  case "$ID" in
    ubuntu|debian) docker_os="$ID" ;;
    *) echo "This installer supports Ubuntu and Debian." >&2; exit 1 ;;
  esac

  $SUDO apt-get install -y ca-certificates curl gnupg
  $SUDO install -m 0755 -d /etc/apt/keyrings
  curl -fsSL "https://download.docker.com/linux/$docker_os/gpg" | $SUDO gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  $SUDO chmod a+r /etc/apt/keyrings/docker.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$docker_os $VERSION_CODENAME stable" | $SUDO tee /etc/apt/sources.list.d/docker.list >/dev/null
  $SUDO apt-get update
  $SUDO apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
}

ask() {
  local label="$1" default="${2:-}" value
  read -r -p "$label${default:+ [$default]}: " value </dev/tty
  printf '%s' "${value:-$default}"
}

validate_name() {
  [[ "$1" =~ ^[a-zA-Z0-9_.-]+$ ]] || { echo "Invalid value: $1" >&2; exit 1; }
}

configure() {
  if [[ -f "$ENV_FILE" ]]; then
    return
  fi

  local domain email app_password db_password mysql_root_password
  local mysql_network mysql_container db_name db_user
  domain="$(ask 'Application domain (example: app.example.com)')"
  [[ -n "$domain" ]] || { echo "The domain is required." >&2; exit 1; }
  email="$(ask 'Administrator email')"
  [[ -n "$email" ]] || { echo "The email is required." >&2; exit 1; }
  read -r -s -p 'Initial administrator password: ' app_password </dev/tty; echo
  [[ -n "$app_password" ]] || { echo "The password is required." >&2; exit 1; }

  mysql_network="$(ask 'Docker network of the existing MySQL' 'getfy_default')"
  mysql_container="$(ask 'Existing MySQL container' 'getfy-mysql-1')"
  db_name="$(ask 'Database name' 'chega_junto_concurseiro')"
  db_user="$(ask 'Database user' 'chega_junto_concurseiro_app')"
  validate_name "$mysql_network"
  validate_name "$mysql_container"
  validate_name "$db_name"
  validate_name "$db_user"

  docker network inspect "$mysql_network" >/dev/null || { echo "Docker network not found: $mysql_network" >&2; exit 1; }
  docker inspect "$mysql_container" >/dev/null || { echo "MySQL container not found: $mysql_container" >&2; exit 1; }
  read -r -s -p 'Root password of the existing MySQL: ' mysql_root_password </dev/tty; echo
  [[ -n "$mysql_root_password" ]] || { echo "The MySQL root password is required." >&2; exit 1; }

  db_password="$(openssl rand -hex 32)"
  docker exec -e "MYSQL_PWD=$mysql_root_password" "$mysql_container" mysql -uroot -e "CREATE DATABASE IF NOT EXISTS \`$db_name\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; CREATE USER IF NOT EXISTS '$db_user'@'%' IDENTIFIED BY '$db_password'; ALTER USER '$db_user'@'%' IDENTIFIED BY '$db_password'; GRANT ALL PRIVILEGES ON \`$db_name\`.* TO '$db_user'@'%'; FLUSH PRIVILEGES;"

  cat > "$ENV_FILE" <<EOF
MYSQL_DOCKER_NETWORK=$mysql_network
DB_HOST=$mysql_container
DB_PORT=3306
DB_NOME=$db_name
DB_USUARIO=$db_user
DB_SENHA=$db_password
MESTRE_EMAIL=$email
MESTRE_SENHA=$app_password
AMBIENTE=producao
FUSO_HORARIO=America/Fortaleza
ORIGEM_FRONTEND=https://$domain
RATE_GLOBAL_RPS=100
RATE_GLOBAL_BURST=500
RATE_LOGIN_RPS=10
RATE_LOGIN_BURST=50
EOF
  chmod 600 "$ENV_FILE"
}

prepare_application_directory
install_docker
configure
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --build --remove-orphans
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" ps
echo "Application started at http://IP_DA_VPS:8082"
echo "In Nginx Proxy Manager, forward the domain to IP_DA_VPS:8082 and enable SSL."
