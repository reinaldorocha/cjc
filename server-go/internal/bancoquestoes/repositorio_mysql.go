package bancoquestoes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"track-concursos-web/internal/identificador"
)

type repositorioMySQL struct {
	banco *sql.DB
}

func novoRepositorioMySQL(banco *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: banco} }

type Questao struct {
	ID              string    `json:"id"`
	Disciplina      string    `json:"disciplina"`
	Assunto         string    `json:"assunto"`
	Tipo            string    `json:"tipo"`
	Enunciado       string    `json:"enunciado"`
	Alternativas    []string  `json:"alternativas,omitempty"`
	RespostaCorreta string    `json:"respostaCorreta"`
	Explicacao      string    `json:"explicacao,omitempty"`
	Alcance         string    `json:"alcance"`
	ConcursoID      *string   `json:"concursoId,omitempty"`
	CriadoPor       string    `json:"criadoPor"`
	Ativo           bool      `json:"ativo"`
	CriadoEm        time.Time `json:"criadoEm"`
	Respondida      *bool     `json:"respondida,omitempty"`
	UltimaResposta  *string   `json:"ultimaResposta,omitempty"`
	UltimoAcerto    *bool     `json:"ultimoAcerto,omitempty"`
}

type EntradaQuestao struct {
	Disciplina      string   `json:"disciplina"`
	Assunto         string   `json:"assunto"`
	Tipo            string   `json:"tipo"`
	Enunciado       string   `json:"enunciado"`
	Alternativas    []string `json:"alternativas,omitempty"`
	RespostaCorreta string   `json:"respostaCorreta"`
	Explicacao      string   `json:"explicacao,omitempty"`
	Alcance         string   `json:"alcance"`
	ConcursoID      *string  `json:"concursoId,omitempty"`
}

type FiltroQuestao struct {
	Disciplina string
	Assunto    string
	Tipo       string
	ConcursoID string
	AlunoID    string
}

type RespostaQuestaoEntrada struct {
	QuestaoID     string `json:"questaoId"`
	RespostaAluno string `json:"respostaAluno"`
	ConcursoID    string `json:"concursoId,omitempty"`
}

type ResultadoResposta struct {
	Correto         bool   `json:"correto"`
	RespostaCorreta string `json:"respostaCorreta"`
	Explicacao      string `json:"explicacao"`
}

type EstatisticasBancoQuestoes struct {
	TotalRespondidas int                         `json:"totalRespondidas"`
	TotalResolvidas  int                         `json:"totalResolvidas"`
	TotalAcertos     int                         `json:"totalAcertos"`
	TotalErros       int                         `json:"totalErros"`
	TaxaAcerto       float64                     `json:"taxaAcerto"`
	PorDisciplina    []EstatisticaPorDisciplina  `json:"porDisciplina"`
	PorAssunto       []EstatisticaPorAssunto     `json:"porAssunto"`
	HistoricoRecente []RegistroHistoricoQuestoes `json:"historicoRecente"`
}

type EstatisticaPorDisciplina struct {
	Disciplina string  `json:"disciplina"`
	Total      int     `json:"total"`
	Acertos    int     `json:"acertos"`
	Erros      int     `json:"erros"`
	TaxaAcerto float64 `json:"taxaAcerto"`
}

type EstatisticaPorAssunto struct {
	Disciplina string  `json:"disciplina"`
	Assunto    string  `json:"assunto"`
	Total      int     `json:"total"`
	Acertos    int     `json:"acertos"`
	Erros      int     `json:"erros"`
	TaxaAcerto float64 `json:"taxaAcerto"`
}

type RegistroHistoricoQuestoes struct {
	ID            string    `json:"id"`
	QuestaoID     string    `json:"questaoId"`
	Disciplina    string    `json:"disciplina"`
	Assunto       string    `json:"assunto"`
	Tipo          string    `json:"tipo"`
	Enunciado     string    `json:"enunciado"`
	RespostaAluno string    `json:"respostaAluno"`
	Correto       bool      `json:"correto"`
	RespondidoEm  time.Time `json:"respondidoEm"`
}

func jsonNulo(v any) any {
	b, err := json.Marshal(v)
	if err != nil || string(b) == "null" {
		return nil
	}
	return string(b)
}

