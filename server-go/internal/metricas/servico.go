package metricas

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"track-concursos-web/internal/dominio"
)

type Servico struct {
	repositorio repositorio
	fuso        *time.Location
}

func entradaInvalida(mensagem string) error {
	return fmt.Errorf("%w: %s", dominio.ErrEntradaInvalida, mensagem)
}

func Novo(banco *sql.DB, fuso *time.Location) *Servico {
	return &Servico{repositorio: novoRepositorioMySQL(banco, fuso), fuso: fuso}
}
func (s *Servico) Resumo(ctx context.Context, aluno, concurso, inicio, fim string) (Resumo, error) {
	if err := validarPeriodo(concurso, inicio, fim, s.repositorio); err != nil {
		return Resumo{}, err
	}
	if err := s.validarConcurso(ctx, aluno, concurso); err != nil {
		return Resumo{}, err
	}
	periodo, err := criarPeriodo(s.fuso, inicio, fim, false)
	if err != nil {
		return Resumo{}, err
	}
	resultado, err := s.repositorio.Resumo(ctx, aluno, concurso, periodo)
	if err != nil {
		return Resumo{}, err
	}
	if resultado.QuestoesResolvidas > 0 {
		v := arredondar(float64(resultado.Acertos) / float64(resultado.QuestoesResolvidas) * 100)
		resultado.PercentualAcertos = &v
	}
	if resultado.ItensEditalTotal > 0 {
		resultado.PercentualEdital = arredondar(float64(resultado.ItensEditalConcluidos) / float64(resultado.ItensEditalTotal) * 100)
	}
	arredondarPonteiro(resultado.MediaSimulados)
	arredondarPonteiro(resultado.MelhorSimulado)
	if len(resultado.UltimosSimulados) > 0 {
		resultado.UltimoSimulado = &resultado.UltimosSimulados[0]
		arredondarPonteiro(resultado.UltimoSimulado)
	}
	if len(resultado.UltimosSimulados) > 1 {
		v := arredondar(resultado.UltimosSimulados[0] - resultado.UltimosSimulados[1])
		resultado.TendenciaSimulados = &v
	}
	datas, err := s.repositorio.DatasEstudo(ctx, aluno, concurso, periodo.Offset)
	if err != nil {
		return Resumo{}, err
	}
	resultado.SequenciaAtual = calcularSequencia(datas, time.Now().In(s.fuso))
	return resultado, nil
}
func (s *Servico) LinhaDoTempo(ctx context.Context, aluno, concurso, inicio, fim string) ([]Dia, error) {
	if concurso == "" {
		return nil, entradaInvalida("concurso obrigatório")
	}
	if err := s.validarConcurso(ctx, aluno, concurso); err != nil {
		return nil, err
	}
	periodo, err := criarPeriodo(s.fuso, inicio, fim, true)
	if err != nil {
		return nil, err
	}
	lista, err := s.repositorio.LinhaDoTempo(ctx, aluno, concurso, periodo)
	if err != nil {
		return nil, err
	}
	porData := make(map[string]Dia, len(lista))
	for _, dia := range lista {
		porData[dia.Data] = dia
	}
	ordenada := []Dia{}
	for data := *periodo.Inicio; data.Before(*periodo.Fim); data = data.AddDate(0, 0, 1) {
		chave := data.In(s.fuso).Format("2006-01-02")
		dia := porData[chave]
		dia.Data = chave
		dia.DiasAtivo = dia.Segundos > 0 || dia.Questoes > 0
		ordenada = append(ordenada, dia)
	}
	return ordenada, nil
}
func (s *Servico) Materias(ctx context.Context, aluno, concurso, inicio, fim string) ([]Materia, error) {
	if err := validarPeriodo(concurso, inicio, fim, s.repositorio); err != nil {
		return nil, err
	}
	if err := s.validarConcurso(ctx, aluno, concurso); err != nil {
		return nil, err
	}
	periodo, err := criarPeriodo(s.fuso, inicio, fim, false)
	if err != nil {
		return nil, err
	}
	lista, err := s.repositorio.Materias(ctx, aluno, concurso, periodo)
	if err != nil {
		return nil, err
	}
	for i := range lista {
		if lista[i].Questoes > 0 {
			v := arredondar(float64(lista[i].Acertos) / float64(lista[i].Questoes) * 100)
			lista[i].PercentualAcertos = &v
		}
		if lista[i].ItensTotal > 0 {
			lista[i].PercentualEdital = arredondar(float64(lista[i].ItensConcluidos) / float64(lista[i].ItensTotal) * 100)
		}
		arredondarPonteiro(lista[i].MediaSimulados)
		arredondarPonteiro(lista[i].MelhorSimulado)
		arredondarPonteiro(lista[i].UltimoSimulado)
	}
	return lista, nil
}

