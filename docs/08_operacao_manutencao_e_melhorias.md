# Operação e manutenção

Frontend (`web/client-react`): `npm run dev`, `npm run build`, `npm run lint`.

API (`web/server-go`): `go run ./cmd/api`, `go build ./...`, `go vet ./...`.

Faça backups pelo MySQL (por exemplo, `mysqldump`) antes de alterações. As migrações são executadas pela API. Use `GET /api/v1/saude`, os logs do processo Go e o console do navegador para diagnóstico.

Testes automatizados não fazem parte desta etapa; a validação atual é por build, lint, vet e testes manuais.

## Rotina de manutenção

- Conferir `GET /api/v1/saude` e a conexão com MySQL.
- Consultar logs da API quando migrações ou operações falharem.
- Executar build e lint do frontend após alterações React.
- Executar build e vet do backend após alterações Go.
- Fazer backup MySQL antes de alterar o esquema.
- Manter `ORIGEM_FRONTEND` igual à origem usada pelo Vite.

Em produção, use segredos próprios, não exponha o phpMyAdmin e publique a aplicação com HTTPS para proteger o cookie de sessão.

## Checklist de implantação

1. Criar o banco MySQL e usuário com permissões necessárias.
2. Definir todas as variáveis do `.env`.
3. Confirmar a origem do frontend e a porta da API.
4. Iniciar a API e verificar saúde e migrações.
5. Executar o build do frontend.
6. Criar ou confirmar a conta mestre.
7. Criar um mentor, um aluno e validar o vínculo.
8. Validar catálogo, edital, cronograma e login do aluno.

## Observabilidade

Registre o horário, usuário, endpoint e mensagem de erro ao investigar falhas. Não registre senha, token CSRF ou cookie. Em problemas de dados, reproduza primeiro com uma conta de teste e preserve o banco original para análise.

## Evolução

Alterações de domínio devem incluir migração, entidade, repositório, caso de uso, handler e integração no frontend. Atualize esta documentação junto com a rota ou regra alterada. Testes unitários permanecem fora do escopo atual, conforme definido para esta etapa.