func (s *repositorioMySQL) ListarQuestoes(ctx context.Context, f FiltroQuestao) ([]Questao, error) {
	var condicoes []string
	var args []any

	condicoes = append(condicoes, "q.ativo = TRUE")

	if strings.TrimSpace(f.Disciplina) != "" {
		condicoes = append(condicoes, "q.disciplina = ?")
		args = append(args, strings.TrimSpace(f.Disciplina))
	}
	if strings.TrimSpace(f.Assunto) != "" {
		condicoes = append(condicoes, "q.assunto = ?")
		args = append(args, strings.TrimSpace(f.Assunto))
	}
	if strings.TrimSpace(f.Tipo) != "" {
		condicoes = append(condicoes, "q.tipo = ?")
		args = append(args, strings.TrimSpace(f.Tipo))
	}
	if strings.TrimSpace(f.ConcursoID) != "" {
		condicoes = append(condicoes, "(q.alcance = 'global' OR q.concurso_id = ?)")
		args = append(args, strings.TrimSpace(f.ConcursoID))
	}

	whereClause := strings.Join(condicoes, " AND ")

	query := `SELECT q.id, q.disciplina, q.assunto, q.tipo, q.enunciado, q.alternativas, q.resposta_correta,
		COALESCE(q.explicacao, ''), q.alcance, q.concurso_id, q.criado_por, q.ativo, q.criado_em`

	if strings.TrimSpace(f.AlunoID) != "" {
		query += `,
		(r.id IS NOT NULL) AS respondida,
		r.resposta_aluno,
		r.correto`
	}

	query += ` FROM banco_questoes q`

	if strings.TrimSpace(f.AlunoID) != "" {
		query += ` LEFT JOIN respostas_banco_questoes r ON r.id = (
			SELECT r2.id FROM respostas_banco_questoes r2
			WHERE r2.aluno_id = ? AND r2.questao_id = q.id
			ORDER BY r2.respondido_em DESC LIMIT 1
		)`
		args = append([]any{f.AlunoID}, args...)
	}

	query += ` WHERE ` + whereClause + ` ORDER BY q.disciplina ASC, q.assunto ASC, q.criado_em DESC`

	linhas, err := s.banco.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar questoes: %w", err)
	}
	defer linhas.Close()

	var lista []Questao
	for linhas.Next() {
		var q Questao
		var altRaw sql.NullString
		var concID sql.NullString
		var respondida sql.NullBool
		var ultResp sql.NullString
		var ultAcerto sql.NullBool

		var errScan error
		if strings.TrimSpace(f.AlunoID) != "" {
			errScan = linhas.Scan(
				&q.ID, &q.Disciplina, &q.Assunto, &q.Tipo, &q.Enunciado, &altRaw, &q.RespostaCorreta,
				&q.Explicacao, &q.Alcance, &concID, &q.CriadoPor, &q.Ativo, &q.CriadoEm,
				&respondida, &ultResp, &ultAcerto,
			)
		} else {
			errScan = linhas.Scan(
				&q.ID, &q.Disciplina, &q.Assunto, &q.Tipo, &q.Enunciado, &altRaw, &q.RespostaCorreta,
				&q.Explicacao, &q.Alcance, &concID, &q.CriadoPor, &q.Ativo, &q.CriadoEm,
			)
		}
		if errScan != nil {
			return nil, errScan
		}

		if concID.Valid {
			v := concID.String
			q.ConcursoID = &v
		}
		if altRaw.Valid && strings.TrimSpace(altRaw.String) != "" {
			if err := json.Unmarshal([]byte(altRaw.String), &q.Alternativas); err != nil {
				return nil, err
			}
		}
		if respondida.Valid {
			v := respondida.Bool
			q.Respondida = &v
		}
		if ultResp.Valid {
			v := ultResp.String
			q.UltimaResposta = &v
		}
		if ultAcerto.Valid {
			v := ultAcerto.Bool
			q.UltimoAcerto = &v
		}

		lista = append(lista, q)
	}
	if err := linhas.Err(); err != nil {
		return nil, err
	}
	return lista, nil
}

func extrairLetraETexto(alt string) (string, string) {
	s := strings.TrimSpace(alt)
	if len(s) >= 2 {
		r := strings.ToUpper(s[0:1])
		if r >= "A" && r <= "E" {
			sep := s[1]
			if sep == ')' || sep == '.' || sep == '-' || sep == ':' {
				resto := strings.TrimSpace(s[2:])
				return r, resto
			}
		}
	}
	return "", s
}

