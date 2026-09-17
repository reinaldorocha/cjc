# Banco de dados

O backend usa exclusivamente MySQL. Na inicialização, cria o banco informado em `DB_NOME` (se necessário) e aplica as migrações em `web/server-go/migrations`.

As migrações atuais (001 a 016) cobrem usuários e sessões, mentoria, concursos, editais verticalizados, cronogramas, cartões e revisões, sessões de estudo, questões, materiais, cadernos/notas, simulados, métricas e identidade visual.

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USUARIO=root
DB_SENHA=rootpassword
DB_NOME=track_concursos
```

Os nomes e relacionamentos exatos estão nos arquivos SQL. Não há suporte a SQLite nem uso de `localStorage` para dados do sistema.

## Domínios persistidos

Usuários e sessões; vínculo mentor-aluno; concursos e editais verticalizados; atribuições e progresso; cronogramas e itens; baralhos, cartões e revisões; sessões de estudo; questões; materiais; cadernos e notas; configurações de prova; simulados; métricas; e identidade visual.

As migrações devem ser executadas na ordem numérica. Uma alteração estrutural deve ser criada em nova migração, sem editar uma já aplicada.

## Tabelas principais

| Área | Tabelas |
|---|---|
| Identidade | `usuarios`, `sessoes_autenticacao`, `eventos_auditoria` |
| Mentoria | `mentor_alunos` |
| Catálogo | `concursos`, `editais`, `edital_materias`, `edital_topicos`, `edital_subtopicos` |
| Atribuições | `aluno_concursos`, `aluno_editais`, `aluno_progresso_edital` |
| Planejamento | `cronogramas`, `cronograma_itens`, `rotinas_estudo` |
| Cartões | `baralhos_cartoes`, `cartoes_estudo`, `revisoes_cartoes` |
| Desempenho | `sessoes_estudo`, `registros_questoes`, `revisoes_programadas` |
| Avaliações | `configuracoes_prova`, `simulados`, `resultados_materias_simulado` |
| Apoio | tabelas de materiais, arquivos e pastas |
| Organização pessoal | cadernos e notas |
| Personalização | configuração de identidade visual do mentor |

As chaves estrangeiras e índices devem ser consultados nas migrações, pois são a fonte oficial do esquema.

## Principais relacionamentos

`mentor_alunos` liga um mentor a um aluno. `aluno_concursos` e `aluno_editais` representam atribuições. `edital_materias`, `edital_topicos` e `edital_subtopicos` formam a hierarquia do edital. `cronogramas` possui vários `cronograma_itens`. `baralhos_cartoes` possui `cartoes_estudo`, e as revisões referenciam o cartão e o aluno. Registros de estudo, questões e simulados referenciam o aluno para manter as métricas isoladas.

## Migrações

As migrações iniciais criam a identidade e o domínio de estudos. Migrações posteriores adicionam cronogramas versionados, cartões completos, rotinas, índices, concursos, materiais individuais, banco de questões, automação de revisões, telefone, pastas, cadernos/notas e configuração de identidade visual. A aplicação registra quais versões já foram executadas.

## Backup e restauração

Use ferramentas do MySQL, como `mysqldump`, para backup. O backup deve incluir estrutura e dados. Teste a restauração em um banco separado antes de atualizar a aplicação. Não copie arquivos de um suposto banco SQLite, pois esse formato não é utilizado.
