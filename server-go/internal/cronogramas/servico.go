package cronogramas

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
	"track-concursos-web/internal/dominio"
	"track-concursos-web/internal/identificador"
)

type Servico struct {
	repositorio repositorio
	fuso        *time.Location
}

type Configuracao struct {
	Tipo                 string             `json:"tipo"`
	CicloModo            string             `json:"cicloModo"`
	MateriasSelecionadas map[string]bool    `json:"materiasSelecionadas"`
	MateriaAfinidade     map[string]float64 `json:"materiaAfinidade"`
	MateriaPrioridades   map[string]string  `json:"materiaPrioridades"`
	Ritmo                string             `json:"ritmo"`
	Horas                map[string]float64 `json:"horas"`
	MinutosTopico        int                `json:"minutosTopico"`
	MaxTopicosDia        int                `json:"maxTopicosDia"`
	AlertaMetaHabilitado bool               `json:"alertaMetaHabilitado"`
	MetaAcertos          int                `json:"metaAcertos"`
	RepeticoesEdital     int                `json:"repeticoesEdital"`
}
type Item struct {
	ID             string  `json:"id,omitempty"`
	DataPlanejada  *string `json:"dataPlanejada,omitempty"`
	PosicaoCiclo   *int    `json:"posicaoCiclo,omitempty"`
	MateriaID      *string `json:"materiaId,omitempty"`
	TopicoID       *string `json:"topicoId,omitempty"`
	SubtopicoID    *string `json:"subtopicoId,omitempty"`
	Nome           string  `json:"nome,omitempty"`
	DuracaoMinutos int     `json:"duracaoMinutos"`
	Prioridade     int     `json:"prioridade"`
	Ordem          int     `json:"ordem"`
	Situacao       string  `json:"situacao"`
	ConcluidoEm    *string `json:"concluidoEm,omitempty"`
}
type Cronograma struct {
	ID           string       `json:"id"`
	AlunoID      string       `json:"alunoId"`
	ConcursoID   *string      `json:"concursoId,omitempty"`
	Tipo         string       `json:"tipo"`
	CriadoPor    string       `json:"criadoPor"`
	Versao       int          `json:"versao"`
	Estado       string       `json:"estado"`
	Configuracao Configuracao `json:"configuracao"`
	Itens        []Item       `json:"itens"`
	CriadoEm     string       `json:"criadoEm"`
}
type Entrada struct {
	ConcursoID   *string      `json:"concursoId"`
	Tipo         string       `json:"tipo"`
	Configuracao Configuracao `json:"configuracao"`
	Itens        []Item       `json:"itens"`
}
type AlteracaoItem struct {
	DataPlanejada  *string `json:"dataPlanejada"`
	PosicaoCiclo   *int    `json:"posicaoCiclo"`
	DuracaoMinutos *int    `json:"duracaoMinutos"`
	Prioridade     *int    `json:"prioridade"`
	Ordem          *int    `json:"ordem"`
	Situacao       *string `json:"situacao"`
	CriarSessao    *bool   `json:"criarSessao"`
}
type unidade struct {
	materiaID, topicoID, subtopicoID, nome string
	prioridade                             int
}
type materiaInfo struct {
	id, nome                                           string
	total, concluidos, resolvidas, acertos, prioridade int
}

func Novo(banco *sql.DB, fuso *time.Location) *Servico {
	return &Servico{repositorio: novoRepositorioMySQL(banco), fuso: fuso}
}
func (s *Servico) ObterAtivo(ctx context.Context, alunoID, concursoID string) (*Cronograma, error) {
	return s.repositorio.ObterAtivo(ctx, alunoID, concursoID)
}
func (s *Servico) Calendario(ctx context.Context, alunoID, concursoID, inicio, fim string) ([]Item, error) {
	return s.repositorio.Calendario(ctx, alunoID, concursoID, inicio, fim)
}

