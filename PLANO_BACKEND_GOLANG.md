# Plano de implementação — Backend Go, autenticação e perfis Mentor/Aluno

## 1. Objetivo

Substituir a persistência atual por uma API em Go com banco de dados MySQL, autenticação e autorização por usuário.

O sistema terá três tipos de acesso:

- **Mestre:** cria e administra contas de mentores e alunos.
- **Mentor:** também pode cadastrar alunos, que ficam automaticamente vinculados à sua mentoria.
- **Mentor:** acompanha cada aluno individualmente e pode criar, consultar, alterar e excluir os dados desse aluno.
- **Aluno:** acessa os próprios dados, segue o cronograma criado pelo mentor e possui permissões limitadas conforme as regras deste documento.

Depois da migração, o frontend não utilizará `localStorage` como banco de dados. Dados de domínio, sessão e configurações serão carregados e persistidos por endpoints autenticados.

## Andamento da implementação

| Fase | Situação | Observação |
|---|---|---|
| 1 — fundação do backend | concluída | módulo Go, MySQL, migrações, configuração, endpoint de saúde e estrutura modular revisada |
| 2 — autenticação e papéis | concluída | serviços separados de autenticação, usuários, mentoria e auditoria; transporte HTTP isolado |
| 3 — concursos e edital | concluída | serviços, escopo por aluno, atribuição reutilizável, estrutura do edital e progresso individual implementados |
| 4 — cronogramas | concluída | versões, agenda, ciclo inteligente, calendário, execução e permissão individual implementados |
| 5 — cartões de estudo | concluída | baralhos hierárquicos, escopos, cartões completos, revisão espaçada e histórico individual implementados |
| 6 — estudos, revisões e simulados | concluída | sessões, questões, revisões, configuração de prova e simulados transacionais implementados |
| 7 — métricas | concluída | resumo, séries temporais, cobertura, simulados e indicadores por matéria implementados |
| 8 — autenticação e mentor no frontend | concluída | sessão, entrada, proteção por papel, administração mestre, seleção de aluno, configurações e acesso individual a todas as telas do aluno implementados |
| 9 — migração das telas | concluída | todas as 9 telas migradas para a API Go sem uso de localStorage |
| 10 — corte definitivo | concluída | persistência antiga, perfis, backups, snapshots e endpoints temporários removidos; documentação atualizada |
| testes unitários | adiado | não serão criados até aprovação explícita do funcionamento |

## 2. Regras de acesso confirmadas

### 2.1 Mestre

O usuário mestre poderá:

- criar contas de mentor e aluno;
- editar e desativar contas;
- definir o mentor único de cada aluno;
- reativar contas desativadas;
- consultar a situação de acesso e o vínculo das contas;
- redefinir credenciais administrativamente quando necessário.

Não haverá cadastro público. A recuperação automática de senha por e-mail ficará para uma fase posterior.

### 2.2 Mentor

O mentor poderá:

- cadastrar alunos com vínculo automático à própria mentoria;
- cadastrar concursos e editais no catálogo da mentoria sem depender de aluno;
- criar flashcards globais ou vinculados a edital sem depender de aluno;
- atribuir concursos e editais já cadastrados aos alunos;

- visualizar a própria lista de alunos;
- selecionar e acompanhar um aluno individualmente;
- acessar todas as telas e todos os dados do aluno selecionado;
- criar, alterar e excluir qualquer dado pertencente ao aluno;
- atribuir concursos e editais ao aluno;
- montar o cronograma de estudos do aluno;
- habilitar ou desabilitar, por aluno, a criação do cronograma inteligente;
- preencher e alterar a data de expiração do plano do aluno;
- criar cartões de estudo disponibilizados ao aluno;
- consultar métricas, histórico, revisões e simulados do aluno.

### 2.3 Aluno

O aluno poderá:

- acessar o cronograma criado pelo mentor;
- usar o cronograma inteligente somente quando estiver habilitado pelo mentor;
- criar, alterar, estudar e excluir os próprios cartões de estudo;
- consultar e estudar cartões de estudo criados pelo mentor;
- acessar o edital atribuído pelo mentor apenas para consulta e execução do estudo;
- acessar dashboard, revisões, simulados, histórico, métricas e demais telas liberadas;
- registrar suas sessões, questões, revisões, simulados e progresso.
- alterar o próprio nome e a própria senha.

O aluno não poderá:

- criar, substituir, editar, reordenar ou excluir o edital atribuído;
- editar ou excluir cartões de estudo criados pelo mentor;
- alterar a configuração que habilita o cronograma inteligente;
- acessar dados de outro aluno;
- alterar dados administrativos do vínculo com o mentor.

### 2.4 Expiração do plano

- a data de expiração é definida pelo mentor no cadastro do aluno vinculado;
- quando a data atual for igual ou posterior à data de expiração, o aluno perde o acesso ao sistema;
- a verificação ocorre no login e em toda requisição autenticada do aluno;
- sessões já abertas são bloqueadas assim que a data expira;
- o mentor continua podendo acessar o histórico do aluno expirado e alterar a data para restabelecer o acesso;
- o usuário mestre também pode consultar e administrar a situação da conta;
- a expiração não apaga nenhum dado;
- a comparação usa a data civil no fuso horário configurado no servidor;
- a API retorna um erro específico `PLANO_EXPIRADO`, e o frontend exibe uma tela informando que o acesso expirou.

