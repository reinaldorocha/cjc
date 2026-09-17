package cronogramas

import (
	"context"
	"encoding/json"
)

// repositorio delimita toda a persistencia do dominio de cronogramas.
type repositorio interface {
	ObterAtivo(context.Context, string, string) (*Cronograma, error)
	Calendario(context.Context, string, string, string, string) ([]Item, error)
	ConcursoAtribuido(context.Context, string, string) (bool, error)
	ConteudoAtribuido(context.Context, string, string, Item) (bool, error)
	Diagnostico(context.Context, string, string, Configuracao) ([]materiaInfo, map[string][]unidade, error)
	Substituir(context.Context, string, string, Entrada, []Item, json.RawMessage) (string, error)
	ObterItemAtivo(context.Context, string, string) (itemCronograma, error)
	AlterarItem(context.Context, string, string, AlteracaoItem, *sessaoCronograma, bool) error
	CronogramaAtivo(context.Context, string) (cronogramaAtivo, error)
	CriarItem(context.Context, string, Item) error
	ItensPendentes(context.Context, string, string) ([]string, error)
	Reprogramar(context.Context, []reprogramacaoItem) error
	ReplanejarCalendario(context.Context, string, string, string, string) error
}

type itemCronograma struct {
	Situacao       string
	DuracaoMinutos int
	MateriaID      *string
	TopicoID       *string
	SubtopicoID    *string
	ConcursoID     *string
}

type sessaoCronograma struct {
	ID, AlunoID, ItemID string
	ConcursoID          *string
	MateriaID           *string
	TopicoID            *string
	SubtopicoID         *string
	Segundos            int
}

type cronogramaAtivo struct {
	ID           string
	Configuracao json.RawMessage
}

type reprogramacaoItem struct {
	ID, Data string
	Ordem    int
}
