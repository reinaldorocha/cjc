package cursos

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"chega-junto-concurseiro-web/internal/identificador"
	"unicode"
	"unicode/utf8"
)

const LimitePDF = 10 << 20

func validarPDF(nome string, conteudo []byte) (string, error) {
	nome = strings.ReplaceAll(nome, "\\", "/")
	partes := strings.Split(nome, "/")
	nome = strings.TrimSpace(partes[len(partes)-1])
	if nome == "" || utf8.RuneCountInString(nome) > 180 || !strings.HasSuffix(strings.ToLower(nome), ".pdf") || len(conteudo) > LimitePDF || !bytes.HasPrefix(conteudo, []byte("%PDF-")) || strings.IndexFunc(nome, unicode.IsControl) >= 0 {
		return "", ErrEntrada
	}
	return nome, nil
}

func (s *Servico) SalvarPDF(ctx context.Context, mentor, nome string, conteudo []byte) (string, string, error) {
	nome, err := validarPDF(nome, conteudo)
	if err != nil {
		return "", "", err
	}
	id := identificador.UUID()
	_, err = s.db.ExecContext(ctx, `INSERT INTO cursos_pdfs (id,mentor_id,nome,conteudo) VALUES (?,?,?,?)`, id, mentor, nome, conteudo)
	return id, nome, err
}

// O arquivo só é entregue se estiver em uma aula de um curso acessível.
func (s *Servico) ObterPDF(ctx context.Context, usuario string, mentor bool, cursoID, pdfID string) (string, []byte, error) {
	lista, err := s.Listar(ctx, usuario, mentor, cursoID)
	if err != nil {
		return "", nil, err
	}
	if len(lista) != 1 {
		return "", nil, ErrNaoEncontrado
	}
	vinculado := false
	for _, a := range lista[0].Aulas {
		if a.PDFID == pdfID && pdfID != "" {
			vinculado = true
			break
		}
	}
	if !vinculado {
		return "", nil, ErrNaoEncontrado
	}
	var nome string
	var conteudo []byte
	err = s.db.QueryRowContext(ctx, `SELECT p.nome,p.conteudo FROM cursos_pdfs p JOIN cursos c ON c.mentor_id=p.mentor_id WHERE c.id=? AND p.id=? AND c.ativo=TRUE`, cursoID, pdfID).Scan(&nome, &conteudo)
	if err == sql.ErrNoRows {
		err = ErrNaoEncontrado
	}
	return nome, conteudo, err
}
