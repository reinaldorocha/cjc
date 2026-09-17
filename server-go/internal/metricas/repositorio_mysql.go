package metricas

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type repositorioMySQL struct {
	db *sql.DB
}

func novoRepositorioMySQL(db *sql.DB, fuso *time.Location) *repositorioMySQL {
	return &repositorioMySQL{db: db}
}

type Resumo struct {
	SegundosEstudados       int64     `json:"segundosEstudados"`
	SegundosEstudadosTotal  int64     `json:"segundosEstudadosTotal"`
	QuestoesResolvidas      int64     `json:"questoesResolvidas"`
	QuestoesResolvidasTotal int64     `json:"questoesResolvidasTotal"`
	Acertos                 int64     `json:"acertos"`
	Erros                   int64     `json:"erros"`
	PercentualAcertos       *float64  `json:"percentualAcertos,omitempty"`
	DiasAtivos              int       `json:"diasAtivos"`
	SequenciaAtual          int       `json:"sequenciaAtual"`
	MediaSimulados          *float64  `json:"mediaSimulados,omitempty"`
	MelhorSimulado          *float64  `json:"melhorSimulado,omitempty"`
	UltimoSimulado          *float64  `json:"ultimoSimulado,omitempty"`
	TendenciaSimulados      *float64  `json:"tendenciaSimulados,omitempty"`
	SimuladosRealizados     int       `json:"simuladosRealizados"`
	SimuladosTotal          int       `json:"simuladosTotal"`
	ItensEditalConcluidos   int       `json:"itensEditalConcluidos"`
	ItensEditalTotal        int       `json:"itensEditalTotal"`
	PercentualEdital        float64   `json:"percentualEdital"`
	MateriasIniciadas       int       `json:"materiasIniciadas"`
	MateriasTotal           int       `json:"materiasTotal"`
	SegundosHoje            int64     `json:"segundosHoje"`
	QuestoesHoje            int64     `json:"questoesHoje"`
	UltimosSimulados        []float64 `json:"-"`
}

type Dia struct {
	Data      string `json:"data"`
	Segundos  int64  `json:"segundos"`
	Questoes  int64  `json:"questoes"`
	Acertos   int64  `json:"acertos"`
	Erros     int64  `json:"erros"`
	DiasAtivo bool   `json:"diaAtivo"`
}

type Materia struct {
	ID                string   `json:"id"`
	Nome              string   `json:"nome"`
	Ordem             int      `json:"ordem"`
	Segundos          int64    `json:"segundos"`
	Questoes          int64    `json:"questoes"`
	Acertos           int64    `json:"acertos"`
	Erros             int64    `json:"erros"`
	PercentualAcertos *float64 `json:"percentualAcertos,omitempty"`
	ItensConcluidos   int      `json:"itensConcluidos"`
	ItensTotal        int      `json:"itensTotal"`
	PercentualEdital  float64  `json:"percentualEdital"`
	Simulados         int      `json:"simulados"`
	MediaSimulados    *float64 `json:"mediaSimulados,omitempty"`
	MelhorSimulado    *float64 `json:"melhorSimulado,omitempty"`
	UltimoSimulado    *float64 `json:"ultimoSimulado,omitempty"`
}

