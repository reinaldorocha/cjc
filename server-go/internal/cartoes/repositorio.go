package cartoes

import "context"

type repositorio interface {
	ListarBaralhosMentor(context.Context, string) ([]Baralho, error)
	PropriedadeBaralho(context.Context, string) (PropriedadeBaralho, error)
	BaralhoDoCartao(context.Context, string) (string, error)
	BaralhoDescendente(context.Context, string, string) (bool, error)
	EditalDoMentor(context.Context, string, string) (bool, error)
	EditalDoAluno(context.Context, string, string) (bool, error)
	EditaisDoAluno(context.Context, string) ([]string, error)
	CriarBaralhoMentor(context.Context, string, EntradaBaralho) (string, error)
	ListarBaralhos(context.Context, string, string, string, string) ([]Baralho, error)
	CriarBaralho(context.Context, string, string, string, EntradaBaralho) (string, error)
	AlterarBaralho(context.Context, string, string, string, string, AlteracaoBaralho) error
	DesativarBaralho(context.Context, string, string, string, string) error
	CriarCartao(context.Context, string, string, string, string, EntradaCartao) (string, error)
	AlterarCartao(context.Context, string, string, string, string, AlteracaoCartao) error
	DesativarCartao(context.Context, string, string, string, string) error
	UltimaRevisao(context.Context, string, string, string) (EstadoRevisao, error)
	RegistrarRevisao(context.Context, string, string, string, Revisao) error
	ListarPendentes(context.Context, string) ([]CartaoPendente, error)
	ListarHistorico(context.Context, string, string) ([]HistoricoRevisao, error)
}

type PropriedadeBaralho struct {
	AlunoID          *string
	CriadoPor        string
	TipoProprietario string
	EditalID         *string
	Alcance          string
}
