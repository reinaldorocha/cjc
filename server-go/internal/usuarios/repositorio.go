package usuarios

import (
	"context"
	"time"
)

type repositorio interface {
	Listar(context.Context) ([]Resumo, error)
	Criar(context.Context, registroCriacao) error
	CriarAlunoVinculado(context.Context, registroCriacao, string, bool, *time.Time) error
	MentorAtivo(context.Context, string) (bool, error)
	Alterar(context.Context, string, *string, *bool, string) error
	Reativar(context.Context, string) error
	RedefinirSenha(context.Context, string, string) error
	AlterarNomeProprio(context.Context, string, string) error
}

type registroCriacao struct {
	ID        string
	Nome      string
	Email     string
	SenhaHash string
	Papel     string
	Telefone  string
}
