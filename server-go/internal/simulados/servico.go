package simulados

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
	"chega-junto-concurseiro-web/internal/dominio"
	"chega-junto-concurseiro-web/internal/identificador"
)

type Servico struct{ repositorio repositorio }
type ResultadoMateria struct {
	ID         string   `json:"id,omitempty"`
	MateriaID  *string  `json:"materiaId,omitempty"`
	Nome       string   `json:"nome"`
	Questoes   *int     `json:"questoes,omitempty"`
	Acertos    *int     `json:"acertos,omitempty"`
	Erros      *int     `json:"erros,omitempty"`
	Brancos    *int     `json:"brancos,omitempty"`
	Percentual *float64 `json:"percentual,omitempty"`
	Pontos     *float64 `json:"pontos,omitempty"`
}
type Simulado struct {
	ID                 string             `json:"id"`
	AlunoID            string             `json:"alunoId"`
	ConcursoID         *string            `json:"concursoId,omitempty"`
	Nome               string             `json:"nome"`
	Tipo               string             `json:"tipo"`
	RealizadoEm        *string            `json:"realizadoEm,omitempty"`
	Link               *string            `json:"link,omitempty"`
	Observacoes        *string            `json:"observacoes,omitempty"`
	Percentual         *float64           `json:"percentual,omitempty"`
	TempoMinutos       *int               `json:"tempoMinutos,omitempty"`
	QuestoesFeitas     *int               `json:"questoesFeitas,omitempty"`
	Resultado          json.RawMessage    `json:"resultado,omitempty"`
	ConfiguracaoUsada  json.RawMessage    `json:"configuracaoUsada,omitempty"`
	ResultadosMaterias []ResultadoMateria `json:"resultadosMaterias"`
}
type Entrada struct {
	ConcursoID         *string            `json:"concursoId"`
	Nome               string             `json:"nome"`
	Tipo               string             `json:"tipo"`
	RealizadoEm        *string            `json:"realizadoEm"`
	Link               *string            `json:"link"`
	Observacoes        *string            `json:"observacoes"`
	Percentual         *float64           `json:"percentual"`
	TempoMinutos       *int               `json:"tempoMinutos"`
	QuestoesFeitas     *int               `json:"questoesFeitas"`
	Resultado          json.RawMessage    `json:"resultado"`
	ConfiguracaoUsada  json.RawMessage    `json:"configuracaoUsada"`
	ResultadosMaterias []ResultadoMateria `json:"resultadosMaterias"`
}

