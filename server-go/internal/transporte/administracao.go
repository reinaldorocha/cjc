package transporte

import (
	"net/http"
	"chega-junto-concurseiro-web/internal/modelos"
	"chega-junto-concurseiro-web/internal/usuarios"
)

func (s *Servidor) listarUsuarios(w http.ResponseWriter, r *http.Request, _ modelos.ContextoAutenticado) {
	lista, err := s.usuarios.Listar(r.Context())
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível listar usuários.")
		return
	}
	responder(w, 200, map[string]any{"usuarios": lista})
}
func (s *Servidor) criarUsuario(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e usuarios.Criacao
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.usuarios.Criar(r.Context(), e)
	if err != nil {
		responderErro(w, 400, "USUARIO_INVALIDO", "Não foi possível criar o usuário.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, "", "criar", "usuario", id, map[string]any{"papel": e.Papel})
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarUsuario(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e struct {
		Nome  *string `json:"nome"`
		Ativo *bool   `json:"ativo"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.usuarios.Alterar(r.Context(), r.PathValue("usuarioId"), e.Nome, e.Ativo, c.Usuario.ID); err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível alterar o usuário.")
		return
	}
	responder(w, 200, map[string]any{"atualizado": true})
}
func (s *Servidor) reativarUsuario(w http.ResponseWriter, r *http.Request, _ modelos.ContextoAutenticado) {
	if err := s.usuarios.Reativar(r.Context(), r.PathValue("usuarioId")); err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível reativar.")
		return
	}
	responder(w, 200, map[string]any{"reativado": true})
}
func (s *Servidor) redefinirSenha(w http.ResponseWriter, r *http.Request, _ modelos.ContextoAutenticado) {
	var e struct {
		NovaSenha string `json:"novaSenha"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.usuarios.RedefinirSenha(r.Context(), r.PathValue("usuarioId"), e.NovaSenha); err != nil {
		responderErro(w, 400, "SENHA_INVALIDA", "Não foi possível redefinir a senha.")
		return
	}
	responder(w, 200, map[string]any{"redefinida": true})
}
func (s *Servidor) definirMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e struct {
		MentorID string `json:"mentorId"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	aluno := r.PathValue("alunoId")
	if err := s.mentoria.DefinirMentor(r.Context(), aluno, e.MentorID, c.Usuario.ID); err != nil {
		responderErro(w, 400, "VINCULO_INVALIDO", "Não foi possível definir o mentor.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "definir_mentor", "mentor_alunos", aluno, map[string]any{"mentorId": e.MentorID})
	responder(w, 200, map[string]any{"definido": true})
}