### 2.5 Distinção importante no edital

O conteúdo estrutural do edital é controlado pelo mentor:

- matérias;
- tópicos;
- subtópicos;
- ordem;
- pesos e relevâncias;
- materiais e metadados estruturais.

O aluno poderá registrar dados de execução associados ao edital sem alterar sua estrutura, como:

- item estudado;
- tempo de estudo;
- questões e acertos;
- revisões;
- progresso individual.

Para garantir essa separação, o conteúdo do edital e o progresso do aluno serão armazenados em tabelas diferentes.

## 3. Decisões técnicas propostas

Estas decisões devem ser confirmadas durante o início da implementação, mas formam a base recomendada do projeto:

- Go 1.24 ou versão estável disponível no ambiente.
- API HTTP REST em JSON.
- MySQL 8 com `utf8mb4`.
- Migrações de banco versionadas.
- Senhas armazenadas somente como hash Argon2id ou bcrypt.
- Autenticação por cookie de sessão `HttpOnly`, `Secure` em produção e `SameSite=Lax`.
- Token CSRF para operações que alteram dados.
- Sessões revogáveis armazenadas no banco.
- Identificadores UUID.
- Validação de entrada no backend.
- Consultas sempre limitadas pelo usuário autenticado e pelo vínculo mentor–aluno.
- Datas e horários persistidos em UTC e formatados no frontend.
- Logs estruturados sem senhas, tokens ou conteúdo sensível.

O uso de cookie `HttpOnly` é preferível a armazenar tokens no navegador, pois impede que o JavaScript leia a credencial de sessão.

## 4. Arquitetura proposta

```text
web/
├── client/                         # build gerado do React
├── client-react/                   # frontend React + TypeScript
├── server-go/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── autenticacao/           # entrada, saída, sessão e senha
│   │   ├── auditoria/              # registro de ações sensíveis
│   │   ├── configuracao/           # variáveis de ambiente
│   │   ├── banco/                  # conexão e transações
│   │   ├── identificador/          # UUID, segredos e hashes técnicos
│   │   ├── modelos/                # contratos compartilhados entre módulos
│   │   ├── transporte/             # rotas, controladores e intermediários HTTP
│   │   ├── usuarios/               # usuários e papéis
│   │   ├── mentoria/               # vínculos mentor–aluno
│   │   ├── concursos/              # concursos e atribuições
│   │   ├── editais/                # edital e progresso individual
│   │   ├── cronogramas/            # cronogramas e configuração inteligente
│   │   ├── cartoes/                # baralhos, cartões e estudo
│   │   ├── estudos/                # sessões, questões e revisões
│   │   ├── simulados/              # simulados e regras de prova
│   │   └── metricas/               # consultas agregadas
│   ├── migracoes/                  # SQL versionado
│   ├── testes/
│   ├── .env.example
│   ├── go.mod
│   └── go.sum
├── README.md
└── PLANO_BACKEND_GOLANG.md
```

Fluxo geral:

```text
React → API Go → autenticação/autorização → serviço de domínio → MySQL
```

## 5. Modelo de dados MySQL

### 5.1 Identidade e acesso

#### `usuarios`

- `id` UUID, chave primária;
- `nome`;
- `email`, único e normalizado;
- `senha_hash`;
- `papel`: `mestre`, `mentor` ou `aluno`;
- `ativo`;
- `criado_em`;
- `atualizado_em`;
- `ultimo_acesso_em`.

Contas não são excluídas fisicamente. A desativação preenche `desativado_em` e mantém o histórico.

#### `sessoes_autenticacao`

- `id` UUID;
- `usuario_id`;
- `token_hash`;
- `expira_em`;
- `criado_em`;
- `ultimo_uso_em`;
- `revogado_em`;
- informações limitadas de dispositivo, se necessárias para auditoria.

O token puro nunca será salvo no banco; apenas seu hash.

#### `mentor_alunos`

- `id` UUID;
- `mentor_id`;
- `aluno_id`;
- `situacao`: ativo ou inativo;
- `permite_cronograma_inteligente` booleano;
- `data_expiracao_plano` do tipo `DATE`;
- `criado_em`;
- `atualizado_em`.

Restrições:

- mentor e aluno devem possuir os papéis correspondentes;
- cada aluno pode possuir somente um vínculo ativo com mentor;
- todas as operações do mentor sobre um aluno exigem vínculo ativo.

Um índice único condicional será representado de forma compatível com MySQL pela regra transacional do serviço e por uma coluna auxiliar/indexada que impeça dois vínculos ativos para o mesmo `aluno_id`.

### 5.2 Concursos e atribuição

#### `concursos`

- dados gerais do concurso;
- banca, cargo, salário e data da prova;
- situação de pré-edital;
- resultado, classificação, nota e nomeação;
- imagem ou referência da imagem;
- datas de criação e atualização.

#### `aluno_concursos`

Relaciona o concurso ao aluno e armazena informações individuais:

- `aluno_id`;
- `concurso_id`;
- `atribuido_por`;
- grupo: foco, mira ou realizado;
- ordem de exibição;
- inclusão nas estatísticas;
- concurso ativo;
- datas de atribuição e atualização.

