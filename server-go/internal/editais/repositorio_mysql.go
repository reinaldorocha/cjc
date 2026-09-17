package editais

import (
	"context"
	"database/sql"
	"encoding/json"
	"track-concursos-web/internal/identificador"
)

type repositorioMySQL struct {
	db *sql.DB
}
type ProgressoMaterial struct {
	Concluido   bool    `json:"concluido"`
	ConcluidoEm *string `json:"concluidoEm,omitempty"`
}
type Subtopico struct {
	ID                 string                       `json:"id"`
	Nome               string                       `json:"nome"`
	Ordem              int                          `json:"ordem"`
	Peso               *float64                     `json:"peso,omitempty"`
	Relevancia         *int                         `json:"relevancia,omitempty"`
	Observacoes        *string                      `json:"observacoes,omitempty"`
	Materiais          json.RawMessage              `json:"materiais,omitempty"`
	ProgressoMateriais map[string]ProgressoMaterial `json:"progressoMateriais,omitempty"`
	Estudado           bool                         `json:"estudado"`
	ConcluidoEm        *string                      `json:"concluidoEm,omitempty"`
	ObservacoesAluno   *string                      `json:"observacoesAluno,omitempty"`
}
type Topico struct {
	ID                 string                       `json:"id"`
	Nome               string                       `json:"nome"`
	Ordem              int                          `json:"ordem"`
	Peso               *float64                     `json:"peso,omitempty"`
	Relevancia         *int                         `json:"relevancia,omitempty"`
	Observacoes        *string                      `json:"observacoes,omitempty"`
	Materiais          json.RawMessage              `json:"materiais,omitempty"`
	ProgressoMateriais map[string]ProgressoMaterial `json:"progressoMateriais,omitempty"`
	Estudado           bool                         `json:"estudado"`
	ConcluidoEm        *string                      `json:"concluidoEm,omitempty"`
	ObservacoesAluno   *string                      `json:"observacoesAluno,omitempty"`
	Subtopicos         []Subtopico                  `json:"subtopicos"`
}
type Materia struct {
	ID                 string                       `json:"id"`
	Nome               string                       `json:"nome"`
	Ordem              int                          `json:"ordem"`
	Peso               *float64                     `json:"peso,omitempty"`
	Relevancia         *int                         `json:"relevancia,omitempty"`
	Observacoes        *string                      `json:"observacoes,omitempty"`
	Materiais          json.RawMessage              `json:"materiais,omitempty"`
	ProgressoMateriais map[string]ProgressoMaterial `json:"progressoMateriais,omitempty"`
	Topicos            []Topico                     `json:"topicos"`
}
type Edital struct {
	ID         string    `json:"id"`
	ConcursoID string    `json:"concursoId"`
	Nome       string    `json:"nome"`
	Versao     *string   `json:"versao,omitempty"`
	Materias   []Materia `json:"materias"`
}
type Entrada struct {
	EditalID      string    `json:"editalId"`
	ConcursoID    string    `json:"concursoId"`
	Nome          string    `json:"nome"`
	Versao        *string   `json:"versao"`
	Materias      []Materia `json:"materias"`
	Novo          bool      `json:"-"`
	AtualizarNome bool      `json:"-"`
}
type AlteracaoItem struct {
	Tipo        string          `json:"tipo"`
	Nome        *string         `json:"nome"`
	Ordem       *int            `json:"ordem"`
	Peso        *float64        `json:"peso"`
	Relevancia  *int            `json:"relevancia"`
	Observacoes *string         `json:"observacoes"`
	Materiais   json.RawMessage `json:"materiais"`
	Ativo       *bool           `json:"ativo"`
}
type NovoItem struct {
	ID          string          `json:"-"`
	Tipo        string          `json:"tipo"`
	PaiID       string          `json:"paiId"`
	Nome        string          `json:"nome"`
	Ordem       int             `json:"ordem"`
	Peso        *float64        `json:"peso"`
	Relevancia  *int            `json:"relevancia"`
	Observacoes *string         `json:"observacoes"`
	Materiais   json.RawMessage `json:"materiais"`
}
type Progresso struct {
	Tipo        string  `json:"tipo"`
	Estudado    bool    `json:"estudado"`
	Observacoes *string `json:"observacoes"`
}
type EntradaProgressoMaterial struct {
	Tipo      string `json:"tipo"`
	ItemID    string `json:"itemId"`
	Concluido bool   `json:"concluido"`
}
type ItemOrdem struct {
	ID    string `json:"id"`
	Tipo  string `json:"tipo"`
	Ordem int    `json:"ordem"`
}
type ResumoCatalogo struct {
	ID         string  `json:"id"`
	ConcursoID string  `json:"concursoId"`
	Nome       string  `json:"nome"`
	Versao     *string `json:"versao,omitempty"`
}

