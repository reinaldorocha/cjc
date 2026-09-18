package cartoes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
	"chega-junto-concurseiro-web/internal/identificador"
)

type repositorioMySQL struct {
	db   *sql.DB
	fuso *time.Location
}

type Baralho struct {
	ID               string   `json:"id"`
	AlunoID          *string  `json:"alunoId,omitempty"`
	ConcursoID       *string  `json:"concursoId,omitempty"`
	EditalID         *string  `json:"editalId,omitempty"`
	BaralhoPaiID     *string  `json:"baralhoPaiId,omitempty"`
	MateriaID        *string  `json:"materiaId,omitempty"`
	TopicoID         *string  `json:"topicoId,omitempty"`
	SubtopicoID      *string  `json:"subtopicoId,omitempty"`
	CriadoPor        string   `json:"criadoPor"`
	TipoProprietario string   `json:"tipoProprietario"`
	Alcance          string   `json:"alcance"`
	Nome             string   `json:"nome"`
	Descricao        *string  `json:"descricao,omitempty"`
	Icone            *string  `json:"icone,omitempty"`
	Ordem            int      `json:"ordem"`
	SomenteLeitura   bool     `json:"somenteLeitura"`
	TotalCartoes     int      `json:"totalCartoes"`
	Cartoes          []Cartao `json:"cartoes"`
}

type Cartao struct {
	ID                      string          `json:"id"`
	BaralhoID               string          `json:"baralhoId"`
	Tipo                    string          `json:"tipo"`
	Frente                  string          `json:"frente"`
	Verso                   *string         `json:"verso,omitempty"`
	Dica                    *string         `json:"dica,omitempty"`
	Alternativas            json.RawMessage `json:"alternativas,omitempty"`
	RespostaCorreta         *string         `json:"respostaCorreta,omitempty"`
	Explicacao              *string         `json:"explicacao,omitempty"`
	ExplicacoesAlternativas json.RawMessage `json:"explicacoesAlternativas,omitempty"`
	Etiquetas               json.RawMessage `json:"etiquetas,omitempty"`
	TopicoID                *string         `json:"topicoId,omitempty"`
	SubtopicoID             *string         `json:"subtopicoId,omitempty"`
	CriadoPor               string          `json:"criadoPor"`
	Repeticoes              int             `json:"repeticoes"`
	IntervaloDias           int             `json:"intervaloDias"`
	Facilidade              float64         `json:"facilidade"`
	ProximaRevisao          *string         `json:"proximaRevisao,omitempty"`
	UltimaRevisao           *string         `json:"ultimaRevisao,omitempty"`
	Falhas                  int             `json:"falhas"`
}

type HistoricoRevisao struct {
	ID               string    `json:"id"`
	CartaoID         string    `json:"cartaoId"`
	BaralhoID        string    `json:"baralhoId"`
	Qualidade        int       `json:"qualidade"`
	RevisadoEm       time.Time `json:"revisadoEm"`
	AlunoID          *string   `json:"-"`
	EditalID         *string   `json:"-"`
	Alcance          string    `json:"-"`
	TipoProprietario string    `json:"-"`
	ConcursoID       *string   `json:"-"`
	EditalConcursoID *string   `json:"-"`
}

type EntradaBaralho struct {
	Nome         string  `json:"nome"`
	Descricao    *string `json:"descricao"`
	ConcursoID   *string `json:"concursoId"`
	EditalID     *string `json:"editalId"`
	BaralhoPaiID *string `json:"baralhoPaiId"`
	MateriaID    *string `json:"materiaId"`
	TopicoID     *string `json:"topicoId"`
	SubtopicoID  *string `json:"subtopicoId"`
	Alcance      string  `json:"alcance"`
	Icone        *string `json:"icone"`
	Ordem        int     `json:"ordem"`
	DestinoAluno *string `json:"-"`
	Proprietario string  `json:"-"`
}

type AlteracaoBaralho struct {
	Nome         *string `json:"nome"`
	Descricao    *string `json:"descricao"`
	ConcursoID   *string `json:"concursoId"`
	EditalID     *string `json:"editalId"`
	BaralhoPaiID *string `json:"baralhoPaiId"`
	MateriaID    *string `json:"materiaId"`
	TopicoID     *string `json:"topicoId"`
	SubtopicoID  *string `json:"subtopicoId"`
	Alcance      *string `json:"alcance"`
	Icone        *string `json:"icone"`
	Ordem        *int    `json:"ordem"`
	DestinoAluno *string `json:"-"`
}

