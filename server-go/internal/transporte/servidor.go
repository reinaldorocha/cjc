package transporte

import (
	"chega-junto-concurseiro-web/internal/auditoria"
	"chega-junto-concurseiro-web/internal/autenticacao"
	"chega-junto-concurseiro-web/internal/bancoquestoes"
	"chega-junto-concurseiro-web/internal/cadernos"
	"chega-junto-concurseiro-web/internal/cartoes"
	"chega-junto-concurseiro-web/internal/concursos"
	"chega-junto-concurseiro-web/internal/configuracao"
	"chega-junto-concurseiro-web/internal/cronogramas"
	"chega-junto-concurseiro-web/internal/cursos"
	"chega-junto-concurseiro-web/internal/editais"
	"chega-junto-concurseiro-web/internal/estudos"
	"chega-junto-concurseiro-web/internal/materiaisapoio"
	"chega-junto-concurseiro-web/internal/mentoria"
	"chega-junto-concurseiro-web/internal/metricas"
	"chega-junto-concurseiro-web/internal/revisoes"
	"chega-junto-concurseiro-web/internal/simulados"
	"chega-junto-concurseiro-web/internal/usuarios"
	"chega-junto-concurseiro-web/internal/whitelabel"
	"context"
	"database/sql"
	"log/slog"
	"net/http"
)

type Servidor struct {
	cursos         *cursos.Servico
	banco          *sql.DB
	cfg            configuracao.Configuracao
	log            *slog.Logger
	autenticacao   *autenticacao.Servico
	usuarios       *usuarios.Servico
	mentoria       *mentoria.Servico
	concursos      *concursos.Servico
	editais        *editais.Servico
	cronogramas    *cronogramas.Servico
	cartoes        *cartoes.Servico
	bancoQuestoes  *bancoquestoes.Servico
	estudos        *estudos.Servico
	materiaisApoio *materiaisapoio.Servico
	cadernos       *cadernos.Servico
	revisoes       *revisoes.Servico
	simulados      *simulados.Servico
	metricas       *metricas.Servico
	auditoria      *auditoria.Servico
	whitelabel     *whitelabel.Servico
}

