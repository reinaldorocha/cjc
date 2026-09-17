package editais

import "context"

type repositorio interface {
	ConcursoDoMentor(context.Context, string, string) (bool, error)
	ConcursoDoAluno(context.Context, string, string) (bool, error)
	EditalDoMentor(context.Context, string, string, string) (bool, error)
	PaiDoAluno(context.Context, string, string, string) (bool, error)
	ItemDoAluno(context.Context, string, string, string) (bool, error)
	TipoDoItem(context.Context, string) (string, error)
	ContextoRevisao(context.Context, string, string) (ContextoRevisao, error)
	SincronizarRevisoes(context.Context, string, string, string, *ContextoRevisao, []RevisaoAutomatica) error
	Catalogar(context.Context, string, string) ([]ResumoCatalogo, error)
	ListarCatalogo(context.Context, string, string) ([]Edital, error)
	CriarCatalogo(context.Context, string, Entrada) (string, error)
	Listar(context.Context, string, string) ([]Edital, error)
	Atribuir(context.Context, string, string, Entrada) (string, error)
	CriarItem(context.Context, string, NovoItem) (string, error)
	AlterarItem(context.Context, string, string, AlteracaoItem) error
	Reordenar(context.Context, string, []ItemOrdem) error
	AtualizarProgressoMaterial(context.Context, string, string, EntradaProgressoMaterial) error
	AtualizarProgresso(context.Context, string, string, Progresso) error
}

type ContextoRevisao struct {
	ConcursoID  string
	MateriaID   string
	TopicoID    string
	SubtopicoID *string
	Prazos      string
}

type RevisaoAutomatica struct {
	Ciclo       int
	ProximaData string
}
