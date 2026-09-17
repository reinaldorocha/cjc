# Fluxos funcionais

## Mestre

Cria e desativa usuários mentor/aluno, redefine senhas e define o mentor de cada aluno.

## Mentor

Cadastra alunos, concursos e editais no catálogo, cria baralhos e cartões globais ou associados a edital, cadastra materiais e acompanha alunos. Em cada aluno, atribui concurso/edital, configura permissões, data de expiração e monta o cronograma.

## Aluno

Consulta o edital atribuído, cronograma, cartões, questões, materiais, cadernos, sessões, revisões, simulados e métricas. Pode alterar nome e senha. Quando habilitado pelo mentor, pode gerar cronograma inteligente, substituindo o atual.

## Sequência recomendada

1. O mestre cria o mentor e as contas de acesso.
2. O mentor cadastra o aluno, concurso e edital.
3. O mentor atribui o edital, define expiração e permissões.
4. O mentor monta o cronograma e disponibiliza materiais e cartões.
5. O aluno estuda e registra sessões, questões, revisões e simulados.
6. O mentor acompanha a visão geral e as métricas.

## Edital e conteúdos

O edital é criado e ordenado pelo mentor. O aluno pode registrar progresso, mas não criar, excluir ou reorganizar matérias, tópicos e sub tópicos. Baralhos podem ser globais ou vinculados a edital, e um edital pode atender vários alunos.

Todos os fluxos dependem de autenticação e do vínculo correto entre mentor e aluno.

## Procedimentos detalhados

### Criar aluno

O mentor abre a área de alunos, informa os dados solicitados e salva. O aluno passa a aparecer vinculado ao mentor. A atribuição de concurso ou edital é uma etapa posterior e opcional.

### Atribuir edital

O mentor seleciona um edital existente no catálogo, escolhe o aluno e confirma a atribuição. A estrutura do edital continua sendo a mesma do catálogo. O progresso é armazenado no contexto do aluno.

### Montar cronograma

O mentor informa itens, datas e organização do estudo e salva a versão. Se a configuração inteligente do aluno estiver ativa, o aluno pode solicitar geração; a nova versão substitui a anterior conforme a regra da aplicação.

### Replanejamento com revisões de assuntos

Ao marcar um tópico/subtópico como estudado, as revisões automáticas passam a ocupar a mesma carga diária dos assuntos. Cada revisão reserva 30 minutos; assuntos usam sua duração cadastrada. As revisões têm prioridade nas datas previstas e os assuntos pendentes são redistribuídos no espaço restante. Revisões excedentes seguem para o próximo dia disponível, sem antecipar sua data; ciclos do mesmo assunto não são concentrados no mesmo dia.

O replanejamento respeita horas por dia e o limite de atividades (`maxTopicosDia`, contando assuntos e revisões). Após concluir um assunto, começa amanhã e preserva as tarefas de hoje; pendências anteriores a hoje são recuperadas. O botão **Reprogramar Pendentes** começa hoje e considera o concurso selecionado. Atividades concluídas/ignoradas não são movidas; concluídas já planejadas no período descontam sua carga. Cronogramas de ciclo sem datas não são convertidos em calendário.

Criar ou alterar uma revisão também reorganiza os cronogramas agendados ativos do aluno. Concluir uma revisão não marca o tópico novamente, evitando recriar seus ciclos. Se a carga configurada não comportar alguma atividade, a reorganização não aplica mudanças parciais e a interface informa que a alteração principal foi salva, mas o replanejamento precisa ser repetido após ajustar a carga. Não é necessário concluir a atividade novamente.

Testes: `go test ./internal/cronogramas`. Integração MySQL opcional: `$env:CRONOGRAMA_TEST_MYSQL='1'; go test ./internal/cronogramas -count=1`; cria e remove somente um banco temporário `test_replanejar_*`.

### Revisão de cartões

O aluno consulta cartões pendentes, informa o resultado da revisão e o sistema registra o histórico. Baralhos criados pelo mentor aparecem conforme o escopo global ou o edital atribuído.

### Acompanhar desempenho

O aluno registra sessões, respostas e simulados. O mentor abre a visão individual para consultar resumo, linha do tempo e resultados por matéria.

## Estados importantes

Conta desativada não autentica. Aluno com plano expirado não acessa a área protegida. Um edital não atribuído não aparece como edital ativo do aluno. Um cronograma vazio pode ser criado pelo mentor ou gerado quando a configuração permitir.
