package transporte

import (
	"net/http"
	"track-concursos-web/internal/cursos"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) publicarCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e struct {
		Revisao  int  `json:"revisao"`
		Publicar bool `json:"publicar"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.cursos.Publicar(r.Context(), c.Usuario.ID, r.PathValue("cursoId"), e.Revisao, e.Publicar); err != nil {
		s.erroCurso(w, err)
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, "", "publicacao", "curso", r.PathValue("cursoId"), e)
	s.listarCursos(w, r, c)
}
func (s *Servidor) previaCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	curso, err := s.cursos.Previa(r.Context(), c.Usuario.ID, r.PathValue("cursoId"))
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	responder(w, 200, map[string]any{"curso": curso})
}
func (s *Servidor) progressoCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	if c.Usuario.Papel != "aluno" {
		responderErro(w, 403, "ACESSO_NEGADO", "Somente o aluno pode registrar seu progresso.")
		return
	}
	var e cursos.EntradaProgresso
	if !decodificar(w, r, &e) {
		return
	}
	curso, err := s.cursos.SalvarProgresso(r.Context(), c.Usuario.ID, r.PathValue("cursoId"), r.PathValue("aulaId"), e)
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	responder(w, 200, map[string]any{"curso": curso})
}
func (s *Servidor) acompanharCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	alunos, err := s.cursos.Acompanhar(r.Context(), c.Usuario.ID, r.PathValue("cursoId"))
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	responder(w, 200, map[string]any{"alunos": alunos})
}
func (s *Servidor) catalogoQuestoesCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	lista, err := s.cursos.CatalogoQuestoes(r.Context(), c.Usuario.ID, r.URL.Query().Get("busca"))
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	responder(w, 200, map[string]any{"questoes": lista})
}
func (s *Servidor) questoesCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	mentor := c.Usuario.Papel == "mentor"
	if !mentor && c.Usuario.Papel != "aluno" {
		responderErro(w, 403, "ACESSO_NEGADO", "Acesso restrito.")
		return
	}
	lista, err := s.cursos.QuestoesAula(r.Context(), c.Usuario.ID, mentor, r.PathValue("cursoId"), r.PathValue("aulaId"))
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	responder(w, 200, map[string]any{"questoes": lista})
}
func (s *Servidor) responderQuestaoCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	mentor := c.Usuario.Papel == "mentor"
	if !mentor && c.Usuario.Papel != "aluno" {
		responderErro(w, 403, "ACESSO_NEGADO", "Acesso restrito.")
		return
	}
	var e struct {
		Resposta string `json:"resposta"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	resultado, err := s.cursos.ResponderQuestao(r.Context(), c.Usuario.ID, mentor, r.PathValue("cursoId"), r.PathValue("aulaId"), r.PathValue("questaoId"), e.Resposta)
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	responder(w, 200, resultado)
}
