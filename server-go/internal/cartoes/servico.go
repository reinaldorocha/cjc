package cartoes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
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
func naoAutorizado(mensagem string) error {
	return fmt.Errorf("%w: %s", dominio.ErrNaoAutorizado, mensagem)
}

func Novo(b *sql.DB, f *time.Location) *Servico {
	return &Servico{repositorio: novoRepositorioMySQL(b, f), fuso: f}
}
func (s *Servico) ListarBaralhosMentor(ctx context.Context, a string) ([]Baralho, error) {
	lista, err := s.repositorio.ListarBaralhosMentor(ctx, a)
	if err != nil {
		return nil, err
	}
	normalizarFacilidade(lista)
	return lista, nil
}
func (s *Servico) CriarBaralhoMentor(ctx context.Context, a string, e EntradaBaralho) (string, error) {
	if err := validarEntradaBaralho(e); err != nil {
		return "", err
	}
	if e.Alcance == "aluno" {
		return "", entradaInvalida("alcance inválido")
	}
	if e.Alcance == "edital" {
		ok, err := s.repositorio.EditalDoMentor(ctx, a, *e.EditalID)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", dominio.ErrNaoEncontrado
		}
	}
	if e.BaralhoPaiID != nil && *e.BaralhoPaiID != "" {
		if err := s.validarEdicao(ctx, "", a, "mentor", *e.BaralhoPaiID); err != nil {
			return "", err
		}
	}
	e.Proprietario = "mentor"
	e.DestinoAluno = nil
	return s.repositorio.CriarBaralhoMentor(ctx, a, e)
}
func (s *Servico) ListarBaralhos(ctx context.Context, a, b, c, d string) ([]Baralho, error) {
	lista, err := s.repositorio.ListarBaralhos(ctx, a, b, c, d)
	if err != nil {
		return nil, err
	}
	editais, err := s.editaisDoAluno(ctx, a)
	if err != nil {
		return nil, err
	}
	filtrada := make([]Baralho, 0, len(lista))
	for i := range lista {
		if !baralhoVisivel(a, lista[i].AlunoID, lista[i].EditalID, lista[i].Alcance, lista[i].TipoProprietario, editais) {
			continue
		}
		if c == "aluno" {
			lista[i].SomenteLeitura = !(lista[i].TipoProprietario == "aluno" && lista[i].CriadoPor == b)
		} else {
			lista[i].SomenteLeitura = !(lista[i].CriadoPor == b || (lista[i].TipoProprietario == "aluno" && lista[i].AlunoID != nil && *lista[i].AlunoID == a))
		}
		filtrada = append(filtrada, lista[i])
	}
	normalizarFacilidade(filtrada)
	return filtrada, nil
}
func (s *Servico) CriarBaralho(ctx context.Context, a, b, c string, e EntradaBaralho) (string, error) {
	if strings.TrimSpace(e.Nome) == "" {
		return "", entradaInvalida("nome obrigatório")
	}
	if c == "aluno" {
		e.Alcance = "aluno"
		e.EditalID = nil
		e.Proprietario = "aluno"
	} else if e.Alcance == "" {
		e.Alcance = "aluno"
		e.Proprietario = "mentor"
	}
	if e.Proprietario == "" {
		e.Proprietario = "mentor"
	}
	if e.Alcance == "aluno" {
		e.DestinoAluno = &a
	} else {
		e.DestinoAluno = nil
	}
	if e.Alcance == "global" {
		e.EditalID = nil
	}
	if err := validarEntradaBaralho(e); err != nil {
		return "", err
	}
	if e.Alcance == "edital" {
		var ok bool
		var err error
		if a == "" {
			ok, err = s.repositorio.EditalDoMentor(ctx, b, *e.EditalID)
		} else {
			ok, err = s.repositorio.EditalDoAluno(ctx, a, *e.EditalID)
		}
		if err != nil {
			return "", err
		}
		if !ok {
			return "", dominio.ErrNaoEncontrado
		}
	}
	if e.BaralhoPaiID != nil && *e.BaralhoPaiID != "" {
		if err := s.validarEdicao(ctx, a, b, c, *e.BaralhoPaiID); err != nil {
			return "", entradaInvalida("baralho pai inválido")
		}
	}
	return s.repositorio.CriarBaralho(ctx, a, b, c, e)
}
func (s *Servico) AlterarBaralho(ctx context.Context, a, b, c, d string, e AlteracaoBaralho) error {
	if e.Nome != nil && strings.TrimSpace(*e.Nome) == "" {
		return entradaInvalida("nome obrigatório")
	}
	if e.Alcance != nil && *e.Alcance != "aluno" && *e.Alcance != "edital" && *e.Alcance != "global" {
		return entradaInvalida("alcance inválido")
	}
	if c == "aluno" && (e.Alcance != nil || e.EditalID != nil || e.ConcursoID != nil) {
		return naoAutorizado("alteração de alcance não permitida")
	}
	if err := s.validarEdicao(ctx, a, b, c, d); err != nil {
		return err
	}
	if e.BaralhoPaiID != nil {
		if *e.BaralhoPaiID == d {
			return entradaInvalida("um baralho não pode ser pai de si mesmo")
		}
		if *e.BaralhoPaiID != "" {
			if err := s.validarEdicao(ctx, a, b, c, *e.BaralhoPaiID); err != nil {
				return entradaInvalida("baralho pai inválido")
			}
			descendente, err := s.repositorio.BaralhoDescendente(ctx, d, *e.BaralhoPaiID)
			if err != nil {
				return err
			}
			if descendente {
				return dominio.ErrConflito
			}
		}
	}
	if e.Alcance != nil && *e.Alcance == "edital" {
		if e.EditalID == nil || strings.TrimSpace(*e.EditalID) == "" {
			return entradaInvalida("edital obrigatório")
		}
		ok, err := s.repositorio.EditalDoAluno(ctx, a, *e.EditalID)
		if err != nil {
			return err
		}
		if !ok {
			return dominio.ErrNaoEncontrado
		}
	}
	if e.Alcance != nil {
		if *e.Alcance == "aluno" {
			e.DestinoAluno = &a
		} else {
			e.DestinoAluno = nil
		}
		if *e.Alcance == "global" {
			e.EditalID = nil
		}
	}
	return s.repositorio.AlterarBaralho(ctx, a, b, c, d, e)
}
func (s *Servico) DesativarBaralho(ctx context.Context, a, b, c, d string) error {
	if err := s.validarEdicao(ctx, a, b, c, d); err != nil {
		return err
	}
	return s.repositorio.DesativarBaralho(ctx, a, b, c, d)
}
func (s *Servico) CriarCartao(ctx context.Context, a, b, c, d string, e EntradaCartao) (string, error) {
	if err := validarCartaoServico(e.Tipo, e.Frente); err != nil {
		return "", err
	}
	if err := s.validarEdicao(ctx, a, b, c, d); err != nil {
		return "", err
	}
	return s.repositorio.CriarCartao(ctx, a, b, c, d, e)
}
func (s *Servico) AlterarCartao(ctx context.Context, a, b, c, d string, e AlteracaoCartao) error {
	if e.Tipo != nil {
		if err := validarTipoServico(*e.Tipo); err != nil {
			return err
		}
	}
	if e.Frente != nil && strings.TrimSpace(*e.Frente) == "" {
		return entradaInvalida("frente obrigatória")
	}
	baralho, err := s.repositorio.BaralhoDoCartao(ctx, d)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dominio.ErrNaoEncontrado
		}
		return err
	}
	if err = s.validarEdicao(ctx, a, b, c, baralho); err != nil {
		return err
	}
	return s.repositorio.AlterarCartao(ctx, a, b, c, d, e)
}
func (s *Servico) DesativarCartao(ctx context.Context, a, b, c, d string) error {
	baralho, err := s.repositorio.BaralhoDoCartao(ctx, d)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dominio.ErrNaoEncontrado
		}
		return err
	}
	if err = s.validarEdicao(ctx, a, b, c, baralho); err != nil {
		return err
	}
	return s.repositorio.DesativarCartao(ctx, a, b, c, d)
}
func (s *Servico) Revisar(ctx context.Context, a, b string, c int, d string) (Revisao, error) {
	if c < 0 || c > 3 {
		return Revisao{}, entradaInvalida("qualidade deve estar entre zero e três")
	}
	baralhoID, err := s.repositorio.BaralhoDoCartao(ctx, b)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Revisao{}, dominio.ErrNaoEncontrado
		}
		return Revisao{}, err
	}
	propriedade, err := s.repositorio.PropriedadeBaralho(ctx, baralhoID)
	if err != nil {
		return Revisao{}, err
	}
	editais, err := s.editaisDoAluno(ctx, a)
	if err != nil {
		return Revisao{}, err
	}
	if !baralhoVisivel(a, propriedade.AlunoID, propriedade.EditalID, propriedade.Alcance, propriedade.TipoProprietario, editais) {
		return Revisao{}, dominio.ErrNaoEncontrado
	}
	estado, err := s.repositorio.UltimaRevisao(ctx, a, b, d)
	if err != nil {
		return Revisao{}, err
	}
	revisao := calcularRevisao(b, c, estado, time.Now().In(s.fuso))
	if err = s.repositorio.RegistrarRevisao(ctx, a, b, d, revisao); err != nil {
		return Revisao{}, err
	}
	return revisao, nil
}

