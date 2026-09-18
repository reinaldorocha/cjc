package editais

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"chega-junto-concurseiro-web/internal/dominio"
	"chega-junto-concurseiro-web/internal/identificador"
)

type Servico struct{ repositorio repositorio }

func entradaInvalida(mensagem string) error {
	return fmt.Errorf("%w: %s", dominio.ErrEntradaInvalida, mensagem)
}

func Novo(b *sql.DB) *Servico { return &Servico{repositorio: novoRepositorioMySQL(b)} }
func (s *Servico) Catalogar(c context.Context, a, b string) ([]ResumoCatalogo, error) {
	return s.repositorio.Catalogar(c, a, b)
}
func (s *Servico) ListarCatalogo(c context.Context, a, b string) ([]Edital, error) {
	return s.repositorio.ListarCatalogo(c, a, b)
}
func (s *Servico) CriarCatalogo(c context.Context, a string, e Entrada) (string, error) {
	if strings.TrimSpace(e.Nome) == "" || e.ConcursoID == "" {
		return "", entradaInvalida("nome e concurso obrigatórios")
	}
	if err := validarArvore(e.Materias); err != nil {
		return "", err
	}
	e.Nome = strings.TrimSpace(e.Nome)
	normalizarArvore(e.Materias)
	ok, err := s.repositorio.ConcursoDoMentor(c, a, e.ConcursoID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", dominio.ErrNaoEncontrado
	}
	prepararArvore(&e)
	return s.repositorio.CriarCatalogo(c, a, e)
}
func (s *Servico) Listar(c context.Context, a, b string) ([]Edital, error) {
	return s.repositorio.Listar(c, a, b)
}
func (s *Servico) Atribuir(c context.Context, a, b string, e Entrada) (string, error) {
	if e.ConcursoID == "" {
		return "", entradaInvalida("concurso obrigatório")
	}
	if e.EditalID == "" && len(e.Materias) == 0 {
		return "", entradaInvalida("selecione um edital do catálogo ou informe as matérias")
	}
	if e.EditalID == "" && strings.TrimSpace(e.Nome) == "" {
		return "", entradaInvalida("nome obrigatório")
	}
	if err := validarArvore(e.Materias); err != nil {
		return "", err
	}
	e.Nome = strings.TrimSpace(e.Nome)
	e.AtualizarNome = e.EditalID != "" && e.Nome != ""
	normalizarArvore(e.Materias)
	ok, err := s.repositorio.ConcursoDoAluno(c, a, e.ConcursoID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", dominio.ErrNaoEncontrado
	}
	if e.EditalID != "" {
		ok, err = s.repositorio.EditalDoMentor(c, b, e.EditalID, e.ConcursoID)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", dominio.ErrNaoEncontrado
		}
	} else {
		prepararArvore(&e)
	}
	return s.repositorio.Atribuir(c, a, b, e)
}
func (s *Servico) CriarItem(c context.Context, a string, e NovoItem) (string, error) {
	if strings.TrimSpace(e.Nome) == "" || e.PaiID == "" {
		return "", entradaInvalida("nome e pai obrigatórios")
	}
	if e.Tipo != "materia" && e.Tipo != "topico" && e.Tipo != "subtopico" {
		return "", entradaInvalida("tipo inválido")
	}
	ok, err := s.repositorio.PaiDoAluno(c, a, e.Tipo, e.PaiID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", dominio.ErrNaoEncontrado
	}
	e.ID = identificador.UUID()
	e.Nome = strings.TrimSpace(e.Nome)
	return s.repositorio.CriarItem(c, a, e)
}
func (s *Servico) AlterarItem(c context.Context, a, b string, e AlteracaoItem) error {
	if e.Tipo != "materia" && e.Tipo != "topico" && e.Tipo != "subtopico" {
		return entradaInvalida("tipo inválido")
	}
	if e.Nome != nil && strings.TrimSpace(*e.Nome) == "" {
		return entradaInvalida("nome obrigatório")
	}
	if e.Nome != nil {
		v := strings.TrimSpace(*e.Nome)
		e.Nome = &v
	}
	ok, err := s.repositorio.ItemDoAluno(c, a, e.Tipo, b)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return s.repositorio.AlterarItem(c, a, b, e)
}
func (s *Servico) Reordenar(c context.Context, a string, e []ItemOrdem) error {
	if len(e) == 0 {
		return entradaInvalida("itens obrigatórios")
	}
	for _, item := range e {
		if item.ID == "" || item.Ordem < 0 || (item.Tipo != "materia" && item.Tipo != "topico" && item.Tipo != "subtopico") {
			return entradaInvalida("item inválido")
		}
	}
	for _, item := range e {
		ok, err := s.repositorio.ItemDoAluno(c, a, item.Tipo, item.ID)
		if err != nil {
			return err
		}
		if !ok {
			return dominio.ErrNaoEncontrado
		}
	}
	return s.repositorio.Reordenar(c, a, e)
}
func (s *Servico) AtualizarProgressoMaterial(c context.Context, a, b string, e EntradaProgressoMaterial) error {
	if strings.TrimSpace(b) == "" || e.ItemID == "" {
		return entradaInvalida("material e item obrigatórios")
	}
	if e.Tipo != "materia" && e.Tipo != "topico" && e.Tipo != "subtopico" {
		return entradaInvalida("tipo inválido")
	}
	ok, err := s.repositorio.ItemDoAluno(c, a, e.Tipo, e.ItemID)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return s.repositorio.AtualizarProgressoMaterial(c, a, b, e)
}
func (s *Servico) AtualizarProgresso(c context.Context, a, b string, e Progresso) error {
	if b == "" {
		return entradaInvalida("item obrigatório")
	}
	if e.Tipo == "" {
		tipo, err := s.repositorio.TipoDoItem(c, b)
		if errors.Is(err, sql.ErrNoRows) {
			return dominio.ErrNaoEncontrado
		}
		if err != nil {
			return err
		}
		e.Tipo = tipo
	}
	if e.Tipo != "topico" && e.Tipo != "subtopico" {
		return entradaInvalida("tipo inválido")
	}
	ok, err := s.repositorio.ItemDoAluno(c, a, e.Tipo, b)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	if err := s.repositorio.AtualizarProgresso(c, a, b, e); err != nil {
		return err
	}
	if !e.Estudado {
		return s.repositorio.SincronizarRevisoes(c, a, e.Tipo, b, nil, nil)
	}
	contexto, err := s.repositorio.ContextoRevisao(c, e.Tipo, b)
	if err != nil {
		return err
	}
	agendamentos := []RevisaoAutomatica{}
	for i, prazo := range strings.Split(contexto.Prazos, ",") {
		dias, err := strconv.Atoi(strings.TrimSpace(prazo))
		if err != nil || dias <= 0 {
			continue
		}
		agendamentos = append(agendamentos, RevisaoAutomatica{
			Ciclo:       i + 1,
			ProximaData: time.Now().AddDate(0, 0, dias).Format("2006-01-02"),
		})
	}
	return s.repositorio.SincronizarRevisoes(c, a, e.Tipo, b, &contexto, agendamentos)
}

