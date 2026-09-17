package metricas

import "context"
import "time"

type repositorio interface {
	ConcursoAtribuido(context.Context, string, string) (bool, error)
	DatasEstudo(context.Context, string, string, string) ([]string, error)
	Resumo(context.Context, string, string, Periodo) (Resumo, error)
	LinhaDoTempo(context.Context, string, string, Periodo) ([]Dia, error)
	Materias(context.Context, string, string, Periodo) ([]Materia, error)
}

type Periodo struct {
	Inicio     *time.Time
	Fim        *time.Time
	InicioHoje time.Time
	FimHoje    time.Time
	Offset     string
}