type EntradaCartao struct {
	Tipo                    string          `json:"tipo"`
	Frente                  string          `json:"frente"`
	Verso                   *string         `json:"verso"`
	Dica                    *string         `json:"dica"`
	Alternativas            json.RawMessage `json:"alternativas"`
	RespostaCorreta         *string         `json:"respostaCorreta"`
	Explicacao              *string         `json:"explicacao"`
	ExplicacoesAlternativas json.RawMessage `json:"explicacoesAlternativas"`
	Etiquetas               json.RawMessage `json:"etiquetas"`
	TopicoID                *string         `json:"topicoId"`
	SubtopicoID             *string         `json:"subtopicoId"`
}

type AlteracaoCartao struct {
	Tipo                    *string          `json:"tipo"`
	Frente                  *string          `json:"frente"`
	Verso                   *string          `json:"verso"`
	Dica                    *string          `json:"dica"`
	Alternativas            *json.RawMessage `json:"alternativas"`
	RespostaCorreta         *string          `json:"respostaCorreta"`
	Explicacao              *string          `json:"explicacao"`
	ExplicacoesAlternativas *json.RawMessage `json:"explicacoesAlternativas"`
	Etiquetas               *json.RawMessage `json:"etiquetas"`
	TopicoID                *string          `json:"topicoId"`
	SubtopicoID             *string          `json:"subtopicoId"`
}

type Revisao struct {
	CartaoID       string  `json:"cartaoId"`
	Qualidade      int     `json:"qualidade"`
	Repeticoes     int     `json:"repeticoes"`
	IntervaloDias  int     `json:"intervaloDias"`
	Facilidade     float64 `json:"facilidade"`
	ProximaRevisao string  `json:"proximaRevisao"`
}

type CartaoPendente struct {
	Cartao
	BaralhoNome      string  `json:"baralhoNome"`
	Origem           string  `json:"origem"`
	Repeticoes       int     `json:"repeticoes"`
	IntervaloDias    int     `json:"intervaloDias"`
	Facilidade       float64 `json:"facilidade"`
	ProximaRevisao   *string `json:"proximaRevisao,omitempty"`
	AlunoID          *string `json:"-"`
	EditalID         *string `json:"-"`
	Alcance          string  `json:"-"`
	TipoProprietario string  `json:"-"`
}

func novoRepositorioMySQL(db *sql.DB, fuso *time.Location) *repositorioMySQL {
	return &repositorioMySQL{db: db, fuso: fuso}
}

