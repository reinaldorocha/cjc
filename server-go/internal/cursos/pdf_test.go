package cursos

import (
	"bytes"
	"testing"
)

func TestValidarPDF(t *testing.T) {
	if nome, err := validarPDF(`C:\fakepath\Aula 1.pdf`, []byte("%PDF-1.7\n%%EOF")); err != nil || nome != "Aula 1.pdf" {
		t.Fatal(nome, err)
	}
	for _, caso := range []struct {
		nome     string
		conteudo []byte
	}{
		{"aula.pdf", []byte("<html>falso</html>")},
		{"aula.html", []byte("%PDF-1.4")},
		{"aula\r\n.pdf", []byte("%PDF-1.4")},
		{"aula.pdf", nil},
		{"aula.pdf", append([]byte("%PDF-"), bytes.Repeat([]byte("x"), LimitePDF)...)},
	} {
		if _, err := validarPDF(caso.nome, caso.conteudo); err == nil {
			t.Errorf("PDF inválido aceito: %q", caso.nome)
		}
	}
}
