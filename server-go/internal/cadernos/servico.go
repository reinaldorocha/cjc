package cadernos

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"chega-junto-concurseiro-web/internal/dominio"
	"chega-junto-concurseiro-web/internal/identificador"
)

type Caderno struct {
	ID           string  `json:"id"`
	AlunoID      string  `json:"alunoId"`
	Titulo       string  `json:"titulo"`
	Pasta        string  `json:"pasta"`
	Conteudo     string  `json:"conteudo"`
	EditalID     *string `json:"editalId"`
	MateriaID    *string `json:"materiaId"`
	TopicoID     *string `json:"topicoId"`
	Cor          string  `json:"cor"`
	Ativo        bool    `json:"ativo"`
	CriadoEm     string  `json:"criadoEm"`
	AtualizadoEm string  `json:"atualizadoEm"`
}
type EntradaCaderno struct {
	Titulo    string  `json:"titulo"`
	Pasta     *string `json:"pasta"`
	Conteudo  *string `json:"conteudo"`
	EditalID  *string `json:"editalId"`
	MateriaID *string `json:"materiaId"`
	TopicoID  *string `json:"topicoId"`
	Cor       *string `json:"cor"`
}
type Servico struct{ repositorio repositorio }

func Novo(banco *sql.DB) *Servico               { return NovoComRepositorio(novoRepositorioMySQL(banco)) }
func NovoComRepositorio(r repositorio) *Servico { return &Servico{repositorio: r} }
func textoNulo(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	x := v.String
	return &x
}
func (s *Servico) Criar(ctx context.Context, alunoID string, e EntradaCaderno) (string, error) {
	n, err := normalizar(e)
	if err != nil {
		return "", err
	}
	id := identificador.UUID()
	return id, s.repositorio.Criar(ctx, id, alunoID, n)
}
func (s *Servico) Listar(ctx context.Context, alunoID, editalID string) ([]Caderno, error) {
	return s.repositorio.Listar(ctx, alunoID, editalID)
}
func (s *Servico) Obter(ctx context.Context, alunoID, id string) (*Caderno, error) {
	caderno, err := s.repositorio.Obter(ctx, alunoID, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dominio.ErrNaoEncontrado
	}
	return caderno, err
}
func (s *Servico) Alterar(ctx context.Context, alunoID, id string, e EntradaCaderno) error {
	n, err := normalizar(e)
	if err != nil {
		return err
	}
	return s.repositorio.Alterar(ctx, alunoID, id, n)
}
func (s *Servico) Desativar(ctx context.Context, alunoID, id string) error {
	return s.repositorio.Desativar(ctx, alunoID, id)
}
func normalizar(e EntradaCaderno) (EntradaCadernoNormalizada, error) {
	if strings.TrimSpace(e.Titulo) == "" {
		return EntradaCadernoNormalizada{}, errors.New("titulo obrigatorio")
	}
	n := EntradaCadernoNormalizada{Titulo: strings.TrimSpace(e.Titulo), Pasta: "Geral", EditalID: e.EditalID, MateriaID: e.MateriaID, TopicoID: e.TopicoID, Cor: "#4f8ef7"}
	if e.Pasta != nil && strings.TrimSpace(*e.Pasta) != "" {
		n.Pasta = strings.TrimSpace(*e.Pasta)
	}
	if e.Conteudo != nil {
		n.Conteudo = *e.Conteudo
	}
	if e.Cor != nil && strings.TrimSpace(*e.Cor) != "" {
		n.Cor = strings.TrimSpace(*e.Cor)
	}
	return n, nil
}
