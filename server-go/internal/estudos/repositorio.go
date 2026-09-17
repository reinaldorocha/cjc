package estudos

import "context"

type repositorio interface {
	ListarSessoes(context.Context, string, string, string, string) ([]Sessao, error)
	ListarQuestoes(context.Context, string, string, string, string) ([]RegistroQuestoes, error)
	ConcursoAtribuido(context.Context, string, *string) (bool, error)
	CriarSessao(context.Context, string, string, EntradaSessao, any) error
	CriarQuestoes(context.Context, string, string, EntradaQuestoes, any) error
	Existe(context.Context, string, string, string) (bool, error)
	ObterContadores(context.Context, string, string) (int, int, int, error)
	AtualizarSessao(context.Context, string, AlteracaoSessao, any) error
	AtualizarQuestoes(context.Context, string, AlteracaoQuestoes, int, int, int, any) error
	Desativar(context.Context, string, string) error
}