func (s *repositorioMySQL) ListarBaralhosMentor(ctx context.Context, mentorID string) ([]Baralho, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT id,aluno_id,concurso_id,edital_id,baralho_pai_id,materia_id,topico_id,subtopico_id,criado_por,tipo_proprietario,alcance,nome,descricao,icone,ordem FROM baralhos_cartoes WHERE criado_por=? AND tipo_proprietario='mentor' AND ativo=TRUE ORDER BY nome,criado_em`, mentorID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Baralho{}
	for linhas.Next() {
		var b Baralho
		var aluno, concurso, edital, pai, materia, topico, sub, descricao, icone sql.NullString
		if err = linhas.Scan(&b.ID, &aluno, &concurso, &edital, &pai, &materia, &topico, &sub, &b.CriadoPor, &b.TipoProprietario, &b.Alcance, &b.Nome, &descricao, &icone, &b.Ordem); err != nil {
			return nil, err
		}
		b.AlunoID, b.ConcursoID, b.EditalID, b.BaralhoPaiID = texto(aluno), texto(concurso), texto(edital), texto(pai)
		b.MateriaID, b.TopicoID, b.SubtopicoID, b.Descricao, b.Icone = texto(materia), texto(topico), texto(sub), texto(descricao), texto(icone)
		b.Cartoes, err = s.listarCartoes(ctx, b.ID, "", "")
		if err != nil {
			return nil, err
		}
		b.TotalCartoes = len(b.Cartoes)
		lista = append(lista, b)
	}
	return lista, linhas.Err()
}
func (s *repositorioMySQL) CriarBaralhoMentor(ctx context.Context, mentorID string, e EntradaBaralho) (string, error) {
	if e.Alcance == "global" {
		e.EditalID = nil
		e.ConcursoID = nil
	}
	id := identificador.UUID()
	_, err := s.db.ExecContext(ctx, `INSERT INTO baralhos_cartoes(id,aluno_id,concurso_id,edital_id,baralho_pai_id,materia_id,topico_id,subtopico_id,criado_por,tipo_proprietario,alcance,nome,descricao,icone,ordem) VALUES(?,NULL,?,?,?,?,?,?,?,'mentor',?,?,?,?,?)`, id, e.ConcursoID, e.EditalID, e.BaralhoPaiID, e.MateriaID, e.TopicoID, e.SubtopicoID, mentorID, e.Alcance, e.Nome, e.Descricao, e.Icone, e.Ordem)
	return id, err
}

func (s *repositorioMySQL) ListarBaralhos(ctx context.Context, alunoID, usuarioID, papel, concursoID string) ([]Baralho, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT b.id,b.aluno_id,b.concurso_id,b.edital_id,b.baralho_pai_id,b.materia_id,b.topico_id,b.subtopico_id,b.criado_por,b.tipo_proprietario,b.alcance,b.nome,b.descricao,b.icone,b.ordem
		FROM baralhos_cartoes b WHERE b.ativo=TRUE
		ORDER BY b.nome,b.criado_em`)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Baralho{}
	for linhas.Next() {
		var b Baralho
		var aluno, concurso, edital, pai, materia, topico, subtopico, descricao, icone sql.NullString
		if err = linhas.Scan(&b.ID, &aluno, &concurso, &edital, &pai, &materia, &topico, &subtopico, &b.CriadoPor, &b.TipoProprietario, &b.Alcance, &b.Nome, &descricao, &icone, &b.Ordem); err != nil {
			return nil, err
		}
		b.AlunoID, b.ConcursoID, b.EditalID, b.BaralhoPaiID = texto(aluno), texto(concurso), texto(edital), texto(pai)
		b.MateriaID, b.TopicoID, b.SubtopicoID, b.Descricao, b.Icone = texto(materia), texto(topico), texto(subtopico), texto(descricao), texto(icone)
		b.Cartoes, err = s.listarCartoes(ctx, b.ID, alunoID, concursoID)
		if err != nil {
			return nil, err
		}
		b.TotalCartoes = len(b.Cartoes)
		lista = append(lista, b)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) CriarBaralho(ctx context.Context, alunoID, executor, papel string, e EntradaBaralho) (string, error) {
	id := identificador.UUID()
	_, err := s.db.ExecContext(ctx, `INSERT INTO baralhos_cartoes (id,aluno_id,concurso_id,edital_id,baralho_pai_id,materia_id,topico_id,subtopico_id,criado_por,tipo_proprietario,alcance,nome,descricao,icone,ordem) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, id, e.DestinoAluno, e.ConcursoID, e.EditalID, e.BaralhoPaiID, e.MateriaID, e.TopicoID, e.SubtopicoID, executor, e.Proprietario, e.Alcance, e.Nome, e.Descricao, e.Icone, e.Ordem)
	return id, err
}

func (s *repositorioMySQL) AlterarBaralho(ctx context.Context, alunoID, executor, papel, baralhoID string, e AlteracaoBaralho) error {
	if e.Nome != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET nome=? WHERE id=?`, *e.Nome, baralhoID); err != nil {
			return err
		}
	}
	if e.Descricao != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET descricao=NULLIF(?,'') WHERE id=?`, *e.Descricao, baralhoID); err != nil {
			return err
		}
	}
	if e.ConcursoID != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET concurso_id=NULLIF(?,'') WHERE id=?`, *e.ConcursoID, baralhoID); err != nil {
			return err
		}
	}
	if e.BaralhoPaiID != nil {
		if *e.BaralhoPaiID != "" {
		}
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET baralho_pai_id=NULLIF(?,'') WHERE id=?`, *e.BaralhoPaiID, baralhoID); err != nil {
			return err
		}
	}
	if e.MateriaID != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET materia_id=NULLIF(?,'') WHERE id=?`, *e.MateriaID, baralhoID); err != nil {
			return err
		}
	}
	if e.TopicoID != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET topico_id=NULLIF(?,'') WHERE id=?`, *e.TopicoID, baralhoID); err != nil {
			return err
		}
	}
	if e.SubtopicoID != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET subtopico_id=NULLIF(?,'') WHERE id=?`, *e.SubtopicoID, baralhoID); err != nil {
			return err
		}
	}
	if e.Icone != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET icone=NULLIF(?,'') WHERE id=?`, *e.Icone, baralhoID); err != nil {
			return err
		}
	}
	if e.Ordem != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET ordem=? WHERE id=?`, *e.Ordem, baralhoID); err != nil {
			return err
		}
	}
	if e.Alcance != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE baralhos_cartoes SET alcance=?,aluno_id=?,edital_id=? WHERE id=?`, *e.Alcance, e.DestinoAluno, e.EditalID, baralhoID); err != nil {
			return err
		}
	}
	return nil
}

func (s *repositorioMySQL) DesativarBaralho(ctx context.Context, alunoID, executor, papel, baralhoID string) error {
	linhas, err := s.db.QueryContext(ctx, `WITH RECURSIVE arvore AS (SELECT id FROM baralhos_cartoes WHERE id=? AND ativo=TRUE UNION ALL SELECT b.id FROM baralhos_cartoes b JOIN arvore a ON b.baralho_pai_id=a.id WHERE b.ativo=TRUE) SELECT id FROM arvore`, baralhoID)
	if err != nil {
		return err
	}
	ids := []string{}
	for linhas.Next() {
		var id string
		if err = linhas.Scan(&id); err != nil {
			linhas.Close()
			return err
		}
		ids = append(ids, id)
	}
	if err = linhas.Close(); err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, id := range ids {
		if _, err = tx.ExecContext(ctx, `UPDATE cartoes_estudo SET ativo=FALSE WHERE baralho_id=?`, id); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `UPDATE baralhos_cartoes SET ativo=FALSE WHERE id=?`, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *repositorioMySQL) CriarCartao(ctx context.Context, alunoID, executor, papel, baralhoID string, e EntradaCartao) (string, error) {
	id := identificador.UUID()
	_, err := s.db.ExecContext(ctx, `INSERT INTO cartoes_estudo (id,baralho_id,tipo,frente,verso,dica,alternativas,resposta_correta,explicacao,explicacoes_alternativas,etiquetas,topico_id,subtopico_id,criado_por) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, id, baralhoID, e.Tipo, e.Frente, e.Verso, e.Dica, jsonNulo(e.Alternativas), e.RespostaCorreta, e.Explicacao, jsonNulo(e.ExplicacoesAlternativas), jsonNulo(e.Etiquetas), e.TopicoID, e.SubtopicoID, executor)
	return id, err
}

