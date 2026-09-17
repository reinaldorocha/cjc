package transporte

import (
	"net/http"
	"track-concursos-web/internal/cadernos"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) listarCadernos(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	editalID := r.URL.Query().Get("editalId")
	lista, err := s.cadernos.Listar(r.Context(), aluno, editalID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os cadernos.")
		return
	}
	responder(w, 200, map[string]any{"cadernos": lista})
}

func (s *Servidor) criarCaderno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e cadernos.EntradaCaderno
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.cadernos.Criar(r.Context(), aluno, e)
	if err != nil {
		responderErro(w, 400, "CADERNO_INVALIDO", err.Error())
		return
	}
	responder(w, 201, map[string]any{"id": id})
}

func (s *Servidor) obterCaderno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	cadernoID := r.PathValue("cadernoId")
	item, err := s.cadernos.Obter(r.Context(), aluno, cadernoID)
	if err != nil {
		responderErro(w, 404, "CADERNO_NAO_ENCONTRADO", "Caderno não encontrado.")
		return
	}
	responder(w, 200, map[string]any{"caderno": item})
}

func (s *Servidor) alterarCaderno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	cadernoID := r.PathValue("cadernoId")
	var e cadernos.EntradaCaderno
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.cadernos.Alterar(r.Context(), aluno, cadernoID, e); err != nil {
		responderErro(w, 400, "CADERNO_INVALIDO", err.Error())
		return
	}
	responder(w, 200, map[string]any{"atualizado": true})
}

func (s *Servidor) desativarCaderno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	cadernoID := r.PathValue("cadernoId")
	if err := s.cadernos.Desativar(r.Context(), aluno, cadernoID); err != nil {
		responderErro(w, 400, "CADERNO_INVALIDO", "Não foi possível remover o caderno.")
		return
	}
	responder(w, 200, map[string]any{"desativado": true})
}
