package materiaisapoio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"track-concursos-web/internal/arquivos"
	"track-concursos-web/internal/dominio"
	"track-concursos-web/internal/identificador"
)

type Servico struct{ repositorio repositorio }

func entradaInvalida(mensagem string) error {
	return fmt.Errorf("%w: %s", dominio.ErrEntradaInvalida, mensagem)
}
func Novo(b *sql.DB) *Servico { return &Servico{repositorio: novoRepositorioMySQL(b)} }

func (s *Servico) Criar(ctx context.Context, mentor string, e EntradaMaterial) (string, error) {
	if err := validarMaterial(e, false); err != nil {
		return "", err
	}
	if err := s.validarEdital(ctx, mentor, e.Escopo, e.EditalID); err != nil {
		return "", err
	}
	normalizarMaterial(&e)
	id := identificador.UUID()
	if err := s.repositorio.Criar(ctx, mentor, id, e); err != nil {
		return "", err
	}
	return id, nil
}
func (s *Servico) Alterar(ctx context.Context, mentor, material string, e EntradaMaterial) error {
	if strings.TrimSpace(material) == "" {
		return entradaInvalida("material obrigatório")
	}
	if err := s.validarMaterialDoMentor(ctx, mentor, material); err != nil {
		return err
	}
	if err := validarMaterial(e, true); err != nil {
		return err
	}
	if err := s.validarEdital(ctx, mentor, e.Escopo, e.EditalID); err != nil {
		return err
	}
	normalizarMaterial(&e)
	ok, err := s.repositorio.Alterar(ctx, mentor, material, e)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return nil
}
func (s *Servico) Desativar(ctx context.Context, mentor, material string) error {
	if strings.TrimSpace(material) == "" {
		return entradaInvalida("material obrigatório")
	}
	if err := s.validarMaterialDoMentor(ctx, mentor, material); err != nil {
		return err
	}
	ok, err := s.repositorio.Desativar(ctx, mentor, material)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return nil
}
func (s *Servico) ObterParaMentor(ctx context.Context, mentor, material string) (*Material, error) {
	resultado, err := s.repositorio.ObterParaMentor(ctx, mentor, material)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dominio.ErrNaoEncontrado
	}
	return resultado, err
}
func (s *Servico) ObterParaAluno(ctx context.Context, aluno, material string) (*Material, error) {
	resultado, err := s.repositorio.ObterParaAluno(ctx, aluno, material)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dominio.ErrNaoEncontrado
	}
	return resultado, err
}
func (s *Servico) ListarMentor(ctx context.Context, mentor, edital string) ([]Material, error) {
	return s.repositorio.ListarMentor(ctx, mentor, edital)
}
func (s *Servico) ListarAluno(ctx context.Context, aluno, edital string) ([]Material, error) {
	return s.repositorio.ListarAluno(ctx, aluno, edital)
}

func (s *Servico) CriarArquivo(ctx context.Context, mentor string, e EntradaArquivo, conteudo io.Reader, raiz string) (string, error) {
	if err := validarArquivo(e, conteudo != nil); err != nil {
		return "", err
	}
	if err := s.validarEdital(ctx, mentor, e.Escopo, e.EditalID); err != nil {
		return "", err
	}
	if e.Escopo == "global" {
		e.EditalID = nil
	}
	normalizarArquivo(&e)
	id := identificador.UUID()
	caminho, tipoMime, err := salvarArquivo(e.NomeArquivo, conteudo, raiz)
	if err != nil {
		return "", err
	}
	if err = s.repositorio.CriarArquivoRegistro(ctx, mentor, id, e, caminho, tipoMime); err != nil {
		if removeErr := os.Remove(caminho); removeErr != nil {
			return "", fmt.Errorf("%w; falha ao remover arquivo temporário: %v", err, removeErr)
		}
		return "", err
	}
	return id, nil
}
func (s *Servico) AlterarArquivo(ctx context.Context, mentor, material string, e EntradaArquivo, conteudo io.Reader, raiz string) error {
	if strings.TrimSpace(material) == "" {
		return entradaInvalida("material obrigatório")
	}
	if err := s.validarMaterialDoMentor(ctx, mentor, material); err != nil {
		return err
	}
	if err := validarArquivo(e, conteudo != nil); err != nil {
		return err
	}
	if err := s.validarEdital(ctx, mentor, e.Escopo, e.EditalID); err != nil {
		return err
	}
	if e.Escopo == "global" {
		e.EditalID = nil
	}
	normalizarArquivo(&e)
	caminho, tipoMime := "", ""
	var err error
	if conteudo != nil {
		caminho, tipoMime, err = salvarArquivo(e.NomeArquivo, conteudo, raiz)
		if err != nil {
			return err
		}
	}
	anterior, err := s.repositorio.AlterarArquivoRegistro(ctx, mentor, material, e, caminho, tipoMime)
	if err != nil {
		if caminho != "" {
			if removeErr := os.Remove(caminho); removeErr != nil {
				return fmt.Errorf("%w; falha ao remover arquivo temporário: %v", err, removeErr)
			}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return dominio.ErrNaoEncontrado
		}
		return err
	}
	if caminho != "" && anterior != nil && *anterior != "" && *anterior != caminho {
		if err := os.Remove(*anterior); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("arquivo atualizado, mas não foi possível remover o anterior: %w", err)
		}
	}
	return nil
}

