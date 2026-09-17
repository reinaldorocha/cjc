# API Go — Track Concursos

Backend da aplicação web. A API usa MySQL como fonte de dados, autenticação por sessão em cookie `HttpOnly` e autorização para os papéis mestre, mentor e aluno.

## Estado atual

Implementado nesta etapa:

- configuração por `.env`;
- criação automática do banco MySQL;
- migrações versionadas;
- usuário mestre inicial;
- entrada, saída, sessão e alteração de senha;
- token CSRF para requisições de escrita;
- criação, desativação e reativação de mentor e aluno pelo mestre;
- vínculo exclusivo mentor–aluno;
- data de expiração do plano e bloqueio do aluno quando expirada;
- listagem e configuração de alunos pelo mentor;
- cadastro transacional de aluno pelo mentor, com vínculo automático;
- catálogos de concursos, editais e flashcards independentes de aluno;
- concursos atribuídos e desativação lógica;
- logotipo, ordenação transacional e resultados completos dos concursos atribuídos;
- editais reutilizáveis com matérias, tópicos e subtópicos;
- estrutura do edital editável somente pelo mentor;
- progresso individual editável pelo aluno ou mentor;
- materiais em matérias, tópicos e subtópicos, com conclusão individual por aluno;
- cronograma manual criado pelo mentor;
- cronograma agendado e ciclo inteligente gerados no backend;
- permissão individual para o aluno gerar o próprio cronograma;
- substituição versionada do cronograma ativo;
- calendário e execução dos itens pelo aluno;
- baralhos pessoais e do mentor, com alcance por aluno, edital ou global;
- subbaralhos e cartões completos, incluindo dicas, etiquetas, alternativas e explicações;
- repetição espaçada individual com histórico e fila de cartões pendentes;
- sessões de estudo, métricas complementares e registros de questões;
- revisões programadas com ciclos, adiamento e conclusão;
- configuração completa de prova por aluno e concurso;
- simulados pendentes ou realizados, com resultado geral e detalhamento por matéria;
- métricas agregadas de estudo, questões, sequência, simulados e cobertura do edital;
- linha do tempo diária e indicadores consolidados por matéria;
- endpoint de saúde.

O frontend React consome estes domínios exclusivamente pela API. Testes unitários estão intencionalmente adiados até aprovação do funcionamento.

## Requisitos

- Go 1.24 ou superior;
- MySQL 8;
- usuário MySQL com permissão para criar o banco configurado.

## Configuração

Copie `.env.example` para `.env` e preencha as credenciais:

```powershell
Copy-Item .env.example .env
```

Defina uma senha forte para `MESTRE_SENHA`. Na primeira inicialização, o usuário mestre é criado automaticamente se ainda não existir outro mestre no banco.

## Execução

```powershell
go mod tidy
go run ./cmd/api
```

A API responde em `http://127.0.0.1:8080` por padrão.

## Endpoints disponíveis

| Método | Endpoint |
|---|---|
| `GET` | `/api/v1/saude` |
| `POST` | `/api/v1/autenticacao/entrar` |
| `POST` | `/api/v1/autenticacao/sair` |
| `POST` | `/api/v1/autenticacao/sair-de-todas` |
| `GET` | `/api/v1/autenticacao/eu` |
| `POST` | `/api/v1/autenticacao/alterar-senha` |
| `GET`, `POST` | `/api/v1/administracao/usuarios` |
| `PATCH` | `/api/v1/administracao/usuarios/{usuarioId}` |
| `POST` | `/api/v1/administracao/usuarios/{usuarioId}/reativar` |
| `POST` | `/api/v1/administracao/usuarios/{usuarioId}/redefinir-senha` |
| `PUT` | `/api/v1/administracao/alunos/{alunoId}/mentor` |
| `GET` | `/api/v1/mentor/alunos` |
| `POST` | `/api/v1/mentor/alunos` |
| `PATCH` | `/api/v1/mentor/alunos/{alunoId}/configuracoes` |
| `GET`, `POST` | `/api/v1/alunos/{alunoId}/concursos` |
| `PATCH`, `DELETE` | `/api/v1/alunos/{alunoId}/concursos/{concursoId}` |
| `GET`, `PUT` | `/api/v1/alunos/{alunoId}/edital` |
| `GET` | `/api/v1/editais?concursoId={concursoId}` |
| `POST` | `/api/v1/alunos/{alunoId}/edital/itens` |
| `PUT` | `/api/v1/alunos/{alunoId}/edital/ordem` |
| `PATCH` | `/api/v1/alunos/{alunoId}/edital/itens/{itemId}` |
| `PATCH` | `/api/v1/alunos/{alunoId}/edital/progresso/{itemId}` |
| `PATCH` | `/api/v1/alunos/{alunoId}/edital/materiais/{materialId}/progresso` |
| `GET`, `PUT` | `/api/v1/alunos/{alunoId}/cronograma` |
| `POST` | `/api/v1/alunos/{alunoId}/cronograma/gerar` |
| `PATCH` | `/api/v1/alunos/{alunoId}/cronograma/itens/{itemId}` |
| `GET` | `/api/v1/alunos/{alunoId}/cronograma/calendario` |
| `GET`, `POST` | `/api/v1/alunos/{alunoId}/baralhos-cartoes` |
| `PATCH`, `DELETE` | `/api/v1/alunos/{alunoId}/baralhos-cartoes/{baralhoId}` |
| `POST` | `/api/v1/alunos/{alunoId}/baralhos-cartoes/{baralhoId}/cartoes` |
| `PATCH`, `DELETE` | `/api/v1/alunos/{alunoId}/cartoes/{cartaoId}` |
| `POST` | `/api/v1/alunos/{alunoId}/cartoes/{cartaoId}/revisar` |
| `GET` | `/api/v1/alunos/{alunoId}/cartoes/pendentes` |
| `GET` | `/api/v1/alunos/{alunoId}/cartoes/historico` |
| `GET`, `POST` | `/api/v1/alunos/{alunoId}/sessoes-estudo` |
| `PATCH`, `DELETE` | `/api/v1/alunos/{alunoId}/sessoes-estudo/{sessaoId}` |
| `GET`, `POST` | `/api/v1/alunos/{alunoId}/registros-questoes` |
| `PATCH`, `DELETE` | `/api/v1/alunos/{alunoId}/registros-questoes/{registroId}` |
| `GET`, `POST` | `/api/v1/alunos/{alunoId}/revisoes` |
| `PATCH`, `DELETE` | `/api/v1/alunos/{alunoId}/revisoes/{revisaoId}` |
| `GET`, `PUT` | `/api/v1/alunos/{alunoId}/configuracao-prova/{concursoId}` |
| `GET`, `POST` | `/api/v1/alunos/{alunoId}/simulados` |
| `PUT`, `DELETE` | `/api/v1/alunos/{alunoId}/simulados/{simuladoId}` |
| `GET` | `/api/v1/alunos/{alunoId}/metricas/resumo` |
| `GET` | `/api/v1/alunos/{alunoId}/metricas/linha-do-tempo` |
| `GET` | `/api/v1/alunos/{alunoId}/metricas/materias` |

As requisições de escrita autenticadas devem enviar `X-Token-CSRF` com o valor retornado pela entrada no sistema.
