package concursos

import "context"

type repositorio interface {
	Catalogar(context.Context, string) ([]Concurso, error)
	CriarCatalogo(context.Context, string, string, Atribuicao) error
	Listar(context.Context, string) ([]Concurso, error)
	Atribuir(context.Context, string, string, string, bool, Atribuicao) error
	Reordenar(context.Context, string, []Ordem) error
	Alterar(context.Context, string, string, Alteracao) error
	Desativar(context.Context, string, string) (bool, error)
}
