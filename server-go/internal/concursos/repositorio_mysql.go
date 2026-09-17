package concursos

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"track-concursos-web/internal/dominio"
	"track-concursos-web/internal/identificador"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(b *sql.DB) *repositorioMySQL { return &repositorioMySQL{b} }
func (r *repositorioMySQL) Catalogar(ctx context.Context, mentor string) ([]Concurso, error) {
	rows, err := r.banco.QueryContext(ctx, `SELECT id,nome,banca,cargo,logotipo,salario,data_prova,pre_edital,COALESCE(prazos_revisao,'1,7,30') FROM concursos WHERE criado_por=? AND ativo=TRUE ORDER BY atualizado_em DESC,nome`, mentor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Concurso{}
	for rows.Next() {
		var x Concurso
		var cargo, logo sql.NullString
		var sal sql.NullFloat64
		var data sql.NullTime
		if err := rows.Scan(&x.ID, &x.Nome, &x.Banca, &cargo, &logo, &sal, &data, &x.PreEdital, &x.PrazosRevisao); err != nil {
			return nil, err
		}
		x.Cargo = textoNulo(cargo)
		x.Logotipo = textoNulo(logo)
		x.Salario = decimalNulo(sal)
		x.DataProva = dataNula(data)
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *repositorioMySQL) CriarCatalogo(ctx context.Context, mentor, id string, e Atribuicao) error {
	_, err := r.banco.ExecContext(ctx, `INSERT INTO concursos (id,nome,banca,cargo,logotipo,salario,data_prova,pre_edital,prazos_revisao,criado_por) VALUES (?,?,?,NULLIF(?,''),NULLIF(?,''),?,NULLIF(?,''),?,?,?)`, id, strings.TrimSpace(e.Nome), strings.TrimSpace(e.Banca), e.Cargo, e.Logotipo, e.Salario, valor(e.DataProva), e.PreEdital, e.PrazosRevisao, mentor)
	return err
}
func (r *repositorioMySQL) Listar(ctx context.Context, aluno string) ([]Concurso, error) {
	rows, err := r.banco.QueryContext(ctx, `SELECT c.id,c.nome,c.banca,c.cargo,c.logotipo,c.salario,c.data_prova,c.pre_edital,COALESCE(c.prazos_revisao,'1,7,30'),ac.grupo,ac.ordem,ac.resultado,ac.classificacao,ac.nota_final,ac.nomeado,ac.data_nomeacao FROM aluno_concursos ac JOIN concursos c ON c.id=ac.concurso_id WHERE ac.aluno_id=? AND ac.ativo=TRUE AND c.ativo=TRUE ORDER BY ac.ordem,c.nome`, aluno)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Concurso{}
	for rows.Next() {
		var x Concurso
		var cargo, logo, res sql.NullString
		var sal, nota sql.NullFloat64
		var prova, nomeacao sql.NullTime
		var class sql.NullInt64
		if err := rows.Scan(&x.ID, &x.Nome, &x.Banca, &cargo, &logo, &sal, &prova, &x.PreEdital, &x.PrazosRevisao, &x.Grupo, &x.Ordem, &res, &class, &nota, &x.Nomeado, &nomeacao); err != nil {
			return nil, err
		}
		x.Cargo = textoNulo(cargo)
		x.Logotipo = textoNulo(logo)
		x.Salario = decimalNulo(sal)
		x.DataProva = dataNula(prova)
		x.Resultado = textoNulo(res)
		if class.Valid {
			v := int(class.Int64)
			x.Classificacao = &v
		}
		x.NotaFinal = decimalNulo(nota)
		x.DataNomeacao = dataNula(nomeacao)
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *repositorioMySQL) Atribuir(ctx context.Context, aluno, executor, id string, novo bool, e Atribuicao) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if novo {
		if _, err = tx.ExecContext(ctx, `INSERT INTO concursos (id,nome,banca,cargo,logotipo,salario,data_prova,pre_edital,prazos_revisao,criado_por) VALUES (?,?,?,NULLIF(?,''),NULLIF(?,''),?,NULLIF(?,''),?,?,?)`, id, strings.TrimSpace(e.Nome), strings.TrimSpace(e.Banca), e.Cargo, e.Logotipo, e.Salario, valor(e.DataProva), e.PreEdital, e.PrazosRevisao, executor); err != nil {
			return err
		}
	} else {
		var ativo bool
		if err = tx.QueryRowContext(ctx, `SELECT ativo FROM concursos WHERE id=?`, id).Scan(&ativo); err != nil {
			return err
		}
		if !ativo {
			return dominio.ErrNaoEncontrado
		}
		if e.PrazosRevisao != "" {
			if _, err = tx.ExecContext(ctx, `UPDATE concursos SET prazos_revisao=? WHERE id=?`, e.PrazosRevisao, id); err != nil {
				return err
			}
		}
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO aluno_concursos (id,aluno_id,concurso_id,atribuido_por,grupo,ordem,resultado,classificacao,nota_final,nomeado,data_nomeacao) VALUES (?,?,?,?,?,?,NULLIF(?,''),?,?,?,NULLIF(?,''))`, identificador.UUID(), aluno, id, executor, e.Grupo, e.Ordem, e.Resultado, e.Classificacao, e.NotaFinal, e.Nomeado, valor(e.DataNomeacao)); err != nil {
		return err
	}
	var edital string
	err = tx.QueryRowContext(ctx, `SELECT id FROM editais WHERE concurso_id=? AND ativo=TRUE ORDER BY criado_em DESC LIMIT 1`, id).Scan(&edital)
	if err == nil {
		if _, err = tx.ExecContext(ctx, `UPDATE aluno_editais SET ativo=FALSE WHERE aluno_id=? AND edital_id IN (SELECT id FROM editais WHERE concurso_id=?)`, aluno, id); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO aluno_editais (id,aluno_id,edital_id,atribuido_por) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE ativo=TRUE`, identificador.UUID(), aluno, edital, executor); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return tx.Commit()
}
func (r *repositorioMySQL) Reordenar(ctx context.Context, aluno string, lista []Ordem) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, x := range lista {
		var n int
		if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, aluno, x.ID).Scan(&n); err != nil {
			return err
		}
		if n == 0 {
			return dominio.ErrNaoEncontrado
		}
		if _, err = tx.ExecContext(ctx, `UPDATE aluno_concursos SET ordem=? WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, x.Ordem, aluno, x.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}
func (r *repositorioMySQL) Alterar(ctx context.Context, aluno, id string, e Alteracao) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var n int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, aluno, id).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		var cCriado int
		if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM concursos WHERE id=? AND criado_por=? AND ativo=TRUE`, id, aluno).Scan(&cCriado); err != nil {
			return err
		}
		if cCriado == 0 {
			return dominio.ErrNaoEncontrado
		}
	}
	exec := func(q string, args ...any) error { _, x := tx.ExecContext(ctx, q, args...); return x }
	if e.Nome != nil {
		err = exec(`UPDATE concursos SET nome=? WHERE id=?`, strings.TrimSpace(*e.Nome), id)
	}
	if err == nil && e.Banca != nil {
		err = exec(`UPDATE concursos SET banca=? WHERE id=?`, strings.TrimSpace(*e.Banca), id)
	}
	if err == nil && e.Cargo != nil {
		err = exec(`UPDATE concursos SET cargo=NULLIF(?,'') WHERE id=?`, *e.Cargo, id)
	}
	if err == nil && e.Logotipo != nil {
		err = exec(`UPDATE concursos SET logotipo=NULLIF(?,'') WHERE id=?`, *e.Logotipo, id)
	}
	if err == nil && e.Salario != nil {
		err = exec(`UPDATE concursos SET salario=? WHERE id=?`, *e.Salario, id)
	}
	if err == nil && e.LimparSalario {
		err = exec(`UPDATE concursos SET salario=NULL WHERE id=?`, id)
	}
	if err == nil && e.DataProva != nil {
		err = exec(`UPDATE concursos SET data_prova=NULLIF(?,'') WHERE id=?`, *e.DataProva, id)
	}
	if err == nil && e.PreEdital != nil {
		err = exec(`UPDATE concursos SET pre_edital=? WHERE id=?`, *e.PreEdital, id)
	}
	if err == nil && e.PrazosRevisao != nil {
		err = exec(`UPDATE concursos SET prazos_revisao=? WHERE id=?`, strings.TrimSpace(*e.PrazosRevisao), id)
	}
	if err == nil && e.Grupo != nil {
		err = exec(`UPDATE aluno_concursos SET grupo=? WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, *e.Grupo, aluno, id)
	}
	if err == nil && e.Ordem != nil {
		err = exec(`UPDATE aluno_concursos SET ordem=? WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, *e.Ordem, aluno, id)
	}
	if err == nil && e.Resultado != nil {
		err = exec(`UPDATE aluno_concursos SET resultado=NULLIF(?,'') WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, *e.Resultado, aluno, id)
	}
	if err == nil && e.Classificacao != nil {
		err = exec(`UPDATE aluno_concursos SET classificacao=? WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, *e.Classificacao, aluno, id)
	}
	if err == nil && e.LimparClassificacao {
		err = exec(`UPDATE aluno_concursos SET classificacao=NULL WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, aluno, id)
	}
	if err == nil && e.NotaFinal != nil {
		err = exec(`UPDATE aluno_concursos SET nota_final=? WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, *e.NotaFinal, aluno, id)
	}
	if err == nil && e.LimparNotaFinal {
		err = exec(`UPDATE aluno_concursos SET nota_final=NULL WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, aluno, id)
	}
	if err == nil && e.Nomeado != nil {
		err = exec(`UPDATE aluno_concursos SET nomeado=? WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, *e.Nomeado, aluno, id)
	}
	if err == nil && e.DataNomeacao != nil {
		err = exec(`UPDATE aluno_concursos SET data_nomeacao=NULLIF(?,'') WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, *e.DataNomeacao, aluno, id)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
func (r *repositorioMySQL) Desativar(ctx context.Context, aluno, id string) (bool, error) {
	res, err := r.banco.ExecContext(ctx, `UPDATE aluno_concursos SET ativo=FALSE WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, aluno, id)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if n == 0 {
		res2, err2 := r.banco.ExecContext(ctx, `UPDATE concursos SET ativo=FALSE WHERE id=? AND criado_por=? AND ativo=TRUE`, id, aluno)
		if err2 != nil {
			return false, err2
		}
		n2, _ := res2.RowsAffected()
		return n2 > 0, nil
	}
	return n > 0, nil
}