func Novo(banco *sql.DB) *Servico { return &Servico{repositorio: novoRepositorioMySQL(banco)} }
func (s *Servico) Listar(ctx context.Context, alunoID, concursoID string) ([]Simulado, error) {
	return s.repositorio.Listar(ctx, alunoID, concursoID)
}
func (s *Servico) Criar(ctx context.Context, alunoID string, e Entrada) (string, error) {
	if err := validar(e); err != nil {
		return "", err
	}
	if err := validarResultados(e.ResultadosMaterias); err != nil {
		return "", err
	}
	if err := s.validarConcurso(ctx, alunoID, e.ConcursoID); err != nil {
		return "", err
	}
	id := identificador.UUID()
	if err := s.repositorio.Criar(ctx, id, alunoID, e); err != nil {
		return "", err
	}
	return id, nil
}
func (s *Servico) Alterar(ctx context.Context, alunoID, id string, e Entrada) error {
	if err := s.existe(ctx, alunoID, id); err != nil {
		return err
	}
	if err := validar(e); err != nil {
		return err
	}
	if err := validarResultados(e.ResultadosMaterias); err != nil {
		return err
	}
	if err := s.validarConcurso(ctx, alunoID, e.ConcursoID); err != nil {
		return err
	}
	return s.repositorio.Alterar(ctx, id, e)
}
func (s *Servico) Desativar(ctx context.Context, alunoID, id string) error {
	if err := s.existe(ctx, alunoID, id); err != nil {
		return err
	}
	return s.repositorio.Desativar(ctx, id)
}
func (s *Servico) ObterConfiguracao(ctx context.Context, alunoID, concursoID string) (json.RawMessage, error) {
	bruto, existe, err := s.repositorio.ObterConfiguracao(ctx, alunoID, concursoID)
	if err != nil {
		return nil, err
	}
	if !existe {
		return json.RawMessage(`{}`), nil
	}
	return bruto, nil
}
func (s *Servico) SalvarConfiguracao(ctx context.Context, alunoID, concursoID, executor string, config json.RawMessage) error {
	if len(config) == 0 || !json.Valid(config) {
		return errors.New("configuracao invalida")
	}
	ok, err := s.repositorio.ConcursoAtribuido(ctx, alunoID, &concursoID)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoAutorizado
	}
	return s.repositorio.SalvarConfiguracao(ctx, identificador.UUID(), alunoID, concursoID, executor, config)
}
func (s *Servico) existe(ctx context.Context, alunoID, id string) error {
	ok, err := s.repositorio.Existe(ctx, alunoID, id)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return nil
}
func (s *Servico) validarConcurso(ctx context.Context, alunoID string, concursoID *string) error {
	ok, err := s.repositorio.ConcursoAtribuido(ctx, alunoID, concursoID)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoAutorizado
	}
	return nil
}
func validarResultados(lista []ResultadoMateria) error {
	for _, x := range lista {
		if x.Nome == "" {
			return errors.New("materia sem nome")
		}
		if x.Percentual != nil && (*x.Percentual < 0 || *x.Percentual > 100) {
			return errors.New("percentual invalido")
		}
	}
	return nil
}

type scanner interface{ Scan(...any) error }

func lerSimulado(l scanner) (Simulado, error) {
	var x Simulado
	var concurso, link, obs sql.NullString
	var data sql.NullTime
	var pct sql.NullFloat64
	var tempo, questoes sql.NullInt64
	var resultado, config []byte
	err := l.Scan(&x.ID, &x.AlunoID, &concurso, &x.Nome, &x.Tipo, &data, &link, &obs, &pct, &tempo, &questoes, &resultado, &config)
	x.ConcursoID, x.Link, x.Observacoes = texto(concurso), texto(link), texto(obs)
	if data.Valid {
		v := data.Time.Format("2006-01-02")
		x.RealizadoEm = &v
	}
	x.Percentual, x.TempoMinutos, x.QuestoesFeitas = decimal(pct), inteiro(tempo), inteiro(questoes)
	if len(resultado) > 0 {
		x.Resultado = resultado
	}
	if len(config) > 0 {
		x.ConfiguracaoUsada = config
	}
	return x, err
}
func validar(e Entrada) error {
	if e.Nome == "" {
		return errors.New("nome obrigatorio")
	}
	if e.Tipo != "realizado" && e.Tipo != "pendente" {
		return errors.New("tipo invalido")
	}
	if e.Tipo == "pendente" && (e.Link == nil || *e.Link == "") {
		return errors.New("link obrigatorio")
	}
	if e.RealizadoEm != nil && *e.RealizadoEm != "" {
		if _, err := time.Parse("2006-01-02", *e.RealizadoEm); err != nil {
			return err
		}
	}
	if e.Percentual != nil && (*e.Percentual < 0 || *e.Percentual > 100) {
		return errors.New("percentual invalido")
	}
	if e.TempoMinutos != nil && *e.TempoMinutos < 0 {
		return errors.New("tempo invalido")
	}
	if e.QuestoesFeitas != nil && *e.QuestoesFeitas < 0 {
		return errors.New("questoes invalidas")
	}
	return nil
}
func valor(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
func jsonNulo(v json.RawMessage) any {
	if len(v) == 0 || string(v) == "null" {
		return nil
	}
	return []byte(v)
}
func texto(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	x := v.String
	return &x
}
func inteiro(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	x := int(v.Int64)
	return &x
}
func decimal(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	x := v.Float64
	return &x
}
