package cronogramas

import (
	"context"
	"encoding/json"
)

func (r *repositorioMySQL) ReplanejarCalendario(ctx context.Context, aluno, concurso, inicio, hoje string) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	q := `SELECT id,COALESCE(concurso_id,''),configuracao FROM cronogramas WHERE aluno_id=? AND ativo=TRUE AND tipo<>'ciclo_inteligente'`
	args := []any{aluno}
	if concurso != "" {
		q += ` AND concurso_id=?`
		args = append(args, concurso)
	}
	rows, err := tx.QueryContext(ctx, q+` ORDER BY id FOR UPDATE`, args...)
	if err != nil {
		return err
	}
	type agenda struct {
		id, concurso string
		cfg          Configuracao
	}
	agendas := []agenda{}
	for rows.Next() {
		var a agenda
		var raw []byte
		if err = rows.Scan(&a.id, &a.concurso, &raw); err != nil {
			rows.Close()
			return err
		}
		if len(raw) > 0 {
			if err = json.Unmarshal(raw, &a.cfg); err != nil {
				rows.Close()
				return err
			}
		}
		agendas = append(agendas, a)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()
	for _, a := range agendas {
		tarefas := []tarefaPlanejada{}
		fixos := map[string]ocupacaoDia{}
		rows, err = tx.QueryContext(ctx, `SELECT id,COALESCE(DATE_FORMAT(data_planejada,'%Y-%m-%d'),''),COALESCE(duracao_minutos,0),situacao FROM cronograma_itens WHERE cronograma_id=? ORDER BY COALESCE(data_planejada,'9999-12-31'),ordem,id FOR UPDATE`, a.id)
		if err != nil {
			return err
		}
		for rows.Next() {
			var t tarefaPlanejada
			var estado string
			if err = rows.Scan(&t.ID, &t.Data, &t.Minutos, &estado); err != nil {
				rows.Close()
				return err
			}
			if t.Minutos <= 0 {
				t.Minutos = a.cfg.MinutosTopico
				if t.Minutos <= 0 {
					t.Minutos = 60
				}
			}
			if estado == "pendente" && (t.Data == "" || t.Data >= inicio || t.Data < hoje) {
				tarefas = append(tarefas, t)
			} else if estado != "ignorado" && t.Data >= inicio {
				o := fixos[t.Data]
				o.Minutos += t.Minutos
				o.Quantidade++
				fixos[t.Data] = o
			}
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return err
		}
		rows.Close()
		rows, err = tx.QueryContext(ctx, `SELECT id,DATE_FORMAT(proxima_data,'%Y-%m-%d'),concluida,CONCAT(COALESCE(materia_id,''),'/',COALESCE(topico_id,''),'/',COALESCE(subtopico_id,'')) FROM revisoes_programadas WHERE aluno_id=? AND concurso_id <=> NULLIF(?,'') AND ativo=TRUE AND proxima_data IS NOT NULL ORDER BY proxima_data,ciclo_atual,id FOR UPDATE`, aluno, a.concurso)
		if err != nil {
			return err
		}
		for rows.Next() {
			t := tarefaPlanejada{Revisao: true, Minutos: MinutosRevisao}
			var concluida bool
			if err = rows.Scan(&t.ID, &t.Data, &concluida, &t.Grupo); err != nil {
				rows.Close()
				return err
			}
			if !concluida && (t.Data >= inicio || t.Data < hoje) {
				tarefas = append(tarefas, t)
			} else if t.Data >= inicio {
				o := fixos[t.Data]
				o.Minutos += MinutosRevisao
				o.Quantidade++
				fixos[t.Data] = o
			}
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return err
		}
		rows.Close()
		plano, err := distribuirComRevisoes(tarefas, fixos, a.cfg, inicio)
		if err != nil {
			return err
		}
		ordens := map[string]int{}
		for _, t := range plano {
			if t.Revisao {
				_, err = tx.ExecContext(ctx, `UPDATE revisoes_programadas SET proxima_data=? WHERE id=? AND aluno_id=? AND concluida=FALSE AND ativo=TRUE`, t.Data, t.ID, aluno)
			} else {
				ordens[t.Data]++
				_, err = tx.ExecContext(ctx, `UPDATE cronograma_itens SET data_planejada=?,ordem=? WHERE id=? AND cronograma_id=? AND situacao='pendente'`, t.Data, ordens[t.Data], t.ID, a.id)
			}
			if err != nil {
				return err
			}
		}
		// O calendário lê os itens; descartar resumos antigos evita manter datas inconsistentes.
		if _, err = tx.ExecContext(ctx, `DELETE FROM cronograma_resumos_diarios WHERE cronograma_id=?`, a.id); err != nil {
			return err
		}
	}
	return tx.Commit()
}