func (s *repositorioMySQL) AlterarCartao(ctx context.Context, alunoID, executor, papel, cartaoID string, e AlteracaoCartao) error {
	var err error
	if e.Tipo != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET tipo=? WHERE id=?`, *e.Tipo, cartaoID)
		if err != nil {
			return err
		}
	}
	if e.Frente != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET frente=? WHERE id=?`, *e.Frente, cartaoID)
		if err != nil {
			return err
		}
	}
	if e.Verso != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET verso=NULLIF(?,'') WHERE id=?`, *e.Verso, cartaoID)
		if err != nil {
			return err
		}
	}
	if e.Dica != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET dica=NULLIF(?,'') WHERE id=?`, *e.Dica, cartaoID)
		if err != nil {
			return err
		}
	}
	if e.Alternativas != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET alternativas=? WHERE id=?`, jsonNulo(*e.Alternativas), cartaoID)
		if err != nil {
			return err
		}
	}
	if e.RespostaCorreta != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET resposta_correta=NULLIF(?,'') WHERE id=?`, *e.RespostaCorreta, cartaoID)
		if err != nil {
			return err
		}
	}
	if e.Explicacao != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET explicacao=NULLIF(?,'') WHERE id=?`, *e.Explicacao, cartaoID)
		if err != nil {
			return err
		}
	}
	if e.ExplicacoesAlternativas != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET explicacoes_alternativas=? WHERE id=?`, jsonNulo(*e.ExplicacoesAlternativas), cartaoID)
		if err != nil {
			return err
		}
	}
	if e.Etiquetas != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET etiquetas=? WHERE id=?`, jsonNulo(*e.Etiquetas), cartaoID)
		if err != nil {
			return err
		}
	}
	if e.TopicoID != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET topico_id=NULLIF(?,'') WHERE id=?`, *e.TopicoID, cartaoID)
		if err != nil {
			return err
		}
	}
	if e.SubtopicoID != nil {
		_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET subtopico_id=NULLIF(?,'') WHERE id=?`, *e.SubtopicoID, cartaoID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *repositorioMySQL) DesativarCartao(ctx context.Context, alunoID, executor, papel, cartaoID string) error {
	var err error
	_, err = s.db.ExecContext(ctx, `UPDATE cartoes_estudo SET ativo=FALSE WHERE id=?`, cartaoID)
	return err
}

func (s *repositorioMySQL) ListarPendentes(ctx context.Context, alunoID string) ([]CartaoPendente, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT c.id,c.baralho_id,c.tipo,c.frente,c.verso,c.dica,c.alternativas,c.resposta_correta,c.explicacao,c.explicacoes_alternativas,c.etiquetas,c.topico_id,c.subtopico_id,c.criado_por,b.nome,b.tipo_proprietario,b.aluno_id,b.edital_id,b.alcance,COALESCE(r.repeticoes,0),COALESCE(r.intervalo_dias,0),COALESCE(r.facilidade,0),r.proxima_revisao
		FROM cartoes_estudo c JOIN baralhos_cartoes b ON b.id=c.baralho_id
		LEFT JOIN revisoes_cartoes r ON r.id=(SELECT r2.id FROM revisoes_cartoes r2 WHERE r2.aluno_id=? AND r2.cartao_id=c.id ORDER BY r2.revisado_em DESC,r2.id DESC LIMIT 1)
		WHERE c.ativo=TRUE AND b.ativo=TRUE AND (r.id IS NULL OR r.proxima_revisao<=CURRENT_DATE)
		ORDER BY COALESCE(r.proxima_revisao,'1000-01-01'),b.nome,c.criado_em`, alunoID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []CartaoPendente{}
	for linhas.Next() {
		var x CartaoPendente
		var verso, dica, resposta, explicacao, topico, subtopico sql.NullString
		var alternativas, explicacoes, etiquetas []byte
		var alunoBaralho, editalBaralho sql.NullString
		var proxima sql.NullTime
		if err = linhas.Scan(&x.ID, &x.BaralhoID, &x.Tipo, &x.Frente, &verso, &dica, &alternativas, &resposta, &explicacao, &explicacoes, &etiquetas, &topico, &subtopico, &x.CriadoPor, &x.BaralhoNome, &x.TipoProprietario, &alunoBaralho, &editalBaralho, &x.Alcance, &x.Repeticoes, &x.IntervaloDias, &x.Facilidade, &proxima); err != nil {
			return nil, err
		}
		x.Verso, x.Dica, x.RespostaCorreta, x.Explicacao, x.TopicoID, x.SubtopicoID = texto(verso), texto(dica), texto(resposta), texto(explicacao), texto(topico), texto(subtopico)
		if len(alternativas) > 0 {
			x.Alternativas = alternativas
		}
		if len(explicacoes) > 0 {
			x.ExplicacoesAlternativas = explicacoes
		}
		if len(etiquetas) > 0 {
			x.Etiquetas = etiquetas
		}
		x.AlunoID, x.EditalID = texto(alunoBaralho), texto(editalBaralho)
		if proxima.Valid {
			v := proxima.Time.Format("2006-01-02")
			x.ProximaRevisao = &v
		}
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) PropriedadeBaralho(ctx context.Context, baralhoID string) (PropriedadeBaralho, error) {
	var destino, edital sql.NullString
	var resultado PropriedadeBaralho
	err := s.db.QueryRowContext(ctx, `SELECT aluno_id,criado_por,tipo_proprietario,edital_id,alcance FROM baralhos_cartoes WHERE id=? AND ativo=TRUE`, baralhoID).Scan(&destino, &resultado.CriadoPor, &resultado.TipoProprietario, &edital, &resultado.Alcance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resultado, nil
		}
		return resultado, err
	}
	if destino.Valid {
		resultado.AlunoID = &destino.String
	}
	if edital.Valid {
		resultado.EditalID = &edital.String
	}
	return resultado, nil
}

