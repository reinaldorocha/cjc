# Autenticação e segurança

- Login em `POST /api/v1/autenticacao/entrar`.
- Sessão no servidor identificada pelo cookie HttpOnly `sessao_track`.
- Token CSRF no cabeçalho `X-Token-CSRF` para escritas.
- Senhas com bcrypt.
- Middleware verifica sessão, perfil (`mestre`, `mentor`, `aluno`) e vínculo do aluno.
- Data de expiração do plano bloqueia o aluno após o vencimento.
- Usuários são desativados logicamente.
- CORS usa a origem definida em `ORIGEM_FRONTEND`.

O frontend não deve persistir sessão ou dados de negócio em `localStorage`.

## Ciclo de sessão

Após o login, o servidor cria um registro em `sessoes_autenticacao` e envia o cookie de sessão. O frontend mantém o cookie automaticamente nas requisições. Ao sair, o registro é invalidado. A expiração da sessão é verificada no servidor.

## CSRF

O token CSRF é associado à sessão. O frontend o envia em requisições POST, PUT, PATCH e DELETE. Uma requisição sem token válido é recusada. Requisições GET são somente leitura conforme as rotas atuais.

## Autorização

O perfil mestre acessa a administração. O mentor acessa seus próprios alunos e seus catálogos. O aluno acessa o próprio identificador. Mesmo que um identificador seja conhecido, o middleware verifica o vínculo antes de executar a operação.

## Dados sensíveis

Senhas nunca são retornadas em respostas. Credenciais e segredos devem ficar somente no `.env` do ambiente. Em produção, use HTTPS, restrinja CORS e altere a senha mestre padrão.
