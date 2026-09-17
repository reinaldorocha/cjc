# Cursos da mentoria

A aba **Cursos da Mentoria** permite cadastrar cursos com capa HTTPS opcional, categoria, descrição e aulas do YouTube organizadas em módulos. O aluno encontra o catálogo em **Meus Cursos**, mesmo sem concurso atribuído. O catálogo mostra as capas e os nomes dos cursos; ao abrir um curso, a página interna apresenta o banner e os cards de seus módulos.

## Cadastro por módulos

O catálogo sempre exibe um card por curso. Cada módulo pode ter descrição e capa própria, usando a capa do curso quando não informada. Dentro do curso, o aluno encontra os cards dos módulos e o card abre a sala diretamente no módulo escolhido, preservando progresso e bloqueios.

O catálogo usa tema escuro, coleções por categoria com setas e rolagem horizontal, e a coleção de retomada dos alunos. A página interna usa a capa como banner e concentra os módulos e as ações de gestão do mentor. O layout se adapta ao celular e respeita a preferência por movimento reduzido.

A capa pode ser enviada do computador em **JPG ou PNG, até 5 MB e 16 milhões de pixels**, ou informada por link HTTPS. O editor mostra a prévia e permite substituir ou remover a imagem. Recomenda-se proporção horizontal 16:9. Salve o rascunho e publique para disponibilizar a nova capa. Uploads ficam no MySQL (`cursos_capas`, migração 022), com validação real da imagem e acesso por sessão: proprietário ou aluno autorizado ao curso publicado. Imagens substituídas ou uploads cancelados permanecem armazenados para preservar versões publicadas.

1. Preencha os dados do curso e escolha quem terá acesso.
2. Clique em **Criar módulo** e informe o título.
3. Adicione aulas usando a lista lateral. Apenas a aula selecionada abre seus campos de edição.
4. Em cada aula, informe o título, cole o link do YouTube e, opcionalmente, anexe um PDF de até **10 MB**.
5. Clique em **Salvar módulo**. Isso grava o curso e o módulo no banco e retorna à lista compacta de módulos.

O módulo pode ser reaberto para editar, adicionar, remover ou reordenar aulas. Os módulos também podem ser reordenados ou removidos; confirme essas alterações em **Salvar rascunho** e depois **Publicar atualização**. Títulos de módulo devem ser únicos no curso. Cursos existentes são agrupados pelos títulos já cadastrados, preservando aulas e ordem dentro de cada módulo. Há limite de 300 aulas por curso.

Na sala do aluno, os módulos são recolhíveis e o PDF aparece abaixo do vídeo da aula. O anexo pode ser substituído ou removido pelo mentor. O botão **Salvar módulo** também salva alterações pendentes nos dados e acessos do curso.

## PDFs por aula

Os arquivos ficam no MySQL (`cursos_pdfs`), sem URL pública na pasta `uploads`. Upload exige perfil de mentor e CSRF. A API verifica extensão, assinatura PDF e limite de tamanho; ao vincular um PDF à aula, verifica se o arquivo pertence ao mentor. Downloads exigem sessão válida, autorização no curso e vínculo atual do PDF a uma aula desse curso. Remover o anexo ou revogar o acesso impede novos downloads; arquivos já baixados não podem ser recolhidos.

O upload precede o salvamento do módulo: se a edição for cancelada, o arquivo não fica disponível aos alunos. PDFs enviados e não vinculados permanecem armazenados, assim como versões substituídas, e devem ser considerados no dimensionamento e backup do MySQL.

## Disponibilização

- **Todos:** alunos com vínculo ativo com o mentor proprietário do curso.
- **Alunos selecionados:** somente os alunos escolhidos, enquanto vinculados à mentoria.
- **Concursos selecionados:** alunos da mentoria com atribuição ativa a pelo menos um concurso escolhido. O acesso não depende do concurso selecionado no menu naquele momento.

No acompanhamento individual, **Disponibilizar Cursos** abre o catálogo com o aluno pré-selecionado para novos cursos. Para um curso existente, use **Editar curso**, selecione o escopo e os destinatários, salve e publique a atualização. A seleção substitui a anterior na publicação. Desativar remove o acesso sem apagar o registro.

## Rascunhos e publicação