func (s *repositorioMySQL) ConcursoAtribuido(ctx context.Context, alunoID, concursoID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, alunoID, concursoID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) Resumo(ctx context.Context, alunoID, concursoID string, periodo Periodo) (Resumo, error) {
	var resultado Resumo
	de, ate := periodo.Inicio, periodo.Fim
	var err error
	baseSessao := ` FROM sessoes_estudo WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`
	argsPeriodo := []any{alunoID, concursoID}
	condicao, valoresss := periodoSQL("estudado_em", de, ate)
	if err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(segundos),0)`+baseSessao+condicao, append(argsPeriodo, valoresss...)...).Scan(&resultado.SegundosEstudados); err != nil {
		return resultado, err
	}
	if err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(segundos),0)`+baseSessao, alunoID, concursoID).Scan(&resultado.SegundosEstudadosTotal); err != nil {
		return resultado, err
	}
	baseQuestoes := ` FROM registros_questoes WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`
	condicao, valoresss = periodoSQL("registrado_em", de, ate)
	if err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(resolvidas),0),COALESCE(SUM(acertos),0),COALESCE(SUM(erros),0)`+baseQuestoes+condicao, append(argsPeriodo, valoresss...)...).Scan(&resultado.QuestoesResolvidas, &resultado.Acertos, &resultado.Erros); err != nil {
		return resultado, err
	}
	if err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(resolvidas),0)`+baseQuestoes, alunoID, concursoID).Scan(&resultado.QuestoesResolvidasTotal); err != nil {
		return resultado, err
	}
	condSessao, valSessao := periodoSQL("estudado_em", de, ate)
	condQuestoes, valQuestoes := periodoSQL("registrado_em", de, ate)
	argsDias := []any{periodo.Offset, alunoID, concursoID}
	argsDias = append(argsDias, valSessao...)
	argsDias = append(argsDias, periodo.Offset, alunoID, concursoID)
	argsDias = append(argsDias, valQuestoes...)
	consultaDias := `SELECT COUNT(*) FROM (SELECT DATE(CONVERT_TZ(estudado_em,'+00:00',?)) data FROM sessoes_estudo WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE` + condSessao + ` UNION SELECT DATE(CONVERT_TZ(registrado_em,'+00:00',?)) data FROM registros_questoes WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE` + condQuestoes + `) atividade`
	if err = s.db.QueryRowContext(ctx, consultaDias, argsDias...).Scan(&resultado.DiasAtivos); err != nil {
		return resultado, err
	}
	if err = s.carregarSimulados(ctx, alunoID, concursoID, &resultado); err != nil {
		return resultado, err
	}
	resultado.ItensEditalConcluidos, resultado.ItensEditalTotal, err = s.cobertura(ctx, alunoID, concursoID, "")
	if err != nil {
		return resultado, err
	}
	if err = s.db.QueryRowContext(ctx, `SELECT COUNT(DISTINCT em.id) FROM aluno_editais ae JOIN editais e ON e.id=ae.edital_id AND e.ativo=TRUE JOIN edital_materias em ON em.edital_id=e.id AND em.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND e.concurso_id=?`, alunoID, concursoID).Scan(&resultado.MateriasTotal); err != nil {
		return resultado, err
	}
	if err = s.db.QueryRowContext(ctx, `SELECT COUNT(DISTINCT materia_id) FROM sessoes_estudo WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE AND segundos>0 AND materia_id IS NOT NULL`, alunoID, concursoID).Scan(&resultado.MateriasIniciadas); err != nil {
		return resultado, err
	}
	if err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(segundos),0) FROM sessoes_estudo WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE AND estudado_em>=? AND estudado_em<?`, alunoID, concursoID, periodo.InicioHoje, periodo.FimHoje).Scan(&resultado.SegundosHoje); err != nil {
		return resultado, err
	}
	if err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(resolvidas),0) FROM registros_questoes WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE AND registrado_em>=? AND registrado_em<?`, alunoID, concursoID, periodo.InicioHoje, periodo.FimHoje).Scan(&resultado.QuestoesHoje); err != nil {
		return resultado, err
	}
	return resultado, nil
}

func (s *repositorioMySQL) LinhaDoTempo(ctx context.Context, alunoID, concursoID string, periodo Periodo) ([]Dia, error) {
	de, ate := *periodo.Inicio, *periodo.Fim
	var err error
	mapa := map[string]*Dia{}
	if err = s.agruparDias(ctx, `SELECT DATE(CONVERT_TZ(estudado_em,'+00:00',?)),COALESCE(SUM(segundos),0) FROM sessoes_estudo WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE AND estudado_em>=? AND estudado_em<? GROUP BY 1`, mapa, true, alunoID, concursoID, de, ate, periodo.Offset); err != nil {
		return nil, err
	}
	if err = s.agruparDias(ctx, `SELECT DATE(CONVERT_TZ(registrado_em,'+00:00',?)),COALESCE(SUM(resolvidas),0),COALESCE(SUM(acertos),0),COALESCE(SUM(erros),0) FROM registros_questoes WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE AND registrado_em>=? AND registrado_em<? GROUP BY 1`, mapa, false, alunoID, concursoID, de, ate, periodo.Offset); err != nil {
		return nil, err
	}
	lista := make([]Dia, 0, len(mapa))
	for _, x := range mapa {
		lista = append(lista, *x)
	}
	return lista, nil
}

