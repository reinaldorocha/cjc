# Track Concursos Web

AplicaÃ§Ã£o web para gestÃ£o de estudos de concursos, com cliente React, API Go e persistÃªncia MySQL.

## Perfis e permissÃµes

- **Mestre:** cria, edita, desativa e reativa mentores e alunos, redefine senhas e define o mentor de cada aluno.
- **Mentor:** cadastra alunos com vÃ­nculo automÃ¡tico, mantÃ©m catÃ¡logos prÃ³prios de concursos, editais e flashcards sem depender de aluno, e acompanha individualmente cronograma, revisÃµes, histÃ³rico, simulados e mÃ©tricas atravÃ©s de um menu lateral dedicado.
- **Aluno:** acessa os prÃ³prios dados, executa o cronograma, registra estudos e questÃµes, revisa cartÃµes e consulta as demais telas permitidas.

O edital Ã© atribuÃ­do e estruturado somente pelo mentor. O aluno pode gerar um cronograma inteligente apenas quando estiver habilitado; a nova versÃ£o substitui o cronograma ativo. A data de expiraÃ§Ã£o configurada pelo mentor bloqueia o acesso do aluno ao vencer.

Concursos, editais e flashcards sÃ£o criados na Ã¡rea de conteÃºdos da mentoria. Na Ã¡rea individual do aluno, o mentor apenas atribui itens do catÃ¡logo, monta o cronograma e acompanha o desempenho. Baralhos do mentor podem ser globais ou restritos aos alunos que possuem um edital especÃ­fico.

## Funcionalidades

- [cursos da mentoria](CURSOS.md), com aulas do YouTube e liberaÃ§Ã£o por aluno, concurso ou toda a mentoria;

- concursos e resultados;
- edital verticalizado reutilizÃ¡vel, materiais e progresso individual;
- cronograma manual, agendado ou por ciclo, com alternÃ¢ncia de visualizaÃ§Ã£o entre lista e calendÃ¡rio;
- sessÃµes de estudo, questÃµes e histÃ³rico;
- revisÃµes programadas;
- baralhos pessoais, por edital ou globais e repetiÃ§Ã£o espaÃ§ada;
- simulados, configuraÃ§Ã£o da prova e Raio-X;
- dashboard e mÃ©tricas por perÃ­odo e matÃ©ria;
- biblioteca de editais;
- administraÃ§Ã£o de usuÃ¡rios e configuraÃ§Ã£o individual do aluno.

## Tecnologias

- React 19, TypeScript, Vite e React Router;
- Chart.js;
- Go 1.24;
- MySQL 8;
- autenticaÃ§Ã£o por cookie `HttpOnly` e proteÃ§Ã£o CSRF.

## Estrutura

```text
web/
â”œâ”€â”€ client-react/   # cÃ³digo-fonte React
â”œâ”€â”€ client/         # build gerado do frontend
â”œâ”€â”€ server-go/      # API, serviÃ§os e migraÃ§Ãµes MySQL
â””â”€â”€ PLANO_BACKEND_GOLANG.md
```

## ConfiguraÃ§Ã£o do backend

Na pasta `server-go`, copie o arquivo de exemplo e configure o MySQL e a conta mestre inicial:

```powershell
Copy-Item .env.example .env
go run ./cmd/api
```

A API responde por padrÃ£o em `http://127.0.0.1:8080`. As migraÃ§Ãµes sÃ£o aplicadas na inicializaÃ§Ã£o. A lista completa de endpoints estÃ¡ em [server-go/README.md](server-go/README.md).

## ExecuÃ§Ã£o do frontend

Na pasta `client-react`:

```powershell
npm install
npm run dev
```

O servidor Vite encaminha `/api` para a API Go conforme sua configuraÃ§Ã£o de desenvolvimento.

## PersistÃªncia

Todos os dados de domÃ­nio sÃ£o persistidos no MySQL por endpoints autenticados. O frontend nÃ£o utiliza `localStorage` como banco de dados. ExclusÃµes de registros de negÃ³cio sÃ£o lÃ³gicas para preservar o histÃ³rico.

## ValidaÃ§Ã£o

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

Os testes unitÃ¡rios permanecem adiados atÃ© a aprovaÃ§Ã£o funcional do sistema.