type EstadoRevisao struct {
	Repeticoes    int
	IntervaloDias int
	Facilidade    float64
}

func calcularRevisao(cartao string, qualidade int, estado EstadoRevisao, agora time.Time) Revisao {
	intervalo, facilidade := estado.IntervaloDias, estado.Facilidade
	if facilidade == 0 {
		facilidade = 2.5
	}
	switch qualidade {
	case 0:
		facilidade, intervalo = math.Max(1.3, facilidade-0.2), 0
	case 1:
		facilidade, intervalo = math.Max(1.3, facilidade-0.15), maximoServico(1, int(math.Round(float64(intervalo)*1.2)))
	case 3:
		facilidade += 0.1
		if intervalo < 1 {
			intervalo = 4
		} else {
			bom := intervaloBomServico(intervalo, estado.Facilidade)
			intervalo = maximoServico(bom+1, maximoServico(intervalo+1, int(math.Round(float64(intervalo)*facilidade*1.3))))
		}
	default:
		intervalo = intervaloBomServico(intervalo, facilidade)
	}
	return Revisao{CartaoID: cartao, Qualidade: qualidade, Repeticoes: estado.Repeticoes + 1, IntervaloDias: intervalo, Facilidade: facilidade, ProximaRevisao: agora.AddDate(0, 0, intervalo).Format("2006-01-02")}
}

