package estudos

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
type Sessao struct {
	ID          string          `json:"id"`
	AlunoID     string          `json:"alunoId"`
	ConcursoID  *string         `json:"concursoId,omitempty"`
	MateriaID   *string         `json:"materiaId,omitempty"`
	TopicoID    *string         `json:"topicoId,omitempty"`
	SubtopicoID *string         `json:"subtopicoId,omitempty"`
	Segundos    int             `json:"segundos"`
	Modo        *string         `json:"modo,omitempty"`
	Observacoes *string         `json:"observacoes,omitempty"`
	Metricas    json.RawMessage `json:"metricas,omitempty"`
	Origem      *string         `json:"origem,omitempty"`
	EstudadoEm  string          `json:"estudadoEm"`
}
type EntradaSessao struct {
	ConcursoID  *string         `json:"concursoId"`
	MateriaID   *string         `json:"materiaId"`
	TopicoID    *string         `json:"topicoId"`
	SubtopicoID *string         `json:"subtopicoId"`
	Segundos    int             `json:"segundos"`
	Modo        *string         `json:"modo"`
	Observacoes *string         `json:"observacoes"`
	Metricas    json.RawMessage `json:"metricas"`
	Origem      *string         `json:"origem"`
	EstudadoEm  string          `json:"estudadoEm"`
}
type AlteracaoSessao struct {
	ConcursoID  *string          `json:"concursoId"`
	MateriaID   *string          `json:"materiaId"`
	TopicoID    *string          `json:"topicoId"`
	SubtopicoID *string          `json:"subtopicoId"`
	Segundos    *int             `json:"segundos"`
	Modo        *string          `json:"modo"`
	Observacoes *string          `json:"observacoes"`
	Metricas    *json.RawMessage `json:"metricas"`
	Origem      *string          `json:"origem"`
	EstudadoEm  *string          `json:"estudadoEm"`
}
type RegistroQuestoes struct {
	ID           string  `json:"id"`
	AlunoID      string  `json:"alunoId"`
	ConcursoID   *string `json:"concursoId,omitempty"`
	MateriaID    *string `json:"materiaId,omitempty"`
	TopicoID     *string `json:"topicoId,omitempty"`
	SubtopicoID  *string `json:"subtopicoId,omitempty"`
	Resolvidas   int     `json:"resolvidas"`
	Acertos      int     `json:"acertos"`
	Erros        int     `json:"erros"`
	Origem       *string `json:"origem,omitempty"`
	RegistradoEm string  `json:"registradoEm"`
}
type EntradaQuestoes struct {
	ConcursoID   *string `json:"concursoId"`
	MateriaID    *string `json:"materiaId"`
	TopicoID     *string `json:"topicoId"`
	SubtopicoID  *string `json:"subtopicoId"`
	Resolvidas   int     `json:"resolvidas"`
	Acertos      int     `json:"acertos"`
	Erros        int     `json:"erros"`
	Origem       *string `json:"origem"`
	RegistradoEm string  `json:"registradoEm"`
}
type AlteracaoQuestoes struct {
	ConcursoID   *string `json:"concursoId"`
	MateriaID    *string `json:"materiaId"`
	TopicoID     *string `json:"topicoId"`
	SubtopicoID  *string `json:"subtopicoId"`
	Resolvidas   *int    `json:"resolvidas"`
	Acertos      *int    `json:"acertos"`
	Erros        *int    `json:"erros"`
	Origem       *string `json:"origem"`
	RegistradoEm *string `json:"registradoEm"`
}

