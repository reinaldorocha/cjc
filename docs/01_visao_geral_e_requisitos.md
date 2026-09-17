# Visão geral

O Track Concursos é uma aplicação web para organização de estudos. A interface é React e a API é Go; os dados são persistidos em MySQL.

## Perfis

- **Mestre**: administra usuários e vínculos de mentoria.
- **Mentor**: cadastra alunos, concursos, editais, baralhos, cartões e materiais; acompanha alunos e configura seus estudos.
- **Aluno**: acessa seus conteúdos, cronograma, edital atribuído, cartões, questões, sessões, revisões, simulados e métricas.

O acesso do aluno depende de uma data de expiração de plano, quando preenchida.

## Funcionalidades

Autenticação, perfil e sessões; cadastro e vínculo de alunos; catálogo de concursos e editais; edital verticalizado; cronograma e calendário; baralhos e revisões de cartões; banco de questões; materiais de apoio; cadernos e notas; sessões de estudo; revisões programadas; configuração de prova; simulados; métricas; e identidade visual do mentor.

## Regras de acesso

- O mestre cria, desativa, reativa e redefine contas.
- Cada aluno possui um mentor responsável.
- O mentor cadastra alunos, concursos e editais independentemente de atribuições.
- Concursos e editais podem ser reutilizados por vários alunos.
- O edital atribuído é mantido pelo mentor; o aluno apenas consulta e registra progresso.
- Cartões podem ser globais ou associados a um edital.
- O mentor habilita, por aluno, a geração inteligente de cronograma.
- A data de expiração do plano bloqueia o acesso do aluno após o vencimento.

## Requisitos funcionais

### Usuários e mentoria

O sistema mantém nome, e-mail, senha, perfil, status e data de expiração do plano. O mestre é o responsável pela criação inicial das contas. O mentor pode cadastrar alunos e consultar somente alunos vinculados a ele. A desativação preserva o histórico e impede novo acesso.

### Concursos e editais

O mentor mantém um catálogo próprio de concursos e editais. O concurso é o agrupador e o edital contém a árvore de matérias, tópicos e sub tópicos. A ordem dos itens é persistida. A mesma definição pode ser atribuída a vários alunos sem duplicação do catálogo.

### Estudos

O cronograma possui versões e itens planejados. O mentor pode montar itens manualmente; quando a permissão estiver ativa, o aluno pode gerar um cronograma inteligente, que substitui a versão atual. Sessões, questões, revisões e simulados alimentam as métricas.

### Cartões e materiais

O mentor cria baralhos e cartões para todos os alunos ou para alunos que utilizam determinado edital. O aluno mantém seus próprios baralhos e registra as revisões. Materiais de apoio podem possuir pasta, descrição e arquivo.

### Restrições

Não há cadastro público, recuperação de senha por e-mail, múltiplos mentores por aluno ou armazenamento de dados funcionais no navegador.

## Telas e áreas

### Área pública

Tela de entrada, autenticação e recuperação do contexto de sessão. A aplicação consulta o usuário atual pela API antes de liberar as áreas protegidas.

### Área do mentor

Painel, alunos, radar de acompanhamento, concursos, editais, conteúdos (baralhos, cartões e materiais), configurações do aluno e identidade visual. O catálogo é acessível mesmo quando ainda não há aluno vinculado.

### Área do aluno

Resumo, cronograma, calendário, edital verticalizado, cartões, banco de questões, materiais, cadernos, sessões de estudo, revisões, simulados, configuração de prova e métricas.

### Administração

Gestão de usuários, status das contas, redefinição de senha e vínculo mentor-aluno.
