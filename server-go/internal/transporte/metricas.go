package transporte

import (
	"net/http"
	"chega-junto-concurseiro-web/internal/modelos"
)

func (s *Servidor) resumoMetricas(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	dados, err := s.metricas.Resumo(r.Context(), aluno, r.URL.Query().Get("concursoId"), r.URL.Query().Get("inicio"), r.URL.Query().Get("fim"))
	if err != nil {
		responderErro(w, 400, "METRICAS_INVALIDAS", "Não foi possível calcular o resumo de métricas.")
		return
	}
	responder(w, 200, dados)
}
func (s *Servidor) linhaDoTempoMetricas(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	dados, err := s.metricas.LinhaDoTempo(r.Context(), aluno, r.URL.Query().Get("concursoId"), r.URL.Query().Get("inicio"), r.URL.Query().Get("fim"))
	if err != nil {
		responderErro(w, 400, "PERIODO_INVALIDO", "Não foi possível calcular a linha do tempo.")
		return
	}
	responder(w, 200, map[string]any{"dias": dados})
}
func (s *Servidor) materiasMetricas(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	dados, err := s.metricas.Materias(r.Context(), aluno, r.URL.Query().Get("concursoId"), r.URL.Query().Get("inicio"), r.URL.Query().Get("fim"))
	if err != nil {
		responderErro(w, 400, "METRICAS_INVALIDAS", "Não foi possível calcular as métricas por matéria.")
		return
	}
	responder(w, 200, map[string]any{"materias": dados})
}
