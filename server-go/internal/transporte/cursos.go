package transporte

import (
	"errors"
	"net/http"
	"chega-junto-concurseiro-web/internal/cursos"
	"chega-junto-concurseiro-web/internal/modelos"
)

func (s *Servidor) listarCursos(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	mentor := r.PathValue("alunoId") == ""
	usuario := c.Usuario.ID
	if !mentor {
		var ok bool
		usuario, ok = s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
		if !ok {
			return
		}
	}
	lista, err := s.cursos.Listar(r.Context(), usuario, mentor, r.PathValue("cursoId"))
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	if r.PathValue("cursoId") != "" {
		if len(lista) == 0 {
			s.erroCurso(w, cursos.ErrNaoEncontrado)
			return
		}
		responder(w, 200, map[string]any{"curso": lista[0]})
		return
	}
	if pagina, limite, ok := paginacaoSolicitada(r); ok {
		itens, dados := paginar(lista, pagina, limite)
		responder(w, 200, map[string]any{"cursos": itens, "paginacao": dados})
		return
	}
	responder(w, 200, map[string]any{"cursos": lista})
}
func (s *Servidor) salvarCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var entrada cursos.Curso
	if !decodificar(w, r, &entrada) {
		return
	}
	id, err := s.cursos.Salvar(r.Context(), c.Usuario.ID, r.PathValue("cursoId"), entrada)
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, "", "salvar", "curso", id, map[string]any{"escopo": entrada.Escopo})
	status := 200
	if r.Method == http.MethodPost {
		status = 201
	}
	lista, err := s.cursos.Listar(r.Context(), c.Usuario.ID, true, id)
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	if len(lista) == 0 {
		s.erroCurso(w, cursos.ErrNaoEncontrado)
		return
	}
	responder(w, status, map[string]any{"id": id, "curso": lista[0]})
}
func (s *Servidor) desativarCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	id := r.PathValue("cursoId")
	if err := s.cursos.Desativar(r.Context(), c.Usuario.ID, id); err != nil {
		s.erroCurso(w, err)
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, "", "desativar", "curso", id, nil)
	responder(w, 200, map[string]any{"desativado": true})
}
func (s *Servidor) erroCurso(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, cursos.ErrCapa):
		responderErro(w, 400, "CAPA_INVALIDA", err.Error())
	case errors.Is(err, cursos.ErrConflito):
		responderErro(w, 409, "CURSO_CONFLITO", err.Error())
	case errors.Is(err, cursos.ErrBloqueado):
		responderErro(w, 403, "AULA_BLOQUEADA", err.Error())
	case errors.Is(err, cursos.ErrEntrada):
		responderErro(w, 400, "CURSO_INVALIDO", err.Error())
	case errors.Is(err, cursos.ErrNaoEncontrado):
		responderErro(w, 404, "CURSO_NAO_ENCONTRADO", "Curso não encontrado no seu acesso.")
	default:
		s.log.Error("falha em cursos", "erro", err)
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível processar o curso.")
	}
}