O aluno consulta os concursos atribuídos. A atribuição e a remoção são ações do mentor.

### 5.3 Edital e progresso

#### `editais`

- `id`;
- `concurso_id`;
- nome e versão;
- `criado_por`;
- datas de criação e atualização.

Um edital é uma entidade reutilizável. O mesmo `edital_id` pode ser atribuído a vários alunos; cada aluno mantém progresso e cronograma próprios sem duplicar a estrutura.

#### `edital_materias`

- `id`;
- `edital_id`;
- nome;
- ordem;
- peso, relevância e metadados.

#### `edital_topicos`

- `id`;
- `materia_id`;
- nome;
- ordem;
- metadados estruturais.

#### `edital_subtopicos`

- `id`;
- `topico_id`;
- nome;
- ordem;
- metadados estruturais.

#### `aluno_editais`

- `aluno_id`;
- `edital_id`;
- `atribuido_por`;
- `atribuido_em`;
- `ativo`.

#### `aluno_progresso_edital`

- `aluno_id`;
- tipo do item: tópico ou subtópico;
- `item_id`;
- estudado;
- data de conclusão;
- observações pessoais permitidas;
- datas de criação e atualização.

O backend não aceitará escrita estrutural do edital em uma requisição autenticada como aluno, mesmo que o frontend seja manipulado.

### 5.4 Cronogramas

#### `cronogramas`

- `id`;
- `aluno_id`;
- `concurso_id`;
- tipo: manual, agendado ou ciclo inteligente;
- `criado_por`;
- estado e versão;
- configuração normalizada em colunas e JSON somente onde a estrutura for variável;
- datas de criação e atualização.

#### `cronograma_itens`

- `id`;
- `cronograma_id`;
- data ou posição no ciclo;
- matéria, tópico ou subtópico;
- duração planejada;
- prioridade;
- status;
- ordem;
- datas de criação e atualização.

#### `cronograma_resumos_diarios`

- instantâneos diários necessários para preservar a agenda já apresentada;
- aluno, cronograma, data e conteúdo calculado.

Autorização:

- mentor sempre pode montar e alterar o cronograma de aluno vinculado;
- aluno pode consultar e executar;
- aluno só pode criar ou reconfigurar o cronograma inteligente se `permite_cronograma_inteligente = TRUE`;
- a permissão não autoriza alteração do edital.

### 5.5 Cartões de estudo

#### `baralhos_cartoes`

- `id`;
- `aluno_id`;
- `concurso_id` opcional;
- `criado_por`;
- `tipo_proprietario`: `mentor` ou `aluno`;
- `alcance`: `aluno`, `edital` ou `global`;
- `edital_id` obrigatório quando o alcance for `edital`;
- nome e descrição;
- datas de criação e atualização.

#### `cartoes_estudo`

- `id`;
- `baralho_id`;
- tipo: básico, lacuna ou múltipla escolha;
- conteúdo do cartão;
- referência opcional ao edital;
- `criado_por`;
- datas de criação e atualização.

#### `revisoes_cartoes`

- `id`;
- `aluno_id`;
- `cartao_id`;
- resultado da resposta;
- estado do algoritmo de repetição;
- próxima revisão;
- data da tentativa.

Regras:

- cartão criado pelo aluno é pessoal e editável pelo aluno e pelo mentor vinculado;
- cartão criado pelo mentor pode ser destinado a um aluno, a todos os alunos de um edital ou globalmente a todos os alunos do sistema;
- cartão criado pelo mentor é somente leitura para os alunos alcançados;
- alcance por edital disponibiliza o cartão aos alunos que tenham aquele edital atribuído;
- alcance global independe do edital atribuído;
- o histórico de estudo pertence sempre ao aluno.

### 5.6 Desativação lógica

Entidades de negócio não serão apagadas definitivamente. As tabelas aplicáveis terão:

- `ativo`;
- `desativado_em`;
- `desativado_por`.

Operações `DELETE` da API executarão desativação lógica. Consultas comuns ignorarão registros desativados, enquanto consultas administrativas poderão incluí-los explicitamente. Chaves estrangeiras não usarão exclusão em cascata para dados históricos.

### 5.7 Estudos, questões e revisões

#### `sessoes_estudo`

- aluno, concurso, matéria, tópico e subtópico;
- duração, modo, observações e data.

#### `registros_questoes`

- aluno, concurso e conteúdo associado;
- resolvidas, acertos, erros e data.

#### `revisoes_programadas`

- aluno e conteúdo associado;
- ciclo atual;
- próxima data;
- percentual anterior;
- estado ativo/concluído.

### 5.8 Simulados

#### `configuracoes_prova`

- configuração de prova por aluno e concurso;
- regra de pontuação;
- formato e valor do corte;
- configuração de blocos, matérias, pesos, mínimos e redações.

#### `simulados`

- aluno e concurso;
- nome, tipo, data, link e observações;
- percentual, pontos, tempo e aprovação;
- cópia da regra usada no momento do resultado.

#### `resultados_materias_simulado`

- resultado por matéria ou bloco;
- questões, acertos, erros, brancos, percentual e pontos.

