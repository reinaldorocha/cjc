package materiaisapoio

import "context"

type repositorio interface {
	Criar(context.Context, string, string, EntradaMaterial) error
	Alterar(context.Context, string, string, EntradaMaterial) (bool, error)
	Desativar(context.Context, string, string) (bool, error)
	ObterParaMentor(context.Context, string, string) (*Material, error)
	ObterParaAluno(context.Context, string, string) (*Material, error)
	ListarMentor(context.Context, string, string) ([]Material, error)
	ListarAluno(context.Context, string, string) ([]Material, error)
	EditalDoMentor(context.Context, string, string) (bool, error)
	MaterialDoMentor(context.Context, string, string) (bool, error)
	CriarArquivoRegistro(context.Context, string, string, EntradaArquivo, string, string) error
	AlterarArquivoRegistro(context.Context, string, string, EntradaArquivo, string, string) (*string, error)
}
