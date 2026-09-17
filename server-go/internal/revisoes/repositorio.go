package revisoes

import "context"

type repositorio interface {
	Listar(context.Context, string, string, string) ([]Revisao, error)
	ConcursoAtribuido(context.Context, string, *string) (bool, error)
	ConteudoValido(context.Context, string, *string, *string, *string, *string) (bool, error)
	Pendente(context.Context, string, Entrada) (string, bool, error)
	Criar(context.Context, string, string, Entrada) error
	AtualizarPendente(context.Context, string, Entrada) error
	Obter(context.Context, string, string) (Entrada, bool, error)
	Alterar(context.Context, string, Alteracao) error
	Desativar(context.Context, string, string) (bool, error)
}