func encontrarIndiceAlternativa(alternativas []string, valor string) int {
	v := strings.TrimSpace(valor)
	if v == "" {
		return -1
	}
	vUpper := strings.ToUpper(v)

	// Caso seja apenas a letra: "A", "B", "C", "D", "E"
	if len(vUpper) == 1 && vUpper >= "A" && vUpper <= "E" {
		idx := int(vUpper[0] - 'A')
		if idx < len(alternativas) {
			return idx
		}
	}
	// Se tiver formato "A) Texto" ou "A. Texto", etc.
	letraV, textoV := extrairLetraETexto(v)
	if letraV != "" {
		idx := int(letraV[0] - 'A')
		if idx < len(alternativas) {
			return idx
		}
	}

	// Comparação direta com as alternativas cadastradas
	for i, alt := range alternativas {
		altTrim := strings.TrimSpace(alt)
		if strings.EqualFold(altTrim, v) {
			return i
		}
		letraAlt, textoAlt := extrairLetraETexto(altTrim)
		if letraAlt != "" && (strings.EqualFold(textoAlt, v) || strings.EqualFold(textoAlt, textoV)) {
			return i
		}
		if textoV != "" && strings.EqualFold(altTrim, textoV) {
			return i
		}
	}
	return -1
}

func verificarResposta(tipo string, alternativas []string, respAluno, respCorreta string) bool {
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

	// multipla_escolha
	idxAluno := encontrarIndiceAlternativa(alternativas, aluno)
	idxCorreta := encontrarIndiceAlternativa(alternativas, correta)
	if idxAluno >= 0 && idxCorreta >= 0 && idxAluno == idxCorreta {
		return true
	}

	return false
}

func (s *repositorioMySQL) CriarQuestao(ctx context.Context, e EntradaQuestao, mentorID string) (*Questao, error) {
	disciplina := strings.TrimSpace(e.Disciplina)
	assunto := strings.TrimSpace(e.Assunto)
	if assunto == "" {
		assunto = "Geral"
	}
	enunciado := strings.TrimSpace(e.Enunciado)
	respostaCorreta := strings.TrimSpace(e.RespostaCorreta)

	if disciplina == "" || enunciado == "" || respostaCorreta == "" {
		return nil, errors.New("disciplina, enunciado e resposta correta são obrigatórios")
	}

	tipo := strings.TrimSpace(e.Tipo)
	if tipo != "certo_errado" && tipo != "multipla_escolha" {
		tipo = "multipla_escolha"
	}
	alcance := strings.TrimSpace(e.Alcance)
	if alcance != "concurso" {
		alcance = "global"
	}

	id := identificador.UUID()
	_, err := s.banco.ExecContext(ctx, `INSERT INTO banco_questoes (id, disciplina, assunto, tipo, enunciado, alternativas, resposta_correta, explicacao, alcance, concurso_id, criado_por) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, disciplina, assunto, tipo, enunciado,
		jsonNulo(e.Alternativas), respostaCorreta, strings.TrimSpace(e.Explicacao),
		alcance, e.ConcursoID, mentorID,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao cadastrar questao: %w", err)
	}

	return &Questao{
		ID:              id,
		Disciplina:      disciplina,
		Assunto:         assunto,
		Tipo:            tipo,
		Enunciado:       enunciado,
		Alternativas:    e.Alternativas,
		RespostaCorreta: respostaCorreta,
		Explicacao:      strings.TrimSpace(e.Explicacao),
		Alcance:         alcance,
		ConcursoID:      e.ConcursoID,
		CriadoPor:       mentorID,
		Ativo:           true,
		CriadoEm:        time.Now(),
	}, nil
}

func (s *repositorioMySQL) ImportarQuestoes(ctx context.Context, lista []EntradaQuestao, mentorID string) (int, error) {
	if len(lista) == 0 {
		return 0, errors.New("lista de questões vazia")
	}

	tx, err := s.banco.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	inseridas := 0
	for _, e := range lista {
		disciplina := strings.TrimSpace(e.Disciplina)
		enunciado := strings.TrimSpace(e.Enunciado)
		respostaCorreta := strings.TrimSpace(e.RespostaCorreta)
		if disciplina == "" || enunciado == "" || respostaCorreta == "" {
			continue
		}
		assunto := strings.TrimSpace(e.Assunto)
		if assunto == "" {
			assunto = "Geral"
		}
		tipo := strings.TrimSpace(e.Tipo)
		if tipo != "certo_errado" && tipo != "multipla_escolha" {
			tipo = "multipla_escolha"
		}
		alcance := strings.TrimSpace(e.Alcance)
		if alcance != "concurso" {
			alcance = "global"
		}

		id := identificador.UUID()
		_, err := tx.ExecContext(ctx, `INSERT INTO banco_questoes (id, disciplina, assunto, tipo, enunciado, alternativas, resposta_correta, explicacao, alcance, concurso_id, criado_por) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, disciplina, assunto, tipo, enunciado,
			jsonNulo(e.Alternativas), respostaCorreta, strings.TrimSpace(e.Explicacao),
			alcance, e.ConcursoID, mentorID,
		)
		if err != nil {
			return 0, err
		}
		inseridas++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return inseridas, nil
}

