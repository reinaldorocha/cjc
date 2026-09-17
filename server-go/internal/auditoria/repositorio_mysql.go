package auditoria

import (
	"context"
	"database/sql"
	"track-concursos-web/internal/identificador"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(banco *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: banco} }
func (r *repositorioMySQL) Registrar(ctx context.Context, e Evento) error {
	_, err := r.banco.ExecContext(ctx, `INSERT INTO eventos_auditoria (id,executor_id,aluno_id,acao,entidade,entidade_id,detalhes) VALUES (?,?,?,?,?,?,?)`, identificador.UUID(), e.ExecutorID, nulo(e.AlunoID), e.Acao, e.Entidade, nulo(e.EntidadeID), e.Detalhes)
	return err
}