Manter uma cópia da regra evita que a edição futura da configuração altere resultados históricos.

### 5.9 Auditoria

#### `eventos_auditoria`

Registrar ações sensíveis:

- entrada e saída do sistema;
- criação e remoção de vínculo;
- alteração de permissão do cronograma inteligente;
- atribuição ou mudança de edital;
- alteração feita pelo mentor em dados do aluno;
- restaurações ou importações administrativas.

Campos mínimos:

- usuário executor;
- aluno afetado, quando aplicável;
- ação;
- entidade e identificador;
- data;
- metadados sem conteúdo secreto.

## 6. Autenticação e autorização

### 6.1 Fluxo de entrada

1. Usuário envia e-mail e senha.
2. Backend busca o usuário pelo e-mail normalizado.
3. Backend compara a senha com o hash.
4. Backend cria uma sessão revogável com expiração.
5. Identificador secreto da sessão é enviado em cookie `HttpOnly`.
6. Frontend consulta `/api/v1/autenticacao/eu` para carregar usuário, papel e permissões.

Para alunos, a criação da sessão exige conta ativa, vínculo ativo e plano ainda válido. A mesma validação será repetida durante o uso da sessão.

### 6.2 Middleware

Middleware necessários:

- identificação da sessão;
- exigência de autenticação;
- exigência de papel mestre;
- exigência de papel mentor;
- resolução do aluno atual;
- confirmação do vínculo mentor–aluno;
- proteção CSRF;
- limite de tentativas na entrada;
- validação de conta ativa, vínculo e data de expiração do plano;
- log de requisição com identificador de correlação;
- tratamento uniforme de erros.

### 6.3 Regra de escopo

O identificador do aluno enviado pela URL nunca será suficiente para autorizar a operação.

Para cada requisição:

- aluno: o `aluno_id` efetivo vem da sessão;
- mentor: o backend confirma o vínculo antes de usar o `aluno_id` solicitado;
- consultas e alterações incluem o `aluno_id` autorizado na condição SQL;
- respostas não expõem registros fora do escopo.

O mestre atua apenas nos endpoints administrativos. O papel mestre não será tratado implicitamente como mentor de todos os alunos nas rotas acadêmicas.

## 7. Contrato inicial da API

Prefixo: `/api/v1`.

### 7.1 Autenticação

| Método | Endpoint | Acesso | Finalidade |
|---|---|---|---|
| `POST` | `/autenticacao/entrar` | público | iniciar sessão |
| `POST` | `/autenticacao/sair` | autenticado | encerrar sessão atual |
| `POST` | `/autenticacao/sair-de-todas` | autenticado | revogar todas as sessões |
| `GET` | `/autenticacao/eu` | autenticado | usuário, papel e permissões |
| `POST` | `/autenticacao/alterar-senha` | autenticado | alterar a própria senha |

Não haverá cadastro público. Contas de mentor e aluno serão criadas exclusivamente pelo usuário mestre.

Não haverá recuperação automática por e-mail nesta etapa. O mestre poderá redefinir a senha administrativamente, e o usuário autenticado poderá alterar a própria senha.

### 7.2 Administração de contas

| Método | Endpoint | Acesso | Finalidade |
|---|---|---|---|
| `GET` | `/administracao/usuarios` | mestre | listar contas e situações |
| `POST` | `/administracao/usuarios` | mestre | criar mentor ou aluno |
| `PATCH` | `/administracao/usuarios/{usuarioId}` | mestre | editar ou desativar conta |
| `POST` | `/administracao/usuarios/{usuarioId}/reativar` | mestre | reativar conta |
| `POST` | `/administracao/usuarios/{usuarioId}/redefinir-senha` | mestre | definir senha temporária |
| `PUT` | `/administracao/alunos/{alunoId}/mentor` | mestre | definir ou trocar o mentor único |

### 7.3 Alunos do mentor

| Método | Endpoint | Acesso | Finalidade |
|---|---|---|---|
| `GET` | `/mentor/alunos` | mentor | listar alunos vinculados |
| `GET` | `/mentor/alunos/{alunoId}` | mentor | consultar aluno |
| `PATCH` | `/mentor/alunos/{alunoId}/configuracoes` | mentor | alterar permissões do aluno |
| `GET` | `/mentor/alunos/{alunoId}/visao-geral` | mentor | visão agregada do aluno |

`/mentor/alunos/{alunoId}/configuracoes` permite alterar `permite_cronograma_inteligente` e `data_expiracao_plano`.

### 7.4 Contexto do aluno

Os recursos usam um contexto único:

- aluno chama `/alunos/eu/...`;
- mentor chama `/alunos/{alunoId}/...`.

Os controladores podem compartilhar a mesma camada de serviço depois que o intermediário resolver o aluno autorizado.

### 7.5 Concursos e edital

| Método | Endpoint | Regra |
|---|---|---|
| `GET` | `/alunos/{id}/concursos` | aluno próprio ou mentor vinculado |
| `POST` | `/alunos/{id}/concursos` | mentor |
| `PATCH` | `/alunos/{id}/concursos/{concursoId}` | mentor |
| `DELETE` | `/alunos/{id}/concursos/{concursoId}` | mentor |
| `GET` | `/alunos/{id}/edital` | aluno próprio ou mentor vinculado |
| `PUT` | `/alunos/{id}/edital` | mentor |
| `PATCH` | `/alunos/{id}/edital/itens/{itemId}` | mentor para estrutura |
| `PATCH` | `/alunos/{id}/edital/progresso/{itemId}` | aluno próprio ou mentor vinculado |

