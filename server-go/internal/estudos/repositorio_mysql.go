package estudos

import (
	"context"
	"database/sql"
	"time"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(banco *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: banco} }

func (r *repositorioMySQL) ListarSessoes(ctx context.Context, alunoID, concursoID, inicio, fim string) ([]Sessao, error) {
	consulta := `SELECT id,aluno_id,concurso_id,materia_id,topico_id,subtopico_id,segundos,modo,observacoes,metricas,origem,estudado_em FROM sessoes_estudo WHERE aluno_id=? AND ativo=TRUE`
	argumentos := []any{alunoID}
	consulta, argumentos = filtros(consulta, argumentos, concursoID, inicio, fim, "estudado_em")
	consulta += ` ORDER BY estudado_em DESC,criado_em DESC`
	linhas, err := r.banco.QueryContext(ctx, consulta, argumentos...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	resultado := make([]Sessao, 0)
	for linhas.Next() {
		var item Sessao
		var concurso, materia, topico, subtopico, modo, observacoes, origem sql.NullString
		var metricas []byte
		var estudadoEm time.Time
		if err := linhas.Scan(&item.ID, &item.AlunoID, &concurso, &materia, &topico, &subtopico, &item.Segundos, &modo, &observacoes, &metricas, &origem, &estudadoEm); err != nil {
			return nil, err
		}
		item.ConcursoID, item.MateriaID, item.TopicoID, item.SubtopicoID = texto(concurso), texto(materia), texto(topico), texto(subtopico)
		item.Modo, item.Observacoes, item.Origem = texto(modo), texto(observacoes), texto(origem)
		item.Metricas = metricas
		item.EstudadoEm = estudadoEm.UTC().Format(time.RFC3339)
		resultado = append(resultado, item)
	}
	return resultado, linhas.Err()
}

func (r *repositorioMySQL) ListarQuestoes(ctx context.Context, alunoID, concursoID, inicio, fim string) ([]RegistroQuestoes, error) {
	consulta := `SELECT id,aluno_id,concurso_id,materia_id,topico_id,subtopico_id,resolvidas,acertos,erros,origem,registrado_em FROM registros_questoes WHERE aluno_id=? AND ativo=TRUE`
	argumentos := []any{alunoID}
	consulta, argumentos = filtros(consulta, argumentos, concursoID, inicio, fim, "registrado_em")
	consulta += ` ORDER BY registrado_em DESC,criado_em DESC`
	linhas, err := r.banco.QueryContext(ctx, consulta, argumentos...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	resultado := make([]RegistroQuestoes, 0)
	for linhas.Next() {
		var item RegistroQuestoes
		var concurso, materia, topico, subtopico, origem sql.NullString
		var registradoEm time.Time
		if err := linhas.Scan(&item.ID, &item.AlunoID, &concurso, &materia, &topico, &subtopico, &item.Resolvidas, &item.Acertos, &item.Erros, &origem, &registradoEm); err != nil {
			return nil, err
		}
		item.ConcursoID, item.MateriaID, item.TopicoID, item.SubtopicoID, item.Origem = texto(concurso), texto(materia), texto(topico), texto(subtopico), texto(origem)
		item.RegistradoEm = registradoEm.UTC().Format(time.RFC3339)
		resultado = append(resultado, item)
	}
	return resultado, linhas.Err()
}

func (r *repositorioMySQL) ConcursoAtribuido(ctx context.Context, alunoID string, concursoID *string) (bool, error) {
	if concursoID == nil || *concursoID == "" {
		return true, nil
	}
	var quantidade int
	err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, alunoID, *concursoID).Scan(&quantidade)
	return quantidade > 0, err
}

func (r *repositorioMySQL) CriarSessao(ctx context.Context, id, alunoID string, entrada EntradaSessao, estudadoEm any) error {
	_, err := r.banco.ExecContext(ctx, `INSERT INTO sessoes_estudo (id,aluno_id,concurso_id,materia_id,topico_id,subtopico_id,segundos,modo,observacoes,metricas,origem,estudado_em) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, id, alunoID, entrada.ConcursoID, entrada.MateriaID, entrada.TopicoID, entrada.SubtopicoID, entrada.Segundos, entrada.Modo, entrada.Observacoes, jsonNulo(entrada.Metricas), entrada.Origem, estudadoEm)
	return err
}

func (r *repositorioMySQL) CriarQuestoes(ctx context.Context, id, alunoID string, entrada EntradaQuestoes, registradoEm any) error {
	_, err := r.banco.ExecContext(ctx, `INSERT INTO registros_questoes (id,aluno_id,concurso_id,materia_id,topico_id,subtopico_id,resolvidas,acertos,erros,origem,registrado_em) VALUES (?,?,?,?,?,?,?,?,?,?,?)`, id, alunoID, entrada.ConcursoID, entrada.MateriaID, entrada.TopicoID, entrada.SubtopicoID, entrada.Resolvidas, entrada.Acertos, entrada.Erros, entrada.Origem, registradoEm)
	return err
}

func (r *repositorioMySQL) Existe(ctx context.Context, tabela, alunoID, id string) (bool, error) {
	var n int
	err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+tabela+` WHERE id=? AND aluno_id=? AND ativo=TRUE`, id, alunoID).Scan(&n)
	return n > 0, err
}
func (r *repositorioMySQL) ObterContadores(ctx context.Context, alunoID, id string) (int, int, int, error) {
	var resolvidas, acertos, erros int
	err := r.banco.QueryRowContext(ctx, `SELECT resolvidas,acertos,erros FROM registros_questoes WHERE id=? AND aluno_id=? AND ativo=TRUE`, id, alunoID).Scan(&resolvidas, &acertos, &erros)
	return resolvidas, acertos, erros, err
}
func (r *repositorioMySQL) AtualizarSessao(ctx context.Context, id string, e AlteracaoSessao, estudadoEm any) error {
	if e.Segundos != nil {
		if _, err := r.banco.ExecContext(ctx, `UPDATE sessoes_estudo SET segundos=? WHERE id=?`, *e.Segundos, id); err != nil {
			return err
		}
	}
	for _, c := range []struct {
		v      *string
		coluna string
	}{{e.ConcursoID, "concurso_id"}, {e.MateriaID, "materia_id"}, {e.TopicoID, "topico_id"}, {e.SubtopicoID, "subtopico_id"}, {e.Modo, "modo"}, {e.Observacoes, "observacoes"}, {e.Origem, "origem"}} {
		if c.v != nil {
			if _, err := r.banco.ExecContext(ctx, `UPDATE sessoes_estudo SET `+c.coluna+`=NULLIF(?,'') WHERE id=?`, *c.v, id); err != nil {
				return err
			}
		}
	}
	if e.Metricas != nil {
		if _, err := r.banco.ExecContext(ctx, `UPDATE sessoes_estudo SET metricas=? WHERE id=?`, jsonNulo(*e.Metricas), id); err != nil {
			return err
		}
	}
	if estudadoEm != nil {
		_, err := r.banco.ExecContext(ctx, `UPDATE sessoes_estudo SET estudado_em=? WHERE id=?`, estudadoEm, id)
		return err
	}
	return nil
}
func (r *repositorioMySQL) AtualizarQuestoes(ctx context.Context, id string, e AlteracaoQuestoes, resolvidas, acertos, erros int, registradoEm any) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `UPDATE registros_questoes SET resolvidas=?,acertos=?,erros=? WHERE id=?`, resolvidas, acertos, erros, id); err != nil {
		return err
	}
	for _, c := range []struct {
		v      *string
		coluna string
	}{{e.ConcursoID, "concurso_id"}, {e.MateriaID, "materia_id"}, {e.TopicoID, "topico_id"}, {e.SubtopicoID, "subtopico_id"}, {e.Origem, "origem"}} {
		if c.v != nil {
			if _, err = tx.ExecContext(ctx, `UPDATE registros_questoes SET `+c.coluna+`=NULLIF(?,'') WHERE id=?`, *c.v, id); err != nil {
				return err
			}
		}
	}
	if registradoEm != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE registros_questoes SET registrado_em=? WHERE id=?`, registradoEm, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}
func (r *repositorioMySQL) Desativar(ctx context.Context, tabela, id string) error {
	_, err := r.banco.ExecContext(ctx, `UPDATE `+tabela+` SET ativo=FALSE WHERE id=?`, id)
	return err
}
