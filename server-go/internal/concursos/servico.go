package concursos

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"track-concursos-web/internal/dominio"
	"track-concursos-web/internal/identificador"
)

type Servico struct{ repositorio repositorio }
type Concurso struct {
	ID            string   `json:"id"`
	Nome          string   `json:"nome"`
	Banca         string   `json:"banca"`
	Cargo         *string  `json:"cargo,omitempty"`
	Logotipo      *string  `json:"logotipo,omitempty"`
	Salario       *float64 `json:"salario,omitempty"`
	DataProva     *string  `json:"dataProva,omitempty"`
	PreEdital     bool     `json:"preEdital"`
	PrazosRevisao string   `json:"prazosRevisao"`
	Grupo         string   `json:"grupo"`
	Ordem         int      `json:"ordem"`
	Resultado     *string  `json:"resultado,omitempty"`
	Classificacao *int     `json:"classificacao,omitempty"`
	NotaFinal     *float64 `json:"notaFinal,omitempty"`
	Nomeado       bool     `json:"nomeado"`
	DataNomeacao  *string  `json:"dataNomeacao,omitempty"`
}
type Atribuicao struct {
	ConcursoID          string   `json:"concursoId"`
	Nome                string   `json:"nome"`
	Banca               string   `json:"banca"`
	Cargo               *string  `json:"cargo"`
	Logotipo            *string  `json:"logotipo"`
	Salario             *float64 `json:"salario"`
	DataProva           *string  `json:"dataProva"`
	PreEdital           bool     `json:"preEdital"`
	PrazosRevisao       string   `json:"prazosRevisao"`
	Grupo               string   `json:"grupo"`
	Ordem               int      `json:"ordem"`
	Resultado           *string  `json:"resultado"`
	Classificacao       *int     `json:"classificacao"`
	NotaFinal           *float64 `json:"notaFinal"`
	Nomeado             bool     `json:"nomeado"`
	DataNomeacao        *string  `json:"dataNomeacao"`
	LimparSalario       bool     `json:"limparSalario"`
	LimparClassificacao bool     `json:"limparClassificacao"`
	LimparNotaFinal     bool     `json:"limparNotaFinal"`
}
type Ordem struct {
	ID    string `json:"id"`
	Ordem int    `json:"ordem"`
}
type Alteracao struct {
	Nome                *string  `json:"nome"`
	Banca               *string  `json:"banca"`
	Cargo               *string  `json:"cargo"`
	Logotipo            *string  `json:"logotipo"`
	Salario             *float64 `json:"salario"`
	DataProva           *string  `json:"dataProva"`
	PreEdital           *bool    `json:"preEdital"`
	PrazosRevisao       *string  `json:"prazosRevisao"`
	Grupo               *string  `json:"grupo"`
	Ordem               *int     `json:"ordem"`
	Resultado           *string  `json:"resultado"`
	Classificacao       *int     `json:"classificacao"`
	NotaFinal           *float64 `json:"notaFinal"`
	Nomeado             *bool    `json:"nomeado"`
	DataNomeacao        *string  `json:"dataNomeacao"`
	LimparSalario       bool     `json:"limparSalario"`
	LimparClassificacao bool     `json:"limparClassificacao"`
	LimparNotaFinal     bool     `json:"limparNotaFinal"`
}

func Novo(banco *sql.DB) *Servico { return &Servico{repositorio: novoRepositorioMySQL(banco)} }
func (s *Servico) Catalogar(ctx context.Context, mentor string) ([]Concurso, error) {
	return s.repositorio.Catalogar(ctx, mentor)
}
func (s *Servico) CriarCatalogo(ctx context.Context, mentor string, e Atribuicao) (string, error) {
	if err := validarNovo(e); err != nil {
		return "", err
	}
	if strings.TrimSpace(e.PrazosRevisao) == "" {
		e.PrazosRevisao = "1,7,30"
	}
	normalizar(&e)
	id := identificador.UUID()
	return id, s.repositorio.CriarCatalogo(ctx, mentor, id, e)
}
func (s *Servico) Listar(ctx context.Context, aluno string) ([]Concurso, error) {
	return s.repositorio.Listar(ctx, aluno)
}
func (s *Servico) Atribuir(ctx context.Context, aluno, executor string, e Atribuicao) (string, error) {
	novo := e.ConcursoID == ""
	if novo {
		if err := validarNovo(e); err != nil {
			return "", err
		}
		if strings.TrimSpace(e.PrazosRevisao) == "" {
			e.PrazosRevisao = "1,7,30"
		}
	}
	if e.Grupo != "" && e.Grupo != "foco" && e.Grupo != "mira" && e.Grupo != "realizado" {
		return "", errors.New("grupo inválido")
	}
	normalizar(&e)
	id := e.ConcursoID
	if novo {
		id = identificador.UUID()
	}
	if err := s.repositorio.Atribuir(ctx, aluno, executor, id, novo, e); err != nil {
		return "", err
	}
	return id, nil
}
func (s *Servico) Reordenar(ctx context.Context, aluno string, lista []Ordem) error {
	if len(lista) == 0 {
		return errors.New("ordem vazia")
	}
	return s.repositorio.Reordenar(ctx, aluno, lista)
}
func (s *Servico) Alterar(ctx context.Context, aluno, id string, e Alteracao) error {
	if e.Grupo != nil && *e.Grupo != "foco" && *e.Grupo != "mira" && *e.Grupo != "realizado" {
		return errors.New("grupo inválido")
	}
	return s.repositorio.Alterar(ctx, aluno, id, e)
}
func (s *Servico) Desativar(ctx context.Context, aluno, id string) error {
	ok, err := s.repositorio.Desativar(ctx, aluno, id)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return nil
}
func validarNovo(e Atribuicao) error {
	if strings.TrimSpace(e.Nome) == "" || strings.TrimSpace(e.Banca) == "" {
		return errors.New("nome e banca obrigatórios")
	}
	return nil
}
func normalizar(e *Atribuicao) {
	if e.Grupo == "" {
		e.Grupo = "foco"
	}
}
func textoNulo(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}
func decimalNulo(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	return &v.Float64
}
func dataNula(v sql.NullTime) *string {
	if !v.Valid {
		return nil
	}
	x := v.Time.Format("2006-01-02")
	return &x
}
func valor(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
