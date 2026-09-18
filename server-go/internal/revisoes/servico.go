package revisoes

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"chega-junto-concurseiro-web/internal/dominio"
	"chega-junto-concurseiro-web/internal/identificador"
)

type Servico struct{ repositorio repositorio }
type Revisao struct {
	ID                 string   `json:"id"`
	AlunoID            string   `json:"alunoId"`
	ConcursoID         *string  `json:"concursoId,omitempty"`
	MateriaID          *string  `json:"materiaId,omitempty"`
	TopicoID           *string  `json:"topicoId,omitempty"`
	SubtopicoID        *string  `json:"subtopicoId,omitempty"`
	CicloAtual         int      `json:"cicloAtual"`
	ProximaData        *string  `json:"proximaData,omitempty"`
	PercentualAnterior *float64 `json:"percentualAnterior,omitempty"`
	Observacoes        *string  `json:"observacoes,omitempty"`
	Concluida          bool     `json:"concluida"`
	ConcluidaEm        *string  `json:"concluidaEm,omitempty"`
	Materia            string   `json:"materia,omitempty"`
	Topico             string   `json:"topico,omitempty"`
	Subtopico          string   `json:"subtopico,omitempty"`
}
type Entrada struct {
	ConcursoID         *string  `json:"concursoId"`
	MateriaID          *string  `json:"materiaId"`
	TopicoID           *string  `json:"topicoId"`
	SubtopicoID        *string  `json:"subtopicoId"`
	CicloAtual         int      `json:"cicloAtual"`
	ProximaData        *string  `json:"proximaData"`
	PercentualAnterior *float64 `json:"percentualAnterior"`
	Observacoes        *string  `json:"observacoes"`
	Concluida          bool     `json:"concluida"`
}
type Alteracao struct {
	ConcursoID         *string  `json:"concursoId"`
	MateriaID          *string  `json:"materiaId"`
	TopicoID           *string  `json:"topicoId"`
	SubtopicoID        *string  `json:"subtopicoId"`
	CicloAtual         *int     `json:"cicloAtual"`
	ProximaData        *string  `json:"proximaData"`
	PercentualAnterior *float64 `json:"percentualAnterior"`
	Observacoes        *string  `json:"observacoes"`
	Concluida          *bool    `json:"concluida"`
}

func Novo(b *sql.DB) *Servico                   { return NovoComRepositorio(novoRepositorioMySQL(b)) }
func NovoComRepositorio(r repositorio) *Servico { return &Servico{repositorio: r} }
func (s *Servico) Listar(ctx context.Context, a, c, e string) ([]Revisao, error) {
	return s.repositorio.Listar(ctx, a, c, e)
}
func (s *Servico) Criar(ctx context.Context, aluno string, e Entrada) (string, error) {
	if err := validar(e.CicloAtual, e.ProximaData, e.PercentualAnterior); err != nil {
		return "", err
	}
	if err := s.validarConcurso(ctx, aluno, e.ConcursoID); err != nil {
		return "", err
	}
	if err := s.validarConteudo(ctx, aluno, e.ConcursoID, e.MateriaID, e.TopicoID, e.SubtopicoID); err != nil {
		return "", err
	}
	if id, ok, err := s.repositorio.Pendente(ctx, aluno, e); err != nil {
		return "", err
	} else if ok {
		return id, s.repositorio.AtualizarPendente(ctx, id, e)
	}
	id := identificador.UUID()
	return id, s.repositorio.Criar(ctx, id, aluno, e)
}
func (s *Servico) Alterar(ctx context.Context, aluno, id string, e Alteracao) error {
	atual, ok, err := s.repositorio.Obter(ctx, aluno, id)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	if e.CicloAtual != nil && *e.CicloAtual < 0 {
		return errors.New("ciclo inválido")
	}
	if e.ProximaData != nil && *e.ProximaData != "" {
		if _, err := time.Parse("2006-01-02", *e.ProximaData); err != nil {
			return err
		}
	}
	if e.PercentualAnterior != nil && (*e.PercentualAnterior < 0 || *e.PercentualAnterior > 100) {
		return errors.New("percentual inválido")
	}
	if e.ConcursoID != nil {
		if err := s.validarConcurso(ctx, aluno, e.ConcursoID); err != nil {
			return err
		}
	}
	if e.ConcursoID != nil {
		atual.ConcursoID = e.ConcursoID
	}
	if e.MateriaID != nil {
		atual.MateriaID = e.MateriaID
	}
	if e.TopicoID != nil {
		atual.TopicoID = e.TopicoID
	}
	if e.SubtopicoID != nil {
		atual.SubtopicoID = e.SubtopicoID
	}
	if err := s.validarConteudo(ctx, aluno, atual.ConcursoID, atual.MateriaID, atual.TopicoID, atual.SubtopicoID); err != nil {
		return err
	}
	return s.repositorio.Alterar(ctx, id, e)
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
func (s *Servico) validarConcurso(ctx context.Context, a string, id *string) error {
	ok, err := s.repositorio.ConcursoAtribuido(ctx, a, id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("concurso não atribuído")
	}
	return nil
}
func (s *Servico) validarConteudo(ctx context.Context, a string, c, m, t, sub *string) error {
	if c == nil || *c == "" || m == nil || *m == "" {
		return errors.New("concurso e materia obrigatorios")
	}
	if sub != nil && *sub != "" && (t == nil || *t == "") {
		return errors.New("topico obrigatorio")
	}
	ok, err := s.repositorio.ConteudoValido(ctx, a, c, m, t, sub)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("conteudo fora do edital atribuido")
	}
	return nil
}
func validar(c int, d *string, p *float64) error {
	if c < 0 {
		return errors.New("ciclo inválido")
	}
	if d != nil && *d != "" {
		if _, err := time.Parse("2006-01-02", *d); err != nil {
			return err
		}
	}
	if p != nil && (*p < 0 || *p > 100) {
		return errors.New("percentual inválido")
	}
	return nil
}
func valor(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
func texto(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	x := v.String
	return &x
}
