package transporte

import (
	"net/http"

	"track-concursos-web/internal/bancoquestoes"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) listarBancoQuestoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	disciplina := r.URL.Query().Get("disciplina")
	assunto := r.URL.Query().Get("assunto")
	tipo := r.URL.Query().Get("tipo")
	concursoID := r.URL.Query().Get("concursoId")

	alunoID := ""
	if c.Usuario.Papel == "aluno" {
		alunoID = c.Usuario.ID
	} else {
		solicitado := r.URL.Query().Get("alunoId")
		if solicitado != "" {
			var ok bool
			alunoID, ok = s.resolverAluno(r.Context(), w, c, solicitado)
			if !ok {
				return
			}
		}
	}

	questoes, err := s.bancoQuestoes.ListarQuestoes(r.Context(), bancoquestoes.FiltroQuestao{
		Disciplina: disciplina,
		Assunto:    assunto,
		Tipo:       tipo,
		ConcursoID: concursoID,
		AlunoID:    alunoID,
	})
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "ERRO_INTERNO", err.Error())
		return
	}
	if pagina, limite, ok := paginacaoSolicitada(r); ok {
		itens, dados := paginar(questoes, pagina, limite)
		responder(w, http.StatusOK, map[string]any{"questoes": itens, "paginacao": dados})
		return
	}
	responder(w, http.StatusOK, map[string]any{"questoes": questoes})
}

func (s *Servidor) criarQuestaoBanco(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var entrada bancoquestoes.EntradaQuestao
	if !decodificar(w, r, &entrada) {
		return
	}

	q, err := s.bancoQuestoes.CriarQuestao(r.Context(), entrada, c.Usuario.ID)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "DADOS_INVALIDOS", err.Error())
		return
	}
	responder(w, http.StatusCreated, map[string]any{"questao": q})
}

func (s *Servidor) importarQuestoesBanco(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var payload struct {
		Questoes []bancoquestoes.EntradaQuestao `json:"questoes"`
	}
	if !decodificar(w, r, &payload) {
		return
	}

	total, err := s.bancoQuestoes.ImportarQuestoes(r.Context(), payload.Questoes, c.Usuario.ID)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "ERRO_IMPORTACAO", err.Error())
		return
	}
	responder(w, http.StatusOK, map[string]any{"mensagem": "questões importadas com sucesso", "total": total})
}

func (s *Servidor) alterarQuestaoBanco(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	questaoID := r.PathValue("questaoId")
	if questaoID == "" {
		responderErro(w, http.StatusBadRequest, "PARAMETRO_INVALIDO", "ID da questão não informado")
		return
	}

	var entrada bancoquestoes.EntradaQuestao
	if !decodificar(w, r, &entrada) {
		return
	}

	if err := s.bancoQuestoes.AlterarQuestao(r.Context(), questaoID, entrada); err != nil {
		responderErro(w, http.StatusBadRequest, "ERRO_ALTERACAO", err.Error())
		return
	}
	responder(w, http.StatusOK, map[string]any{"mensagem": "questão alterada com sucesso"})
}

func (s *Servidor) desativarQuestaoBanco(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	questaoID := r.PathValue("questaoId")
	if questaoID == "" {
		responderErro(w, http.StatusBadRequest, "PARAMETRO_INVALIDO", "ID da questão não informado")
		return
	}

	if err := s.bancoQuestoes.DesativarQuestao(r.Context(), questaoID); err != nil {
		responderErro(w, http.StatusBadRequest, "ERRO_DESATIVACAO", err.Error())
		return
	}
	responder(w, http.StatusOK, map[string]any{"mensagem": "questão desativada com sucesso"})
}

func (s *Servidor) responderQuestaoBanco(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	alunoID, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}

	var entrada bancoquestoes.RespostaQuestaoEntrada
	if !decodificar(w, r, &entrada) {
		return
	}

	resultado, err := s.bancoQuestoes.Responder(r.Context(), alunoID, entrada)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "ERRO_RESPOSTA", err.Error())
		return
	}
	responder(w, http.StatusOK, resultado)
}

func (s *Servidor) historicoQuestaoBanco(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	alunoID, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}

	historico, err := s.bancoQuestoes.HistoricoQuestao(r.Context(), alunoID, r.PathValue("questaoId"), r.URL.Query().Get("concursoId"))
	if err != nil {
		responderErro(w, http.StatusBadRequest, "ERRO_HISTORICO", err.Error())
		return
	}
	responder(w, http.StatusOK, map[string]any{"historico": historico})
}

func (s *Servidor) estatisticasBancoQuestoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	alunoID, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}

	concursoID := r.URL.Query().Get("concursoId")
	stats, err := s.bancoQuestoes.ObterEstatisticas(r.Context(), alunoID, concursoID)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "ERRO_INTERNO", err.Error())
		return
	}
	responder(w, http.StatusOK, stats)
}
