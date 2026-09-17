package transporte

import (
	"errors"
	"net/http"
	"time"
	"track-concursos-web/internal/autenticacao"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) saude(w http.ResponseWriter, _ *http.Request) {
	if err := s.banco.Ping(); err != nil {
		responderErro(w, 503, "BANCO_INDISPONIVEL", "Banco de dados indisponível.")
		return
	}
	responder(w, 200, map[string]any{"situacao": "ok"})
}
func (s *Servidor) entrar(w http.ResponseWriter, r *http.Request) {
	var e struct {
		Email string `json:"email"`
		Senha string `json:"senha"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	sessao, err := s.autenticacao.Entrar(r.Context(), e.Email, e.Senha)
	if err != nil {
		if errors.Is(err, autenticacao.ErrPlanoExpirado) {
			responderErro(w, 403, "PLANO_EXPIRADO", "O acesso deste aluno expirou.")
			return
		}
		responderErro(w, 401, "CREDENCIAIS_INVALIDAS", "E-mail ou senha inválidos.")
		return
	}
	expira := time.Now().UTC().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{Name: nomeCookieSessao, Value: sessao.Token, Path: "/", HttpOnly: true, Secure: s.cfg.Ambiente == "producao", SameSite: http.SameSiteLaxMode, Expires: expira})
	responder(w, 200, map[string]any{"usuario": sessao.Usuario, "tokenCsrf": sessao.TokenCSRF})
}
func (s *Servidor) sair(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	if err := s.autenticacao.Revogar(r.Context(), valorCookie(r)); err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível encerrar a sessão.")
		return
	}
	limparCookie(w)
	s.registrarAuditoria(r.Context(), c.Usuario.ID, "", "sair", "sessao", "", nil)
	responder(w, 200, map[string]any{"encerrada": true})
}
func (s *Servidor) sairDeTodas(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	if err := s.autenticacao.RevogarTodas(r.Context(), c.Usuario.ID); err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível encerrar as sessões.")
		return
	}
	limparCookie(w)
	responder(w, 200, map[string]any{"encerradas": true})
}
func (s *Servidor) eu(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	tokenCSRF, err := s.autenticacao.RenovarCSRF(r.Context(), valorCookie(r))
	if err != nil {
		responderErro(w, 401, "NAO_AUTENTICADO", "Sessão inválida ou expirada.")
		return
	}
	responder(w, 200, map[string]any{"usuario": c.Usuario, "tokenCsrf": tokenCSRF})
}
func (s *Servidor) alterarSenha(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e struct {
		SenhaAtual string `json:"senhaAtual"`
		NovaSenha  string `json:"novaSenha"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.autenticacao.AlterarSenha(r.Context(), c.Usuario, e.SenhaAtual, e.NovaSenha); err != nil {
		responderErro(w, 400, "SENHA_INVALIDA", "Não foi possível alterar a senha.")
		return
	}
	responder(w, 200, map[string]any{"alterada": true})
}
func (s *Servidor) alterarPerfil(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e struct {
		Nome string `json:"nome"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.usuarios.AlterarNomeProprio(r.Context(), c.Usuario.ID, e.Nome); err != nil {
		responderErro(w, 400, "NOME_INVALIDO", "Nome inválido.")
		return
	}
	responder(w, 200, map[string]any{"atualizado": true})
}
