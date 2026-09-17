package cursos

import (
	"context"
	"encoding/json"
	"strings"
)

type Questao struct {
	ID           string   `json:"id"`
	Disciplina   string   `json:"disciplina"`
	Assunto      string   `json:"assunto"`
	Tipo         string   `json:"tipo"`
	Enunciado    string   `json:"enunciado"`
	Alternativas []string `json:"alternativas"`
}
type ResultadoQuestao struct {
	Correta         bool   `json:"correta"`
	RespostaCorreta string `json:"respostaCorreta"`
	Explicacao      string `json:"explicacao"`
}

func (s *Servico) CatalogoQuestoes(ctx context.Context, mentor, busca string) ([]Questao, error) {
	return s.lerQuestoes(ctx, `SELECT id,disciplina,assunto,tipo,enunciado,alternativas FROM banco_questoes WHERE criado_por=? AND ativo=TRUE AND (enunciado LIKE ? OR disciplina LIKE ? OR assunto LIKE ?) ORDER BY disciplina,assunto,criado_em DESC LIMIT 100`, mentor, "%"+busca+"%", "%"+busca+"%", "%"+busca+"%")
}

func (s *Servico) lerQuestoes(ctx context.Context, q string, args ...any) ([]Questao, error) {
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Questao{}
	for rows.Next() {
		var q Questao
		var alternativas []byte
		if err = rows.Scan(&q.ID, &q.Disciplina, &q.Assunto, &q.Tipo, &q.Enunciado, &alternativas); err != nil {
			return nil, err
		}
		q.Alternativas = []string{}
		if len(alternativas) > 0 {
			if err = json.Unmarshal(alternativas, &q.Alternativas); err != nil {
				return nil, err
			}
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

func (s *Servico) QuestoesAula(ctx context.Context, usuario string, mentor bool, curso, aula string) ([]Questao, error) {
	c, a, err := s.aulaAcessivel(ctx, usuario, mentor, curso, aula)
	if err != nil {
		return nil, err
	}
	if len(a.Questoes) == 0 {
		return []Questao{}, nil
	}
	if mentor {
		previa, err := s.Previa(ctx, usuario, curso)
		if err != nil {
			return nil, err
		}
		for _, p := range previa.Aulas {
			if p.ID == aula && p.Bloqueada {
				return nil, ErrBloqueado
			}
		}
	}
	ids, err := json.Marshal(a.Questoes)
	if err != nil {
		return nil, err
	}
	return s.lerQuestoes(ctx, `SELECT id,disciplina,assunto,tipo,enunciado,alternativas FROM banco_questoes WHERE criado_por=? AND ativo=TRUE AND JSON_CONTAINS(?,JSON_QUOTE(id)) ORDER BY FIELD(id,`+strings.TrimRight(strings.Repeat("?,", len(a.Questoes)), ",")+`)`, append([]any{c.MentorID, ids}, stringsAny(a.Questoes)...)...)
}
func stringsAny(values []string) []any {
	out := make([]any, len(values))
	for i, v := range values {
		out[i] = v
	}
	return out
}

func (s *Servico) ResponderQuestao(ctx context.Context, usuario string, mentor bool, curso, aula, id, resposta string) (ResultadoQuestao, error) {
	questoes, err := s.QuestoesAula(ctx, usuario, mentor, curso, aula)
	if err != nil {
		return ResultadoQuestao{}, err
	}
	var encontrada *Questao
	for i := range questoes {
		if questoes[i].ID == id {
			encontrada = &questoes[i]
			break
		}
	}
	if encontrada == nil {
		return ResultadoQuestao{}, ErrNaoEncontrado
	}
	resposta = strings.TrimSpace(resposta)
	if len(resposta) == 0 || len(resposta) > 255 {
		return ResultadoQuestao{}, ErrEntrada
	}
	if encontrada.Tipo == "certo_errado" {
		ehCerto := strings.EqualFold(resposta, "Certo") || strings.EqualFold(resposta, "C")
		ehErrado := strings.EqualFold(resposta, "Errado") || strings.EqualFold(resposta, "E")
		if !ehCerto && !ehErrado {
			return ResultadoQuestao{}, ErrEntrada
		}
	} else if encontrarIndiceCurso(encontrada.Alternativas, resposta) < 0 {
		return ResultadoQuestao{}, ErrEntrada
	}
	var result ResultadoQuestao
	if err = s.db.QueryRowContext(ctx, `SELECT resposta_correta,COALESCE(explicacao,'') FROM banco_questoes WHERE id=? AND ativo=TRUE`, id).Scan(&result.RespostaCorreta, &result.Explicacao); err != nil {
		return result, err
	}
	result.Correta = verificarRespostaCurso(encontrada.Tipo, encontrada.Alternativas, resposta, result.RespostaCorreta)
	if !mentor {
		_, err = s.db.ExecContext(ctx, `INSERT INTO cursos_respostas (aluno_id,curso_id,aula_id,questao_id,resposta,correta) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE resposta=VALUES(resposta),correta=VALUES(correta),tentativas=tentativas+1,atualizado_em=UTC_TIMESTAMP(6)`, usuario, curso, aula, id, resposta, result.Correta)
	}
	return result, err
}

func verificarRespostaCurso(tipo string, alternativas []string, respAluno, respCorreta string) bool {
	aluno := strings.TrimSpace(respAluno)
	correta := strings.TrimSpace(respCorreta)

	if strings.EqualFold(aluno, correta) {
		return true
	}

	if tipo == "certo_errado" {
		ehCertoAluno := strings.EqualFold(aluno, "certo") || strings.EqualFold(aluno, "c")
		ehCertoGabarito := strings.EqualFold(correta, "certo") || strings.EqualFold(correta, "c")
		if ehCertoAluno && ehCertoGabarito {
			return true
		}

		ehErradoAluno := strings.EqualFold(aluno, "errado") || strings.EqualFold(aluno, "e")
		ehErradoGabarito := strings.EqualFold(correta, "errado") || strings.EqualFold(correta, "e")
		if ehErradoAluno && ehErradoGabarito {
			return true
		}
		return false
	}

	idxAluno := encontrarIndiceCurso(alternativas, aluno)
	idxCorreta := encontrarIndiceCurso(alternativas, correta)
	if idxAluno >= 0 && idxCorreta >= 0 && idxAluno == idxCorreta {
		return true
	}

	return false
}

func encontrarIndiceCurso(alternativas []string, valor string) int {
	v := strings.TrimSpace(valor)
	if v == "" {
		return -1
	}
	vUpper := strings.ToUpper(v)

	if len(vUpper) == 1 && vUpper >= "A" && vUpper <= "E" {
		idx := int(vUpper[0] - 'A')
		if idx < len(alternativas) {
			return idx
		}
	}
	if len(vUpper) >= 2 && vUpper[0] >= 'A' && vUpper[0] <= 'E' {
		sep := vUpper[1]
		if sep == ')' || sep == '.' || sep == '-' || sep == ':' {
			idx := int(vUpper[0] - 'A')
			if idx < len(alternativas) {
				return idx
			}
		}
	}

	for i, alt := range alternativas {
		altTrim := strings.TrimSpace(alt)
		if strings.EqualFold(altTrim, v) {
			return i
		}
		if len(altTrim) >= 2 {
			r := strings.ToUpper(altTrim[0:1])
			if r >= "A" && r <= "E" {
				sep := altTrim[1]
				if sep == ')' || sep == '.' || sep == '-' || sep == ':' {
					texto := strings.TrimSpace(altTrim[2:])
					if strings.EqualFold(texto, v) {
						return i
					}
				}
			}
		}
	}
	return -1
}
