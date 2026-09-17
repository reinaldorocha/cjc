package transporte

import (
	"encoding/json"
	"net/http"
	"track-concursos-web/internal/modelos"
	"track-concursos-web/internal/simulados"
)

func (s *Servidor) listarSimulados(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.simulados.Listar(r.Context(), aluno, r.URL.Query().Get("concursoId"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os simulados.")
		return
	}
	responder(w, 200, map[string]any{"simulados": lista})
}
func (s *Servidor) criarSimulado(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e simulados.Entrada
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.simulados.Criar(r.Context(), aluno, e)
	if err != nil {
		responderErro(w, 400, "SIMULADO_INVALIDO", "Não foi possível registrar o simulado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "criar", "simulado", id, e)
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarSimulado(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e simulados.Entrada
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.simulados.Alterar(r.Context(), aluno, r.PathValue("simuladoId"), e); err != nil {
		responderErro(w, 400, "SIMULADO_INVALIDO", "Não foi possível alterar o simulado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "alterar", "simulado", r.PathValue("simuladoId"), e)
	responder(w, 200, map[string]any{"atualizado": true})
}
func (s *Servidor) desativarSimulado(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	if err := s.simulados.Desativar(r.Context(), aluno, r.PathValue("simuladoId")); err != nil {
		responderErro(w, 404, "SIMULADO_NAO_ENCONTRADO", "Simulado não encontrado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "desativar", "simulado", r.PathValue("simuladoId"), nil)
	responder(w, 200, map[string]any{"desativado": true})
}
func (s *Servidor) obterConfiguracaoProva(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	config, err := s.simulados.ObterConfiguracao(r.Context(), aluno, r.PathValue("concursoId"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar a configuração da prova.")
		return
	}
	responder(w, 200, map[string]any{"configuracao": config})
}
func (s *Servidor) salvarConfiguracaoProva(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e struct {
		Configuracao json.RawMessage `json:"configuracao"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.simulados.SalvarConfiguracao(r.Context(), aluno, r.PathValue("concursoId"), c.Usuario.ID, e.Configuracao); err != nil {
		responderErro(w, 400, "CONFIGURACAO_INVALIDA", "Não foi possível salvar a configuração da prova.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "salvar", "configuracao_prova", r.PathValue("concursoId"), nil)
	responder(w, 200, map[string]any{"atualizada": true})
}
