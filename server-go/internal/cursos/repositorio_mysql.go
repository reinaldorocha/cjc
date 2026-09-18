package cursos

import (
	"context"
	"database/sql"
	"encoding/json"
	"chega-junto-concurseiro-web/internal/identificador"
)

const acessoAluno = `EXISTS (SELECT 1 FROM mentor_alunos ma JOIN usuarios u ON u.id=ma.mentor_id JOIN usuarios aluno ON aluno.id=ma.aluno_id
 WHERE ma.aluno_id=? AND ma.mentor_id=c.mentor_id AND ma.ativo=TRUE AND u.ativo=TRUE AND aluno.ativo=TRUE
 AND (ma.data_expiracao_plano IS NULL OR DATE(ma.data_expiracao_plano)>?))
 AND (c.publicado->>'$.escopo'='global' OR (c.publicado->>'$.escopo'='alunos' AND JSON_CONTAINS(c.publicado->'$.destinatarios',JSON_QUOTE(?)))
 OR (c.publicado->>'$.escopo'='concursos' AND EXISTS (SELECT 1 FROM aluno_concursos ac JOIN concursos co ON co.id=ac.concurso_id
 WHERE ac.aluno_id=? AND ac.ativo=TRUE AND co.ativo=TRUE AND co.criado_por=c.mentor_id AND JSON_CONTAINS(c.publicado->'$.destinatarios',JSON_QUOTE(ac.concurso_id)))))`