### 7.6 Cronograma

| Método | Endpoint | Regra |
|---|---|---|
| `GET` | `/alunos/{id}/cronograma` | aluno próprio ou mentor vinculado |
| `PUT` | `/alunos/{id}/cronograma` | mentor; aluno somente se habilitado |
| `POST` | `/alunos/{id}/cronograma/gerar` | mentor; aluno somente se habilitado |
| `PATCH` | `/alunos/{id}/cronograma/itens/{itemId}` | execução pelo aluno; edição pelo mentor |
| `GET` | `/alunos/{id}/cronograma/calendario` | aluno próprio ou mentor vinculado |

Quando o aluno habilitado gerar um cronograma inteligente, o novo cronograma substitui o cronograma ativo. A versão anterior será desativada e preservada para histórico.

### 7.7 Cartões de estudo

| Método | Endpoint | Regra |
|---|---|---|
| `GET` | `/alunos/{id}/baralhos-cartoes` | lista pessoais e do mentor |
| `POST` | `/alunos/{id}/baralhos-cartoes` | aluno cria pessoal; mentor cria para aluno |
| `PATCH` | `/alunos/{id}/baralhos-cartoes/{baralhoId}` | criador ou mentor vinculado |
| `DELETE` | `/alunos/{id}/baralhos-cartoes/{baralhoId}` | criador ou mentor vinculado |
| `POST` | `/alunos/{id}/baralhos-cartoes/{baralhoId}/cartoes` | cria cartão em baralho editável |
| `PATCH` | `/alunos/{id}/cartoes/{cartaoId}` | criador ou mentor vinculado |
| `DELETE` | `/alunos/{id}/cartoes/{cartaoId}` | criador ou mentor vinculado |
| `POST` | `/alunos/{id}/cartoes/{cartaoId}/revisar` | aluno próprio |
| `GET` | `/alunos/{id}/cartoes/pendentes` | aluno próprio ou mentor vinculado |

O mentor poderá definir o alcance do baralho como `aluno`, `edital` ou `global`.

### 7.8 Estudos, revisões, simulados e métricas

| Método | Endpoint | Regra |
|---|---|---|
| `GET`, `POST` | `/alunos/{id}/sessoes-estudo` | consulta e registra sessões |
| `PATCH`, `DELETE` | `/alunos/{id}/sessoes-estudo/{sessaoId}` | altera ou desativa logicamente |
| `GET`, `POST` | `/alunos/{id}/registros-questoes` | consulta e registra questões |
| `PATCH`, `DELETE` | `/alunos/{id}/registros-questoes/{registroId}` | altera ou desativa logicamente |
| `GET`, `POST` | `/alunos/{id}/revisoes` | consulta e programa revisões |
| `PATCH`, `DELETE` | `/alunos/{id}/revisoes/{revisaoId}` | conclui, adia, altera ou desativa |
| `GET`, `PUT` | `/alunos/{id}/configuracao-prova/{concursoId}` | configuração completa da regra de prova |
| `GET`, `POST` | `/alunos/{id}/simulados` | consulta e registra simulados |
| `PUT`, `DELETE` | `/alunos/{id}/simulados/{simuladoId}` | substitui transacionalmente ou desativa |
| `GET` | `/alunos/{id}/metricas/resumo` | totais, período, sequência, simulados e cobertura |
| `GET` | `/alunos/{id}/metricas/linha-do-tempo` | atividade diária no intervalo solicitado |
| `GET` | `/alunos/{id}/metricas/materias` | tempo, questões, edital e simulados por matéria |

Aluno próprio e mentor vinculado podem consultar. Escritas seguem as regras de cada domínio; o mentor pode alterar qualquer dado do aluno.

### 7.9 Respostas e erros

Formato de sucesso:

```json
{
  "dados": {},
  "metadados": {}
}
```

Formato de erro:

```json
{
  "erro": {
    "codigo": "ACESSO_NEGADO",
    "message": "Ação não permitida para este usuário.",
    "requisicaoId": "..."
  }
}
```

Códigos HTTP principais:

- `200` consulta ou alteração concluída;
- `201` recurso criado;
- `204` exclusão concluída;
- `400` entrada inválida;
- `401` sessão ausente ou expirada;
- `403` usuário sem permissão;
- `404` recurso inexistente dentro do escopo autorizado;
- `409` conflito de versão ou duplicidade;
- `422` regra de domínio não atendida;
- `429` excesso de tentativas;
- `500` erro interno sem exposição de detalhes sensíveis.

## 8. Mudanças no frontend

### 8.1 Infraestrutura

- substituir o `DataContext` atual por uma camada de sessão e cache de API;
- remover leitura e escrita de coleções no `localStorage`;
- criar cliente HTTP tipado;
- enviar cookies com `credentials: "include"`;
- tratar `401`, `403`, indisponibilidade e conflitos;
- adicionar estados de carregamento, erro e tentativa novamente;
- invalidar ou atualizar o cache após cada mutação;
- manter apenas preferências visuais não sensíveis no navegador, caso desejado.

