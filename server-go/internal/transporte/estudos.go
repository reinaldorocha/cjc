package transporte

import (
	"net/http"
	"chega-junto-concurseiro-web/internal/estudos"
	"chega-junto-concurseiro-web/internal/modelos"
)

func (s *Servidor) listarSessoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.estudos.ListarSessoes(r.Context(), aluno, r.URL.Query().Get("concursoId"), r.URL.Query().Get("inicio"), r.URL.Query().Get("fim"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar as sessões de estudo.")
		return
	}
	responder(w, 200, map[string]any{"sessoes": lista})
}
func (s *Servidor) criarSessao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e estudos.EntradaSessao
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.estudos.CriarSessao(r.Context(), aluno, e)
	if err != nil {
		responderErro(w, 400, "SESSAO_INVALIDA", "Não foi possível registrar a sessão de estudo.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "criar", "sessao_estudo", id, e)
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarSessao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e estudos.AlteracaoSessao
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.estudos.AlterarSessao(r.Context(), aluno, r.PathValue("sessaoId"), e); err != nil {
		responderErro(w, 400, "SESSAO_INVALIDA", "Não foi possível alterar a sessão de estudo.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "alterar", "sessao_estudo", r.PathValue("sessaoId"), e)
	responder(w, 200, map[string]any{"atualizada": true})
}
func (s *Servidor) desativarSessao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	if err := s.estudos.Desativar(r.Context(), aluno, "sessoes_estudo", r.PathValue("sessaoId")); err != nil {
		responderErro(w, 404, "SESSAO_NAO_ENCONTRADA", "Sessão de estudo não encontrada.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "desativar", "sessao_estudo", r.PathValue("sessaoId"), nil)
	responder(w, 200, map[string]any{"desativada": true})
}
func (s *Servidor) listarQuestoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.estudos.ListarQuestoes(r.Context(), aluno, r.URL.Query().Get("concursoId"), r.URL.Query().Get("inicio"), r.URL.Query().Get("fim"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os registros de questões.")
		return
	}
	responder(w, 200, map[string]any{"registros": lista})
}
func (s *Servidor) criarQuestoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e estudos.EntradaQuestoes
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.estudos.CriarQuestoes(r.Context(), aluno, e)
	if err != nil {
		responderErro(w, 400, "QUESTOES_INVALIDAS", "Não foi possível registrar as questões.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "criar", "registro_questoes", id, e)
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarQuestoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e estudos.AlteracaoQuestoes
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.estudos.AlterarQuestoes(r.Context(), aluno, r.PathValue("registroId"), e); err != nil {
		responderErro(w, 400, "QUESTOES_INVALIDAS", "Não foi possível alterar o registro de questões.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "alterar", "registro_questoes", r.PathValue("registroId"), e)
	responder(w, 200, map[string]any{"atualizado": true})
}
func (s *Servidor) desativarQuestoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	if err := s.estudos.Desativar(r.Context(), aluno, "registros_questoes", r.PathValue("registroId")); err != nil {
		responderErro(w, 404, "REGISTRO_NAO_ENCONTRADO", "Registro de questões não encontrado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "desativar", "registro_questoes", r.PathValue("registroId"), nil)
	responder(w, 200, map[string]any{"desativado": true})
}
