package auditoria

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Servico struct{ repositorio repositorio }

func Novo(banco *sql.DB) *Servico               { return NovoComRepositorio(novoRepositorioMySQL(banco)) }
func NovoComRepositorio(r repositorio) *Servico { return &Servico{repositorio: r} }
func (s *Servico) Registrar(ctx context.Context, executor, aluno, acao, entidade, entidadeID string, detalhes any) error {
	var bruto []byte
	var err error
	if detalhes != nil {
		bruto, err = json.Marshal(detalhes)
		if err != nil {
			return err
		}
	}
	return s.repositorio.Registrar(ctx, Evento{ExecutorID: executor, AlunoID: aluno, Acao: acao, Entidade: entidade, EntidadeID: entidadeID, Detalhes: bruto})
}
func nulo(v string) any {
	if v == "" {
		return nil
	}
	return v
}
