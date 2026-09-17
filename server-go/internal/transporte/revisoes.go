package transporte

import (
	"net/http"
	"track-concursos-web/internal/modelos"
	"track-concursos-web/internal/revisoes"
)

func (s *Servidor) listarRevisoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.revisoes.Listar(r.Context(), aluno, r.URL.Query().Get("concursoId"), r.URL.Query().Get("estado"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar as revisões.")
		return
	}
	responder(w, 200, map[string]any{"revisoes": lista})
}
func (s *Servidor) criarRevisao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e revisoes.Entrada
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.revisoes.Criar(r.Context(), aluno, e)
	if err != nil {
		responderErro(w, 400, "REVISAO_INVALIDA", "Não foi possível programar a revisão.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "criar", "revisao_programada", id, e)
	resultado := s.replanejarAposRevisoes(r, aluno, false)
	resultado["id"] = id
	responder(w, 201, resultado)
}
func (s *Servidor) alterarRevisao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e revisoes.Alteracao
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.revisoes.Alterar(r.Context(), aluno, r.PathValue("revisaoId"), e); err != nil {
		responderErro(w, 400, "REVISAO_INVALIDA", "Não foi possível alterar a revisão.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "alterar", "revisao_programada", r.PathValue("revisaoId"), e)
	resultado := s.replanejarAposRevisoes(r, aluno, false)
	resultado["atualizada"] = true
	responder(w, 200, resultado)
}
func (s *Servidor) desativarRevisao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	if err := s.revisoes.Desativar(r.Context(), aluno, r.PathValue("revisaoId")); err != nil {
		responderErro(w, 404, "REVISAO_NAO_ENCONTRADA", "Revisão não encontrada.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "desativar", "revisao_programada", r.PathValue("revisaoId"), nil)
	responder(w, 200, map[string]any{"desativada": true})
}