### 8.2 Autenticação

Novas telas e fluxos:

- entrada no sistema;
- sessão expirada;
- troca de senha;
- saída do sistema;
- bloqueio por plano expirado;
- proteção de rotas;
- redirecionamento conforme o papel.

### 8.3 Experiência do mestre

- painel administrativo de usuários;
- criação de contas de mentor e aluno;
- definição ou troca do mentor único do aluno;
- desativação e reativação de contas;
- redefinição administrativa de senha;
- visualização da situação e da expiração do plano.

### 8.4 Experiência do mentor

- tela inicial com lista de alunos;
- seleção persistente do aluno acompanhado durante a sessão;
- cabeçalho indicando claramente o aluno atual;
- visão geral individual;
- acesso às telas existentes no contexto do aluno;
- controle “Aluno pode montar cronograma inteligente”;
- campo de data de expiração do plano;
- edição integral do edital e do cronograma;
- criação de cartões de estudo para o aluno;
- identificação da autoria dos cartões de estudo.

### 8.5 Experiência do aluno

- remoção dos controles estruturais do edital;
- edital apresentado em modo somente leitura estrutural;
- manutenção dos controles de progresso e estudo;
- cronograma do mentor disponível para execução;
- construtor inteligente visível somente quando autorizado;
- cartões de estudo pessoais editáveis;
- cartões de estudo do mentor identificados e somente leitura;
- demais telas consumindo dados do usuário autenticado;
- tela de plano expirado sem acesso às rotas acadêmicas;
- edição apenas do próprio nome e senha.

### 8.6 Matriz de interface

| Recurso | Mestre | Mentor | Aluno |
|---|---|---|---|
| Criar mentor e aluno | sim | não | não |
| Definir mentor do aluno | sim | não | não |
| Selecionar aluno para acompanhamento | não | sim | não |
| Definir expiração do plano | não | sim | não |
| Editar edital | não | sim | não |
| Atualizar progresso do edital | não | sim | sim |
| Montar cronograma manual | não | sim | execução apenas |
| Gerar cronograma inteligente | não | sim | conforme permissão |
| Alterar permissão inteligente | não | sim | não |
| Criar cartão de estudo pessoal do aluno | não | sim | sim |
| Criar cartões por edital ou globais | não | sim | não |
| Editar cartão de estudo criado pelo mentor | não | sim | não |
| Estudar qualquer cartão disponível | não | sim | sim |
| Registrar sessões e questões | não | sim | sim |
| Alterar nome e senha próprios | sim | sim | sim |

Ocultar um botão não substitui autorização: todas as regras serão repetidas e aplicadas pela API.

## 9. Migração dos dados existentes

A remoção do `localStorage` deve acontecer somente depois que a importação para o MySQL estiver validada.

### 9.1 Estratégia

1. Congelar e documentar o formato atual do payload.
2. Criar tabelas e migrações no MySQL.
3. Criar um importador Go idempotente.
4. Associar cada payload importado a um usuário aluno.
5. Converter concursos, edital, progresso, sessões, questões, revisões, simulados e cartões de estudo.
6. Preservar identificadores quando forem válidos ou manter tabela de correspondência.
7. Validar contagens e totais antes e depois.
8. Gerar relatório de importação com itens convertidos, ignorados e inválidos.
9. Liberar o frontend baseado na API.
10. Remover definitivamente o código de persistência em `localStorage`.

### 9.2 Compatibilidade temporária

Durante a transição, não haverá escrita simultânea em dois bancos por tempo indeterminado. O período de compatibilidade terá uma finalidade única: importar os dados e validar a nova API.

Depois da confirmação:

- o frontend deixa de chamar `getArray` e `saveArray`;
- o payload antigo deixa de ser a fonte da verdade;
- todas as alterações passam pela API Go;
- dados antigos no navegador não são utilizados pela aplicação.

## 10. Ordem de implementação

### Fase 0 — decisões e contrato

- registrar as decisões de produto consolidadas;
- definir o usuário mestre inicial de forma segura;
- definir o fuso horário usado na expiração dos planos;
- fechar contrato OpenAPI inicial;
- definir critérios de aceite.

### Fase 1 — fundação do backend — concluída

- criar módulo Go;
- configuração por ambiente;
- conexão MySQL;
- sistema de migrações;
- tratamento de erros;
- registros e verificações de integridade;
- testes unitários e de integração.

Entrega: API inicial sobe, conecta ao banco e executa migrações.

### Fase 2 — autenticação e RBAC — concluída

- usuários e papéis;
- papel mestre e administração de contas;
- hash de senha;
- entrada, saída e sessão;
- CSRF e limite de tentativas;
- middleware de autenticação e autorização;
- vínculo mentor–aluno;
- mentor único e expiração do plano;
- auditoria das ações sensíveis.

Entrega: mentor e aluno autenticam, e acessos indevidos são bloqueados por testes.

### Fase 3 — concursos e edital — concluída

- concursos atribuídos;
- estrutura normalizada do edital;
- progresso separado por aluno;
- endpoints de leitura e escrita;
- testes de imutabilidade do edital para aluno.