Novos cursos começam como rascunho, podendo ser salvos sem aulas ou destinatários. Para publicar, preencha pelo menos uma aula válida e os destinatários do escopo escolhido. **Salvar módulo** grava o rascunho; **Publicar curso / atualização** disponibiliza uma cópia completa aos alunos. Editar um curso publicado preserva a versão anterior até a próxima publicação, inclusive seus acessos e PDFs. **Retirar publicação** suspende o acesso e preserva conteúdo e progresso. Cursos anteriores à migração 021 continuam publicados.

O editor avisa antes de descartar alterações. Um número de revisão impede sobrescritas por edições concorrentes: se outra janela salvar, reabra o editor. **Prévia como aluno** apresenta o rascunho com as regras de um aluno iniciante, sem gravar progresso ou respostas.

## Aprendizagem e acompanhamento

- **Retomada:** a posição é salva durante reprodução a cada 15 segundos, ao pausar, trocar aula e sair. Fechar a aba tenta enviar a última posição; uma interrupção de rede ou encerramento abrupto pode impedir o último envio. Erros permitem tentar novamente.
- **Conclusão:** o aluno marca ou desmarca cada aula. A porcentagem corresponde às aulas marcadas, não comprova tempo assistido. IDs estáveis preservam progresso ao renomear ou reordenar; substituir o vídeo da aula reinicia seu progresso na exibição.
- **Liberação de módulos:** data opcional no fuso configurado na API e exigência de conclusão de todos os módulos anteriores. Ambas devem ser atendidas quando configuradas. A API omite vídeo, PDF e questões de aulas bloqueadas e nega os acessos diretos correspondentes. Reabra o curso para atualizar uma liberação por data.
- **Exercícios:** selecione até 50 questões ativas do banco do próprio mentor por aula, com busca. O aluno recebe alternativas sem gabarito; a resposta é corrigida no servidor, com explicação. Pode tentar novamente. As questões continuam vinculadas ao banco: alterações e desativações ali afetam os exercícios disponíveis.
- **Acompanhar alunos:** mostra alunos atualmente autorizados, conclusão, última atividade e acertos. Os estados são Não iniciado, Em andamento, Sem atividade há 7 dias e Concluído. Acertos usam a última resposta de cada questão por aula; repetições não multiplicam o total. Alunos sem atividade alguma aparecem como Não iniciado. A busca aceita nome, email e situação.

A coleção **Continue de onde parou** e os filtros de conclusão ajudam o aluno a retomar os estudos. O mentor pode filtrar publicados e rascunhos.

A API valida propriedade do curso e dos destinatários. Listagem e abertura direta aplicam a mesma regra; curso global de outro mentor não é liberado. O middleware existente continua exigindo sessão, CSRF nas escritas e plano válido para o aluno.

## Player e links

Aceita links `watch`, `youtu.be`, `shorts`, `live`, `embed` e o ID do vídeo. O backend armazena somente IDs de vídeos validados. A sala oferece reprodução, pausa, busca na linha do tempo, volume, tela cheia e navegação entre aulas. Não exige chave de API do YouTube.

A interface não apresenta URL nem botão para copiar ou compartilhar; o iframe não recebe cliques ou foco e o menu de contexto da área é desabilitado. **Isso não é DRM:** o ID permanece inspecionável no navegador e nas requisições. O YouTube pode exibir marca, anúncios e recomendações. Vídeos privados, removidos ou com incorporação desabilitada não são reproduzíveis. Para impedir acesso externo de forma mais forte, seria necessária outra hospedagem com controles de distribuição próprios.