func (s *Servico) Salvar(ctx context.Context, alunoID, executor string, e Entrada) (string, error) {
	if e.Tipo != "manual" && e.Tipo != "agendado" && e.Tipo != "ciclo_inteligente" {
		return "", errors.New("tipo inválido")
	}
	return s.substituir(ctx, alunoID, executor, e, e.Itens)
}
func (s *Servico) Gerar(ctx context.Context, alunoID, executor string, e Entrada) (string, error) {
	if e.ConcursoID == nil || *e.ConcursoID == "" {
		return "", errors.New("concurso obrigatório")
	}
	ok, err := s.repositorio.ConcursoAtribuido(ctx, alunoID, *e.ConcursoID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", dominio.ErrNaoAutorizado
	}
	materias, unidades, err := s.repositorio.Diagnostico(ctx, alunoID, *e.ConcursoID, e.Configuracao)
	if err != nil {
		return "", err
	}
	if e.Configuracao.Tipo == "ciclo" || e.Tipo == "ciclo_inteligente" {
		e.Tipo, e.Configuracao.Tipo, e.Itens = "ciclo_inteligente", "ciclo", gerarCiclo(materias)
	} else {
		e.Tipo, e.Configuracao.Tipo, e.Itens = "agendado", "agendado", gerarAgenda(materias, unidades, e.Configuracao, s.fuso)
	}
	return s.substituir(ctx, alunoID, executor, e, e.Itens)
}
func (s *Servico) AlterarItem(ctx context.Context, alunoID, itemID string, e AlteracaoItem, podePlanejar bool) error {
	atual, err := s.repositorio.ObterItemAtivo(ctx, alunoID, itemID)
	if err != nil {
		return err
	}
	if !podePlanejar && (e.PosicaoCiclo != nil || e.DuracaoMinutos != nil || e.Prioridade != nil || e.Ordem != nil) {
		return dominio.ErrNaoAutorizado
	}
	if err := validarAlteracao(e); err != nil {
		return err
	}
	var sessao *sessaoCronograma
	deveCriarSessao := e.CriarSessao == nil || *e.CriarSessao
	if deveCriarSessao && e.Situacao != nil && *e.Situacao == "concluido" && atual.Situacao != "concluido" {
		duracao := atual.DuracaoMinutos
		if e.DuracaoMinutos != nil {
			duracao = *e.DuracaoMinutos
		}
		dados, err := json.Marshal(map[string]any{"cronograma_item_id": itemID})
		if err != nil {
			return err
		}
		_ = dados
		sessao = &sessaoCronograma{ID: identificador.UUID(), AlunoID: alunoID, ItemID: itemID, ConcursoID: atual.ConcursoID, MateriaID: atual.MateriaID, TopicoID: atual.TopicoID, SubtopicoID: atual.SubtopicoID, Segundos: duracao * 60}
	}
	return s.repositorio.AlterarItem(ctx, alunoID, itemID, e, sessao, e.Situacao != nil && *e.Situacao != "concluido" && atual.Situacao == "concluido")
}
func (s *Servico) CriarItemNoCronograma(ctx context.Context, alunoID string, item Item) error {
	return s.repositorio.CriarItem(ctx, alunoID, item)
}
func (s *Servico) ReprogramarPendentes(ctx context.Context, alunoID, aPartirDe string, concurso ...string) error {
	fuso := s.fuso
	if fuso == nil {
		fuso = time.Local
	}
	hoje := time.Now().In(fuso).Format("2006-01-02")
	if aPartirDe == "" {
		aPartirDe = hoje
	}
	if _, err := time.Parse("2006-01-02", aPartirDe); err != nil {
		return errors.New("data de replanejamento inválida")
	}
	if aPartirDe < hoje {
		aPartirDe = hoje
	}
	id := ""
	if len(concurso) > 0 {
		id = concurso[0]
	}
	return s.repositorio.ReplanejarCalendario(ctx, alunoID, id, aPartirDe, hoje)
}
func (s *Servico) substituir(ctx context.Context, alunoID, executor string, e Entrada, itens []Item) (string, error) {
	if e.ConcursoID == nil || *e.ConcursoID == "" {
		return "", errors.New("concurso obrigatorio")
	}
	for _, item := range itens {
		if err := s.validarItem(ctx, alunoID, *e.ConcursoID, item); err != nil {
			return "", err
		}
	}
	cfg, err := json.Marshal(e.Configuracao)
	if err != nil {
		return "", err
	}
	return s.repositorio.Substituir(ctx, alunoID, executor, e, itens, cfg)
}
func (s *Servico) validarItem(ctx context.Context, alunoID, concursoID string, item Item) error {
	if item.MateriaID == nil || *item.MateriaID == "" || item.DuracaoMinutos <= 0 || item.Prioridade < 0 || item.Prioridade > 100 {
		return errors.New("item invalido")
	}
	if item.Situacao != "" && item.Situacao != "pendente" && item.Situacao != "concluido" && item.Situacao != "ignorado" {
		return errors.New("situacao invalida")
	}
	if item.DataPlanejada != nil && *item.DataPlanejada != "" {
		if _, err := time.Parse("2006-01-02", *item.DataPlanejada); err != nil {
			return errors.New("data invalida")
		}
	}
	if item.SubtopicoID != nil && *item.SubtopicoID != "" && (item.TopicoID == nil || *item.TopicoID == "") {
		return errors.New("topico obrigatorio")
	}
	ok, err := s.repositorio.ConteudoAtribuido(ctx, alunoID, concursoID, item)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoAutorizado
	}
	return nil
}
func validarAlteracao(e AlteracaoItem) error {
	if e.DataPlanejada != nil && *e.DataPlanejada != "" {
		if _, err := time.Parse("2006-01-02", *e.DataPlanejada); err != nil {
			return errors.New("data invalida")
		}
	}
	if e.DuracaoMinutos != nil && *e.DuracaoMinutos <= 0 {
		return errors.New("duracao invalida")
	}
	if e.Prioridade != nil && (*e.Prioridade < 0 || *e.Prioridade > 100) {
		return errors.New("prioridade invalida")
	}
	if e.Situacao != nil && *e.Situacao != "pendente" && *e.Situacao != "concluido" && *e.Situacao != "ignorado" {
		return errors.New("situacao invalida")
	}
	return nil
}
func (s *Servico) programarPendentes(ids []string, cfg Configuracao, inicio string) ([]reprogramacaoItem, error) {
	fuso := s.fuso
	if fuso == nil {
		fuso = time.Local
	}
	data, err := time.ParseInLocation("2006-01-02", inicio, fuso)
	if err != nil {
		data = time.Now().In(fuso)
	}
	minutos := cfg.MinutosTopico
	if minutos <= 0 {
		minutos = 60
	}
	horas := cfg.Horas
	if horas == nil {
		horas = map[string]float64{"seg": 2, "ter": 2, "qua": 2, "qui": 2, "sex": 2}
	}
	chaves := []string{"dom", "seg", "ter", "qua", "qui", "sex", "sab"}
	saida := []reprogramacaoItem{}
	indice := 0
	for dia := 0; dia < 365 && indice < len(ids); dia++ {
		atual := data.AddDate(0, 0, dia)
		capacidade := int(horas[chaves[int(atual.Weekday())]]*60) / minutos
		if horas[chaves[int(atual.Weekday())]] > 0 && capacidade < 1 {
			capacidade = 1
		}
		if cfg.MaxTopicosDia > 0 && capacidade > cfg.MaxTopicosDia {
			capacidade = cfg.MaxTopicosDia
		}
		for ordem := 1; ordem <= capacidade && indice < len(ids); ordem++ {
			saida = append(saida, reprogramacaoItem{ID: ids[indice], Data: atual.Format("2006-01-02"), Ordem: ordem})
			indice++
		}
	}
	return saida, nil
}
