package cadernos

import "context"

// repositorio is deliberately small: persistence stays behind this boundary.
type repositorio interface {
	Criar(context.Context, string, string, EntradaCadernoNormalizada) error
	Listar(context.Context, string, string) ([]Caderno, error)
	Obter(context.Context, string, string) (*Caderno, error)
	Alterar(context.Context, string, string, EntradaCadernoNormalizada) error
	Desativar(context.Context, string, string) error
}

type EntradaCadernoNormalizada struct {
	Titulo    string
	Pasta     string
	Conteudo  string
	EditalID  *string
	MateriaID *string
	TopicoID  *string
	Cor       string
}