func (s *Servico) validarEdital(ctx context.Context, mentor, escopo string, edital *string) error {
	if escopo != "edital" {
		return nil
	}
	ok, err := s.repositorio.EditalDoMentor(ctx, mentor, *edital)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: edital inválido", dominio.ErrNaoEncontrado)
	}
	return nil
}

func (s *Servico) validarMaterialDoMentor(ctx context.Context, mentor, material string) error {
	ok, err := s.repositorio.MaterialDoMentor(ctx, mentor, material)
	if err != nil {
		return err
	}
	if !ok {
		return dominio.ErrNaoEncontrado
	}
	return nil
}
func validarMaterial(e EntradaMaterial, permitirArquivo bool) error {
	if strings.TrimSpace(e.Titulo) == "" || !tiposValidos[e.Tipo] || (!permitirArquivo && e.Tipo == "arquivo") || !escoposValidos[e.Escopo] {
		return entradaInvalida("dados de material inválidos")
	}
	if e.Escopo == "edital" && (e.EditalID == nil || strings.TrimSpace(*e.EditalID) == "") {
		return entradaInvalida("edital obrigatório")
	}
	if e.Tipo == "texto" && (e.Texto == nil || strings.TrimSpace(*e.Texto) == "") {
		return entradaInvalida("texto obrigatório")
	}
	if e.Tipo != "texto" && e.Tipo != "arquivo" && (e.Url == nil || strings.TrimSpace(*e.Url) == "") {
		return entradaInvalida("url obrigatória")
	}
	return nil
}
func validarArquivo(e EntradaArquivo, conteudo bool) error {
	if strings.TrimSpace(e.Titulo) == "" || e.Tipo != "arquivo" || !escoposValidos[e.Escopo] || (conteudo && strings.TrimSpace(e.NomeArquivo) == "") {
		return entradaInvalida("dados de arquivo inválidos")
	}
	if e.Escopo == "edital" && (e.EditalID == nil || strings.TrimSpace(*e.EditalID) == "") {
		return entradaInvalida("edital obrigatório")
	}
	return nil
}
func normalizarMaterial(e *EntradaMaterial) {
	e.Titulo = strings.TrimSpace(e.Titulo)
	normalizarTexto(&e.Descricao)
	normalizarTexto(&e.Url)
	normalizarTexto(&e.Texto)
	normalizarTexto(&e.Pasta)
	normalizarTexto(&e.EditalID)
	if e.Escopo == "global" {
		e.EditalID = nil
	}
	if e.Tipo == "texto" {
		e.Url = nil
	} else {
		e.Texto = nil
	}
}

func normalizarArquivo(e *EntradaArquivo) {
	e.Titulo = strings.TrimSpace(e.Titulo)
	e.NomeArquivo = strings.TrimSpace(e.NomeArquivo)
	normalizarTexto(&e.Descricao)
	normalizarTexto(&e.Pasta)
	normalizarTexto(&e.EditalID)
}

func normalizarTexto(v **string) {
	if *v == nil {
		return
	}
	normalizado := strings.TrimSpace(**v)
	*v = &normalizado
}
func salvarArquivo(nome string, conteudo io.Reader, raiz string) (string, string, error) {
	return arquivos.SalvarMaterial(raiz, nome, conteudo)
}
