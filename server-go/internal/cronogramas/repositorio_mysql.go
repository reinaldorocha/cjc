package cronogramas

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"time"
	"chega-junto-concurseiro-web/internal/dominio"
	"chega-junto-concurseiro-web/internal/identificador"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(banco *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: banco} }

func (r *repositorioMySQL) ObterAtivo(ctx context.Context, alunoID, concursoID string) (*Cronograma, error) {
	var c Cronograma
	var concurso sql.NullString
	var cfg []byte
	var criado time.Time
	q := `SELECT id,aluno_id,concurso_id,tipo,criado_por,versao,estado,configuracao,criado_em FROM cronogramas WHERE aluno_id=? AND ativo=TRUE`
	args := []any{alunoID}
	if concursoID != "" {
		q += ` AND concurso_id=?`
		args = append(args, concursoID)
	}
	err := r.banco.QueryRowContext(ctx, q+` ORDER BY criado_em DESC LIMIT 1`, args...).Scan(&c.ID, &c.AlunoID, &concurso, &c.Tipo, &c.CriadoPor, &c.Versao, &c.Estado, &cfg, &criado)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if concurso.Valid {
		c.ConcursoID = &concurso.String
	}
	if len(cfg) > 0 {
		if err := json.Unmarshal(cfg, &c.Configuracao); err != nil {
			return nil, err
		}
	}
	c.CriadoEm = criado.UTC().Format(time.RFC3339)
	itens, err := r.ListarItens(ctx, c.ID, "", "")
	if err != nil {
		return nil, err
	}
	c.Itens = itens
	return &c, nil
}
func (r *repositorioMySQL) Calendario(ctx context.Context, alunoID, concursoID, inicio, fim string) ([]Item, error) {
	var id string
	q := `SELECT id FROM cronogramas WHERE aluno_id=? AND ativo=TRUE`
	args := []any{alunoID}
	if concursoID != "" {
		q += ` AND concurso_id=?`
		args = append(args, concursoID)
	}
	err := r.banco.QueryRowContext(ctx, q+` ORDER BY criado_em DESC LIMIT 1`, args...).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return []Item{}, nil
	}
	if err != nil {
		return nil, err
	}
	return r.ListarItens(ctx, id, inicio, fim)
}
func (r *repositorioMySQL) ListarItens(ctx context.Context, id, inicio, fim string) ([]Item, error) {
	q := `SELECT ci.id,ci.data_planejada,ci.posicao_ciclo,ci.materia_id,ci.topico_id,ci.subtopico_id,COALESCE(es.nome,et.nome,em.nome,''),ci.duracao_minutos,ci.prioridade,ci.ordem,ci.situacao,ci.concluido_em FROM cronograma_itens ci LEFT JOIN edital_materias em ON em.id=ci.materia_id LEFT JOIN edital_topicos et ON et.id=ci.topico_id LEFT JOIN edital_subtopicos es ON es.id=ci.subtopico_id WHERE ci.cronograma_id=?`
	args := []any{id}
	if inicio != "" {
		q += ` AND ci.data_planejada>=?`
		args = append(args, inicio)
	}
	if fim != "" {
		q += ` AND ci.data_planejada<=?`
		args = append(args, fim)
	}
	rows, err := r.banco.QueryContext(ctx, q+` ORDER BY COALESCE(ci.data_planejada,'9999-12-31'),COALESCE(ci.posicao_ciclo,999999),ci.ordem`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Item{}
	for rows.Next() {
		var x Item
		var data, done sql.NullTime
		var pos sql.NullInt64
		var mat, top, sub sql.NullString
		if err := rows.Scan(&x.ID, &data, &pos, &mat, &top, &sub, &x.Nome, &x.DuracaoMinutos, &x.Prioridade, &x.Ordem, &x.Situacao, &done); err != nil {
			return nil, err
		}
		x.DataPlanejada = dataPtr(data)
		x.PosicaoCiclo = intPtr(pos)
		x.MateriaID = stringPtr(mat)
		x.TopicoID = stringPtr(top)
		x.SubtopicoID = stringPtr(sub)
		if done.Valid {
			v := done.Time.UTC().Format(time.RFC3339)
			x.ConcluidoEm = &v
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *repositorioMySQL) ConcursoAtribuido(ctx context.Context, aluno, concurso string) (bool, error) {
	var n int
	err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, aluno, concurso).Scan(&n)
	return n > 0, err
}
func (r *repositorioMySQL) ConteudoAtribuido(ctx context.Context, aluno, concurso string, item Item) (bool, error) {
	q := `SELECT COUNT(*) FROM aluno_editais ae JOIN editais e ON e.id=ae.edital_id JOIN edital_materias em ON em.edital_id=e.id LEFT JOIN edital_topicos et ON et.materia_id=em.id LEFT JOIN edital_subtopicos es ON es.topico_id=et.id WHERE ae.aluno_id=? AND ae.ativo=TRUE AND e.concurso_id=? AND e.ativo=TRUE AND em.ativo=TRUE AND em.id=?`
	args := []any{aluno, concurso, *item.MateriaID}
	if item.TopicoID != nil && *item.TopicoID != "" {
		q += ` AND et.id=? AND et.ativo=TRUE`
		args = append(args, *item.TopicoID)
	}
	if item.SubtopicoID != nil && *item.SubtopicoID != "" {
		q += ` AND es.id=? AND es.ativo=TRUE`
		args = append(args, *item.SubtopicoID)
	}
	var n int
	err := r.banco.QueryRowContext(ctx, q, args...).Scan(&n)
	return n > 0, err
}
func (r *repositorioMySQL) Diagnostico(ctx context.Context, aluno, concurso string, cfg Configuracao) ([]materiaInfo, map[string][]unidade, error) {
	rows, err := r.banco.QueryContext(ctx, `SELECT em.id,em.nome FROM aluno_editais ae JOIN editais e ON e.id=ae.edital_id JOIN edital_materias em ON em.edital_id=e.id WHERE ae.aluno_id=? AND ae.ativo=TRUE AND e.concurso_id=? AND e.ativo=TRUE AND em.ativo=TRUE ORDER BY em.ordem`, aluno, concurso)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	desmarcada := false
	for _, v := range cfg.MateriasSelecionadas {
		if !v {
			desmarcada = true
			break
		}
	}
	materias := []materiaInfo{}
	for rows.Next() {
		var m materiaInfo
		if err := rows.Scan(&m.id, &m.nome); err != nil {
			return nil, nil, err
		}
		if desmarcada {
			if v, ok := cfg.MateriasSelecionadas[m.id]; ok && !v {
				continue
			}
		}
		if err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM edital_topicos WHERE materia_id=? AND ativo=TRUE`, m.id).Scan(&m.total); err != nil {
			return nil, nil, err
		}
		if err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_progresso_edital ap JOIN edital_topicos et ON et.id=ap.item_id WHERE ap.aluno_id=? AND ap.tipo_item='topico' AND ap.estudado=TRUE AND et.materia_id=?`, aluno, m.id).Scan(&m.concluidos); err != nil {
			return nil, nil, err
		}
		if err := r.banco.QueryRowContext(ctx, `SELECT COALESCE(SUM(resolvidas),0),COALESCE(SUM(acertos),0) FROM registros_questoes WHERE aluno_id=? AND materia_id=?`, aluno, m.id).Scan(&m.resolvidas, &m.acertos); err != nil {
			return nil, nil, err
		}
		restante := 0.0
		if m.total > 0 {
			restante = (1 - float64(m.concluidos)/float64(m.total)) * 40
		}
		penalidade := 0.0
		if m.resolvidas > 0 && float64(m.acertos)/float64(m.resolvidas)*100 < 65 {
			penalidade = 15
		}
		m.prioridade = int(50 + restante + penalidade)
		if p, ok := map[string]int{"baixa": 32, "media": 47, "alta": 62, "maxima": 78}[cfg.MateriaPrioridades[m.id]]; ok {
			m.prioridade = p
		}
		materias = append(materias, m)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	sort.SliceStable(materias, func(i, j int) bool { return materias[i].prioridade > materias[j].prioridade })
	unidades := map[string][]unidade{}
	for _, m := range materias {
		linhas, err := r.banco.QueryContext(ctx, `SELECT et.id,et.nome,es.id,es.nome,COALESCE(CASE WHEN es.id IS NULL THEN apt.estudado ELSE aps.estudado END,FALSE) FROM edital_topicos et LEFT JOIN edital_subtopicos es ON es.topico_id=et.id AND es.ativo=TRUE LEFT JOIN aluno_progresso_edital apt ON apt.aluno_id=? AND apt.tipo_item='topico' AND apt.item_id=et.id LEFT JOIN aluno_progresso_edital aps ON aps.aluno_id=? AND aps.tipo_item='subtopico' AND aps.item_id=es.id WHERE et.materia_id=? AND et.ativo=TRUE ORDER BY et.ordem,es.ordem`, aluno, aluno, m.id)
		if err != nil {
			return nil, nil, err
		}
		lista := []unidade{}
		for linhas.Next() {
			var tid, tn string
			var sid, sn sql.NullString
			var estudado bool
			if err := linhas.Scan(&tid, &tn, &sid, &sn, &estudado); err != nil {
				linhas.Close()
				return nil, nil, err
			}
			if !estudado {
				u := unidade{materiaID: m.id, topicoID: tid, nome: tn, prioridade: m.prioridade}
				if sid.Valid {
					u.subtopicoID = sid.String
					u.nome = sn.String
				}
				lista = append(lista, u)
			}
		}
		if err := linhas.Err(); err != nil {
			linhas.Close()
			return nil, nil, err
		}
		if err := linhas.Close(); err != nil {
			return nil, nil, err
		}
		if len(lista) == 0 {
			lista = append(lista, unidade{materiaID: m.id, nome: m.nome, prioridade: m.prioridade})
		}
		unidades[m.id] = lista
	}
	return materias, unidades, nil
}
func (r *repositorioMySQL) Substituir(ctx context.Context, aluno, executor string, e Entrada, itens []Item, cfg json.RawMessage) (string, error) {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	var versao int
	if err = tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(versao),0)+1 FROM cronogramas WHERE aluno_id=?`, aluno).Scan(&versao); err != nil {
		return "", err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE cronogramas SET ativo=FALSE,estado='concluido',desativado_em=UTC_TIMESTAMP() WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, aluno, *e.ConcursoID); err != nil {
		return "", err
	}
	id := identificador.UUID()
	if _, err = tx.ExecContext(ctx, `INSERT INTO cronogramas (id,aluno_id,concurso_id,tipo,criado_por,versao,estado,configuracao) VALUES (?,?,?,?,?,?,'ativo',?)`, id, aluno, e.ConcursoID, e.Tipo, executor, versao, []byte(cfg)); err != nil {
		return "", err
	}
	porDia := map[string][]Item{}
	for i, item := range itens {
		situacao := item.Situacao
		if situacao == "" {
			situacao = "pendente"
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO cronograma_itens (id,cronograma_id,data_planejada,posicao_ciclo,materia_id,topico_id,subtopico_id,duracao_minutos,prioridade,ordem,situacao) VALUES (?,?,NULLIF(?,''),?,?,?,?,?,?,?,?)`, identificador.UUID(), id, valor(item.DataPlanejada), item.PosicaoCiclo, item.MateriaID, item.TopicoID, item.SubtopicoID, item.DuracaoMinutos, item.Prioridade, ordem(item.Ordem, i), situacao); err != nil {
			return "", err
		}
		if item.DataPlanejada != nil && *item.DataPlanejada != "" {
			porDia[*item.DataPlanejada] = append(porDia[*item.DataPlanejada], item)
		}
	}
	for data, lista := range porDia {
		conteudo, err := json.Marshal(lista)
		if err != nil {
			return "", err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO cronograma_resumos_diarios (id,cronograma_id,aluno_id,data_referencia,conteudo) VALUES (?,?,?,?,?)`, identificador.UUID(), id, aluno, data, conteudo); err != nil {
			return "", err
		}
	}
	return id, tx.Commit()
}
func (r *repositorioMySQL) ObterItemAtivo(ctx context.Context, aluno, id string) (itemCronograma, error) {
	var x itemCronograma
	var duracao sql.NullInt64
	var mat, top, sub, concurso sql.NullString
	err := r.banco.QueryRowContext(ctx, `SELECT ci.situacao,ci.duracao_minutos,ci.materia_id,ci.topico_id,ci.subtopico_id,c.concurso_id FROM cronograma_itens ci JOIN cronogramas c ON c.id=ci.cronograma_id WHERE c.aluno_id=? AND c.ativo=TRUE AND ci.id=?`, aluno, id).Scan(&x.Situacao, &duracao, &mat, &top, &sub, &concurso)
	if errors.Is(err, sql.ErrNoRows) {
		return x, dominio.ErrNaoEncontrado
	}
	if err != nil {
		return x, err
	}
	x.DuracaoMinutos = int(duracao.Int64)
	x.MateriaID = stringPtr(mat)
	x.TopicoID = stringPtr(top)
	x.SubtopicoID = stringPtr(sub)
	x.ConcursoID = stringPtr(concurso)
	return x, nil
}
func (r *repositorioMySQL) AlterarItem(ctx context.Context, aluno, id string, e AlteracaoItem, sessao *sessaoCronograma, desativar bool) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if e.DataPlanejada != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE cronograma_itens SET data_planejada=NULLIF(?,'') WHERE id=?`, *e.DataPlanejada, id); err != nil {
			return err
		}
	}
	if e.PosicaoCiclo != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE cronograma_itens SET posicao_ciclo=? WHERE id=?`, *e.PosicaoCiclo, id); err != nil {
			return err
		}
	}
	if e.DuracaoMinutos != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE cronograma_itens SET duracao_minutos=? WHERE id=?`, *e.DuracaoMinutos, id); err != nil {
			return err
		}
	}
	if e.Prioridade != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE cronograma_itens SET prioridade=? WHERE id=?`, *e.Prioridade, id); err != nil {
			return err
		}
	}
	if e.Ordem != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE cronograma_itens SET ordem=? WHERE id=?`, *e.Ordem, id); err != nil {
			return err
		}
	}
	if e.Situacao != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE cronograma_itens SET situacao=?,concluido_em=CASE WHEN ?='concluido' THEN UTC_TIMESTAMP() ELSE NULL END WHERE id=?`, *e.Situacao, *e.Situacao, id); err != nil {
			return err
		}
	}
	if sessao != nil {
		metricas, err := json.Marshal(map[string]any{"cronograma_item_id": sessao.ItemID})
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO sessoes_estudo (id,aluno_id,concurso_id,materia_id,topico_id,subtopico_id,segundos,modo,origem,metricas,estudado_em) VALUES (?,?,?,?,?,?,?,'cronograma','cronograma',?,UTC_TIMESTAMP())`, sessao.ID, sessao.AlunoID, sessao.ConcursoID, sessao.MateriaID, sessao.TopicoID, sessao.SubtopicoID, sessao.Segundos, metricas); err != nil {
			return err
		}
	}
	if desativar {
		if _, err = tx.ExecContext(ctx, `UPDATE sessoes_estudo SET ativo=FALSE WHERE aluno_id=? AND origem='cronograma' AND (JSON_UNQUOTE(JSON_EXTRACT(metricas,'$.cronograma_item_id'))=? OR metricas LIKE ?)`, aluno, id, "%"+id+"%"); err != nil {
			return err
		}
	}
	return tx.Commit()
}
func (r *repositorioMySQL) CronogramaAtivo(ctx context.Context, aluno string) (cronogramaAtivo, error) {
	var x cronogramaAtivo
	err := r.banco.QueryRowContext(ctx, `SELECT id,configuracao FROM cronogramas WHERE aluno_id=? AND ativo=TRUE ORDER BY criado_em DESC LIMIT 1`, aluno).Scan(&x.ID, &x.Configuracao)
	if errors.Is(err, sql.ErrNoRows) {
		return x, dominio.ErrNaoEncontrado
	}
	return x, err
}
func (r *repositorioMySQL) CriarItem(ctx context.Context, aluno string, item Item) error {
	var id string
	err := r.banco.QueryRowContext(ctx, `SELECT id FROM cronogramas WHERE aluno_id=? AND ativo=TRUE ORDER BY criado_em DESC LIMIT 1`, aluno).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.ErrNaoEncontrado
	}
	if err != nil {
		return err
	}
	situacao := item.Situacao
	if situacao == "" {
		situacao = "pendente"
	}
	duracao := item.DuracaoMinutos
	if duracao <= 0 {
		duracao = 60
	}
	if _, err = r.banco.ExecContext(ctx, `INSERT INTO cronograma_itens (id,cronograma_id,data_planejada,posicao_ciclo,materia_id,topico_id,subtopico_id,duracao_minutos,prioridade,ordem,situacao) VALUES (?,?,NULLIF(?,''),?,?,?,?,?,?,?,?)`, identificador.UUID(), id, valor(item.DataPlanejada), item.PosicaoCiclo, item.MateriaID, item.TopicoID, item.SubtopicoID, duracao, item.Prioridade, 999, situacao); err != nil {
		return err
	}
	return nil
}
func (r *repositorioMySQL) ItensPendentes(ctx context.Context, cronograma, inicio string) ([]string, error) {
	rows, err := r.banco.QueryContext(ctx, `SELECT id FROM cronograma_itens WHERE cronograma_id=? AND situacao='pendente' ORDER BY CASE WHEN data_planejada IS NULL THEN '9999-12-31' ELSE data_planejada END,ordem,id`, cronograma)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
func (r *repositorioMySQL) Reprogramar(ctx context.Context, itens []reprogramacaoItem) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, item := range itens {
		if _, err = tx.ExecContext(ctx, `UPDATE cronograma_itens SET data_planejada=?,ordem=? WHERE id=?`, item.Data, item.Ordem, item.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}