func novoRepositorioMySQL(db *sql.DB) *repositorioMySQL { return &repositorioMySQL{db: db} }

func (s *repositorioMySQL) ConcursoDoMentor(ctx context.Context, mentorID, concursoID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM concursos WHERE id=? AND criado_por=? AND ativo=TRUE`, concursoID, mentorID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) ConcursoDoAluno(ctx context.Context, alunoID, concursoID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND concurso_id=? AND ativo=TRUE`, alunoID, concursoID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) EditalDoMentor(ctx context.Context, mentorID, editalID, concursoID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM editais WHERE id=? AND criado_por=? AND concurso_id=? AND ativo=TRUE`, editalID, mentorID, concursoID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) ItemDoAluno(ctx context.Context, alunoID, tipo, itemID string) (bool, error) {
	var consulta string
	switch tipo {
	case "materia":
		consulta = `SELECT COUNT(*) FROM aluno_editais ae JOIN edital_materias em ON em.edital_id=ae.edital_id AND em.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND em.id=?`
	case "topico":
		consulta = `SELECT COUNT(*) FROM aluno_editais ae JOIN edital_materias em ON em.edital_id=ae.edital_id JOIN edital_topicos et ON et.materia_id=em.id AND et.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND et.id=?`
	case "subtopico":
		consulta = `SELECT COUNT(*) FROM aluno_editais ae JOIN edital_materias em ON em.edital_id=ae.edital_id JOIN edital_topicos et ON et.materia_id=em.id JOIN edital_subtopicos es ON es.topico_id=et.id AND es.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND es.id=?`
	default:
		return false, nil
	}
	var n int
	err := s.db.QueryRowContext(ctx, consulta, alunoID, itemID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) PaiDoAluno(ctx context.Context, alunoID, tipo, paiID string) (bool, error) {
	var consulta string
	switch tipo {
	case "materia":
		consulta = `SELECT COUNT(*) FROM aluno_editais WHERE aluno_id=? AND edital_id=? AND ativo=TRUE`
	case "topico":
		consulta = `SELECT COUNT(*) FROM aluno_editais ae JOIN edital_materias em ON em.edital_id=ae.edital_id AND em.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND em.id=?`
	case "subtopico":
		consulta = `SELECT COUNT(*) FROM aluno_editais ae JOIN edital_materias em ON em.edital_id=ae.edital_id JOIN edital_topicos et ON et.materia_id=em.id AND et.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND et.id=?`
	default:
		return false, nil
	}
	var n int
	err := s.db.QueryRowContext(ctx, consulta, alunoID, paiID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) TipoDoItem(ctx context.Context, itemID string) (string, error) {
	var tipo string
	err := s.db.QueryRowContext(ctx, `SELECT tipo FROM (SELECT id,'topico' tipo FROM edital_topicos WHERE ativo=TRUE UNION ALL SELECT id,'subtopico' tipo FROM edital_subtopicos WHERE ativo=TRUE) itens WHERE id=? LIMIT 1`, itemID).Scan(&tipo)
	return tipo, err
}

func (s *repositorioMySQL) Catalogar(ctx context.Context, mentorID, concursoID string) ([]ResumoCatalogo, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT id,concurso_id,nome,versao FROM editais WHERE criado_por=? AND concurso_id=? AND ativo=TRUE ORDER BY atualizado_em DESC,nome`, mentorID, concursoID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []ResumoCatalogo{}
	for linhas.Next() {
		var x ResumoCatalogo
		var versao sql.NullString
		if err = linhas.Scan(&x.ID, &x.ConcursoID, &x.Nome, &versao); err != nil {
			return nil, err
		}
		x.Versao = texto(versao)
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) ListarCatalogo(ctx context.Context, mentorID, concursoID string) ([]Edital, error) {
	q := `SELECT id,concurso_id,nome,versao FROM editais WHERE criado_por=? AND ativo=TRUE`
	args := []any{mentorID}
	if concursoID != "" {
		q += ` AND concurso_id=?`
		args = append(args, concursoID)
	}
	q += ` ORDER BY atualizado_em DESC,nome`
	linhas, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Edital{}
	for linhas.Next() {
		var e Edital
		var v sql.NullString
		if err = linhas.Scan(&e.ID, &e.ConcursoID, &e.Nome, &v); err != nil {
			return nil, err
		}
		e.Versao = texto(v)
		e.Materias, err = s.carregarMaterias(ctx, "", e.ID)
		if err != nil {
			return nil, err
		}
		lista = append(lista, e)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) CriarCatalogo(ctx context.Context, mentorID string, e Entrada) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	id := e.EditalID
	if _, err = tx.ExecContext(ctx, `INSERT INTO editais(id,concurso_id,nome,versao,criado_por) VALUES(?,?,?,?,?)`, id, e.ConcursoID, e.Nome, e.Versao, mentorID); err != nil {
		return "", err
	}
	for _, m := range e.Materias {
		mid := m.ID
		if _, err = tx.ExecContext(ctx, `INSERT INTO edital_materias(id,edital_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES(?,?,?,?,?,?,?,?)`, mid, id, m.Nome, m.Ordem, m.Peso, m.Relevancia, m.Observacoes, jsonNulo(m.Materiais)); err != nil {
			return "", err
		}
		for _, t := range m.Topicos {
			tid := t.ID
			if _, err = tx.ExecContext(ctx, `INSERT INTO edital_topicos(id,materia_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES(?,?,?,?,?,?,?,?)`, tid, mid, t.Nome, t.Ordem, t.Peso, t.Relevancia, t.Observacoes, jsonNulo(t.Materiais)); err != nil {
				return "", err
			}
			for _, sub := range t.Subtopicos {
				if _, err = tx.ExecContext(ctx, `INSERT INTO edital_subtopicos(id,topico_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES(?,?,?,?,?,?,?,?)`, sub.ID, tid, sub.Nome, sub.Ordem, sub.Peso, sub.Relevancia, sub.Observacoes, jsonNulo(sub.Materiais)); err != nil {
					return "", err
				}
			}
		}
	}
	return id, tx.Commit()
}

func (s *repositorioMySQL) Listar(ctx context.Context, alunoID, concursoID string) ([]Edital, error) {
	consulta := `SELECT e.id,e.concurso_id,e.nome,e.versao FROM aluno_editais ae JOIN editais e ON e.id=ae.edital_id JOIN aluno_concursos ac ON ac.aluno_id=ae.aluno_id AND ac.concurso_id=e.concurso_id AND ac.ativo=TRUE WHERE ae.aluno_id=? AND ae.ativo=TRUE AND e.ativo=TRUE`
	args := []any{alunoID}
	if concursoID != "" {
		consulta += ` AND e.concurso_id=?`
		args = append(args, concursoID)
	}
	consulta += ` ORDER BY e.criado_em DESC`
	linhas, err := s.db.QueryContext(ctx, consulta, args...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Edital{}
	for linhas.Next() {
		var e Edital
		var versao sql.NullString
		if err = linhas.Scan(&e.ID, &e.ConcursoID, &e.Nome, &versao); err != nil {
			return nil, err
		}
		e.Versao = texto(versao)
		e.Materias, err = s.carregarMaterias(ctx, alunoID, e.ID)
		if err != nil {
			return nil, err
		}
		lista = append(lista, e)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) Atribuir(ctx context.Context, alunoID, executor string, e Entrada) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	id := e.EditalID
	if !e.Novo {
		if e.AtualizarNome {
			if _, err = tx.ExecContext(ctx, `UPDATE editais SET nome=?,versao=? WHERE id=?`, e.Nome, e.Versao, id); err != nil {
				return "", err
			}
		}
	} else {
		_, err = tx.ExecContext(ctx, `INSERT INTO editais (id,concurso_id,nome,versao,criado_por) VALUES (?,?,?,?,?)`, id, e.ConcursoID, e.Nome, e.Versao, executor)
		if err != nil {
			return "", err
		}
		for _, m := range e.Materias {
			materiaID := m.ID
			_, err = tx.ExecContext(ctx, `INSERT INTO edital_materias (id,edital_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES (?,?,?,?,?,?,?,?)`, materiaID, id, m.Nome, m.Ordem, m.Peso, m.Relevancia, m.Observacoes, jsonNulo(m.Materiais))
			if err != nil {
				return "", err
			}
			for _, t := range m.Topicos {
				topicoID := t.ID
				_, err = tx.ExecContext(ctx, `INSERT INTO edital_topicos (id,materia_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES (?,?,?,?,?,?,?,?)`, topicoID, materiaID, t.Nome, t.Ordem, t.Peso, t.Relevancia, t.Observacoes, jsonNulo(t.Materiais))
				if err != nil {
					return "", err
				}
				for _, sub := range t.Subtopicos {
					_, err = tx.ExecContext(ctx, `INSERT INTO edital_subtopicos (id,topico_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES (?,?,?,?,?,?,?,?)`, sub.ID, topicoID, sub.Nome, sub.Ordem, sub.Peso, sub.Relevancia, sub.Observacoes, jsonNulo(sub.Materiais))
					if err != nil {
						return "", err
					}
				}
			}
		}
	}
	var jaAtribuido int
	if _, err = tx.ExecContext(ctx, `UPDATE aluno_editais ae JOIN editais atual ON atual.id=ae.edital_id SET ae.ativo=FALSE WHERE ae.aluno_id=? AND atual.concurso_id=? AND ae.edital_id<>? AND ae.ativo=TRUE`, alunoID, e.ConcursoID, id); err != nil {
		return "", err
	}
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM aluno_editais WHERE aluno_id=? AND edital_id=? AND ativo=TRUE`, alunoID, id).Scan(&jaAtribuido); err != nil {
		return "", err
	}
	if jaAtribuido == 0 {
		_, err = tx.ExecContext(ctx, `INSERT INTO aluno_editais (id,aluno_id,edital_id,atribuido_por) VALUES (?,?,?,?)`, identificador.UUID(), alunoID, id, executor)
	}
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (s *repositorioMySQL) CriarItem(ctx context.Context, alunoID string, e NovoItem) (string, error) {
	id := e.ID
	var consulta string
	var args []any
	switch e.Tipo {
	case "materia":
		consulta = `INSERT INTO edital_materias (id,edital_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES (?,?,?,?,?,?,?,?)`
		args = []any{id, e.PaiID, e.Nome, e.Ordem, e.Peso, e.Relevancia, e.Observacoes, jsonNulo(e.Materiais)}
	case "topico":
		consulta = `INSERT INTO edital_topicos (id,materia_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES (?,?,?,?,?,?,?,?)`
		args = []any{id, e.PaiID, e.Nome, e.Ordem, e.Peso, e.Relevancia, e.Observacoes, jsonNulo(e.Materiais)}
	case "subtopico":
		consulta = `INSERT INTO edital_subtopicos (id,topico_id,nome,ordem,peso,relevancia,observacoes,materiais) VALUES (?,?,?,?,?,?,?,?)`
		args = []any{id, e.PaiID, e.Nome, e.Ordem, e.Peso, e.Relevancia, e.Observacoes, jsonNulo(e.Materiais)}
	}
	_, err := s.db.ExecContext(ctx, consulta, args...)
	return id, err
}

func (s *repositorioMySQL) AlterarItem(ctx context.Context, alunoID, itemID string, e AlteracaoItem) error {
	tabela := ""
	switch e.Tipo {
	case "materia":
		tabela = "edital_materias"
	case "topico":
		tabela = "edital_topicos"
	case "subtopico":
		tabela = "edital_subtopicos"
	}
	if e.Nome != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE `+tabela+` SET nome=? WHERE id=?`, *e.Nome, itemID); err != nil {
			return err
		}
	}
	if e.Ordem != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE `+tabela+` SET ordem=? WHERE id=?`, *e.Ordem, itemID); err != nil {
			return err
		}
	}
	if e.Peso != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE `+tabela+` SET peso=? WHERE id=?`, *e.Peso, itemID); err != nil {
			return err
		}
	}
	if e.Relevancia != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE `+tabela+` SET relevancia=? WHERE id=?`, *e.Relevancia, itemID); err != nil {
			return err
		}
	}
	if e.Observacoes != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE `+tabela+` SET observacoes=NULLIF(?,'') WHERE id=?`, *e.Observacoes, itemID); err != nil {
			return err
		}
	}
	if e.Ativo != nil {
		if _, err := s.db.ExecContext(ctx, `UPDATE `+tabela+` SET ativo=? WHERE id=?`, *e.Ativo, itemID); err != nil {
			return err
		}
	}
	if len(e.Materiais) > 0 {
		if _, err := s.db.ExecContext(ctx, `UPDATE `+tabela+` SET materiais=? WHERE id=?`, e.Materiais, itemID); err != nil {
			return err
		}
	}
	return nil
}