func (s *Servico) validarConcurso(ctx context.Context, aluno, concurso string) error {
	ok, err := s.repositorio.ConcursoAtribuido(ctx, aluno, concurso)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return nil
}

// validarPeriodo mantém no serviço as regras de entrada; o repositório recebe
// somente filtros já validados e executa as leituras necessárias.
func validarPeriodo(concurso, inicio, fim string, _ repositorio) error {
	if concurso == "" {
		return entradaInvalida("concurso obrigatório")
	}
	if inicio != "" {
		if _, err := time.Parse("2006-01-02", inicio); err != nil {
			return entradaInvalida("início inválido")
		}
	}
	if fim != "" {
		if _, err := time.Parse("2006-01-02", fim); err != nil {
			return entradaInvalida("fim inválido")
		}
	}
	if inicio != "" && fim != "" && inicio > fim {
		return entradaInvalida("período inválido")
	}
	return nil
}

func calcularSequencia(datas []string, agora time.Time) int {
	ativas := make(map[string]bool, len(datas))
	for _, data := range datas {
		ativas[data] = true
	}
	cursor := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	if !ativas[cursor.Format("2006-01-02")] {
		cursor = cursor.AddDate(0, 0, -1)
	}
	total := 0
	for ativas[cursor.Format("2006-01-02")] {
		total++
		cursor = cursor.AddDate(0, 0, -1)
	}
	return total
}

func arredondar(v float64) float64 { return float64(int(v*10+0.5)) / 10 }
func arredondarPonteiro(v *float64) {
	if v != nil {
		*v = arredondar(*v)
	}
}

func criarPeriodo(fuso *time.Location, inicio, fim string, obrigatorio bool) (Periodo, error) {
	agora := time.Now().In(fuso)
	hoje := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, fuso)
	resultado := Periodo{InicioHoje: hoje.UTC(), FimHoje: hoje.AddDate(0, 0, 1).UTC(), Offset: agora.Format("-07:00")}
	if inicio == "" && fim == "" && !obrigatorio {
		return resultado, nil
	}
	if fim == "" {
		fim = agora.Format("2006-01-02")
	}
	if inicio == "" {
		ultimo, err := time.ParseInLocation("2006-01-02", fim, fuso)
		if err != nil {
			return Periodo{}, entradaInvalida("fim inválido")
		}
		inicio = ultimo.AddDate(0, 0, -29).Format("2006-01-02")
	}
	de, err := time.ParseInLocation("2006-01-02", inicio, fuso)
	if err != nil {
		return Periodo{}, entradaInvalida("início inválido")
	}
	ultimo, err := time.ParseInLocation("2006-01-02", fim, fuso)
	if err != nil {
		return Periodo{}, entradaInvalida("fim inválido")
	}
	ate := ultimo.AddDate(0, 0, 1)
	if !de.Before(ate) {
		return Periodo{}, entradaInvalida("período inválido")
	}
	deUTC, ateUTC := de.UTC(), ate.UTC()
	resultado.Inicio, resultado.Fim = &deUTC, &ateUTC
	return resultado, nil
}