func Novo(banco *sql.DB) *Servico               { return NovoComRepositorio(novoRepositorioMySQL(banco)) }
func NovoComRepositorio(r repositorio) *Servico { return &Servico{repositorio: r} }
func (s *Servico) ListarSessoes(ctx context.Context, a, c, i, f string) ([]Sessao, error) {
	return s.repositorio.ListarSessoes(ctx, a, c, i, f)
}
func (s *Servico) ListarQuestoes(ctx context.Context, a, c, i, f string) ([]RegistroQuestoes, error) {
	return s.repositorio.ListarQuestoes(ctx, a, c, i, f)
}
func (s *Servico) CriarSessao(ctx context.Context, aluno string, e EntradaSessao) (string, error) {
	if e.Segundos < 0 || (e.Segundos == 0 && len(e.Metricas) == 0) {
		return "", errors.New("sessão vazia")
	}
	d, err := dataHora(e.EstudadoEm)
	if err != nil {
		return "", err
	}
	if d.After(time.Now().UTC().Add(5 * time.Minute)) {
		return "", errors.New("a data não pode estar no futuro")
	}
	if err = s.validarConcurso(ctx, aluno, e.ConcursoID); err != nil {
		return "", err
	}
	id := identificador.UUID()
	return id, s.repositorio.CriarSessao(ctx, id, aluno, e, d)
}
func (s *Servico) AlterarSessao(ctx context.Context, aluno, id string, e AlteracaoSessao) error {
	if err := s.existe(ctx, "sessoes_estudo", aluno, id); err != nil {
		return err
	}
	if e.Segundos != nil && *e.Segundos < 0 {
		return errors.New("tempo inválido")
	}
	if e.ConcursoID != nil {
		if err := s.validarConcurso(ctx, aluno, e.ConcursoID); err != nil {
			return err
		}
	}
	var d any
	if e.EstudadoEm != nil {
		x, err := dataHora(*e.EstudadoEm)
		if err != nil {
			return err
		}
		if x.After(time.Now().UTC().Add(5 * time.Minute)) {
			return errors.New("a data não pode estar no futuro")
		}
		d = x
	}
	return s.repositorio.AtualizarSessao(ctx, id, e, d)
}
func (s *Servico) CriarQuestoes(ctx context.Context, aluno string, e EntradaQuestoes) (string, error) {
	if err := validarQuestoes(e.Resolvidas, e.Acertos, e.Erros); err != nil {
		return "", err
	}
	if err := s.validarConcurso(ctx, aluno, e.ConcursoID); err != nil {
		return "", err
	}
	d, err := dataHora(e.RegistradoEm)
	if err != nil {
		return "", err
	}
	if d.After(time.Now().UTC().Add(5 * time.Minute)) {
		return "", errors.New("a data não pode estar no futuro")
	}
	id := identificador.UUID()
	return id, s.repositorio.CriarQuestoes(ctx, id, aluno, e, d)
}
func (s *Servico) AlterarQuestoes(ctx context.Context, aluno, id string, e AlteracaoQuestoes) error {
	if err := s.existe(ctx, "registros_questoes", aluno, id); err != nil {
		return err
	}
	r, a, er, err := s.repositorio.ObterContadores(ctx, aluno, id)
	if err != nil {
		return err
	}
	if e.Resolvidas != nil {
		r = *e.Resolvidas
	}
	if e.Acertos != nil {
		a = *e.Acertos
	}
	if e.Erros != nil {
		er = *e.Erros
	}
	if err = validarQuestoes(r, a, er); err != nil {
		return err
	}
	if e.ConcursoID != nil {
		if err = s.validarConcurso(ctx, aluno, e.ConcursoID); err != nil {
			return err
		}
	}
	var d any
	if e.RegistradoEm != nil {
		x, err := dataHora(*e.RegistradoEm)
		if err != nil {
			return err
		}
		if x.After(time.Now().UTC().Add(5 * time.Minute)) {
			return errors.New("a data não pode estar no futuro")
		}
		d = x
	}
	return s.repositorio.AtualizarQuestoes(ctx, id, e, r, a, er, d)
}
func (s *Servico) Desativar(ctx context.Context, aluno, tabela, id string) error {
	if tabela != "sessoes_estudo" && tabela != "registros_questoes" {
		return errors.New("recurso inválido")
	}
	if err := s.existe(ctx, tabela, aluno, id); err != nil {
		return err
	}
	return s.repositorio.Desativar(ctx, tabela, id)
}
func (s *Servico) existe(ctx context.Context, tabela, aluno, id string) error {
	ok, err := s.repositorio.Existe(ctx, tabela, aluno, id)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return nil
}
func (s *Servico) validarConcurso(ctx context.Context, aluno string, id *string) error {
	ok, err := s.repositorio.ConcursoAtribuido(ctx, aluno, id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("concurso não atribuído")
	}
	return nil
}
func validarQuestoes(r, a, e int) error {
	if r < 0 || a < 0 || e < 0 || a+e > r {
		return errors.New("quantidades de questões inválidas")
	}
	return nil
}
func dataHora(v string) (time.Time, error) {
	if v == "" {
		return time.Now().UTC(), nil
	}
	if d, err := time.Parse(time.RFC3339, v); err == nil {
		return d.UTC(), nil
	}
	return time.Parse("2006-01-02", v)
}
func filtros(q string, args []any, concurso, inicio, fim, coluna string) (string, []any) {
	if concurso != "" {
		q += ` AND concurso_id=?`
		args = append(args, concurso)
	}
	if inicio != "" {
		q += ` AND ` + coluna + `>=?`
		args = append(args, inicio)
	}
	if fim != "" {
		q += ` AND ` + coluna + `<?`
		if d, err := time.Parse("2006-01-02", fim); err == nil {
			args = append(args, d.AddDate(0, 0, 1))
		} else {
			args = append(args, fim)
		}
	}
	return q, args
}
func texto(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	x := v.String
	return &x
}
func jsonNulo(v json.RawMessage) any {
	if len(v) == 0 || string(v) == "null" {
		return nil
	}
	return []byte(v)
}
