package cursos

import (
	"context"
	"math"
	"time"
)

func (s *Servico) carregarProgressos(ctx context.Context, aluno, curso string) (map[string]map[string]Progresso, error) {
	q := `SELECT curso_id,aula_id,video_id,posicao,duracao,concluida,atualizado_em FROM cursos_progresso WHERE aluno_id=?`
	args := []any{aluno}
	if curso != "" {
		q += ` AND curso_id=?`
		args = append(args, curso)
	}
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]map[string]Progresso{}
	for rows.Next() {
		var cursoID string
		var p Progresso
		var data time.Time
		if err = rows.Scan(&cursoID, &p.AulaID, &p.VideoID, &p.Posicao, &p.Duracao, &p.Concluida, &data); err != nil {
			return nil, err
		}
		p.AtualizadoEm = data.UTC().Format(time.RFC3339Nano)
		if result[cursoID] == nil {
			result[cursoID] = map[string]Progresso{}
		}
		result[cursoID][p.AulaID] = p
	}
	return result, rows.Err()
}

type EntradaProgresso struct {
	Posicao   float64 `json:"posicao"`
	Duracao   float64 `json:"duracao"`
	Concluida *bool   `json:"concluida"`
}

func (s *Servico) aulaAcessivel(ctx context.Context, usuario string, mentor bool, curso, aula string) (Curso, Aula, error) {
	lista, err := s.Listar(ctx, usuario, mentor, curso)
	if err != nil {
		return Curso{}, Aula{}, err
	}
	if len(lista) != 1 {
		return Curso{}, Aula{}, ErrNaoEncontrado
	}
	for _, a := range lista[0].Aulas {
		if a.ID == aula {
			if a.Bloqueada {
				return Curso{}, Aula{}, ErrBloqueado
			}
			return lista[0], a, nil
		}
	}
	return Curso{}, Aula{}, ErrNaoEncontrado
}

func (s *Servico) SalvarProgresso(ctx context.Context, aluno, curso, aula string, e EntradaProgresso) (Curso, error) {
	if math.IsNaN(e.Posicao) || math.IsNaN(e.Duracao) || math.IsInf(e.Posicao, 0) || math.IsInf(e.Duracao, 0) || e.Posicao < 0 || e.Duracao < 0 || e.Duracao > 86400 || e.Posicao > 86400 {
		return Curso{}, ErrEntrada
	}
	_, a, err := s.aulaAcessivel(ctx, aluno, false, curso, aula)
	if err != nil {
		return Curso{}, err
	}
	if e.Duracao > 0 && e.Posicao > e.Duracao {
		e.Posicao = e.Duracao
	}
	concluida := false
	if e.Concluida != nil {
		concluida = *e.Concluida
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO cursos_progresso (aluno_id,curso_id,aula_id,video_id,posicao,duracao,concluida,atualizado_em) VALUES (?,?,?,?,?,?,?,UTC_TIMESTAMP(6))
 ON DUPLICATE KEY UPDATE concluida=IF(?,VALUES(concluida),IF(video_id=VALUES(video_id),concluida,FALSE)),video_id=VALUES(video_id),posicao=VALUES(posicao),duracao=VALUES(duracao),atualizado_em=UTC_TIMESTAMP(6)`, aluno, curso, aula, a.VideoID, e.Posicao, e.Duracao, concluida, e.Concluida != nil)
	if err != nil {
		return Curso{}, err
	}
	lista, err := s.Listar(ctx, aluno, false, curso)
	if err != nil {
		return Curso{}, err
	}
	if len(lista) == 0 {
		return Curso{}, ErrNaoEncontrado
	}
	return lista[0], nil
}

func (s *Servico) Previa(ctx context.Context, mentor, curso string) (Curso, error) {
	lista, err := s.Listar(ctx, mentor, true, curso)
	if err != nil {
		return Curso{}, err
	}
	if len(lista) != 1 {
		return Curso{}, ErrNaoEncontrado
	}
	c := lista[0]
	aplicarAprendizagem(&c, nil, s.agora())
	return c, nil
}
