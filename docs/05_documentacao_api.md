# API HTTP

Base: `/api/v1`. A API é implementada em `web/server-go/internal/transporte/servidor.go`.

## Grupos de rotas

- `GET /saude`
- `/autenticacao`: entrar, sair, sair de todas, usuário atual, alterar senha e perfil
- `/administracao`: usuários e vínculo mentor/aluno
- `/mentor`: alunos, radar, concursos, editais, baralhos/cartões, materiais e identidade visual
- `/alunos/{alunoId}`: concursos, edital, cronograma, cartões, questões, materiais, cadernos, sessões, registros, revisões, simulados, métricas e configuração de prova
- `GET /editais`: catálogo de editais
- `/mentores/{mentorId}/whitelabel` e `/alunos/{alunoId}/whitelabel`: identidade visual

Rotas protegidas exigem sessão. Operações de escrita exigem `X-Token-CSRF`. Métodos, parâmetros e permissões estão definidos nos handlers e no roteador.

## Autenticação

`POST /autenticacao/entrar` cria a sessão. `GET /autenticacao/eu` retorna o usuário atual. `POST /autenticacao/sair` encerra a sessão atual e `POST /autenticacao/sair-de-todas` encerra todas as sessões. Há endpoints para alteração de senha e perfil.

## Mentor e aluno

O mentor mantém alunos, concursos, editais, cartões, materiais e configurações em `/mentor`. As rotas `/alunos/{alunoId}` validam o próprio aluno ou o mentor responsável e cobrem edital, cronograma, cartões, questões, materiais, cadernos, sessões, revisões, simulados e métricas. O aluno não possui rota para alterar a estrutura do edital.

Os handlers retornam JSON e rejeitam autenticação, vínculo, método ou dados inválidos conforme a regra aplicável.

## Endpoints de conteúdo

| Grupo | Operações disponíveis |
|---|---|
| Concursos | listar, cadastrar, ordenar e atualizar atribuições |
| Edital | consultar catálogo, atribuir, listar itens, ordenar e atualizar progresso |
| Cronograma | consultar, substituir, gerar, criar item, reprogramar e consultar calendário |
| Cartões | listar baralhos, criar/editar/excluir cartões e revisar cartões pendentes |
| Questões | listar banco, responder e consultar estatísticas |
| Estudo | criar e acompanhar sessões, registros de questões e revisões |
| Simulados | criar, consultar, atualizar e excluir simulados |
| Métricas | resumo, linha do tempo e consolidação por matéria |

Os identificadores são enviados nos caminhos conforme o roteador. O cliente deve tratar respostas não autorizadas encerrando o contexto local e consultando novamente `/autenticacao/eu`.

## Rotas completas por área

### Administração

`GET/POST /administracao/usuarios`, `PATCH /administracao/usuarios/{usuarioId}`, `POST /administracao/usuarios/{usuarioId}/reativar`, `POST /administracao/usuarios/{usuarioId}/redefinir-senha` e `PUT /administracao/alunos/{alunoId}/mentor`.

### Catálogo do mentor

`GET/POST /mentor/concursos`, `GET/POST /mentor/editais`, `GET/POST/PATCH/DELETE /mentor/baralhos-cartoes`, `POST /mentor/baralhos-cartoes/{baralhoId}/cartoes`, `PATCH/DELETE /mentor/cartoes/{cartaoId}`, `GET/POST/PATCH/DELETE /mentor/materiais-apoio` e `GET /mentor/materiais-apoio/{materialId}/arquivo`.

### Gestão de alunos

`GET/POST /mentor/alunos`, `GET /mentor/radar-alunos`, `GET /mentor/alunos/{alunoId}`, `GET /mentor/alunos/{alunoId}/visao-geral` e `PATCH /mentor/alunos/{alunoId}/configuracoes`.

### Edital e concurso do aluno

`GET/POST /alunos/{alunoId}/concursos`, `PUT /alunos/{alunoId}/concursos/ordem`, `PATCH/DELETE /alunos/{alunoId}/concursos/{concursoId}`, `GET /alunos/{alunoId}/edital`, `PUT /alunos/{alunoId}/edital`, `POST /alunos/{alunoId}/edital/itens`, `PUT /alunos/{alunoId}/edital/ordem` e `PATCH /alunos/{alunoId}/edital/itens/{itemId}`.

### Cronograma

`GET/PUT /alunos/{alunoId}/cronograma`, `POST /alunos/{alunoId}/cronograma/gerar`, `POST /alunos/{alunoId}/cronograma/itens`, `POST /alunos/{alunoId}/cronograma/reprogramar`, `PATCH /alunos/{alunoId}/cronograma/itens/{itemId}` e `GET /alunos/{alunoId}/cronograma/calendario`.

### Cartões e desempenho

`GET/POST/PATCH/DELETE /alunos/{alunoId}/baralhos-cartoes`, `POST /alunos/{alunoId}/cartoes/{cartaoId}/revisar`, `GET /alunos/{alunoId}/cartoes/pendentes` e `GET /alunos/{alunoId}/cartoes/historico`.

### Banco de questões, estudo e simulados

`GET /banco-questoes`, operações de manutenção do banco para mentor, `POST /alunos/{alunoId}/banco-questoes/responder`, `GET /alunos/{alunoId}/banco-questoes/estatisticas`; operações CRUD de `/sessoes-estudo`, `/registros-questoes` e `/revisoes`; `GET/PUT /configuracao-prova/{concursoId}`; e `GET/POST /simulados`, `PUT/DELETE /simulados/{simuladoId}`.

### Métricas e identidade visual

`GET /alunos/{alunoId}/metricas/resumo`, `/linha-do-tempo`, `/materias`; `GET/POST /mentor/whitelabel`; `GET /mentores/{mentorId}/whitelabel`; `GET /alunos/{alunoId}/whitelabel`.
