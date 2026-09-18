package transporte

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"chega-junto-concurseiro-web/internal/autenticacao"
	"chega-junto-concurseiro-web/internal/identificador"
	"chega-junto-concurseiro-web/internal/modelos"
	"chega-junto-concurseiro-web/internal/ratelimit"

	"golang.org/x/time/rate"
)

const nomeCookieSessao = "sessao_chega_junto_concurseiro"

type manipuladorAutenticado func(http.ResponseWriter, *http.Request, modelos.ContextoAutenticado)

func (s *Servidor) comAutenticacao(proximo manipuladorAutenticado) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contextoAutenticado, err := s.autenticacao.Autenticar(r.Context(), valorCookie(r))
		if err != nil {
			if errors.Is(err, autenticacao.ErrPlanoExpirado) {
				responderErro(w, http.StatusForbidden, "PLANO_EXPIRADO", "O acesso deste aluno expirou.")
				return
			}
			responderErro(w, http.StatusUnauthorized, "NAO_AUTENTICADO", "Sessão inválida ou expirada.")
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions && s.autenticacao.HashCSRF(r.Header.Get("X-Token-CSRF")) != contextoAutenticado.CSRFHash {
			responderErro(w, http.StatusForbidden, "CSRF_INVALIDO", "Token CSRF inválido.")
			return
		}
		proximo(w, r, contextoAutenticado)
	}
}

func (s *Servidor) somenteMestre(proximo manipuladorAutenticado) http.HandlerFunc {
	return s.comAutenticacao(func(w http.ResponseWriter, r *http.Request, contexto modelos.ContextoAutenticado) {
		if contexto.Usuario.Papel != "mestre" {
			responderErro(w, http.StatusForbidden, "ACESSO_NEGADO", "Acesso restrito ao mestre.")
			return
		}
		proximo(w, r, contexto)
	})
}

func (s *Servidor) somenteMentor(proximo manipuladorAutenticado) http.HandlerFunc {
	return s.comAutenticacao(func(w http.ResponseWriter, r *http.Request, contexto modelos.ContextoAutenticado) {
		if contexto.Usuario.Papel != "mentor" {
			responderErro(w, http.StatusForbidden, "ACESSO_NEGADO", "Acesso restrito ao mentor.")
			return
		}
		proximo(w, r, contexto)
	})
}

func (s *Servidor) cors(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.cfg.OrigemFrontend)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Token-CSRF, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
		w.Header().Set("Vary", "Origin")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		proximo.ServeHTTP(w, r)
	})
}

type respostaRegistrada struct {
	http.ResponseWriter
	status  int
	tamanho int
}

func (r *respostaRegistrada) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *respostaRegistrada) Write(conteudo []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(conteudo)
	r.tamanho += n
	return n, err
}

func (s *Servidor) registrar(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		idRequisicao := r.Header.Get("X-Request-ID")
		if idRequisicao == "" {
			idRequisicao = identificador.UUID()
		}
		w.Header().Set("X-Request-ID", idRequisicao)
		contexto, cancelar := context.WithTimeout(r.Context(), 25*time.Second)
		defer cancelar()
		resposta := &respostaRegistrada{ResponseWriter: w}
		proximo.ServeHTTP(resposta, r.WithContext(contexto))
		if resposta.status == 0 {
			resposta.status = http.StatusOK
		}
		s.log.Info("requisicao concluida", "request_id", idRequisicao, "metodo", r.Method, "caminho", r.URL.Path, "status", resposta.status, "bytes", resposta.tamanho, "duracao_ms", time.Since(inicio).Milliseconds())
	})
}

func valorCookie(r *http.Request) string {
	cookie, err := r.Cookie(nomeCookieSessao)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func limparCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: nomeCookieSessao, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
}

// ipCliente extrai o IP real do cliente, respeitando X-Forwarded-For em produção.
func ipCliente(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		partes := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(partes[0])
	}
	addr := r.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i != -1 {
		return addr[:i]
	}
	return addr
}

// onLimitado escreve a resposta HTTP 429 padrão com header Retry-After.
func onLimitado(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Retry-After", "60")
	responderErro(w, http.StatusTooManyRequests, "LIMITE_EXCEDIDO", "Muitas requisições. Tente novamente em breve.")
}

// limitarGlobal aplica rate limiting por IP a todas as rotas.
func (s *Servidor) limitarGlobal() func(http.Handler) http.Handler {
	l := ratelimit.Novo(ratelimit.Config{
		Requisicoes: rate.Limit(s.cfg.RateGlobalRPS),
		Burst:       s.cfg.RateGlobalBurst,
	}, 10*time.Minute)
	return l.Middleware(ipCliente, onLimitado)
}

// limitarLogin aplica rate limiting mais restritivo por IP para a rota de login.
func (s *Servidor) limitarLogin() func(http.Handler) http.Handler {
	l := ratelimit.Novo(ratelimit.Config{
		Requisicoes: rate.Limit(s.cfg.RateLoginRPS),
		Burst:       s.cfg.RateLoginBurst,
	}, 10*time.Minute)
	return l.Middleware(ipCliente, onLimitado)
}