Referências consultadas: [player oficial do YouTube](https://developers.google.com/youtube/player_parameters) e [exemplo de catálogo Netflix em React](https://github.com/karlhadwen/netflix). Implementação própria, sem copiar o projeto ou adicionar dependências.

## Instalação e manutenção

Reinicie a API Go para aplicar automaticamente `019_cursos.sql`, `020_cursos_pdfs.sql` e `021_cursos_aprendizagem.sql`. A migração 021 depende do banco de questões (017). Gere o frontend com `npm run build` em `web/client-react`; a saída vai para `web/client`. Em desenvolvimento, a configuração atual do Vite usa `5174` e encaminha `/api` a `8082`; configure a API com `PORTA=8082` e `ORIGEM_FRONTEND=http://127.0.0.1:5174`.

Os metadados, destinatários, aulas e versão publicada são persistidos no MySQL, tabela `cursos`. Aulas e destinatários usam colunas JSON; as alterações são transacionais. `cursos_progresso` registra posição/conclusão por aluno e aula; `cursos_respostas` registra última resposta e tentativas. Capas externas e vídeos dependem de conexão à internet.

| Método | Rota | Acesso |
| --- | --- | --- |
| GET / POST | `/api/v1/mentor/cursos` | Catálogo / criação pelo mentor |
| GET / PUT / DELETE | `/api/v1/mentor/cursos/{cursoId}` | Consulta / edição / desativação pelo proprietário |
| GET | `/api/v1/alunos/{alunoId}/cursos` | Cursos disponíveis ao aluno |
| GET | `/api/v1/alunos/{alunoId}/cursos/{cursoId}` | Abertura autorizada do curso |
| POST | `/api/v1/mentor/cursos/pdfs` | Upload multipart, campo `arquivo` |
| GET | `/api/v1/mentor/cursos/{cursoId}/pdfs/{pdfId}` | Download pelo proprietário |
| GET | `/api/v1/alunos/{alunoId}/cursos/{cursoId}/pdfs/{pdfId}` | Download autorizado para aluno |
| POST | `/api/v1/mentor/cursos/{cursoId}/publicacao` | Publicar/retirar: `{revisao, publicar}` |
| GET | `/api/v1/mentor/cursos/{cursoId}/previa` | Rascunho como aluno iniciante |
| GET | `/api/v1/mentor/cursos/{cursoId}/acompanhamento` | Progresso e acertos dos alunos autorizados |
| GET | `/api/v1/mentor/cursos/questoes?busca=` | Questões ativas do mentor, até 100 resultados |
| PUT | `/api/v1/cursos/{cursoId}/aulas/{aulaId}/progresso` | Aluno: `{posicao, duracao, concluida?}` |
| GET | `/api/v1/cursos/{cursoId}/aulas/{aulaId}/questoes` | Questões da aula autorizada |
| POST | `/api/v1/cursos/{cursoId}/aulas/{aulaId}/questoes/{questaoId}/responder` | `{resposta}`; prévia do mentor não persiste |

`alunoId=eu` representa o usuário autenticado. O mentor também pode consultar os cursos disponíveis de um aluno vinculado pela rota de aluno.

## Verificação

Frontend: `npm run build` e `npm run lint`. Backend: `go test ./...` e `go vet ./...`.

O teste `TestIntegracaoAcessoCursos` é opcional: em `web/server-go`, defina `$env:CURSOS_TEST_MYSQL='1'` e execute `go test ./internal/cursos -count=1`. Ele usa a conexão do `.env`, cria um banco temporário exclusivo `test_cursos_*`, aplica apenas o esquema necessário ao domínio (`001`, `002`, `017`, `019`, `020`, `021`) e remove somente esse banco ao finalizar. Verifica isolamento entre mentores, três escopos, PDFs, rascunhos, publicação, conflito de revisão, progresso, exercícios, prévia, acompanhamento, expiração e revogação. Testes unitários verificam liberação no fuso horário e identidade ao reordenar. Esse teste não valida o executor de todas as migrações da aplicação.

Para testar agrupamento e preservação de 100 aulas e PDFs: em `web/client-react`, execute `node --experimental-strip-types --test tests/cursos-modulos.test.mjs`.

### Resultado desta implementação

- `npm run build`: aprovado; aviso de bundle acima de 500 kB.
- `npm run lint`: aprovado com cinco avisos preexistentes em `WhiteLabelContext.tsx` e `NovoConcursoPage.tsx`; nenhum aviso nos arquivos novos.
- `go test ./...` e `go vet ./...`: aprovados, com integração MySQL desabilitada por padrão.
- Teste MySQL com opt-in e esquema específico do domínio: aprovado, incluindo publicação, progresso, questões, prévia, acompanhamento, autorização e PDFs. A migração preexistente `009_isolamento_concursos.sql` continua com erro no executor geral ao criar um banco vazio; ela não faz parte desse teste de domínio. A atualização do banco local existente aplicou a migração 021 ao iniciar a nova API.
- Três testes de módulos aprovados, incluindo preservação da ordem e dos PDFs de 100 aulas distribuídas em 10 módulos, IDs, questões e regras após reordenar/renomear.
- Validação visual e reprodução real: pendentes; nenhum navegador estava conectado no ambiente.
