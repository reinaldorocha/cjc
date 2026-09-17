package simulados

import (
	"context"
	"encoding/json"
)

// repositorio delimita a persistencia de simulados e configuracoes de prova.
type repositorio interface {
	Listar(context.Context, string, string) ([]Simulado, error)
	Criar(context.Context, string, string, Entrada) error
	Alterar(context.Context, string, Entrada) error
	Desativar(context.Context, string) error
	Existe(context.Context, string, string) (bool, error)
	ConcursoAtribuido(context.Context, string, *string) (bool, error)
	ObterConfiguracao(context.Context, string, string) (json.RawMessage, bool, error)
	SalvarConfiguracao(context.Context, string, string, string, string, json.RawMessage) error
}
