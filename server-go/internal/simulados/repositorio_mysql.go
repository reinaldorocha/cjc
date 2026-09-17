package simulados

import (
	"context"
	"database/sql"
	"encoding/json"
	"track-concursos-web/internal/identificador"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(b *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: b} }

func (r *repositorioMySQL) Listar(ctx context.Context, aluno, concurso string) ([]Simulado, error) {
	q := `SELECT id,aluno_id,concurso_id,nome,tipo,realizado_em,link,observacoes,percentual,tempo_minutos,questoes_feitas,resultado,configuracao_usada FROM simulados WHERE aluno_id=? AND ativo=TRUE`
	args := []any{aluno}
	if concurso != "" {
		q += ` AND concurso_id=?`
		args = append(args, concurso)
	}
	linhas, err := r.banco.QueryContext(ctx, q+` ORDER BY COALESCE(realizado_em,'9999-12-31'),criado_em`, args...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Simulado{}
	for linhas.Next() {
		x, err := lerSimulado(linhas)
		if err != nil {
			return nil, err
		}
		rs, err := r.resultados(ctx, x.ID)
		if err != nil {
			return nil, err
		}
		x.ResultadosMaterias = rs
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}

func (r *repositorioMySQL) Criar(ctx context.Context, id, alunoID string, e Entrada) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `INSERT INTO simulados (id,aluno_id,concurso_id,nome,tipo,realizado_em,link,observacoes,percentual,tempo_minutos,questoes_feitas,resultado,configuracao_usada) VALUES (?,?,?,?,?,NULLIF(?,''),?,?,?,?,?,?,?)`, id, alunoID, e.ConcursoID, e.Nome, e.Tipo, valor(e.RealizadoEm), e.Link, e.Observacoes, e.Percentual, e.TempoMinutos, e.QuestoesFeitas, jsonNulo(e.Resultado), jsonNulo(e.ConfiguracaoUsada)); err != nil {
		return err
	}
	if err = r.salvarResultados(ctx, tx, id, e.ResultadosMaterias); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *repositorioMySQL) Alterar(ctx context.Context, id string, e Entrada) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `UPDATE simulados SET concurso_id=?,nome=?,tipo=?,realizado_em=NULLIF(?,''),link=?,observacoes=?,percentual=?,tempo_minutos=?,questoes_feitas=?,resultado=?,configuracao_usada=? WHERE id=?`, e.ConcursoID, e.Nome, e.Tipo, valor(e.RealizadoEm), e.Link, e.Observacoes, e.Percentual, e.TempoMinutos, e.QuestoesFeitas, jsonNulo(e.Resultado), jsonNulo(e.ConfiguracaoUsada), id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE resultados_materias_simulado SET ativo=FALSE,desativado_em=UTC_TIMESTAMP() WHERE simulado_id=? AND ativo=TRUE`, id); err != nil {
		return err
	}
	if err = r.salvarResultados(ctx, tx, id, e.ResultadosMaterias); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *repositorioMySQL) Desativar(ctx context.Context, id string) error {
	_, err := r.banco.ExecContext(ctx, `UPDATE simulados SET ativo=FALSE WHERE id=?`, id)
	return err
}
func (r *repositorioMySQL) Existe(ctx context.Context, alunoID, id string) (bool, error) {
	var n int
	err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM simulados WHERE id=? AND aluno_id=? AND ativo=TRUE`, id, alunoID).Scan(&n)
	return n > 0, err
}
func (r *repositorioMySQL) ConcursoAtribuido(ctx context.Context, alunoID string, concursoID *string) (bool, error) {
	if concursoID == nil || *concursoID == "" {
		return true, nil
	}
	var n int
	err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, alunoID, *concursoID).Scan(&n)
	return n > 0, err
}
func (r *repositorioMySQL) ObterConfiguracao(ctx context.Context, alunoID, concursoID string) (json.RawMessage, bool, error) {
	var bruto []byte
	err := r.banco.QueryRowContext(ctx, `SELECT configuracao FROM configuracoes_prova WHERE aluno_id=? AND concurso_id=?`, alunoID, concursoID).Scan(&bruto)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	return bruto, true, err
}
func (r *repositorioMySQL) SalvarConfiguracao(ctx context.Context, id, alunoID, concursoID, executor string, config json.RawMessage) error {
	_, err := r.banco.ExecContext(ctx, `INSERT INTO configuracoes_prova (id,aluno_id,concurso_id,configuracao,criado_por) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE configuracao=VALUES(configuracao),criado_por=VALUES(criado_por)`, id, alunoID, concursoID, []byte(config), executor)
	return err
}

func (r *repositorioMySQL) salvarResultados(ctx context.Context, tx *sql.Tx, simuladoID string, lista []ResultadoMateria) error {
	for _, x := range lista {
		if _, err := tx.ExecContext(ctx, `INSERT INTO resultados_materias_simulado (id,simulado_id,materia_id,nome,questoes,acertos,erros,brancos,percentual,pontos) VALUES (?,?,?,?,?,?,?,?,?,?)`, identificador.UUID(), simuladoID, x.MateriaID, x.Nome, x.Questoes, x.Acertos, x.Erros, x.Brancos, x.Percentual, x.Pontos); err != nil {
			return err
		}
	}
	return nil
}
func (r *repositorioMySQL) resultados(ctx context.Context, id string) ([]ResultadoMateria, error) {
	linhas, err := r.banco.QueryContext(ctx, `SELECT id,materia_id,nome,questoes,acertos,erros,brancos,percentual,pontos FROM resultados_materias_simulado WHERE simulado_id=? AND ativo=TRUE ORDER BY nome`, id)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []ResultadoMateria{}
	for linhas.Next() {
		var x ResultadoMateria
		var m sql.NullString
		var q, a, e, b sql.NullInt64
		var p, pt sql.NullFloat64
		if err := linhas.Scan(&x.ID, &m, &x.Nome, &q, &a, &e, &b, &p, &pt); err != nil {
			return nil, err
		}
		x.MateriaID, x.Questoes, x.Acertos, x.Erros, x.Brancos, x.Percentual, x.Pontos = texto(m), inteiro(q), inteiro(a), inteiro(e), inteiro(b), decimal(p), decimal(pt)
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}