Entrega: mentor administra o edital e aluno apenas consulta e registra progresso.

### Fase 4 — cronogramas — concluída

- cronograma montado pelo mentor;
- itens e calendário;
- configuração por aluno;
- geração inteligente no backend;
- bloqueio ou liberação para o aluno.

Entrega: cronograma funciona conforme a permissão individual.

### Fase 5 — cartões de estudo — concluída

- baralhos e cartões;
- autoria e origem;
- permissões de edição;
- repetição espaçada e histórico individual.

Também foram implementados os recursos já utilizados pela interface: subbaralhos, ícone, ordem, vínculos com matéria/tópico/subtópico, dicas, etiquetas, alternativas e explicações. As avaliações da revisão são `0` (errei), `1` (difícil), `2` (bom) e `3` (fácil).

Entrega: aluno distingue e estuda cartões pessoais e do mentor.

### Fase 6 — estudos, revisões e simulados — concluída

- sessões e questões;
- revisões programadas;
- configuração da prova;
- simulados e resultados por matéria;
- endpoints transacionais.

As exclusões são lógicas. Os resultados por matéria substituídos durante a edição de um simulado também são preservados como inativos. A configuração de prova permanece em JSON para comportar os modos tradicional, Cebraspe e média por áreas, grupos objetivos e redações já existentes na interface.

Entrega: todas as rotinas de estudo usam exclusivamente a API.

### Fase 7 — métricas — concluída

- consultas agregadas;
- séries temporais;
- indicadores por matéria, concurso e aluno;
- índices necessários para desempenho.

Os filtros aceitam `concursoId`, `inicio` e `fim` no formato `AAAA-MM-DD`. As datas são consolidadas no fuso horário configurado pela aplicação. A linha do tempo inclui dias sem lançamentos para alimentar diretamente gráficos, mapas de atividade e retrospectivas.

Entrega: dashboard, histórico e Raio-X sem cálculos dependentes de coleções locais completas.

### Fase 8 — frontend de autenticação e mentor — concluída

- sessão global;
- entrada e rotas protegidas;
- seleção e visão individual do aluno;
- permissões refletidas na interface;
- cliente HTTP e cache.

O token da sessão permanece em cookie `HttpOnly` e o token CSRF fica somente em memória. O endpoint de sessão renova o CSRF após recarregar a página, e o cliente refaz uma escrita uma vez quando detectar rotação concorrente. Mestre, mentor e aluno possuem entradas e rotas separadas. O aluno sem permissão não visualiza o construtor inteligente, mas continua acessando o cronograma ativo publicado.

### Fase 9 — migração de todas as telas

Migrar uma tela por vez nesta ordem:

1. concursos — concluída;
2. edital — concluída;
3. cronograma — concluída;
4. sessões e histórico — concluídos;
5. revisões — concluída;
6. cartões de estudo — concluída;
7. simulados — concluída;
8. dashboard e métricas — concluída;
9. ajuda e configurações de perfil — concluída.

Cada tela só é considerada migrada quando não lê nem grava dados de domínio no `localStorage`.

Na migração do edital verticalizado foram concluídos: atribuição e reutilização pelo mentor, importação JSON, manutenção de matérias, tópicos e subtópicos, ordenação, materiais em qualquer nível, progresso de conteúdo e de materiais separado por aluno, lançamentos manuais de estudo e questões, cronômetro persistido pela API e programação de revisões. O aluno possui leitura da estrutura e escrita apenas no próprio progresso e histórico; as alterações estruturais permanecem exclusivas do mentor.

Na migração do cronograma foram concluídos: editor manual exclusivo do mentor, geração agendada, ciclo inteligente, configuração de disponibilidade semanal, seleção e prioridade de matérias, ritmo, limite diário, alertas, estimativa, calendário produzido pelo backend, versões substituídas com preservação histórica, alternância de visualização entre lista e calendário na tela do aluno e execução das atividades pelo aluno. O aluno só visualiza o construtor quando habilitado e, quando gera um novo plano, substitui o cronograma ativo conforme a regra definida. Nenhuma configuração ou ciclo dessa tela permanece no `localStorage`.

Na migração de sessões e histórico foram concluídos: carregamento autenticado de sessões, questões e simulados; edição por endpoint; exclusão lógica com preservação do registro; cálculo de períodos, sequência, distribuição por matéria, desempenho, cobertura, mapas de atividade e retrospectiva anual; comparação com a disponibilidade do cronograma armazenado; exportação visual do relatório; e acesso do mentor ao histórico individual do aluno. A tela e o cronômetro não leem nem gravam sessões ou questões no `localStorage`.

Na migração de revisões foram concluídos: agenda de sete dias, classificação em atrasadas, atuais, futuras e concluídas, ciclos de 1, 7, 30 e 90 dias, percentual obtido, conclusão, encerramento do ciclo, adiamento, exclusão lógica, alertas de desempenho por questões e navegação direta ao conteúdo do edital. O backend retorna os nomes dos conteúdos e valida que toda revisão criada ou remanejada pertence ao edital atribuído ao aluno. O mentor possui acesso à agenda individual e o aluno pode iniciar o cronômetro vinculado à revisão.

