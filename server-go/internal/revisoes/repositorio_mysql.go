package revisoes

import (
	"context"
	"database/sql"
	"time"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(b *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: b} }
func (r *repositorioMySQL) Listar(ctx context.Context, aluno, concurso, estado string) ([]Revisao, error) {
	q := `SELECT rp.id,rp.aluno_id,rp.concurso_id,rp.materia_id,rp.topico_id,rp.subtopico_id,rp.ciclo_atual,rp.proxima_data,rp.percentual_anterior,rp.observacoes,rp.concluida,rp.concluida_em,COALESCE(em.nome,''),COALESCE(et.nome,''),COALESCE(es.nome,'') FROM revisoes_programadas rp LEFT JOIN edital_materias em ON em.id=rp.materia_id LEFT JOIN edital_topicos et ON et.id=rp.topico_id LEFT JOIN edital_subtopicos es ON es.id=rp.subtopico_id WHERE rp.aluno_id=? AND rp.ativo=TRUE`
	args := []any{aluno}
	if concurso != "" {
		q += ` AND rp.concurso_id=?`
		args = append(args, concurso)
	}
	if estado == "pendentes" {
		q += ` AND rp.concluida=FALSE`
	} else if estado == "concluidas" {
		q += ` AND rp.concluida=TRUE`
	}
	linhas, err := r.banco.QueryContext(ctx, q+` ORDER BY rp.concluida,COALESCE(rp.proxima_data,'9999-12-31'),rp.criado_em`, args...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Revisao{}
	for linhas.Next() {
		var x Revisao
		var concurso, materia, topico, subtopico, obs sql.NullString
		var prox, feito sql.NullTime
		var pct sql.NullFloat64
		if err := linhas.Scan(&x.ID, &x.AlunoID, &concurso, &materia, &topico, &subtopico, &x.CicloAtual, &prox, &pct, &obs, &x.Concluida, &feito, &x.Materia, &x.Topico, &x.Subtopico); err != nil {
			return nil, err
		}
		x.ConcursoID, x.MateriaID, x.TopicoID, x.SubtopicoID, x.Observacoes = texto(concurso), texto(materia), texto(topico), texto(subtopico), texto(obs)
		if prox.Valid {
			v := prox.Time.Format("2006-01-02")
			x.ProximaData = &v
		}
		if pct.Valid {
			x.PercentualAnterior = &pct.Float64
		}
		if feito.Valid {
			v := feito.Time.UTC().Format(time.RFC3339)
			x.ConcluidaEm = &v
		}
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}
func (r *repositorioMySQL) ConcursoAtribuido(ctx context.Context, aluno string, id *string) (bool, error) {
	if id == nil || *id == "" {
		return true, nil
	}
	var n int
	err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, aluno, *id).Scan(&n)
	return n > 0, err
}
func (r *repositorioMySQL) ConteudoValido(ctx context.Context, aluno string, concurso, materia, topico, sub *string) (bool, error) {
	if concurso == nil || *concurso == "" || materia == nil || *materia == "" {
		return false, nil
	}
	q := `SELECT COUNT(*) FROM aluno_editais ae JOIN editais e ON e.id=ae.edital_id JOIN edital_materias em ON em.edital_id=e.id LEFT JOIN edital_topicos et ON et.materia_id=em.id LEFT JOIN edital_subtopicos es ON es.topico_id=et.id WHERE ae.aluno_id=? AND ae.ativo=TRUE AND e.ativo=TRUE AND e.concurso_id=? AND em.id=? AND em.ativo=TRUE`
	args := []any{aluno, *concurso, *materia}
	if topico != nil && *topico != "" {
		q += ` AND et.id=? AND et.ativo=TRUE`
		args = append(args, *topico)
	}
	if sub != nil && *sub != "" {
		q += ` AND es.id=? AND es.ativo=TRUE`
		args = append(args, *sub)
	}
	var n int
	err := r.banco.QueryRowContext(ctx, q, args...).Scan(&n)
	return n > 0, err
}
func (r *repositorioMySQL) Pendente(ctx context.Context, aluno string, e Entrada) (string, bool, error) {
	var id string
	err := r.banco.QueryRowContext(ctx, `SELECT id FROM revisoes_programadas WHERE aluno_id=? AND concurso_id <=> ? AND materia_id <=> ? AND topico_id <=> ? AND subtopico_id <=> ? AND concluida=FALSE AND ativo=TRUE ORDER BY criado_em DESC LIMIT 1`, aluno, e.ConcursoID, e.MateriaID, e.TopicoID, e.SubtopicoID).Scan(&id)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	return id, err == nil, err
}
func (r *repositorioMySQL) AtualizarPendente(ctx context.Context, id string, e Entrada) error {
	_, err := r.banco.ExecContext(ctx, `UPDATE revisoes_programadas SET ciclo_atual=?,proxima_data=NULLIF(?,''),percentual_anterior=?,observacoes=? WHERE id=?`, e.CicloAtual, valor(e.ProximaData), e.PercentualAnterior, e.Observacoes, id)
	return err
}
func (r *repositorioMySQL) Criar(ctx context.Context, id, aluno string, e Entrada) error {
	var em any
	if e.Concluida {
		em = time.Now().UTC()
	}
	_, err := r.banco.ExecContext(ctx, `INSERT INTO revisoes_programadas (id,aluno_id,concurso_id,materia_id,topico_id,subtopico_id,ciclo_atual,proxima_data,percentual_anterior,observacoes,concluida,concluida_em) VALUES (?,?,?,?,?,?,?,NULLIF(?,''),?,?,?,?)`, id, aluno, e.ConcursoID, e.MateriaID, e.TopicoID, e.SubtopicoID, e.CicloAtual, valor(e.ProximaData), e.PercentualAnterior, e.Observacoes, e.Concluida, em)
	return err
}
func (r *repositorioMySQL) Obter(ctx context.Context, aluno, id string) (Entrada, bool, error) {
	var e Entrada
	var c, m, t, s sql.NullString
	err := r.banco.QueryRowContext(ctx, `SELECT concurso_id,materia_id,topico_id,subtopico_id,ciclo_atual,proxima_data,percentual_anterior,observacoes,concluida FROM revisoes_programadas WHERE id=? AND aluno_id=? AND ativo=TRUE`, id, aluno).Scan(&c, &m, &t, &s, &e.CicloAtual, &e.ProximaData, &e.PercentualAnterior, &e.Observacoes, &e.Concluida)
	if err == sql.ErrNoRows {
		return Entrada{}, false, nil
	}
	e.ConcursoID, e.MateriaID, e.TopicoID, e.SubtopicoID = texto(c), texto(m), texto(t), texto(s)
	return e, err == nil, err
}
func (r *repositorioMySQL) Alterar(ctx context.Context, id string, e Alteracao) error {
	for _, c := range []struct {
		v      *string
		coluna string
	}{{e.ConcursoID, "concurso_id"}, {e.MateriaID, "materia_id"}, {e.TopicoID, "topico_id"}, {e.SubtopicoID, "subtopico_id"}, {e.ProximaData, "proxima_data"}, {e.Observacoes, "observacoes"}} {
		if c.v != nil {
			if _, err := r.banco.ExecContext(ctx, `UPDATE revisoes_programadas SET `+c.coluna+`=NULLIF(?,'') WHERE id=?`, *c.v, id); err != nil {
				return err
			}
		}
	}
	if e.CicloAtual != nil {
		if _, err := r.banco.ExecContext(ctx, `UPDATE revisoes_programadas SET ciclo_atual=? WHERE id=?`, *e.CicloAtual, id); err != nil {
			return err
		}
	}
	if e.PercentualAnterior != nil {
		if _, err := r.banco.ExecContext(ctx, `UPDATE revisoes_programadas SET percentual_anterior=? WHERE id=?`, *e.PercentualAnterior, id); err != nil {
			return err
		}
	}
	if e.Concluida != nil {
		em := "NULL"
		if *e.Concluida {
			em = "UTC_TIMESTAMP()"
		}
		_, err := r.banco.ExecContext(ctx, `UPDATE revisoes_programadas SET concluida=?,concluida_em=`+em+` WHERE id=?`, *e.Concluida, id)
		return err
	}
	return nil
}
func (r *repositorioMySQL) Desativar(ctx context.Context, aluno, id string) (bool, error) {
	res, err := r.banco.ExecContext(ctx, `UPDATE revisoes_programadas SET ativo=FALSE WHERE id=? AND aluno_id=? AND ativo=TRUE`, id, aluno)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}
