package cursos

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrEntrada = errors.New("curso inválido: confira título, aulas, links e destinatários")
var ErrNaoEncontrado = errors.New("curso não encontrado")
var ErrConflito = errors.New("o curso foi alterado em outra janela; reabra o editor para carregar a versão atual")
var ErrBloqueado = errors.New("aula ainda não liberada")

type Aula struct {
	ID              string     `json:"id"`
	ModuloID        string     `json:"moduloId"`
	ModuloCapaURL   string     `json:"moduloCapaUrl,omitempty"`
	ModuloDescricao string     `json:"moduloDescricao,omitempty"`
	LiberarEm       string     `json:"liberarEm,omitempty"`
	ExigeAnterior   bool       `json:"exigeAnterior,omitempty"`
	Questoes        []string   `json:"questoes"`
	Bloqueada       bool       `json:"bloqueada,omitempty"`
	MotivoBloqueio  string     `json:"motivoBloqueio,omitempty"`
	Progresso       *Progresso `json:"progresso,omitempty"`
	Titulo          string     `json:"titulo"`
	Modulo          string     `json:"modulo"`
	VideoID         string     `json:"videoId"`
	PDFID           string     `json:"pdfId,omitempty"`
	PDFNome         string     `json:"pdfNome,omitempty"`
}
type Curso struct {
	Revisao             int      `json:"revisao"`
	Publicado           bool     `json:"publicado"`
	AlteracoesPendentes bool     `json:"alteracoesPendentes"`
	Resumo              Resumo   `json:"resumo"`
	MentorID            string   `json:"-"`
	ID                  string   `json:"id"`
	Titulo              string   `json:"titulo"`
	Descricao           string   `json:"descricao"`
	Categoria           string   `json:"categoria"`
	CapaURL             string   `json:"capaUrl"`
	ModoExibicao        string   `json:"modoExibicao,omitempty"`
	Escopo              string   `json:"escopo"`
	Destinatarios       []string `json:"destinatarios"`
	Aulas               []Aula   `json:"aulas"`
}
type Servico struct {
	db    *sql.DB
	agora func() time.Time
}

func Novo(db *sql.DB) *Servico { return &Servico{db: db, agora: time.Now} }
func NovoComFuso(db *sql.DB, fuso *time.Location) *Servico {
	return &Servico{db: db, agora: func() time.Time { return time.Now().In(fuso) }}
}

var videoValido = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)
var uuidValido = regexp.MustCompile(`^[a-fA-F0-9-]{36}$`)

// Somente IDs e URLs de hosts oficiais são convertidos para embeds.
func videoID(valor string) (string, bool) {
	valor = strings.TrimSpace(valor)
	if videoValido.MatchString(valor) {
		return valor, true
	}
	u, err := url.Parse(valor)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.User != nil {
		return "", false
	}
	host := strings.ToLower(u.Host)
	partes := strings.Split(strings.Trim(u.Path, "/"), "/")
	id := ""
	switch host {
	case "youtu.be":
		if len(partes) == 1 {
			id = partes[0]
		}
	case "youtube.com", "www.youtube.com", "m.youtube.com", "youtube-nocookie.com", "www.youtube-nocookie.com":
		if u.Path == "/watch" {
			id = u.Query().Get("v")
		} else if len(partes) == 2 && (partes[0] == "embed" || partes[0] == "shorts" || partes[0] == "live") {
			id = partes[1]
		}
	}
	return id, videoValido.MatchString(id)
}
func validar(c *Curso) error {
	c.Titulo = strings.TrimSpace(c.Titulo)
	c.Descricao = strings.TrimSpace(c.Descricao)
	c.Categoria = strings.TrimSpace(c.Categoria)
	c.CapaURL = strings.TrimSpace(c.CapaURL)
	if c.ModoExibicao == "" {
		c.ModoExibicao = "curso"
	}
	if c.ModoExibicao != "curso" && c.ModoExibicao != "modulos" {
		return ErrEntrada
	}
	if c.Categoria == "" {
		c.Categoria = "Geral"
	}
	if c.Titulo == "" || utf8.RuneCountInString(c.Titulo) > 180 || utf8.RuneCountInString(c.Categoria) > 80 || len(c.Descricao) > 16000 || len(c.Aulas) > 300 || len(c.Destinatarios) > 1000 {
		return ErrEntrada
	}
	if c.CapaURL != "" && idCapa(c.CapaURL) == "" {
		u, err := url.Parse(c.CapaURL)
		if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || len(c.CapaURL) > 2048 {
			return ErrEntrada
		}
	}
	switch c.Escopo {
	case "global":
		c.Destinatarios = []string{}
	case "alunos", "concursos":
	default:
		return ErrEntrada
	}
	vistos := map[string]bool{}
	ids := []string{}
	for _, id := range c.Destinatarios {
		if !uuidValido.MatchString(id) {
			return ErrEntrada
		}
		if !vistos[id] {
			ids = append(ids, id)
			vistos[id] = true
		}
	}
	c.Destinatarios = ids
	for i := range c.Aulas {
		a := &c.Aulas[i]
		a.Titulo = strings.TrimSpace(a.Titulo)
		a.Modulo = strings.TrimSpace(a.Modulo)
		a.ModuloCapaURL = strings.TrimSpace(a.ModuloCapaURL)
		a.ModuloDescricao = strings.TrimSpace(a.ModuloDescricao)
		if len(a.ModuloDescricao) > 4000 {
			return ErrEntrada
		}
		if a.ModuloCapaURL != "" && idCapa(a.ModuloCapaURL) == "" {
			u, err := url.Parse(a.ModuloCapaURL)
			if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || len(a.ModuloCapaURL) > 2048 {
				return ErrEntrada
			}
		}
		id, ok := videoID(a.VideoID)
		if !ok || a.Titulo == "" || utf8.RuneCountInString(a.Titulo) > 180 || utf8.RuneCountInString(a.Modulo) > 100 {
			return ErrEntrada
		}
		a.VideoID = id
		if a.PDFID != "" && !uuidValido.MatchString(a.PDFID) {
			return ErrEntrada
		}
		if a.PDFID == "" {
			a.PDFNome = ""
		}
		if a.Modulo == "" {
			a.Modulo = "Aulas"
		}
		if a.LiberarEm != "" {
			if _, err := time.Parse("2006-01-02", a.LiberarEm); err != nil {
				return ErrEntrada
			}
		}
		if len(a.Questoes) > 50 {
			return ErrEntrada
		}
		vistosQ := map[string]bool{}
		for _, q := range a.Questoes {
			if !uuidValido.MatchString(q) || vistosQ[q] {
				return ErrEntrada
			}
			vistosQ[q] = true
		}
		a.Progresso = nil
		a.Bloqueada = false
		a.MotivoBloqueio = ""
	}
	return normalizarAulas(c)
}

func validarPublicacao(c Curso) error {
	if len(c.Aulas) == 0 || (c.Escopo != "global" && len(c.Destinatarios) == 0) {
		return ErrEntrada
	}
	return validar(&c)
}

func (s *Servico) Salvar(ctx context.Context, mentor, id string, c Curso) (string, error) {
	if err := validar(&c); err != nil {
		return "", err
	}
	return s.salvar(ctx, mentor, id, c)
}
