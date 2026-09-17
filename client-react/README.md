# Track Concursos — cliente React

Interface React + TypeScript da versão web do Track Concursos.

## Desenvolvimento

1. Configure e inicie a API Go em `web/server-go` na porta 8080.
2. Nesta pasta, execute `npm run dev`.

O Vite encaminha `/api` à API Go durante o desenvolvimento. A compilação de produção (`npm run build`) é gravada em `web/client`.

## Autenticação

- sessão em cookie `HttpOnly`;
- CSRF mantido somente em memória;
- rotas separadas para mestre, mentor e aluno;
- mestre administra contas e vínculos;
- mentor seleciona e configura cada aluno individualmente;
- aluno acessa o cronograma publicado e só vê o construtor inteligente quando autorizado.

## Migração das telas

- Concursos: usa exclusivamente os endpoints Go e estado React em memória; aluno consulta e mentor administra no contexto individual.
- Demais telas de domínio: serão conectadas sequencialmente conforme o plano de implementação.

## Verificações

- `npm run lint`
- `npm run build`