func (s *repositorioMySQL) AlterarQuestao(ctx context.Context, id string, e EntradaQuestao) error {
	tipo := strings.TrimSpace(e.Tipo)
	if tipo != "certo_errado" && tipo != "multipla_escolha" {
		tipo = "multipla_escolha"
	}
	alcance := strings.TrimSpace(e.Alcance)
	if alcance != "concurso" {
		alcance = "global"
	}

	_, err := s.banco.ExecContext(ctx, `UPDATE banco_questoes SET disciplina=?, assunto=?, tipo=?, enunciado=?, alternativas=?, resposta_correta=?, explicacao=?, alcance=?, concurso_id=? WHERE id=? AND ativo=TRUE`,
		strings.TrimSpace(e.Disciplina), strings.TrimSpace(e.Assunto), tipo, strings.TrimSpace(e.Enunciado),
		jsonNulo(e.Alternativas), strings.TrimSpace(e.RespostaCorreta), strings.TrimSpace(e.Explicacao),
		alcance, e.ConcursoID, id,
	)
	return err
}

func (s *repositorioMySQL) DesativarQuestao(ctx context.Context, id string) error {
	_, err := s.banco.ExecContext(ctx, `UPDATE banco_questoes SET ativo=FALSE WHERE id=?`, id)
	return err
}

func (s *repositorioMySQL) Responder(ctx context.Context, alunoID string, e RespostaQuestaoEntrada) (*ResultadoResposta, error) {
	if strings.TrimSpace(e.QuestaoID) == "" || strings.TrimSpace(e.RespostaAluno) == "" {
		return nil, errors.New("questaoId e respostaAluno são obrigatórios")
	}

	var respCorreta, explicacao, tipo string
	var altRaw sql.NullString
	err := s.banco.QueryRowContext(ctx, `SELECT tipo, alternativas, resposta_correta, COALESCE(explicacao, '') FROM banco_questoes WHERE id=? AND ativo=TRUE`, e.QuestaoID).Scan(&tipo, &altRaw, &respCorreta, &explicacao)
	if err != nil {
		return nil, errors.New("questão não encontrada")
	}

	var alternativas []string
	if altRaw.Valid && strings.TrimSpace(altRaw.String) != "" {
		_ = json.Unmarshal([]byte(altRaw.String), &alternativas)
	}

	correto := verificarResposta(tipo, alternativas, e.RespostaAluno, respCorreta)

	var concID *string
	if strings.TrimSpace(e.ConcursoID) != "" {
		v := strings.TrimSpace(e.ConcursoID)
		concID = &v
	}

	id := identificador.UUID()
	_, err = s.banco.ExecContext(ctx, `INSERT INTO respostas_banco_questoes (id, aluno_id, questao_id, concurso_id, resposta_aluno, correto) VALUES (?, ?, ?, ?, ?, ?)`,
		id, alunoID, e.QuestaoID, concID, strings.TrimSpace(e.RespostaAluno), correto,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao registrar resposta: %w", err)
	}

	return &ResultadoResposta{
		Correto:         correto,
		RespostaCorreta: respCorreta,
		Explicacao:      explicacao,
	}, nil
}

func (s *repositorioMySQL) HistoricoQuestao(ctx context.Context, alunoID, questaoID, concursoID string) ([]RegistroHistoricoQuestoes, error) {
	if strings.TrimSpace(questaoID) == "" {
		return nil, errors.New("questaoId e obrigat\u00f3rio")
	}

	query := `SELECT r.id, r.questao_id, q.disciplina, q.assunto, q.tipo, q.enunciado, r.resposta_aluno, r.correto, r.respondido_em
		FROM respostas_banco_questoes r
		JOIN banco_questoes q ON q.id = r.questao_id
		WHERE r.aluno_id = ? AND r.questao_id = ?`
	args := []any{alunoID, strings.TrimSpace(questaoID)}
	if strings.TrimSpace(concursoID) != "" {
		query += ` AND r.concurso_id = ?`
		args = append(args, strings.TrimSpace(concursoID))
	}
	query += ` ORDER BY r.respondido_em DESC, r.id DESC`

	linhas, err := s.banco.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar hist\u00f3rico da quest\u00e3o: %w", err)
	}
	defer linhas.Close()

	var historico []RegistroHistoricoQuestoes
	for linhas.Next() {
		var registro RegistroHistoricoQuestoes
		if err := linhas.Scan(&registro.ID, &registro.QuestaoID, &registro.Disciplina, &registro.Assunto, &registro.Tipo, &registro.Enunciado, &registro.RespostaAluno, &registro.Correto, &registro.RespondidoEm); err != nil {
			return nil, err
		}
		historico = append(historico, registro)
	}
	return historico, linhas.Err()
}