func validarEntradaBaralho(e EntradaBaralho) error {
	if strings.TrimSpace(e.Nome) == "" {
		return entradaInvalida("nome obrigatório")
	}
	if e.Alcance != "aluno" && e.Alcance != "edital" && e.Alcance != "global" {
		return entradaInvalida("alcance inválido")
	}
	if e.Alcance == "edital" && (e.EditalID == nil || strings.TrimSpace(*e.EditalID) == "") {
		return entradaInvalida("edital obrigatório")
	}
	return nil
}
func validarCartaoServico(tipo, frente string) error {
	if strings.TrimSpace(frente) == "" {
		return entradaInvalida("frente obrigatória")
	}
	return validarTipoServico(tipo)
}
func validarTipoServico(tipo string) error {
	if tipo != "basico" && tipo != "lacuna" && tipo != "multipla_escolha" && tipo != "certo_errado" {
		return entradaInvalida("tipo inválido")
	}
	return nil
}
func intervaloBomServico(anterior int, facilidade float64) int {
	if anterior < 1 {
		return 1
	}
	if anterior < 6 {
		return 6
	}
	return maximoServico(1, int(math.Round(float64(anterior)*facilidade)))
}
func maximoServico(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *Servico) validarEdicao(ctx context.Context, aluno, executor, papel, baralho string) error {
	propriedade, err := s.repositorio.PropriedadeBaralho(ctx, baralho)
	if err != nil {
		return err
	}
	if propriedade.CriadoPor == "" {
		return dominio.ErrNaoEncontrado
	}
	ok := false
	if papel == "aluno" {
		ok = propriedade.AlunoID != nil && *propriedade.AlunoID == aluno && propriedade.TipoProprietario == "aluno" && propriedade.CriadoPor == executor
	} else {
		ok = propriedade.CriadoPor == executor || (propriedade.TipoProprietario == "aluno" && propriedade.AlunoID != nil && *propriedade.AlunoID == aluno)
	}
	if !ok {
		return dominio.ErrNaoAutorizado
	}
	return nil
}
func (s *Servico) ListarPendentes(ctx context.Context, a string) ([]CartaoPendente, error) {
	lista, err := s.repositorio.ListarPendentes(ctx, a)
	if err != nil {
		return nil, err
	}
	editais, err := s.editaisDoAluno(ctx, a)
	if err != nil {
		return nil, err
	}
	filtrada := make([]CartaoPendente, 0, len(lista))
	for _, cartao := range lista {
		if !baralhoVisivel(a, cartao.AlunoID, cartao.EditalID, cartao.Alcance, cartao.TipoProprietario, editais) {
			continue
		}
		if cartao.TipoProprietario == "aluno" {
			cartao.Origem = "pessoal"
		} else {
			cartao.Origem = "mentor"
		}
		if cartao.Facilidade == 0 {
			cartao.Facilidade = 2.5
		}
		filtrada = append(filtrada, cartao)
	}
	return filtrada, nil
}
func (s *Servico) ListarHistorico(ctx context.Context, a, b string) ([]HistoricoRevisao, error) {
	lista, err := s.repositorio.ListarHistorico(ctx, a, b)
	if err != nil {
		return nil, err
	}
	editais, err := s.editaisDoAluno(ctx, a)
	if err != nil {
		return nil, err
	}
	filtrada := make([]HistoricoRevisao, 0, len(lista))
	for _, revisao := range lista {
		if !baralhoVisivel(a, revisao.AlunoID, revisao.EditalID, revisao.Alcance, revisao.TipoProprietario, editais) {
			continue
		}
		if b != "" && revisao.Alcance != "global" && (revisao.ConcursoID == nil || *revisao.ConcursoID != b) && (revisao.EditalConcursoID == nil || *revisao.EditalConcursoID != b) {
			continue
		}
		filtrada = append(filtrada, revisao)
	}
	return filtrada, nil
}

func (s *Servico) editaisDoAluno(ctx context.Context, aluno string) (map[string]bool, error) {
	ids, err := s.repositorio.EditaisDoAluno(ctx, aluno)
	if err != nil {
		return nil, err
	}
	resultado := make(map[string]bool, len(ids))
	for _, id := range ids {
		resultado[id] = true
	}
	return resultado, nil
}

func baralhoVisivel(aluno string, destino, edital *string, alcance, proprietario string, editais map[string]bool) bool {
	switch alcance {
	case "aluno":
		return destino != nil && *destino == aluno
	case "global":
		return proprietario == "mentor"
	case "edital":
		return edital != nil && editais[*edital]
	default:
		return false
	}
}

func normalizarFacilidade(baralhos []Baralho) {
	for i := range baralhos {
		for j := range baralhos[i].Cartoes {
			if baralhos[i].Cartoes[j].Facilidade == 0 {
				baralhos[i].Cartoes[j].Facilidade = 2.5
			}
		}
	}
}
