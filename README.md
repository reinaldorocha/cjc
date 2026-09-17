# Track Concursos Web

Aplicação web para gestão de estudos de concursos, com cliente React, API Go e persistência MySQL.

## Perfis e permissões

- **Mestre:** cria, edita, desativa e reativa mentores e alunos, redefine senhas e define o mentor de cada aluno.
- **Mentor:** cadastra alunos com vínculo automático, mantém catálogos próprios de concursos, editais e flashcards sem depender de aluno, e acompanha individualmente cronograma, revisões, histórico, simulados e métricas através de um menu lateral dedicado.
- **Aluno:** acessa os próprios dados, executa o cronograma, registra estudos e questões, revisa cartões e consulta as demais telas permitidas.

O edital é atribuído e estruturado somente pelo mentor. O aluno pode gerar um cronograma inteligente apenas quando estiver habilitado; a nova versão substitui o cronograma ativo. A data de expiração configurada pelo mentor bloqueia o acesso do aluno ao vencer.

Concursos, editais e flashcards são criados na área de conteúdos da mentoria. Na área individual do aluno, o mentor apenas atribui itens do catálogo, monta o cronograma e acompanha o desempenho. Baralhos do mentor podem ser globais ou restritos aos alunos que possuem um edital específico.

## Funcionalidades

- [cursos da mentoria](CURSOS.md), com aulas do YouTube e liberação por aluno, concurso ou toda a mentoria;

- concursos e resultados;
- edital verticalizado reutilizável, materiais e progresso individual;
- cronograma manual, agendado ou por ciclo, com alternância de visualização entre lista e calendário;
- sessões de estudo, questões e histórico;
- revisões programadas;
- baralhos pessoais, por edital ou globais e repetição espaçada;
- simulados, configuração da prova e Raio-X;
- dashboard e métricas por período e matéria;
- biblioteca de editais;
- administração de usuários e configuração individual do aluno.

## Tecnologias

- React 19, TypeScript, Vite e React Router;
- Chart.js;
- Go 1.24;
- MySQL 8;
- autenticação por cookie `HttpOnly` e proteção CSRF.

## Estrutura

```text
web/
├── client-react/   # código-fonte React
├── client/         # build gerado do frontend
├── server-go/      # API, serviços e migrações MySQL
└── PLANO_BACKEND_GOLANG.md
```

## Configuração do backend

Na pasta `server-go`, copie o arquivo de exemplo e configure o MySQL e a conta mestre inicial:

```powershell
Copy-Item .env.example .env
go run ./cmd/api
```

A API responde por padrão em `http://127.0.0.1:8080`. As migrações são aplicadas na inicialização. A lista completa de endpoints está em [server-go/README.md](server-go/README.md).

## Execução do frontend

Na pasta `client-react`:

```powershell
npm install
npm run dev
```

O servidor Vite encaminha `/api` para a API Go conforme sua configuração de desenvolvimento.

## Persistência

Todos os dados de domínio são persistidos no MySQL por endpoints autenticados. O frontend não utiliza `localStorage` como banco de dados. Exclusões de registros de negócio são lógicas para preservar o histórico.

## Validação

Frontend:

```powershell
npm run build
npm run lint
```

Backend:

```powershell
go build ./...
go vet ./...
```

Os testes unitários permanecem adiados até a aprovação funcional do sistema.