func (s *repositorioMySQL) Reordenar(ctx context.Context, alunoID string, itens []ItemOrdem) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, item := range itens {
		tabela := tabelaItem(item.Tipo)
		if _, err = tx.ExecContext(ctx, `UPDATE `+tabela+` SET ordem=? WHERE id=?`, item.Ordem, item.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *repositorioMySQL) AtualizarProgressoMaterial(ctx context.Context, alunoID, materialID string, e EntradaProgressoMaterial) error {
	concluido := "NULL"
	if e.Concluido {
		concluido = "UTC_TIMESTAMP()"
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO aluno_progresso_materiais (id,aluno_id,tipo_item,item_id,material_id,concluido,concluido_em) VALUES (?,?,?,?,?, ?,`+concluido+`) ON DUPLICATE KEY UPDATE concluido=VALUES(concluido),concluido_em=`+concluido, identificador.UUID(), alunoID, e.Tipo, e.ItemID, materialID, e.Concluido)
	return err
}

func tabelaItem(tipo string) string {
	switch tipo {
	case "materia":
		return "edital_materias"
	case "topico":
		return "edital_topicos"
	case "subtopico":
		return "edital_subtopicos"
	default:
		return ""
	}
}

func (s *repositorioMySQL) AtualizarProgresso(ctx context.Context, alunoID, itemID string, e Progresso) error {
	var err error
	concluido := "NULL"
	if e.Estudado {
		concluido = "UTC_TIMESTAMP()"
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO aluno_progresso_edital (id,aluno_id,tipo_item,item_id,estudado,concluido_em,observacoes) VALUES (?,?,?,?,?,`+concluido+`,?) ON DUPLICATE KEY UPDATE estudado=VALUES(estudado),concluido_em=`+concluido+`,observacoes=VALUES(observacoes)`, identificador.UUID(), alunoID, e.Tipo, itemID, e.Estudado, e.Observacoes)
	if err != nil {
		return err
	}

	return nil
}

func (s *repositorioMySQL) ContextoRevisao(ctx context.Context, tipo, itemID string) (ContextoRevisao, error) {
	var resultado ContextoRevisao
	if tipo == "topico" {
		resultado.TopicoID = itemID
		err := s.db.QueryRowContext(ctx, `SELECT e.concurso_id,em.id,COALESCE(c.prazos_revisao,'1,7,30') FROM edital_topicos et JOIN edital_materias em ON em.id=et.materia_id JOIN editais e ON e.id=em.edital_id JOIN concursos c ON c.id=e.concurso_id WHERE et.id=?`, itemID).Scan(&resultado.ConcursoID, &resultado.MateriaID, &resultado.Prazos)
		return resultado, err
	}
	resultado.SubtopicoID = &itemID
	err := s.db.QueryRowContext(ctx, `SELECT e.concurso_id,em.id,et.id,COALESCE(c.prazos_revisao,'1,7,30') FROM edital_subtopicos es JOIN edital_topicos et ON et.id=es.topico_id JOIN edital_materias em ON em.id=et.materia_id JOIN editais e ON e.id=em.edital_id JOIN concursos c ON c.id=e.concurso_id WHERE es.id=?`, itemID).Scan(&resultado.ConcursoID, &resultado.MateriaID, &resultado.TopicoID, &resultado.Prazos)
	return resultado, err
}

func (s *repositorioMySQL) SincronizarRevisoes(ctx context.Context, alunoID, tipo, itemID string, contexto *ContextoRevisao, revisoes []RevisaoAutomatica) error {
	if contexto == nil {
		if tipo == "topico" {
			_, err := s.db.ExecContext(ctx, `UPDATE revisoes_programadas SET ativo=FALSE WHERE aluno_id=? AND topico_id=? AND subtopico_id IS NULL AND concluida=FALSE`, alunoID, itemID)
			return err
		}
		_, err := s.db.ExecContext(ctx, `UPDATE revisoes_programadas SET ativo=FALSE WHERE aluno_id=? AND subtopico_id=? AND concluida=FALSE`, alunoID, itemID)
		return err
	}
	for _, revisao := range revisoes {
		var existente string
		err := s.db.QueryRowContext(ctx, `SELECT id FROM revisoes_programadas WHERE aluno_id=? AND concurso_id<=>? AND materia_id<=>? AND topico_id<=>? AND subtopico_id<=>? AND ciclo_atual=? AND ativo=TRUE LIMIT 1`, alunoID, contexto.ConcursoID, contexto.MateriaID, contexto.TopicoID, contexto.SubtopicoID, revisao.Ciclo).Scan(&existente)
		if err == nil {
			if _, err = s.db.ExecContext(ctx, `UPDATE revisoes_programadas SET proxima_data=?,concluida=FALSE,ativo=TRUE WHERE id=?`, revisao.ProximaData, existente); err != nil {
				return err
			}
			continue
		}
		if err != sql.ErrNoRows {
			return err
		}
		if _, err = s.db.ExecContext(ctx, `INSERT INTO revisoes_programadas (id,aluno_id,concurso_id,materia_id,topico_id,subtopico_id,ciclo_atual,proxima_data,concluida) VALUES (?,?,?,?,?,?,?,?,FALSE)`, identificador.UUID(), alunoID, contexto.ConcursoID, contexto.MateriaID, contexto.TopicoID, contexto.SubtopicoID, revisao.Ciclo, revisao.ProximaData); err != nil {
			return err
		}
	}
	return nil
}

func (s *repositorioMySQL) carregarMaterias(ctx context.Context, alunoID, editalID string) ([]Materia, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT id,nome,ordem,peso,relevancia,observacoes,materiais FROM edital_materias WHERE edital_id=? AND ativo=TRUE ORDER BY ordem,nome`, editalID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Materia{}
	for linhas.Next() {
		var m Materia
		var peso sql.NullFloat64
		var rel sql.NullInt64
		var obs, materiais sql.NullString
		if err = linhas.Scan(&m.ID, &m.Nome, &m.Ordem, &peso, &rel, &obs, &materiais); err != nil {
			return nil, err
		}
		m.Peso = decimal(peso)
		m.Relevancia = inteiro(rel)
		m.Observacoes = texto(obs)
		m.Materiais = jsonBruto(materiais)
		m.ProgressoMateriais, err = s.carregarProgressoMateriais(ctx, alunoID, "materia", m.ID)
		if err != nil {
			return nil, err
		}
		m.Topicos, err = s.carregarTopicos(ctx, alunoID, m.ID)
		if err != nil {
			return nil, err
		}
		lista = append(lista, m)
	}
	return lista, linhas.Err()
}
func (s *repositorioMySQL) carregarTopicos(ctx context.Context, alunoID, materiaID string) ([]Topico, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT et.id,et.nome,et.ordem,et.peso,et.relevancia,et.observacoes,et.materiais,COALESCE(ap.estudado,FALSE),ap.concluido_em,ap.observacoes FROM edital_topicos et LEFT JOIN aluno_progresso_edital ap ON ap.aluno_id=? AND ap.tipo_item='topico' AND ap.item_id=et.id WHERE et.materia_id=? AND et.ativo=TRUE ORDER BY et.ordem,et.nome`, alunoID, materiaID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Topico{}
	for linhas.Next() {
		var t Topico
		var peso sql.NullFloat64
		var rel sql.NullInt64
		var obs, materiais, obsAluno sql.NullString
		var concluido sql.NullTime
		if err = linhas.Scan(&t.ID, &t.Nome, &t.Ordem, &peso, &rel, &obs, &materiais, &t.Estudado, &concluido, &obsAluno); err != nil {
			return nil, err
		}
		t.Peso = decimal(peso)
		t.Relevancia = inteiro(rel)
		t.Observacoes = texto(obs)
		t.Materiais = jsonBruto(materiais)
		t.ProgressoMateriais, err = s.carregarProgressoMateriais(ctx, alunoID, "topico", t.ID)
		if err != nil {
			return nil, err
		}
		t.ConcluidoEm = dataHora(concluido)
		t.ObservacoesAluno = texto(obsAluno)
		t.Subtopicos, err = s.carregarSubtopicos(ctx, alunoID, t.ID)
		if err != nil {
			return nil, err
		}
		lista = append(lista, t)
	}
	return lista, linhas.Err()
}
func (s *repositorioMySQL) carregarSubtopicos(ctx context.Context, alunoID, topicoID string) ([]Subtopico, error) {
	linhas, err := s.db.QueryContext(ctx, `SELECT es.id,es.nome,es.ordem,es.peso,es.relevancia,es.observacoes,es.materiais,COALESCE(ap.estudado,FALSE),ap.concluido_em,ap.observacoes FROM edital_subtopicos es LEFT JOIN aluno_progresso_edital ap ON ap.aluno_id=? AND ap.tipo_item='subtopico' AND ap.item_id=es.id WHERE es.topico_id=? AND es.ativo=TRUE ORDER BY es.ordem,es.nome`, alunoID, topicoID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Subtopico{}
	for linhas.Next() {
		var x Subtopico
		var peso sql.NullFloat64
		var rel sql.NullInt64
		var obs, materiais, obsAluno sql.NullString
		var concluido sql.NullTime
		if err = linhas.Scan(&x.ID, &x.Nome, &x.Ordem, &peso, &rel, &obs, &materiais, &x.Estudado, &concluido, &obsAluno); err != nil {
			return nil, err
		}
		x.Peso = decimal(peso)
		x.Relevancia = inteiro(rel)
		x.Observacoes = texto(obs)
		x.Materiais = jsonBruto(materiais)
		x.ProgressoMateriais, err = s.carregarProgressoMateriais(ctx, alunoID, "subtopico", x.ID)
		if err != nil {
			return nil, err
		}
		x.ConcluidoEm = dataHora(concluido)
		x.ObservacoesAluno = texto(obsAluno)
		lista = append(lista, x)
	}
	return lista, linhas.Err()
}
func (s *repositorioMySQL) carregarProgressoMateriais(ctx context.Context, alunoID, tipo, itemID string) (map[string]ProgressoMaterial, error) {
	if alunoID == "" {
		return map[string]ProgressoMaterial{}, nil
	}
	linhas, err := s.db.QueryContext(ctx, `SELECT material_id,concluido,concluido_em FROM aluno_progresso_materiais WHERE aluno_id=? AND tipo_item=? AND item_id=?`, alunoID, tipo, itemID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	resultado := map[string]ProgressoMaterial{}
	for linhas.Next() {
		var id string
		var p ProgressoMaterial
		var data sql.NullTime
		if err = linhas.Scan(&id, &p.Concluido, &data); err != nil {
			return nil, err
		}
		p.ConcluidoEm = dataHora(data)
		resultado[id] = p
	}
	return resultado, linhas.Err()
}
func texto(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}
func decimal(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	return &v.Float64
}
func inteiro(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	x := int(v.Int64)
	return &x
}
func dataHora(v sql.NullTime) *string {
	if !v.Valid {
		return nil
	}
	x := v.Time.UTC().Format("2006-01-02T15:04:05Z")
	return &x
}
func jsonBruto(v sql.NullString) json.RawMessage {
	if !v.Valid {
		return nil
	}
	return json.RawMessage(v.String)
}
func jsonNulo(v json.RawMessage) any {
	if len(v) == 0 {
		return nil
	}
	return v
}
