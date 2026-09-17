package bancoquestoes

import "context"

type repositorio interface {
	ListarQuestoes(context.Context, FiltroQuestao) ([]Questao, error)
	CriarQuestao(context.Context, EntradaQuestao, string) (*Questao, error)
	ImportarQuestoes(context.Context, []EntradaQuestao, string) (int, error)
	AlterarQuestao(context.Context, string, EntradaQuestao) error
	DesativarQuestao(context.Context, string) error
	Responder(context.Context, string, RespostaQuestaoEntrada) (*ResultadoResposta, error)
	HistoricoQuestao(context.Context, string, string, string) ([]RegistroHistoricoQuestoes, error)
	ObterEstatisticas(context.Context, string, string) (*EstatisticasBancoQuestoes, error)
}