func Novo(banco *sql.DB, cfg configuracao.Configuracao, log *slog.Logger) *Servidor {
	return &Servidor{
		cursos:         cursos.NovoComFuso(banco, cfg.FusoHorario),
		banco:          banco,
		cfg:            cfg,
		log:            log,
		autenticacao:   autenticacao.Novo(banco, cfg),
		usuarios:       usuarios.Novo(banco),
		mentoria:       mentoria.Novo(banco),
		concursos:      concursos.Novo(banco),
		editais:        editais.Novo(banco),
		cronogramas:    cronogramas.Novo(banco, cfg.FusoHorario),
		cartoes:        cartoes.Novo(banco, cfg.FusoHorario),
		bancoQuestoes:  bancoquestoes.Novo(banco),
		estudos:        estudos.Novo(banco),
		materiaisApoio: materiaisapoio.Novo(banco),
		cadernos:       cadernos.Novo(banco),
		revisoes:       revisoes.Novo(banco),
		simulados:      simulados.Novo(banco),
		metricas:       metricas.Novo(banco, cfg.FusoHorario),
		auditoria:      auditoria.Novo(banco),
		whitelabel:     whitelabel.NovoServico(banco),
	}
}
func (s *Servidor) SemearMestre() error { return s.autenticacao.SemearMestre() }
func (s *Servidor) registrarAuditoria(ctx context.Context, executor, aluno, acao, entidade, entidadeID string, detalhes any) {
	if err := s.auditoria.Registrar(ctx, executor, aluno, acao, entidade, entidadeID, detalhes); err != nil {
		s.log.Error("falha ao registrar auditoria", "erro", err, "acao", acao, "entidade", entidade, "entidadeId", entidadeID)
	}
}
func (s *Servidor) Rotas() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/mentor/cursos/capas", s.somenteMentor(s.enviarCapaCurso))
	mux.HandleFunc("GET /api/v1/cursos/capas/{capaId}", s.comAutenticacao(s.exibirCapaCurso))
	mux.HandleFunc("POST /api/v1/mentor/cursos/{cursoId}/publicacao", s.somenteMentor(s.publicarCurso))
	mux.HandleFunc("GET /api/v1/mentor/cursos/{cursoId}/previa", s.somenteMentor(s.previaCurso))
	mux.HandleFunc("GET /api/v1/mentor/cursos/{cursoId}/acompanhamento", s.somenteMentor(s.acompanharCurso))
	mux.HandleFunc("GET /api/v1/mentor/cursos/questoes", s.somenteMentor(s.catalogoQuestoesCurso))
	mux.HandleFunc("PUT /api/v1/cursos/{cursoId}/aulas/{aulaId}/progresso", s.comAutenticacao(s.progressoCurso))
	mux.HandleFunc("GET /api/v1/cursos/{cursoId}/aulas/{aulaId}/questoes", s.comAutenticacao(s.questoesCurso))
	mux.HandleFunc("POST /api/v1/cursos/{cursoId}/aulas/{aulaId}/questoes/{questaoId}/responder", s.comAutenticacao(s.responderQuestaoCurso))
	mux.HandleFunc("GET /api/v1/mentor/cursos", s.somenteMentor(s.listarCursos))
	mux.HandleFunc("POST /api/v1/mentor/cursos/pdfs", s.somenteMentor(s.enviarPDFCurso))
	mux.HandleFunc("GET /api/v1/mentor/cursos/{cursoId}/pdfs/{pdfId}", s.somenteMentor(s.baixarPDFCurso))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cursos/{cursoId}/pdfs/{pdfId}", s.comAutenticacao(s.baixarPDFCurso))
	mux.HandleFunc("GET /api/v1/mentor/cursos/{cursoId}", s.somenteMentor(s.listarCursos))
	mux.HandleFunc("POST /api/v1/mentor/cursos", s.somenteMentor(s.salvarCurso))
	mux.HandleFunc("PUT /api/v1/mentor/cursos/{cursoId}", s.somenteMentor(s.salvarCurso))
	mux.HandleFunc("DELETE /api/v1/mentor/cursos/{cursoId}", s.somenteMentor(s.desativarCurso))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cursos", s.comAutenticacao(s.listarCursos))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cursos/{cursoId}", s.comAutenticacao(s.listarCursos))
	mux.HandleFunc("GET /api/v1/saude", s.saude)
	mux.HandleFunc("GET /api/v1/mentor/whitelabel", s.somenteMentor(s.obterWhiteLabelMentor))
	mux.HandleFunc("GET /api/v1/mentor/whitelabel/logo", s.somenteMentor(s.exibirLogoWhiteLabelMentor))
	mux.HandleFunc("GET /api/v1/mentor/whitelabel/banner", s.somenteMentor(s.exibirBannerWhiteLabelMentor))
	mux.HandleFunc("POST /api/v1/mentor/whitelabel", s.somenteMentor(s.salvarWhiteLabelMentor))
	mux.HandleFunc("POST /api/v1/mentor/whitelabel/logo", s.somenteMentor(s.uploadLogoWhiteLabel))
	mux.HandleFunc("POST /api/v1/mentor/whitelabel/banner", s.somenteMentor(s.uploadBannerWhiteLabel))
	mux.HandleFunc("GET /api/v1/mentores/{mentorId}/whitelabel", s.comAutenticacao(s.obterWhiteLabelMentorPorID))
	mux.HandleFunc("GET /api/v1/mentores/{mentorId}/whitelabel/logo", s.comAutenticacao(s.exibirLogoWhiteLabelMentorPorID))
	mux.HandleFunc("GET /api/v1/mentores/{mentorId}/whitelabel/banner", s.comAutenticacao(s.exibirBannerWhiteLabelMentorPorID))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/whitelabel", s.comAutenticacao(s.obterWhiteLabelAluno))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/whitelabel/logo", s.comAutenticacao(s.exibirLogoWhiteLabelAluno))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/whitelabel/banner", s.comAutenticacao(s.exibirBannerWhiteLabelAluno))

	// Rota de login com limite mais restritivo (anti brute-force)
	limitLogin := s.limitarLogin()
	mux.Handle("POST /api/v1/autenticacao/entrar", limitLogin(http.HandlerFunc(s.entrar)))

	mux.HandleFunc("POST /api/v1/autenticacao/sair", s.comAutenticacao(s.sair))
	mux.HandleFunc("POST /api/v1/autenticacao/sair-de-todas", s.comAutenticacao(s.sairDeTodas))
	mux.HandleFunc("GET /api/v1/autenticacao/eu", s.comAutenticacao(s.eu))
	mux.HandleFunc("POST /api/v1/autenticacao/alterar-senha", s.comAutenticacao(s.alterarSenha))
	mux.HandleFunc("PATCH /api/v1/autenticacao/perfil", s.comAutenticacao(s.alterarPerfil))
	mux.HandleFunc("GET /api/v1/administracao/usuarios", s.somenteMestre(s.listarUsuarios))
	mux.HandleFunc("POST /api/v1/administracao/usuarios", s.somenteMestre(s.criarUsuario))
	mux.HandleFunc("PATCH /api/v1/administracao/usuarios/{usuarioId}", s.somenteMestre(s.alterarUsuario))
	mux.HandleFunc("POST /api/v1/administracao/usuarios/{usuarioId}/reativar", s.somenteMestre(s.reativarUsuario))
	mux.HandleFunc("POST /api/v1/administracao/usuarios/{usuarioId}/redefinir-senha", s.somenteMestre(s.redefinirSenha))
	mux.HandleFunc("PUT /api/v1/administracao/alunos/{alunoId}/mentor", s.somenteMestre(s.definirMentor))
	mux.HandleFunc("GET /api/v1/mentor/alunos", s.somenteMentor(s.listarAlunos))
	mux.HandleFunc("POST /api/v1/mentor/alunos", s.somenteMentor(s.criarAlunoPeloMentor))
	mux.HandleFunc("GET /api/v1/mentor/radar-alunos", s.somenteMentor(s.radarAlunos))
	mux.HandleFunc("GET /api/v1/mentor/concursos", s.somenteMentor(s.catalogarConcursosMentor))
	mux.HandleFunc("POST /api/v1/mentor/concursos", s.somenteMentor(s.criarConcursoMentor))
	mux.HandleFunc("GET /api/v1/mentor/editais", s.somenteMentor(s.catalogarEditaisMentor))
	mux.HandleFunc("POST /api/v1/mentor/editais", s.somenteMentor(s.criarEditalMentor))
	mux.HandleFunc("GET /api/v1/mentor/baralhos-cartoes", s.somenteMentor(s.listarBaralhosMentor))
	mux.HandleFunc("POST /api/v1/mentor/baralhos-cartoes", s.somenteMentor(s.criarBaralhoMentor))
	mux.HandleFunc("PATCH /api/v1/mentor/baralhos-cartoes/{baralhoId}", s.somenteMentor(s.alterarBaralhoMentor))
	mux.HandleFunc("DELETE /api/v1/mentor/baralhos-cartoes/{baralhoId}", s.somenteMentor(s.desativarBaralhoMentor))
	mux.HandleFunc("POST /api/v1/mentor/baralhos-cartoes/{baralhoId}/cartoes", s.somenteMentor(s.criarCartaoMentor))
	mux.HandleFunc("PATCH /api/v1/mentor/cartoes/{cartaoId}", s.somenteMentor(s.alterarCartaoMentor))
	mux.HandleFunc("DELETE /api/v1/mentor/cartoes/{cartaoId}", s.somenteMentor(s.desativarCartaoMentor))
	mux.HandleFunc("GET /api/v1/mentor/alunos/{alunoId}", s.somenteMentor(s.obterAluno))
	mux.HandleFunc("GET /api/v1/mentor/alunos/{alunoId}/visao-geral", s.somenteMentor(s.visaoGeralAluno))
	mux.HandleFunc("PATCH /api/v1/mentor/alunos/{alunoId}/configuracoes", s.somenteMentor(s.configurarAluno))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/concursos", s.comAutenticacao(s.listarConcursos))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/concursos", s.somenteMentor(s.atribuirConcurso))
	mux.HandleFunc("PUT /api/v1/alunos/{alunoId}/concursos/ordem", s.somenteMentor(s.reordenarConcursos))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/concursos/{concursoId}", s.somenteMentor(s.alterarConcurso))
	mux.HandleFunc("DELETE /api/v1/alunos/{alunoId}/concursos/{concursoId}", s.somenteMentor(s.desativarConcurso))
	mux.HandleFunc("GET /api/v1/banco-questoes", s.comAutenticacao(s.listarBancoQuestoes))
	mux.HandleFunc("POST /api/v1/banco-questoes", s.somenteMentor(s.criarQuestaoBanco))
	mux.HandleFunc("POST /api/v1/banco-questoes/importar", s.somenteMentor(s.importarQuestoesBanco))
	mux.HandleFunc("PATCH /api/v1/banco-questoes/{questaoId}", s.somenteMentor(s.alterarQuestaoBanco))
	mux.HandleFunc("DELETE /api/v1/banco-questoes/{questaoId}", s.somenteMentor(s.desativarQuestaoBanco))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/banco-questoes/responder", s.comAutenticacao(s.responderQuestaoBanco))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/banco-questoes/{questaoId}/historico", s.comAutenticacao(s.historicoQuestaoBanco))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/banco-questoes/estatisticas", s.comAutenticacao(s.estatisticasBancoQuestoes))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/edital", s.comAutenticacao(s.listarEditais))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/materiais-apoio", s.comAutenticacao(s.listarMateriaisApoio))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cadernos", s.comAutenticacao(s.listarCadernos))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/cadernos", s.comAutenticacao(s.criarCaderno))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cadernos/{cadernoId}", s.comAutenticacao(s.obterCaderno))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/cadernos/{cadernoId}", s.comAutenticacao(s.alterarCaderno))
	mux.HandleFunc("DELETE /api/v1/alunos/{alunoId}/cadernos/{cadernoId}", s.comAutenticacao(s.desativarCaderno))
	mux.HandleFunc("GET /api/v1/mentor/materiais-apoio", s.somenteMentor(s.listarMateriaisApoioMentor))
	mux.HandleFunc("GET /api/v1/editais", s.somenteMentor(s.catalogarEditais))
	mux.HandleFunc("PUT /api/v1/alunos/{alunoId}/edital", s.somenteMentor(s.atribuirEdital))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/edital/itens", s.somenteMentor(s.criarItemEdital))
	mux.HandleFunc("PUT /api/v1/alunos/{alunoId}/edital/ordem", s.somenteMentor(s.reordenarItensEdital))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/edital/itens/{itemId}", s.somenteMentor(s.alterarItemEdital))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/edital/progresso/{itemId}", s.comAutenticacao(s.atualizarProgressoEdital))
	mux.HandleFunc("POST /api/v1/mentor/materiais-apoio", s.somenteMentor(s.criarMaterialApoioMentor))
	mux.HandleFunc("PATCH /api/v1/mentor/materiais-apoio/{materialId}", s.somenteMentor(s.alterarMaterialApoioMentor))
	mux.HandleFunc("DELETE /api/v1/mentor/materiais-apoio/{materialId}", s.somenteMentor(s.desativarMaterialApoioMentor))
	mux.HandleFunc("GET /api/v1/mentor/materiais-apoio/{materialId}/arquivo", s.somenteMentor(s.baixarArquivoMaterialApoio))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/materiais-apoio/{materialId}/arquivo", s.comAutenticacao(s.baixarArquivoMaterialApoio))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/edital/materiais/{materialId}/progresso", s.comAutenticacao(s.atualizarProgressoMaterial))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cronograma", s.comAutenticacao(s.obterCronograma))
	mux.HandleFunc("PUT /api/v1/alunos/{alunoId}/cronograma", s.comAutenticacao(s.salvarCronograma))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/cronograma/gerar", s.comAutenticacao(s.gerarCronograma))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/cronograma/itens", s.comAutenticacao(s.criarItemCronograma))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/cronograma/reprogramar", s.comAutenticacao(s.reprogramarCronograma))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/cronograma/itens/{itemId}", s.comAutenticacao(s.alterarItemCronograma))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cronograma/calendario", s.comAutenticacao(s.calendarioCronograma))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/baralhos-cartoes", s.comAutenticacao(s.listarBaralhos))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/baralhos-cartoes", s.comAutenticacao(s.criarBaralho))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/baralhos-cartoes/{baralhoId}", s.comAutenticacao(s.alterarBaralho))
	mux.HandleFunc("DELETE /api/v1/alunos/{alunoId}/baralhos-cartoes/{baralhoId}", s.comAutenticacao(s.desativarBaralho))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/baralhos-cartoes/{baralhoId}/cartoes", s.comAutenticacao(s.criarCartao))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/cartoes/{cartaoId}", s.comAutenticacao(s.alterarCartao))
	mux.HandleFunc("DELETE /api/v1/alunos/{alunoId}/cartoes/{cartaoId}", s.comAutenticacao(s.desativarCartao))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/cartoes/{cartaoId}/revisar", s.comAutenticacao(s.revisarCartao))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cartoes/pendentes", s.comAutenticacao(s.listarCartoesPendentes))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/cartoes/historico", s.comAutenticacao(s.listarHistoricoCartoes))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/sessoes-estudo", s.comAutenticacao(s.listarSessoes))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/sessoes-estudo", s.comAutenticacao(s.criarSessao))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/sessoes-estudo/{sessaoId}", s.comAutenticacao(s.alterarSessao))
	mux.HandleFunc("DELETE /api/v1/alunos/{alunoId}/sessoes-estudo/{sessaoId}", s.comAutenticacao(s.desativarSessao))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/registros-questoes", s.comAutenticacao(s.listarQuestoes))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/registros-questoes", s.comAutenticacao(s.criarQuestoes))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/registros-questoes/{registroId}", s.comAutenticacao(s.alterarQuestoes))
	mux.HandleFunc("DELETE /api/v1/alunos/{alunoId}/registros-questoes/{registroId}", s.comAutenticacao(s.desativarQuestoes))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/revisoes", s.comAutenticacao(s.listarRevisoes))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/revisoes", s.comAutenticacao(s.criarRevisao))
	mux.HandleFunc("PATCH /api/v1/alunos/{alunoId}/revisoes/{revisaoId}", s.comAutenticacao(s.alterarRevisao))
	mux.HandleFunc("DELETE /api/v1/alunos/{alunoId}/revisoes/{revisaoId}", s.comAutenticacao(s.desativarRevisao))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/configuracao-prova/{concursoId}", s.comAutenticacao(s.obterConfiguracaoProva))
	mux.HandleFunc("PUT /api/v1/alunos/{alunoId}/configuracao-prova/{concursoId}", s.comAutenticacao(s.salvarConfiguracaoProva))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/simulados", s.comAutenticacao(s.listarSimulados))
	mux.HandleFunc("POST /api/v1/alunos/{alunoId}/simulados", s.comAutenticacao(s.criarSimulado))
	mux.HandleFunc("PUT /api/v1/alunos/{alunoId}/simulados/{simuladoId}", s.comAutenticacao(s.alterarSimulado))
	mux.HandleFunc("DELETE /api/v1/alunos/{alunoId}/simulados/{simuladoId}", s.comAutenticacao(s.desativarSimulado))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/metricas/resumo", s.comAutenticacao(s.resumoMetricas))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/metricas/linha-do-tempo", s.comAutenticacao(s.linhaDoTempoMetricas))
	mux.HandleFunc("GET /api/v1/alunos/{alunoId}/metricas/materias", s.comAutenticacao(s.materiasMetricas))

	// Cadeia: CORS → Rate Limit Global (por IP) → Logger → Rotas
	if s.cfg.FrontendPath != "" {
		mux.HandleFunc("GET /{caminho...}", s.servirFrontend)
	}
	return s.cors(s.limitarGlobal()(s.registrar(mux)))
}