func (s *Servico) Listar(ctx context.Context, usuario string, mentor bool, id string) ([]Curso, error) {
	q := `SELECT c.id,c.mentor_id,c.titulo,c.descricao,c.categoria,c.capa_url,c.modo_exibicao,c.escopo,c.destinatarios,c.aulas,c.publicado,c.revisao,c.revisao_publicada FROM cursos c WHERE c.ativo=TRUE AND `
	args := []any{}
	if mentor {
		q += `c.mentor_id=?`
		args = append(args, usuario)
	} else {
		q += `c.publicado IS NOT NULL AND ` + acessoAluno
		args = append(args, usuario, s.agora().Format("2006-01-02"), usuario, usuario)
	}
	if id != "" {
		q += ` AND c.id=?`
		args = append(args, id)
	}
	q += ` ORDER BY c.atualizado_em DESC,c.titulo`
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lista := []Curso{}
	for rows.Next() {
		var c Curso
		var destinos, aulas, publicado []byte
		var revPub int
		if err = rows.Scan(&c.ID, &c.MentorID, &c.Titulo, &c.Descricao, &c.Categoria, &c.CapaURL, &c.ModoExibicao, &c.Escopo, &destinos, &aulas, &publicado, &c.Revisao, &revPub); err != nil {
			return nil, err
		}
		if c.ModoExibicao == "" {
			c.ModoExibicao = "curso"
		}
		if mentor {
			if err = json.Unmarshal(aulas, &c.Aulas); err != nil {
				return nil, err
			}
			if err = json.Unmarshal(destinos, &c.Destinatarios); err != nil {
				return nil, err
			}
		} else {
			id, mentorID, rev := c.ID, c.MentorID, c.Revisao
			// Publicações antigas não têm esta propriedade: nunca herdar o rascunho.
			c.ModoExibicao = "curso"
			if err = json.Unmarshal(publicado, &c); err != nil {
				return nil, err
			}
			c.ID = id
			c.MentorID = mentorID
			c.Revisao = rev
		}
		c.Publicado = len(publicado) > 0
		c.AlteracoesPendentes = c.Revisao != revPub
		if err = normalizarAulas(&c); err != nil {
			return nil, err
		}
		c.Resumo = Resumo{Total: len(c.Aulas)}
		lista = append(lista, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	if !mentor {
		todos, err := s.carregarProgressos(ctx, usuario, id)
		if err != nil {
			return nil, err
		}
		for i := range lista {
			aplicarAprendizagem(&lista[i], todos[lista[i].ID], s.agora())
		}
	}
	return lista, nil
}

func (s *Servico) salvar(ctx context.Context, mentor, id string, c Curso) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	novo := id == ""
	if novo {
		id = identificador.UUID()
	} else {
		var rev int
		err = tx.QueryRowContext(ctx, `SELECT revisao FROM cursos WHERE id=? AND mentor_id=? AND ativo=TRUE FOR UPDATE`, id, mentor).Scan(&rev)
		if err == sql.ErrNoRows {
			return "", ErrNaoEncontrado
		}
		if err != nil {
			return "", err
		}
		if rev != c.Revisao {
			return "", ErrConflito
		}
	}
	if err = validarReferencias(ctx, tx, mentor, c); err != nil {
		return "", err
	}
	for i := range c.Aulas {
		a := &c.Aulas[i]
		if a.PDFID != "" {
			if err = tx.QueryRowContext(ctx, `SELECT nome FROM cursos_pdfs WHERE id=? AND mentor_id=?`, a.PDFID, mentor).Scan(&a.PDFNome); err != nil {
				return "", err
			}
		}
	}
	destinos, err := json.Marshal(c.Destinatarios)
	if err != nil {
		return "", err
	}
	aulas, err := json.Marshal(c.Aulas)
	if err != nil {
		return "", err
	}
	if c.ModoExibicao == "" {
		c.ModoExibicao = "curso"
	}
	if novo {
		_, err = tx.ExecContext(ctx, `INSERT INTO cursos (id,mentor_id,titulo,descricao,categoria,capa_url,modo_exibicao,escopo,destinatarios,aulas) VALUES (?,?,?,?,?,?,?,?,?,?)`, id, mentor, c.Titulo, c.Descricao, c.Categoria, c.CapaURL, c.ModoExibicao, c.Escopo, destinos, aulas)
	} else {
		_, err = tx.ExecContext(ctx, `UPDATE cursos SET titulo=?,descricao=?,categoria=?,capa_url=?,modo_exibicao=?,escopo=?,destinatarios=?,aulas=?,revisao=revisao+1 WHERE id=? AND mentor_id=?`, c.Titulo, c.Descricao, c.Categoria, c.CapaURL, c.ModoExibicao, c.Escopo, destinos, aulas, id, mentor)
	}
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func validarReferencias(ctx context.Context, tx *sql.Tx, mentor string, c Curso) error {
	verificar := func(q string, args ...any) error {
		var id string
		err := tx.QueryRowContext(ctx, q, args...).Scan(&id)
		if err == sql.ErrNoRows {
			return ErrEntrada
		}
		return err
	}
	if id := idCapa(c.CapaURL); id != "" {
		if err := verificar(`SELECT id FROM cursos_capas WHERE id=? AND mentor_id=?`, id, mentor); err != nil {
			return err
		}
	}
	for _, destino := range c.Destinatarios {
		q := `SELECT aluno_id FROM mentor_alunos WHERE aluno_id=? AND mentor_id=? AND ativo=TRUE FOR UPDATE`
		if c.Escopo == "concursos" {
			q = `SELECT id FROM concursos WHERE id=? AND criado_por=? AND ativo=TRUE FOR UPDATE`
		}
		if err := verificar(q, destino, mentor); err != nil {
			return err
		}
	}
	pdfs, questoes := map[string]bool{}, map[string]bool{}
	for _, a := range c.Aulas {
		if id := idCapa(a.ModuloCapaURL); id != "" {
			if err := verificar(`SELECT id FROM cursos_capas WHERE id=? AND mentor_id=?`, id, mentor); err != nil {
				return err
			}
		}
		if a.PDFID != "" && !pdfs[a.PDFID] {
			if err := verificar(`SELECT id FROM cursos_pdfs WHERE id=? AND mentor_id=?`, a.PDFID, mentor); err != nil {
				return err
			}
			pdfs[a.PDFID] = true
		}
		for _, id := range a.Questoes {
			if !questoes[id] {
				if err := verificar(`SELECT id FROM banco_questoes WHERE id=? AND criado_por=? AND ativo=TRUE`, id, mentor); err != nil {
					return err
				}
				questoes[id] = true
			}
		}
	}
	return nil
}

func (s *Servico) Publicar(ctx context.Context, mentor, id string, revisao int, publicar bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var c Curso
	var aulas, destinos []byte
	err = tx.QueryRowContext(ctx, `SELECT id,titulo,descricao,categoria,capa_url,modo_exibicao,escopo,destinatarios,aulas,revisao FROM cursos WHERE id=? AND mentor_id=? AND ativo=TRUE FOR UPDATE`, id, mentor).Scan(&c.ID, &c.Titulo, &c.Descricao, &c.Categoria, &c.CapaURL, &c.ModoExibicao, &c.Escopo, &destinos, &aulas, &c.Revisao)
	if err == sql.ErrNoRows {
		return ErrNaoEncontrado
	}
	if err != nil {
		return err
	}
	if c.ModoExibicao == "" {
		c.ModoExibicao = "curso"
	}
	if c.Revisao != revisao {
		return ErrConflito
	}
	if publicar {
		if err = json.Unmarshal(aulas, &c.Aulas); err != nil {
			return err
		}
		if err = json.Unmarshal(destinos, &c.Destinatarios); err != nil {
			return err
		}
		if err = validarPublicacao(c); err != nil {
			return err
		}
		if err = normalizarAulas(&c); err != nil {
			return err
		}
		if err = validarReferencias(ctx, tx, mentor, c); err != nil {
			return err
		}
		snapshot, err := json.Marshal(c)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `UPDATE cursos SET publicado=?,revisao=revisao+1,revisao_publicada=revisao WHERE id=?`, snapshot, id)
		if err != nil {
			return err
		}
	} else {
		if _, err = tx.ExecContext(ctx, `UPDATE cursos SET publicado=NULL,revisao=revisao+1,revisao_publicada=0 WHERE id=?`, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Servico) Desativar(ctx context.Context, mentor, id string) error {
	res, err := s.db.ExecContext(ctx, `UPDATE cursos SET ativo=FALSE WHERE id=? AND mentor_id=? AND ativo=TRUE`, id, mentor)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
