package cursos

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type AcompanhamentoAluno struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Resumo    Resumo `json:"resumo"`
	Situacao  string `json:"situacao"`
	Respostas int    `json:"respostas"`
	Acertos   int    `json:"acertos"`
}

func (s *Servico) Acompanhar(ctx context.Context, mentor, cursoID string) ([]AcompanhamentoAluno, error) {
	var snapshot []byte
	err := s.db.QueryRowContext(ctx, `SELECT publicado FROM cursos WHERE id=? AND mentor_id=? AND ativo=TRUE`, cursoID, mentor).Scan(&snapshot)
	if err == sql.ErrNoRows {
		return nil, ErrNaoEncontrado
	}
	if err != nil {
		return nil, err
	}
	if len(snapshot) == 0 {
		return []AcompanhamentoAluno{}, nil
	}
	var base Curso
	if err = json.Unmarshal(snapshot, &base); err != nil {
		return nil, err
	}
	base.ID = cursoID
	if err = normalizarAulas(&base); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT DISTINCT u.id,u.nome,u.email FROM usuarios u JOIN mentor_alunos ma ON ma.aluno_id=u.id JOIN cursos c ON c.mentor_id=ma.mentor_id
 WHERE c.id=? AND c.mentor_id=? AND ma.ativo=TRUE AND `+acessoRelatorio, cursoID, mentor, s.agora().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	alunos := []AcompanhamentoAluno{}
	for rows.Next() {
		var a AcompanhamentoAluno
		if err = rows.Scan(&a.ID, &a.Nome, &a.Email); err != nil {
			rows.Close()
			return nil, err
		}
		alunos = append(alunos, a)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()
	progressos := map[string]map[string]Progresso{}
	rows, err = s.db.QueryContext(ctx, `SELECT aluno_id,aula_id,video_id,posicao,duracao,concluida,atualizado_em FROM cursos_progresso WHERE curso_id=?`, cursoID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var aluno string
		var p Progresso
		var data time.Time
		if err = rows.Scan(&aluno, &p.AulaID, &p.VideoID, &p.Posicao, &p.Duracao, &p.Concluida, &data); err != nil {
			rows.Close()
			return nil, err
		}
		p.AtualizadoEm = data.UTC().Format(time.RFC3339Nano)
		if progressos[aluno] == nil {
			progressos[aluno] = map[string]Progresso{}
		}
		progressos[aluno][p.AulaID] = p
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()
	type respostaStats struct {
		total, acertos int
		ultima         string
	}
	stats := map[string]respostaStats{}
	vinculos := map[string]bool{}
	for _, a := range base.Aulas {
		for _, q := range a.Questoes {
			vinculos[a.ID+"/"+q] = true
		}
	}
	rows, err = s.db.QueryContext(ctx, `SELECT aluno_id,aula_id,questao_id,correta,atualizado_em FROM cursos_respostas WHERE curso_id=?`, cursoID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var aluno, aula, q string
		var correta bool
		var data time.Time
		if err = rows.Scan(&aluno, &aula, &q, &correta, &data); err != nil {
			rows.Close()
			return nil, err
		}
		if vinculos[aula+"/"+q] {
			st := stats[aluno]
			st.total++
			if correta {
				st.acertos++
			}
			d := data.UTC().Format(time.RFC3339Nano)
			if d > st.ultima {
				st.ultima = d
			}
			stats[aluno] = st
		}
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()
	for i := range alunos {
		a := &alunos[i]
		c := base
		c.Aulas = append([]Aula{}, base.Aulas...)
		aplicarAprendizagem(&c, progressos[a.ID], s.agora())
		a.Resumo = c.Resumo
		st := stats[a.ID]
		a.Respostas = st.total
		a.Acertos = st.acertos
		if st.ultima > a.Resumo.UltimaAtividade {
			a.Resumo.UltimaAtividade = st.ultima
		}
		a.Situacao = "Não iniciado"
		if a.Resumo.UltimaAtividade != "" {
			a.Situacao = "Em andamento"
			ultima, _ := time.Parse(time.RFC3339Nano, a.Resumo.UltimaAtividade)
			if s.agora().Sub(ultima) >= 7*24*time.Hour {
				a.Situacao = "Sem atividade há 7 dias"
			}
		}
		if a.Resumo.Total > 0 && a.Resumo.Concluidas == a.Resumo.Total {
			a.Situacao = "Concluído"
		}
	}
	return alunos, nil
}

const acessoRelatorio = `u.ativo=TRUE AND (ma.data_expiracao_plano IS NULL OR DATE(ma.data_expiracao_plano)>?)
 AND (c.publicado->>'$.escopo'='global' OR (c.publicado->>'$.escopo'='alunos' AND JSON_CONTAINS(c.publicado->'$.destinatarios',JSON_QUOTE(u.id)))
 OR (c.publicado->>'$.escopo'='concursos' AND EXISTS (SELECT 1 FROM aluno_concursos ac JOIN concursos co ON co.id=ac.concurso_id
 WHERE ac.aluno_id=u.id AND ac.ativo=TRUE AND co.ativo=TRUE AND co.criado_por=c.mentor_id AND JSON_CONTAINS(c.publicado->'$.destinatarios',JSON_QUOTE(ac.concurso_id))))) ORDER BY u.nome`
