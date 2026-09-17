package whitelabel

import "context"

type repositorio interface {
	ObterPorMentor(context.Context, string) (Config, bool, error)
	MentorDoAluno(context.Context, string) (string, bool, error)
	Salvar(context.Context, Config) error
}
