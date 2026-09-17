package cursos

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"
)

func imagemTeste(t *testing.T) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := png.Encode(&b, image.NewRGBA(image.Rect(0, 0, 16, 9))); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestValidarCapa(t *testing.T) {
	pngDados := imagemTeste(t)
	if tipo, err := validarCapa(pngDados); err != nil || tipo != "image/png" {
		t.Fatal(tipo, err)
	}
	var b bytes.Buffer
	if err := jpeg.Encode(&b, image.NewRGBA(image.Rect(0, 0, 16, 9)), nil); err != nil {
		t.Fatal(err)
	}
	if tipo, err := validarCapa(b.Bytes()); err != nil || tipo != "image/jpeg" {
		t.Fatal(tipo, err)
	}
	for _, dados := range [][]byte{nil, []byte("<svg onload='alert(1)'/>"), pngDados[:24], make([]byte, LimiteCapa+1)} {
		if _, err := validarCapa(dados); err != ErrCapa {
			t.Fatal("imagem inválida aceita", err)
		}
	}
}
