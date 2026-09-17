package mentoria

import (
	"context"
	"time"
)

type repositorio interface {
	ResolverAluno(context.Context, string, string) (bool, error)
	DefinirMentor(context.Context, string, string, string) error
	ListarAlunos(context.Context, string) ([]Aluno, error)
	ObterAluno(context.Context, string, string) (Aluno, error)
	VisaoGeral(context.Context, string) (map[string]int64, error)
	Configurar(context.Context, string, string, *bool, bool, *time.Time, *string, *string, *string) error
	Radar(context.Context, string) ([]dadosRadar, error)
}

type dadosRadar struct {
	AlunoID            string
	AlunoNome          string
	AlunoEmail         string
	Telefone           *string
	ConcursoNome       string
	UltimoEstudoEm     *time.Time
	SegundosSemana     int64
	SegundosMes        int64
	QuestoesResolvidas int64
	QuestoesAcertos    int64
	RevisoesPendentes  int64
}
