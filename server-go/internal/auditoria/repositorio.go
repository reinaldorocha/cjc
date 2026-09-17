package auditoria

import "context"

type repositorio interface {
	Registrar(context.Context, Evento) error
}

type Evento struct {
	ExecutorID, AlunoID, Acao, Entidade, EntidadeID string
	Detalhes                                        []byte
}
