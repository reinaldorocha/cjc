# Instalação e configuração

Dependências: Docker Compose, Go 1.24 e Node.js com npm.

Na raiz: `docker compose up -d mysql_db phpmyadmin_web`. O MySQL fica na porta 3306 e o phpMyAdmin em `http://localhost:8080`.

```bash
cd web/server-go
copy .env.example .env
go run ./cmd/api
```

Configure `PORTA`, `DB_HOST`, `DB_PORT`, `DB_USUARIO`, `DB_SENHA`, `DB_NOME`, `MESTRE_EMAIL`, `MESTRE_SENHA` e `ORIGEM_FRONTEND`. A configuração local atual usa API em 8082 e frontend em 5174.

```bash
cd web/client-react
npm install
npm run dev
```

Verifique `GET /api/v1/saude`.

## Variáveis de ambiente

| Variável | Finalidade |
|---|---|
| `AMBIENTE` | ambiente de execução |
| `PORTA` | porta HTTP da API |
| `FUSO_HORARIO` | fuso de datas e cronogramas |
| `DB_HOST`, `DB_PORT` | endereço e porta do MySQL |
| `DB_USUARIO`, `DB_SENHA`, `DB_NOME` | credenciais e banco |
| `MESTRE_EMAIL`, `MESTRE_SENHA` | conta mestre inicial |
| `ORIGEM_FRONTEND` | origem permitida pelo CORS |

Não versione `.env`; use `.env.example` como modelo. Na primeira execução, aguarde a aplicação das migrações e entre com a conta mestre configurada.

## Docker Compose

O serviço `mysql_db` utiliza MySQL 8, publica a porta 3306 e cria o banco `chega_junto_concurseiro`. O serviço `phpmyadmin_web` utiliza o MySQL como servidor e publica a porta 8080. Para parar os serviços sem remover dados:

```bash
docker compose stop
```

Para iniciar novamente:

```bash
docker compose start
```

## Execução do frontend

O frontend deve ser executado a partir de `web/client-react`. O comando de desenvolvimento mantém atualização automática. Para uma verificação de produção, `npm run build` gera os arquivos estáticos e reporta erros de TypeScript e bundling.

## Problemas comuns

- **API indisponível**: confirme que a porta configurada no Vite é a mesma porta da API.
- **Erro de conexão**: confirme se o container MySQL está saudável, o host é `127.0.0.1` e a senha coincide com o Compose.
- **CORS**: ajuste `ORIGEM_FRONTEND` para a origem exata, incluindo porta.
- **Login inicial**: confira `MESTRE_EMAIL` e `MESTRE_SENHA` antes da primeira inicialização.
