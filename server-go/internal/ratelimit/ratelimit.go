package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Config define os parâmetros de um limitador.
type Config struct {
	// Requisicoes é a taxa de tokens por segundo (use rate.Every para frações).
	Requisicoes rate.Limit
	// Burst é o número máximo de tokens acumulados (pico instantâneo permitido).
	Burst int
}

type entrada struct {
	limiter    *rate.Limiter
	ultimoUso  time.Time
	mu         sync.Mutex
}

// Limitador mantém um Token Bucket por chave (ex.: IP do cliente).
type Limitador struct {
	entradas sync.Map
	cfg      Config
	ttl      time.Duration
}

// Novo cria um Limitador com limpeza automática de entradas inativas.
func Novo(cfg Config, ttl time.Duration) *Limitador {
	l := &Limitador{cfg: cfg, ttl: ttl}
	go l.limpar()
	return l
}

// Permitir retorna true se a chave ainda está dentro do limite.
func (l *Limitador) Permitir(chave string) bool {
	v, _ := l.entradas.LoadOrStore(chave, &entrada{
		limiter: rate.NewLimiter(l.cfg.Requisicoes, l.cfg.Burst),
	})
	e := v.(*entrada)
	e.mu.Lock()
	e.ultimoUso = time.Now()
	ok := e.limiter.Allow()
	e.mu.Unlock()
	return ok
}

// limpar remove entradas que não foram usadas recentemente.
func (l *Limitador) limpar() {
	ticker := time.NewTicker(l.ttl)
	defer ticker.Stop()
	for range ticker.C {
		agora := time.Now()
		l.entradas.Range(func(chave, valor any) bool {
			e := valor.(*entrada)
			e.mu.Lock()
			inativo := agora.Sub(e.ultimoUso) > l.ttl
			e.mu.Unlock()
			if inativo {
				l.entradas.Delete(chave)
			}
			return true
		})
	}
}

// Middleware retorna um handler HTTP que aplica rate limiting por chave.
// A função extrairChave deve retornar a chave de identificação da requisição (ex.: IP).
// onLimitado é chamado quando o limite é excedido; deve escrever a resposta e retornar.
func (l *Limitador) Middleware(
	extrairChave func(*http.Request) string,
	onLimitado func(http.ResponseWriter, *http.Request),
) func(http.Handler) http.Handler {
	return func(proximo http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			chave := extrairChave(r)
			if !l.Permitir(chave) {
				onLimitado(w, r)
				return
			}
			proximo.ServeHTTP(w, r)
		})
	}
}