func prepararArvore(e *Entrada) {
	e.Novo = true
	e.EditalID = identificador.UUID()
	for mi := range e.Materias {
		e.Materias[mi].ID = identificador.UUID()
		if e.Materias[mi].Ordem == 0 {
			e.Materias[mi].Ordem = mi + 1
		}
		for ti := range e.Materias[mi].Topicos {
			e.Materias[mi].Topicos[ti].ID = identificador.UUID()
			if e.Materias[mi].Topicos[ti].Ordem == 0 {
				e.Materias[mi].Topicos[ti].Ordem = ti + 1
			}
			for si := range e.Materias[mi].Topicos[ti].Subtopicos {
				e.Materias[mi].Topicos[ti].Subtopicos[si].ID = identificador.UUID()
				if e.Materias[mi].Topicos[ti].Subtopicos[si].Ordem == 0 {
					e.Materias[mi].Topicos[ti].Subtopicos[si].Ordem = si + 1
				}
			}
		}
	}
}

func normalizarArvore(materias []Materia) {
	for mi := range materias {
		materias[mi].Nome = strings.TrimSpace(materias[mi].Nome)
		for ti := range materias[mi].Topicos {
			materias[mi].Topicos[ti].Nome = strings.TrimSpace(materias[mi].Topicos[ti].Nome)
			for si := range materias[mi].Topicos[ti].Subtopicos {
				materias[mi].Topicos[ti].Subtopicos[si].Nome = strings.TrimSpace(materias[mi].Topicos[ti].Subtopicos[si].Nome)
			}
		}
	}
}

func validarArvore(materias []Materia) error {
	for _, materia := range materias {
		if strings.TrimSpace(materia.Nome) == "" {
			return entradaInvalida("nome da matéria obrigatório")
		}
		for _, topico := range materia.Topicos {
			if strings.TrimSpace(topico.Nome) == "" {
				return entradaInvalida("nome do tópico obrigatório")
			}
			for _, subtopico := range topico.Subtopicos {
				if strings.TrimSpace(subtopico.Nome) == "" {
					return entradaInvalida("nome do subtópico obrigatório")
				}
			}
		}
	}
	return nil
}
