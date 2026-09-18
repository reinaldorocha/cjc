package arquivos

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"chega-junto-concurseiro-web/internal/identificador"
)

const (
	LimiteImagem   = 5 << 20
	LimiteMaterial = 15 << 20
)

var ErrImagemInvalida = errors.New("envie uma imagem JPG ou PNG de até 5 MB")
var ErrMaterialInvalido = errors.New("envie um arquivo PDF, imagem, Office, TXT ou CSV de até 15 MB")

type Imagem struct {
	Conteudo []byte
	Extensao string
	Mime     string
}

func LerImagem(origem io.Reader) (Imagem, error) {
	conteudo, err := io.ReadAll(io.LimitReader(origem, LimiteImagem+1))
	if err != nil || len(conteudo) == 0 || len(conteudo) > LimiteImagem {
		return Imagem{}, ErrImagemInvalida
	}
	config, formato, err := image.DecodeConfig(bytes.NewReader(conteudo))
	if err != nil || (formato != "jpeg" && formato != "png") || config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 16_000_000 {
		return Imagem{}, ErrImagemInvalida
	}
	if _, _, err = image.Decode(bytes.NewReader(conteudo)); err != nil {
		return Imagem{}, ErrImagemInvalida
	}
	if formato == "jpeg" {
		return Imagem{Conteudo: conteudo, Extensao: ".jpg", Mime: "image/jpeg"}, nil
	}
	return Imagem{Conteudo: conteudo, Extensao: ".png", Mime: "image/png"}, nil
}

func SalvarImagem(raiz, pasta string, imagem Imagem) (string, error) {
	if imagem.Extensao == "" || len(imagem.Conteudo) == 0 {
		return "", ErrImagemInvalida
	}
	dir := filepath.Join(raiz, pasta)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("criar diretório de upload: %w", err)
	}
	return salvarAtomico(dir, identificador.UUID()+imagem.Extensao, bytes.NewReader(imagem.Conteudo), LimiteImagem)
}

func SalvarMaterial(raiz, nome string, origem io.Reader) (string, string, error) {
	extensao := strings.ToLower(filepath.Ext(filepath.Base(nome)))
	mimes := map[string]string{
		".pdf": "application/pdf", ".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
		".doc": "application/msword", ".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls": "application/vnd.ms-excel", ".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt": "application/vnd.ms-powerpoint", ".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt": "text/plain; charset=utf-8", ".csv": "text/csv; charset=utf-8",
	}
	tipo, permitido := mimes[extensao]
	if !permitido {
		return "", "", ErrMaterialInvalido
	}
	if err := os.MkdirAll(raiz, 0o755); err != nil {
		return "", "", fmt.Errorf("criar diretório de upload: %w", err)
	}
	caminho, err := salvarAtomico(raiz, identificador.UUID()+extensao, origem, LimiteMaterial)
	if err != nil {
		return "", "", err
	}
	if err := validarMaterialSalvo(caminho, extensao); err != nil {
		_ = os.Remove(caminho)
		return "", "", err
	}
	return caminho, tipo, nil
}

func validarMaterialSalvo(caminho, extensao string) error {
	if extensao == ".pdf" {
		arquivo, err := os.Open(caminho)
		if err != nil {
			return err
		}
		defer arquivo.Close()
		cabecalho := make([]byte, 5)
		n, err := io.ReadFull(arquivo, cabecalho)
		if err != nil || n != 5 || string(cabecalho) != "%PDF-" {
			return ErrMaterialInvalido
		}
	}
	if extensao == ".jpg" || extensao == ".jpeg" || extensao == ".png" {
		arquivo, err := os.Open(caminho)
		if err != nil {
			return err
		}
		defer arquivo.Close()
		if _, err = LerImagem(arquivo); err != nil {
			return ErrMaterialInvalido
		}
	}
	return nil
}

func salvarAtomico(dir, nome string, origem io.Reader, limite int64) (string, error) {
	temporario, err := os.CreateTemp(dir, ".upload-*")
	if err != nil {
		return "", err
	}
	caminhoTemporario := temporario.Name()
	defer os.Remove(caminhoTemporario)
	n, err := io.Copy(temporario, io.LimitReader(origem, limite+1))
	if fechar := temporario.Close(); err == nil {
		err = fechar
	}
	if err != nil || n == 0 || n > limite {
		return "", ErrMaterialInvalido
	}
	caminho := filepath.Join(dir, nome)
	if err := os.Rename(caminhoTemporario, caminho); err != nil {
		return "", err
	}
	return caminho, nil
}
