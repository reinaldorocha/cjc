package cursos

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"chega-junto-concurseiro-web/internal/identificador"
)

const LimiteCapa = 5 << 20
const prefixoCapa = "/api/v1/cursos/capas/"

var ErrCapa = errors.New("envie uma imagem JPG ou PNG de até 5 MB e 16 milhões de pixels")

func idCapa(url string) string {
	if !strings.HasPrefix(url, prefixoCapa) {
		return ""
	}
	id := strings.TrimPrefix(url, prefixoCapa)
	if !uuidValido.MatchString(id) {
		return ""
	}
	return id
}

func validarCapa(conteudo []byte) (string, error) {
	if len(conteudo) == 0 || len(conteudo) > LimiteCapa {
		return "", ErrCapa
	}
	config, formato, err := image.DecodeConfig(bytes.NewReader(conteudo))
	if err != nil || (formato != "jpeg" && formato != "png") || config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 16000000 {
		return "", ErrCapa
	}
	if _, _, err = image.Decode(bytes.NewReader(conteudo)); err != nil {
		return "", ErrCapa
	}
	return "image/" + formato, nil
}

func (s *Servico) SalvarCapa(ctx context.Context, mentor string, conteudo []byte) (string, error) {
	tipo, err := validarCapa(conteudo)
	if err != nil {
		return "", err
	}
	id := identificador.UUID()
	_, err = s.db.ExecContext(ctx, `INSERT INTO cursos_capas (id,mentor_id,tipo,conteudo) VALUES (?,?,?,?)`, id, mentor, tipo, conteudo)
	return prefixoCapa + id, err
}

func (s *Servico) ObterCapa(ctx context.Context, usuario string, mentor bool, id string) (string, []byte, error) {
	if !uuidValido.MatchString(id) {
		return "", nil, ErrNaoEncontrado
	}
	var dono, tipo string
	var conteudo []byte
	err := s.db.QueryRowContext(ctx, `SELECT mentor_id,tipo,conteudo FROM cursos_capas WHERE id=?`, id).Scan(&dono, &tipo, &conteudo)
	if err == sql.ErrNoRows {
		return "", nil, ErrNaoEncontrado
	}
	if err != nil {
		return "", nil, err
	}
	if mentor && dono == usuario {
		return tipo, conteudo, nil
	}
	if !mentor {
		lista, err := s.Listar(ctx, usuario, false, "")
		if err != nil {
			return "", nil, err
		}
		for _, c := range lista {
			if c.MentorID == dono {
				for _, a := range c.Aulas {
					if a.ModuloCapaURL == prefixoCapa+id {
						return tipo, conteudo, nil
					}
				}
			}
			if c.CapaURL == prefixoCapa+id && c.MentorID == dono {
				return tipo, conteudo, nil
			}
		}
	}
	return "", nil, ErrNaoEncontrado
}
