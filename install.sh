#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$ROOT_DIR/deploy"
ENV_FILE="$DEPLOY_DIR/.env.production"
COMPOSE_FILE="$DEPLOY_DIR/compose.production.yml"

if [[ "${EUID}" -ne 0 ]]; then
  SUDO="sudo"
else
  SUDO=""
fi

install_docker() {
  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    return
  fi

  if ! command -v apt-get >/dev/null 2>&1; then
    echo "Instale Docker Engine e Docker Compose Plugin e execute este script novamente." >&2
    exit 1
  fi

  $SUDO apt-get update
  . /etc/os-release
  case "$ID" in
    ubuntu|debian) docker_os="$ID" ;;
    *)
      echo "Este instalador suporta Ubuntu e Debian." >&2
      exit 1
      ;;
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
  read -r -p "$label${default:+ [$default]}: " value
  printf '%s' "${value:-$default}"
}

configure() {
  if [[ -f "$ENV_FILE" ]]; then
    return
  fi

  local domain email app_password db_password root_password
  domain="$(ask 'Domínio da aplicação (ex.: app.exemplo.com)')"
  [[ -n "$domain" ]] || { echo 'O domínio é obrigatório.' >&2; exit 1; }
  email="$(ask 'E-mail do administrador')"
  [[ -n "$email" ]] || { echo 'O e-mail é obrigatório.' >&2; exit 1; }
  read -r -s -p 'Senha inicial do administrador: ' app_password; echo
  [[ -n "$app_password" ]] || { echo 'A senha é obrigatória.' >&2; exit 1; }
  db_password="$(openssl rand -hex 32)"
  root_password="$(openssl rand -hex 32)"

  cat > "$ENV_FILE" <<EOF
DOMAIN=$domain
DB_NOME=track_concursos
DB_USUARIO=track_app
DB_SENHA=$db_password
MYSQL_ROOT_PASSWORD=$root_password
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

install_docker
configure
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --build
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" ps
echo "Aplicação disponível em https://$(grep '^DOMAIN=' "$ENV_FILE" | cut -d= -f2)"
