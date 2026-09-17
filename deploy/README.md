# Deploy em VPS

1. Aponte o DNS do domínio para o IP da VPS e libere as portas 80 e 443.
2. Instale Docker Engine e o plugin Docker Compose na VPS.
3. Envie este repositório para a VPS e execute:

```bash
cd deploy
cp .env.production.example .env.production
nano .env.production
docker compose --env-file .env.production -f compose.production.yml up -d --build
docker compose --env-file .env.production -f compose.production.yml ps
```

O Caddy emite e renova o certificado HTTPS automaticamente. A API e o MySQL não expõem portas públicas.

Para atualizar depois de enviar uma nova versão:

```bash
cd deploy
docker compose --env-file .env.production -f compose.production.yml up -d --build
```

Backup do banco:

```bash
docker compose --env-file .env.production -f compose.production.yml exec -T mysql sh -c 'mysqldump -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' > backup.sql
```
