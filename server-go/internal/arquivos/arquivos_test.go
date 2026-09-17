package arquivos

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
)

func pngTeste(t *testing.T) []byte {
	t.Helper()
	imagem := image.NewRGBA(image.Rect(0, 0, 2, 2))
	imagem.Set(0, 0, color.RGBA{R: 255, A: 255})
	var b bytes.Buffer
	if err := png.Encode(&b, imagem); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestLerImagemValidaConteudoENormalizaExtensao(t *testing.T) {
	imagem, err := LerImagem(bytes.NewReader(pngTeste(t)))
	if err != nil {
		t.Fatal(err)
	}
	if imagem.Extensao != ".png" || imagem.Mime != "image/png" {
		t.Fatalf("tipo inesperado: %#v", imagem)
	}
	if _, err := LerImagem(strings.NewReader("<html>arquivo falso</html>")); err != ErrImagemInvalida {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestSalvarMaterialRestringeFormatoEGravaEmCaminhoPrivado(t *testing.T) {
	raiz := t.TempDir()
	if _, _, err := SalvarMaterial(raiz, "perigoso.exe", strings.NewReader("x")); err != ErrMaterialInvalido {
		t.Fatalf("deveria rejeitar extensão: %v", err)
	}
	if _, _, err := SalvarMaterial(raiz, "falso.pdf", strings.NewReader("<html>falso</html>")); err != ErrMaterialInvalido {
		t.Fatalf("deveria rejeitar conteúdo PDF inválido: %v", err)
	}
	caminho, tipo, err := SalvarMaterial(raiz, "aula.pdf", strings.NewReader("%PDF-1.7\nconteúdo"))
	if err != nil {
		t.Fatal(err)
	}
	if tipo != "application/pdf" {
		t.Fatalf("mime inesperado: %s", tipo)
	}
	if _, err := os.Stat(caminho); err != nil {
		t.Fatal(err)
	}
}