func (s *repositorioMySQL) BaralhoDescendente(ctx context.Context, baralhoID, candidatoID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `WITH RECURSIVE arvore AS (SELECT id FROM baralhos_cartoes WHERE baralho_pai_id=? AND ativo=TRUE UNION ALL SELECT b.id FROM baralhos_cartoes b JOIN arvore a ON b.baralho_pai_id=a.id WHERE b.ativo=TRUE) SELECT COUNT(*) FROM arvore WHERE id=?`, baralhoID, candidatoID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) EditalDoMentor(ctx context.Context, mentorID, editalID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM editais WHERE id=? AND criado_por=? AND ativo=TRUE`, editalID, mentorID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) EditalDoAluno(ctx context.Context, alunoID, editalID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_editais WHERE aluno_id=? AND edital_id=? AND ativo=TRUE`, alunoID, editalID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) EditaisDoAluno(ctx context.Context, alunoID string) ([]string, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT edital_id FROM aluno_editais WHERE aluno_id=? AND ativo=TRUE`, alunoID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	resultado := []string{}
	for linhas.Next() {
		var id string
		if err = linhas.Scan(&id); err != nil {
			return nil, err
		}
		resultado = append(resultado, id)
	}
	return resultado, linhas.Err()
}

func (s *repositorioMySQL) listarCartoes(ctx context.Context, baralhoID, alunoID, concursoID string) ([]Cartao, error) {
	consulta := `SELECT c.id,c.baralho_id,c.tipo,c.frente,c.verso,c.dica,c.alternativas,c.resposta_correta,c.explicacao,c.explicacoes_alternativas,c.etiquetas,c.topico_id,c.subtopico_id,c.criado_por,
		COALESCE(r.repeticoes,0),COALESCE(r.intervalo_dias,0),COALESCE(r.facilidade,0),r.proxima_revisao,r.revisado_em,
		(SELECT COUNT(*) FROM revisoes_cartoes rf WHERE rf.aluno_id=? AND rf.cartao_id=c.id AND rf.qualidade=0`
	args := []any{alunoID}
	if concursoID != "" {
		consulta += ` AND rf.concurso_id=?`
		args = append(args, concursoID)
	}
	consulta += `) FROM cartoes_estudo c LEFT JOIN revisoes_cartoes r ON r.id=(SELECT r2.id FROM revisoes_cartoes r2 WHERE r2.aluno_id=? AND r2.cartao_id=c.id`
	args = append(args, alunoID)
	if concursoID != "" {
		consulta += ` AND r2.concurso_id=?`
		args = append(args, concursoID)
	}
	consulta += ` ORDER BY r2.revisado_em DESC,r2.id DESC LIMIT 1) WHERE c.baralho_id=? AND c.ativo=TRUE ORDER BY c.criado_em`
	args = append(args, baralhoID)

	linhas, err := s.db.QueryContext(ctx, consulta, args...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Cartao{}
	for linhas.Next() {
		var x Cartao
		var verso, dica, resposta, explicacao, topico, subtopico sql.NullString
		var alternativas, explicacoes, etiquetas []byte
		var proxima, ultima sql.NullTime
		if err = linhas.Scan(&x.ID, &x.BaralhoID, &x.Tipo, &x.Frente, &verso, &dica, &alternativas, &resposta, &explicacao, &explicacoes, &etiquetas, &topico, &subtopico, &x.CriadoPor, &x.Repeticoes, &x.IntervaloDias, &x.Facilidade, &proxima, &ultima, &x.Falhas); err != nil {
			return nil, err
		}
		x.Verso, x.Dica, x.RespostaCorreta, x.Explicacao, x.TopicoID, x.SubtopicoID = texto(verso), texto(dica), texto(resposta), texto(explicacao), texto(topico), texto(subtopico)
		if proxima.Valid {
			v := proxima.Time.Format("2006-01-02")
			x.ProximaRevisao = &v
		}
		if ultima.Valid {
			v := ultima.Time.Format(time.RFC3339)
			x.UltimaRevisao = &v
		}
		if len(alternativas) > 0 {
			x.Alternativas = alternativas
		}
		if len(explicacoes) > 0 {
			x.ExplicacoesAlternativas = explicacoes
		}
		if len(etiquetas) > 0 {
			x.Etiquetas = etiquetas
		}
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) ListarHistorico(ctx context.Context, alunoID, concursoID string) ([]HistoricoRevisao, error) {
	consulta := `SELECT r.id,r.cartao_id,c.baralho_id,r.qualidade,r.revisado_em,b.aluno_id,b.edital_id,b.alcance,b.tipo_proprietario,b.concurso_id,e.concurso_id FROM revisoes_cartoes r
		JOIN cartoes_estudo c ON c.id=r.cartao_id JOIN baralhos_cartoes b ON b.id=c.baralho_id
		LEFT JOIN editais e ON e.id=b.edital_id
		WHERE r.aluno_id=? AND c.ativo=TRUE AND b.ativo=TRUE`
	args := []any{alunoID}
	if concursoID != "" {
		consulta += ` AND r.concurso_id=?`
		args = append(args, concursoID)
	}
	consulta += ` ORDER BY r.revisado_em DESC,r.id DESC`
	linhas, err := s.db.QueryContext(ctx, consulta, args...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []HistoricoRevisao{}
	for linhas.Next() {
		var x HistoricoRevisao
		var alunoBaralho, editalBaralho, concursoBaralho, concursoEdital sql.NullString
		if err = linhas.Scan(&x.ID, &x.CartaoID, &x.BaralhoID, &x.Qualidade, &x.RevisadoEm, &alunoBaralho, &editalBaralho, &x.Alcance, &x.TipoProprietario, &concursoBaralho, &concursoEdital); err != nil {
			return nil, err
		}
		x.AlunoID, x.EditalID = texto(alunoBaralho), texto(editalBaralho)
		x.ConcursoID, x.EditalConcursoID = texto(concursoBaralho), texto(concursoEdital)
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) BaralhoDoCartao(ctx context.Context, cartaoID string) (string, error) {
	var id string
	err := s.db.QueryRowContext(ctx, `SELECT baralho_id FROM cartoes_estudo WHERE id=? AND ativo=TRUE`, cartaoID).Scan(&id)
	return id, err
}
func (s *repositorioMySQL) UltimaRevisao(ctx context.Context, alunoID, cartaoID, concursoID string) (EstadoRevisao, error) {
	consulta := `SELECT repeticoes,intervalo_dias,facilidade FROM revisoes_cartoes WHERE aluno_id=? AND cartao_id=?`
	args := []any{alunoID, cartaoID}
	if concursoID != "" {
		consulta += ` AND concurso_id=?`
		args = append(args, concursoID)
	}
	consulta += ` ORDER BY revisado_em DESC,id DESC LIMIT 1`
	var estado EstadoRevisao
	err := s.db.QueryRowContext(ctx, consulta, args...).Scan(&estado.Repeticoes, &estado.IntervaloDias, &estado.Facilidade)
	if errors.Is(err, sql.ErrNoRows) {
		return EstadoRevisao{}, nil
	}
	return estado, err
}
func (s *repositorioMySQL) RegistrarRevisao(ctx context.Context, alunoID, cartaoID, concursoID string, revisao Revisao) error {
	var concurso any
	if concursoID != "" {
		concurso = concursoID
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO revisoes_cartoes (id,aluno_id,concurso_id,cartao_id,qualidade,repeticoes,intervalo_dias,facilidade,proxima_revisao) VALUES (?,?,?,?,?,?,?,?,?)`, identificador.UUID(), alunoID, concurso, cartaoID, revisao.Qualidade, revisao.Repeticoes, revisao.IntervaloDias, revisao.Facilidade, revisao.ProximaRevisao)
	return err
}
func texto(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	x := v.String
	return &x
}
func jsonNulo(v json.RawMessage) any {
	if len(v) == 0 || string(v) == "null" {
		return nil
	}
	return []byte(v)
}