func (s *repositorioMySQL) Materias(ctx context.Context, alunoID, concursoID string, periodo Periodo) ([]Materia, error) {
	de, ate := periodo.Inicio, periodo.Fim
	var err error
	linhas, err := s.db.QueryContext(ctx, `SELECT DISTINCT em.id,em.nome,em.ordem FROM aluno_editais ae JOIN editais e ON e.id=ae.edital_id AND e.ativo=TRUE JOIN edital_materias em ON em.edital_id=e.id AND em.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND e.concurso_id=? ORDER BY em.ordem,em.nome`, alunoID, concursoID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Materia{}
	for linhas.Next() {
		var x Materia
		if err = linhas.Scan(&x.ID, &x.Nome, &x.Ordem); err != nil {
			return nil, err
		}
		cond, valoresss := periodoSQL("estudado_em", de, ate)
		args := []any{alunoID, concursoID, x.ID}
		args = append(args, valoresss...)
		if err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(segundos),0) FROM sessoes_estudo WHERE aluno_id=? AND concurso_id=? AND materia_id=? AND ativo=TRUE`+cond, args...).Scan(&x.Segundos); err != nil {
			return nil, err
		}
		cond, valoresss = periodoSQL("registrado_em", de, ate)
		args = []any{alunoID, concursoID, x.ID}
		args = append(args, valoresss...)
		if err = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(resolvidas),0),COALESCE(SUM(acertos),0),COALESCE(SUM(erros),0) FROM registros_questoes WHERE aluno_id=? AND concurso_id=? AND materia_id=? AND ativo=TRUE`+cond, args...).Scan(&x.Questoes, &x.Acertos, &x.Erros); err != nil {
			return nil, err
		}
		x.ItensConcluidos, x.ItensTotal, err = s.cobertura(ctx, alunoID, concursoID, x.ID)
		if err != nil {
			return nil, err
		}
		if err = s.metricasSimuladoMateria(ctx, alunoID, concursoID, &x); err != nil {
			return nil, err
		}
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) carregarSimulados(ctx context.Context, alunoID, concursoID string, r *Resumo) error {
	var media, melhor sql.NullFloat64
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*),COALESCE(SUM(tipo='realizado' AND percentual IS NOT NULL),0),AVG(CASE WHEN tipo='realizado' THEN percentual END),MAX(CASE WHEN tipo='realizado' THEN percentual END) FROM simulados WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, alunoID, concursoID).Scan(&r.SimuladosTotal, &r.SimuladosRealizados, &media, &melhor); err != nil {
		return err
	}
	r.MediaSimulados = decimal(media)
	r.MelhorSimulado = decimal(melhor)
	linhas, err := s.db.QueryContext(ctx, `SELECT percentual FROM simulados WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE AND tipo='realizado' AND percentual IS NOT NULL ORDER BY COALESCE(realizado_em,DATE(criado_em)) DESC,criado_em DESC LIMIT 2`, alunoID, concursoID)
	if err != nil {
		return err
	}
	defer linhas.Close()
	valoresss := []float64{}
	for linhas.Next() {
		var v float64
		if err = linhas.Scan(&v); err != nil {
			return err
		}
		valoresss = append(valoresss, v)
	}
	r.UltimosSimulados = valoresss
	return linhas.Err()
}
func (s *repositorioMySQL) metricasSimuladoMateria(ctx context.Context, alunoID, concursoID string, m *Materia) error {
	var media, melhor sql.NullFloat64
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*),AVG(r.percentual),MAX(r.percentual) FROM resultados_materias_simulado r JOIN simulados si ON si.id=r.simulado_id WHERE si.aluno_id=? AND si.concurso_id=? AND si.ativo=TRUE AND si.tipo='realizado' AND r.ativo=TRUE AND r.materia_id=?`, alunoID, concursoID, m.ID).Scan(&m.Simulados, &media, &melhor); err != nil {
		return err
	}
	m.MediaSimulados = decimal(media)
	m.MelhorSimulado = decimal(melhor)
	var ultima sql.NullFloat64
	err := s.db.QueryRowContext(ctx, `SELECT r.percentual FROM resultados_materias_simulado r JOIN simulados si ON si.id=r.simulado_id WHERE si.aluno_id=? AND si.concurso_id=? AND si.ativo=TRUE AND si.tipo='realizado' AND r.ativo=TRUE AND r.materia_id=? AND r.percentual IS NOT NULL ORDER BY COALESCE(si.realizado_em,DATE(si.criado_em)) DESC,si.criado_em DESC LIMIT 1`, alunoID, concursoID, m.ID).Scan(&ultima)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	m.UltimoSimulado = decimal(ultima)
	return err
}
func (s *repositorioMySQL) cobertura(ctx context.Context, alunoID, concursoID, materiaID string) (int, int, error) {
	filtro := ""
	args := []any{alunoID, concursoID}
	if materiaID != "" {
		filtro = ` AND em.id=?`
		args = append(args, materiaID)
	}
	q := `SELECT COALESCE(SUM(CASE WHEN COALESCE(ap.estudado,FALSE) THEN 1 ELSE 0 END),0),COUNT(*) FROM (SELECT et.id,'topico' tipo FROM aluno_editais ae JOIN editais e ON e.id=ae.edital_id AND e.ativo=TRUE JOIN edital_materias em ON em.edital_id=e.id AND em.ativo=TRUE JOIN edital_topicos et ON et.materia_id=em.id AND et.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND e.concurso_id=?` + filtro + ` AND NOT EXISTS(SELECT 1 FROM edital_subtopicos es WHERE es.topico_id=et.id AND es.ativo=TRUE) UNION ALL SELECT es.id,'subtopico' tipo FROM aluno_editais ae JOIN editais e ON e.id=ae.edital_id AND e.ativo=TRUE JOIN edital_materias em ON em.edital_id=e.id AND em.ativo=TRUE JOIN edital_topicos et ON et.materia_id=em.id AND et.ativo=TRUE JOIN edital_subtopicos es ON es.topico_id=et.id AND es.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND e.concurso_id=?` + filtro + `) itens LEFT JOIN aluno_progresso_edital ap ON ap.aluno_id=? AND ap.tipo_item=itens.tipo AND ap.item_id=itens.id`
	todos := append([]any{}, args...)
	todos = append(todos, args...)
	todos = append(todos, alunoID)
	var concluidos, total int
	err := s.db.QueryRowContext(ctx, q, todos...).Scan(&concluidos, &total)
	return concluidos, total, err
}
func (s *repositorioMySQL) DatasEstudo(ctx context.Context, alunoID, concursoID, offset string) ([]string, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT DISTINCT DATE(CONVERT_TZ(estudado_em,'+00:00',?)) data FROM sessoes_estudo WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE ORDER BY data DESC`, offset, alunoID, concursoID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	datas := []string{}
	for linhas.Next() {
		var d time.Time
		if err = linhas.Scan(&d); err != nil {
			return nil, err
		}
		datas = append(datas, d.Format("2006-01-02"))
	}
	return datas, linhas.Err()
}
func (s *repositorioMySQL) agruparDias(ctx context.Context, q string, mapa map[string]*Dia, sessoes bool, aluno, concurso string, de, ate time.Time, offset string) error {
	linhas, err := s.db.QueryContext(ctx, q, offset, aluno, concurso, de.UTC(), ate.UTC())
	if err != nil {
		return err
	}
	defer linhas.Close()
	for linhas.Next() {
		var data time.Time
		if sessoes {
			var segundos int64
			if err = linhas.Scan(&data, &segundos); err != nil {
				return err
			}
			chave := data.Format("2006-01-02")
			x := mapa[chave]
			if x == nil {
				x = &Dia{Data: chave}
				mapa[chave] = x
			}
			x.Segundos = segundos
		} else {
			var q, a, e int64
			if err = linhas.Scan(&data, &q, &a, &e); err != nil {
				return err
			}
			chave := data.Format("2006-01-02")
			x := mapa[chave]
			if x == nil {
				x = &Dia{Data: chave}
				mapa[chave] = x
			}
			x.Questoes = q
			x.Acertos = a
			x.Erros = e
		}
	}
	return linhas.Err()
}
func periodoSQL(coluna string, de, ate *time.Time) (string, []any) {
	if de == nil || ate == nil {
		return "", nil
	}
	return fmt.Sprintf(" AND %s>=? AND %s<?", coluna, coluna), []any{de.UTC(), ate.UTC()}
}
func decimal(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	x := v.Float64
	return &x
}
