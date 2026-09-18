package transporte

import (
	"net/http"
	"chega-junto-concurseiro-web/internal/concursos"
	"chega-junto-concurseiro-web/internal/modelos"
)

func (s *Servidor) listarConcursos(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.concursos.Listar(r.Context(), aluno)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível listar concursos.")
		return
	}
	responder(w, 200, map[string]any{"concursos": lista})
}
func (s *Servidor) atribuirConcurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e concursos.Atribuicao
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.concursos.Atribuir(r.Context(), aluno, c.Usuario.ID, e)
	if err != nil {
		responderErro(w, 400, "CONCURSO_INVALIDO", "Não foi possível atribuir o concurso.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "atribuir", "concurso", id, nil)
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarConcurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e concursos.Alteracao
	if !decodificar(w, r, &e) {
		return
	}
	id := r.PathValue("concursoId")
	if err := s.concursos.Alterar(r.Context(), aluno, id, e); err != nil {
		responderErro(w, 400, "CONCURSO_INVALIDO", "Não foi possível alterar o concurso.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "alterar", "concurso", id, nil)
	responder(w, 200, map[string]any{"atualizado": true})
}
func (s *Servidor) desativarConcurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	id := r.PathValue("concursoId")
	if err := s.concursos.Desativar(r.Context(), aluno, id); err != nil {
		responderErro(w, 404, "CONCURSO_NAO_ENCONTRADO", "Concurso não encontrado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "desativar", "concurso", id, nil)
	responder(w, 200, map[string]any{"desativado": true})
}
func (s *Servidor) reordenarConcursos(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e struct {
		Itens []concursos.Ordem `json:"itens"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.concursos.Reordenar(r.Context(), aluno, e.Itens); err != nil {
		responderErro(w, 400, "ORDEM_INVALIDA", "Não foi possível reordenar os concursos.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "reordenar", "concursos", aluno, e)
	responder(w, 200, map[string]any{"atualizado": true})
}
