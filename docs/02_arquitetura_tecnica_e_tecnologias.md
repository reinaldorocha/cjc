# Arquitetura e tecnologias

- `web/client-react`: frontend React 19 com TypeScript, Vite, React Router e Tailwind.
- `web/server-go`: API Go (módulo Go 1.24), organizada em domínio, aplicação, banco, autenticação e transporte HTTP.
- `web/server-go/migrations`: migrações SQL versionadas.

O frontend chama a API pelo prefixo `/api`; no desenvolvimento o Vite faz o proxy para a API. A autenticação usa cookie de sessão HttpOnly e o cabeçalho CSRF `X-Token-CSRF` nas escritas. O único banco suportado é MySQL 8.

## Organização do código

O comando da API fica em `server-go/cmd/api`. O backend separa domínio, casos de uso, repositórios MySQL, autenticação e transporte HTTP. O roteador central registra as rotas públicas, administrativas, de mentor e de aluno. As migrações SQL são aplicadas na inicialização.

No frontend, páginas ficam em `client-react/src/pages`, componentes compartilhados em `src/components` e chamadas HTTP em `src/services`. O estado de sessão é obtido pela API; não há persistência de dados de negócio no navegador.

## Portas e proxy

O Vite escuta em `127.0.0.1:5174` na configuração local atual e encaminha `/api` para `127.0.0.1:8082`. Esses valores não são fixos para produção: devem ser mantidos sincronizados com `PORTA` e `ORIGEM_FRONTEND`.

## Camada HTTP

A API configura CORS, cookies, limite de requisição e roteamento. Handlers convertem JSON em comandos de aplicação e convertem erros de domínio em respostas HTTP. A autorização é aplicada antes dos repositórios para impedir acesso cruzado entre alunos.

## Camada de dados

O repositório abre a conexão MySQL, configura o pool e executa migrações pendentes. Transações são utilizadas nas operações que alteram mais de uma tabela, como atribuição de edital, geração de cronograma e registro de respostas.