Na migração de simulados e Raio-X foram concluídos: listagem, cadastro, edição e exclusão lógica de simulados transacionais via endpoints Go; carregamento e persistência da configuração completa da prova (regras objetivas, penalidade por erro, cortes, blocos e discursivas) no MySQL; cálculo automático e detalhamento por matéria sem uso de `localStorage`.

Na migração do Dashboard e Métricas foram concluídos: carregamento de métricas agregadas via API Go; consolidação de tempo de estudo, questões, taxa de acerto, último simulado e progresso do edital; remoção da persistência local de regras e matérias.

Na migração de Editais Premium, Ajuda e Perfis foram concluídos: criação de concursos e atribuição de editais estruturados pelo mentor via API Go; alteração de dados de perfil autenticado sem fallback em `localStorage`.

### Fase 10 — corte definitivo — concluída

- remoção completa da persistência em `localStorage` no frontend React (`DataContext.tsx` e telas);
- validação de leitura e escrita exclusivas na API Go com banco MySQL;
- preservação das regras de escopo por papel (mestre, mentor e aluno).
- remover endpoints substituídos;
- executar testes de regressão;
- preparar procedimento de rollback do banco.

## 11. Estratégia de testes

### Backend

- testes unitários de regras de domínio;
- testes de repositório com MySQL isolado;
- testes HTTP dos endpoints;
- testes de autenticação, expiração e revogação;
- testes de matriz de permissão;
- testes de concorrência e transações;
- testes do importador e idempotência.

Casos obrigatórios de autorização:

- aluno não consulta outro aluno;
- mentor não consulta aluno sem vínculo;
- aluno não altera edital por nenhum endpoint;
- aluno sem permissão não gera cronograma inteligente;
- aluno habilitado gera apenas o próprio cronograma;
- aluno não edita cartão de estudo criado pelo mentor;
- mentor altera qualquer dado do aluno vinculado;
- aluno com plano expirado não entra e não reutiliza uma sessão aberta;
- mentor pode renovar o plano e restabelecer o acesso;
- mestre cria, desativa e reativa contas, mas não acessa implicitamente os dados acadêmicos;
- segundo vínculo ativo de mentor para o mesmo aluno é rejeitado;
- exclusões preservam registros desativados;
- sessão revogada não pode ser reutilizada.

### Frontend

- testes dos estados de autenticação;
- testes de proteção de rotas;
- testes de visibilidade por papel e permissão;
- testes de formulário e erros da API;
- testes de integração dos fluxos principais;
- testes ponta a ponta para mentor e aluno.

## 12. Segurança e operação

- HTTPS obrigatório fora do ambiente local.
- Cookies `Secure` em produção.
- CORS limitado à origem do frontend.
- Credenciais somente em variáveis de ambiente.
- Usuário MySQL com permissões mínimas.
- Backups automáticos do MySQL.
- Migrações aplicadas antes da nova versão da API.
- Limite de tamanho de requisições e imagens.
- Sanitização de conteúdo rico dos cartões de estudo e materiais.
- Cabeçalhos de segurança.
- Política de expiração e rotação de sessões.
- Auditoria de alterações realizadas pelo mentor.
- Monitoramento de erros, latência e disponibilidade.

## 13. Critérios de conclusão

A mudança será considerada concluída quando:

- todos os usuários precisarem autenticar;
- mestre, mentor e aluno tiverem permissões aplicadas no backend;
- somente o mestre puder criar contas;
- cada aluno possuir no máximo um mentor ativo;
- plano expirado bloquear imediatamente o acesso do aluno;
- mentor conseguir selecionar e administrar cada aluno vinculado;
- aluno não conseguir alterar a estrutura do edital;
- permissão do cronograma inteligente funcionar individualmente;
- cartões de estudo pessoais e do mentor respeitarem autoria e edição;
- cartões do mentor funcionarem com alcance individual, por edital e global;
- editais forem reutilizáveis, mantendo progresso separado por aluno;
- exclusões forem lógicas e preservarem o histórico;
- todas as telas lerem e gravarem pela API Go;
- nenhuma coleção de domínio for salva em `localStorage`;
- MySQL for a fonte única da verdade;
- importação dos dados anteriores tiver relatório e validação;
- testes de autorização e fluxos principais estiverem aprovados;
- build do frontend e testes do backend passarem;
- documentação da API e operação estiver atualizada.

## 14. Decisões de produto consolidadas

1. Contas de mentor e aluno são criadas exclusivamente pelo usuário mestre; não existe cadastro público.
2. Cada aluno possui somente um mentor ativo.
3. Recuperação de senha por e-mail fica para uma fase posterior.
4. O aluno pode alterar apenas o próprio nome e a própria senha.
5. Exclusões são lógicas: registros ficam desativados para preservação do histórico.
6. Cartões criados pelo mentor podem ter alcance individual, por edital ou global para todos os alunos do sistema.
7. O mesmo edital é reutilizado por vários alunos; somente o progresso é individual.
8. Quando autorizado, o cronograma inteligente gerado pelo aluno substitui o cronograma ativo, preservando a versão anterior como desativada.
9. A data de expiração do plano é preenchida pelo mentor e bloqueia o acesso do aluno quando a data é alcançada.

Essas regras passam a fazer parte do modelo de dados, da autorização da API, dos testes e dos critérios de aceite.