func (s *repositorioMySQL) ObterEstatisticas(ctx context.Context, alunoID string, concursoID string) (*EstatisticasBancoQuestoes, error) {
	var cond []string
	var args []any

	cond = append(cond, "r.aluno_id = ?")
	args = append(args, alunoID)

	if strings.TrimSpace(concursoID) != "" {
		cond = append(cond, "r.concurso_id = ?")
		args = append(args, strings.TrimSpace(concursoID))
	}

	whereClause := strings.Join(cond, " AND ")

	var stats EstatisticasBancoQuestoes

	// Totais
	err := s.banco.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(SUM(CASE WHEN r.correto THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN NOT r.correto THEN 1 ELSE 0 END), 0) FROM respostas_banco_questoes r WHERE `+whereClause, args...).Scan(&stats.TotalRespondidas, &stats.TotalAcertos, &stats.TotalErros)
	if err != nil {
		return nil, err
	}
	stats.TotalResolvidas = stats.TotalRespondidas
	if stats.TotalRespondidas > 0 {
		stats.TaxaAcerto = float64(stats.TotalAcertos) / float64(stats.TotalRespondidas) * 100.0
	}

	// Por Disciplina
	discRows, err := s.banco.QueryContext(ctx, `SELECT q.disciplina, COUNT(r.id), SUM(CASE WHEN r.correto THEN 1 ELSE 0 END), SUM(CASE WHEN NOT r.correto THEN 1 ELSE 0 END) FROM respostas_banco_questoes r JOIN banco_questoes q ON q.id=r.questao_id WHERE `+whereClause+` GROUP BY q.disciplina ORDER BY COUNT(r.id) DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer discRows.Close()
	for discRows.Next() {
		var d EstatisticaPorDisciplina
		if err := discRows.Scan(&d.Disciplina, &d.Total, &d.Acertos, &d.Erros); err != nil {
			return nil, err
		}
		if d.Total > 0 {
			d.TaxaAcerto = float64(d.Acertos) / float64(d.Total) * 100.0
		}
		stats.PorDisciplina = append(stats.PorDisciplina, d)
	}
	if err := discRows.Err(); err != nil {
		return nil, err
	}

	// Por Assunto
	assuntoRows, err := s.banco.QueryContext(ctx, `SELECT q.disciplina, q.assunto, COUNT(r.id), SUM(CASE WHEN r.correto THEN 1 ELSE 0 END), SUM(CASE WHEN NOT r.correto THEN 1 ELSE 0 END) FROM respostas_banco_questoes r JOIN banco_questoes q ON q.id=r.questao_id WHERE `+whereClause+` GROUP BY q.disciplina, q.assunto ORDER BY COUNT(r.id) DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer assuntoRows.Close()
	for assuntoRows.Next() {
		var a EstatisticaPorAssunto
		if err := assuntoRows.Scan(&a.Disciplina, &a.Assunto, &a.Total, &a.Acertos, &a.Erros); err != nil {
			return nil, err
		}
		if a.Total > 0 {
			a.TaxaAcerto = float64(a.Acertos) / float64(a.Total) * 100.0
		}
		stats.PorAssunto = append(stats.PorAssunto, a)
	}
	if err := assuntoRows.Err(); err != nil {
		return nil, err
	}

	// Histórico Recente
	histRows, err := s.banco.QueryContext(ctx, `SELECT r.id, r.questao_id, q.disciplina, q.assunto, q.tipo, q.enunciado, r.resposta_aluno, r.correto, r.respondido_em FROM respostas_banco_questoes r JOIN banco_questoes q ON q.id=r.questao_id WHERE `+whereClause+` ORDER BY r.respondido_em DESC LIMIT 20`, args...)
	if err != nil {
		return nil, err
	}
	defer histRows.Close()
	for histRows.Next() {
		var h RegistroHistoricoQuestoes
		if err := histRows.Scan(&h.ID, &h.QuestaoID, &h.Disciplina, &h.Assunto, &h.Tipo, &h.Enunciado, &h.RespostaAluno, &h.Correto, &h.RespondidoEm); err != nil {
			return nil, err
		}
		stats.HistoricoRecente = append(stats.HistoricoRecente, h)
	}
	if err := histRows.Err(); err != nil {
		return nil, err
	}

	return &stats, nil
}
