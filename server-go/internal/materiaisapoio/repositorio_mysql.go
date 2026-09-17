package materiaisapoio

import (
	"context"
	"database/sql"
)

type Material struct {
	ID             string  `json:"id"`
	Titulo         string  `json:"titulo"`
	Descricao      *string `json:"descricao,omitempty"`
	Tipo           string  `json:"tipo"`
	Url            *string `json:"url,omitempty"`
	Texto          *string `json:"texto,omitempty"`
	Escopo         string  `json:"escopo"`
	EditalID       *string `json:"editalId,omitempty"`
	ArquivoNome    *string `json:"arquivoNome,omitempty"`
	ArquivoMime    *string `json:"arquivoMime,omitempty"`
	ArquivoCaminho *string `json:"-"`
	Pasta          *string `json:"pasta,omitempty"`
}

type EntradaMaterial struct {
	Titulo    string  `json:"titulo"`
	Descricao *string `json:"descricao"`
	Tipo      string  `json:"tipo"`
	Url       *string `json:"url"`
	Texto     *string `json:"texto"`
	Escopo    string  `json:"escopo"`
	EditalID  *string `json:"editalId"`
	Pasta     *string `json:"pasta"`
}

var tiposValidos = map[string]bool{"arquivo": true, "youtube": true, "texto": true, "link": true}
var escoposValidos = map[string]bool{"global": true, "edital": true}

type repositorioMySQL struct {
	db *sql.DB
}

func novoRepositorioMySQL(db *sql.DB) *repositorioMySQL { return &repositorioMySQL{db: db} }

func (s *repositorioMySQL) EditalDoMentor(ctx context.Context, mentorID, editalID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM editais WHERE id=? AND criado_por=? AND ativo=TRUE`, editalID, mentorID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) MaterialDoMentor(ctx context.Context, mentorID, materialID string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM materiais_apoio WHERE id=? AND mentor_id=? AND ativo=TRUE`, materialID, mentorID).Scan(&n)
	return n > 0, err
}

func (s *repositorioMySQL) CriarArquivoRegistro(ctx context.Context, mentorID, id string, e EntradaArquivo, caminho, tipoMime string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO materiais_apoio (id,mentor_id,titulo,descricao,tipo,url,texto,escopo,edital_id,arquivo_nome,arquivo_caminho,arquivo_mime,pasta) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, id, mentorID, e.Titulo, textoNuloRaw(e.Descricao), e.Tipo, nil, nil, e.Escopo, e.EditalID, e.NomeArquivo, caminho, tipoMime, textoNuloRaw(e.Pasta))
	return err
}

func (s *repositorioMySQL) AlterarArquivoRegistro(ctx context.Context, mentorID, materialID string, e EntradaArquivo, caminho, tipoMime string) (*string, error) {
	var anterior sql.NullString
	if err := s.db.QueryRowContext(ctx, `SELECT arquivo_caminho FROM materiais_apoio WHERE id=? AND mentor_id=? AND ativo=TRUE`, materialID, mentorID).Scan(&anterior); err != nil {
		return nil, err
	}
	if caminho == "" {
		_, err := s.db.ExecContext(ctx, `UPDATE materiais_apoio SET titulo=?,descricao=?,tipo=?,url=NULL,texto=NULL,escopo=?,edital_id=?,pasta=? WHERE id=?`, e.Titulo, textoNuloRaw(e.Descricao), e.Tipo, e.Escopo, e.EditalID, textoNuloRaw(e.Pasta), materialID)
		return textoNulo(anterior), err
	}
	_, err := s.db.ExecContext(ctx, `UPDATE materiais_apoio SET titulo=?,descricao=?,tipo=?,url=NULL,texto=NULL,escopo=?,edital_id=?,arquivo_nome=?,arquivo_caminho=?,arquivo_mime=?,pasta=? WHERE id=?`, e.Titulo, textoNuloRaw(e.Descricao), e.Tipo, e.Escopo, e.EditalID, e.NomeArquivo, caminho, tipoMime, textoNuloRaw(e.Pasta), materialID)
	return textoNulo(anterior), err
}

func (s *repositorioMySQL) ListarMentor(ctx context.Context, mentorID, editalID string) ([]Material, error) {
	consulta := `SELECT id,titulo,descricao,tipo,url,texto,escopo,edital_id,arquivo_nome,arquivo_caminho,arquivo_mime,pasta FROM materiais_apoio WHERE mentor_id=? AND ativo=TRUE`
	args := []any{mentorID}
	if editalID != "" {
		consulta += ` AND (escopo='global' OR edital_id=?)`
		args = append(args, editalID)
	}
	consulta += ` ORDER BY criado_em DESC,titulo`
	return s.listar(ctx, consulta, args...)
}

func (s *repositorioMySQL) ListarAluno(ctx context.Context, alunoID, editalID string) ([]Material, error) {
	consulta := `SELECT m.id,m.titulo,m.descricao,m.tipo,m.url,m.texto,m.escopo,m.edital_id,m.arquivo_nome,m.arquivo_caminho,m.arquivo_mime,m.pasta FROM materiais_apoio m WHERE m.ativo=TRUE AND (m.escopo='global'`
	args := []any{}
	if editalID != "" {
		consulta += ` OR (m.escopo='edital' AND m.edital_id=? AND EXISTS (SELECT 1 FROM aluno_editais ae WHERE ae.aluno_id=? AND ae.edital_id=m.edital_id AND ae.ativo=TRUE)))`
		args = append(args, editalID, alunoID)
	} else {
		consulta += `)`
	}
	consulta += ` ORDER BY m.criado_em DESC,m.titulo`
	return s.listar(ctx, consulta, args...)
}

func (s *repositorioMySQL) listar(ctx context.Context, consulta string, args ...any) ([]Material, error) {
	linhas, err := s.db.QueryContext(ctx, consulta, args...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Material{}
	for linhas.Next() {
		var m Material
		var descricao, url, texto, editalID, arquivoNome, arquivoCaminho, arquivoMime, pasta sql.NullString
		if err = linhas.Scan(&m.ID, &m.Titulo, &descricao, &m.Tipo, &url, &texto, &m.Escopo, &editalID, &arquivoNome, &arquivoCaminho, &arquivoMime, &pasta); err != nil {
			return nil, err
		}
		m.Descricao, m.Url, m.Texto, m.EditalID = textoNulo(descricao), textoNulo(url), textoNulo(texto), textoNulo(editalID)
		m.ArquivoNome, m.ArquivoCaminho, m.ArquivoMime, m.Pasta = textoNulo(arquivoNome), textoNulo(arquivoCaminho), textoNulo(arquivoMime), textoNulo(pasta)
		lista = append(lista, m)
	}
	return lista, linhas.Err()
}

func (s *repositorioMySQL) Criar(ctx context.Context, mentorID, id string, e EntradaMaterial) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO materiais_apoio (id,mentor_id,titulo,descricao,tipo,url,texto,escopo,edital_id,pasta) VALUES (?,?,?,?,?,?,?,?,?,?)`, id, mentorID, e.Titulo, textoNuloRaw(e.Descricao), e.Tipo, textoNuloRaw(e.Url), textoNuloRaw(e.Texto), e.Escopo, e.EditalID, textoNuloRaw(e.Pasta))
	return err
}

func (s *repositorioMySQL) ObterParaMentor(ctx context.Context, mentorID, materialID string) (*Material, error) {
	var m Material
	var descricao, url, texto, editalID, arquivoNome, arquivoCaminho, arquivoMime, pasta sql.NullString
	if err := s.db.QueryRowContext(ctx, `SELECT id,titulo,descricao,tipo,url,texto,escopo,edital_id,arquivo_nome,arquivo_caminho,arquivo_mime,pasta FROM materiais_apoio WHERE id=? AND mentor_id=? AND ativo=TRUE`, materialID, mentorID).Scan(&m.ID, &m.Titulo, &descricao, &m.Tipo, &url, &texto, &m.Escopo, &editalID, &arquivoNome, &arquivoCaminho, &arquivoMime, &pasta); err != nil {
		return nil, err
	}
	m.Descricao = textoNulo(descricao)
	m.Url = textoNulo(url)
	m.Texto = textoNulo(texto)
	m.EditalID = textoNulo(editalID)
	m.ArquivoNome = textoNulo(arquivoNome)
	m.ArquivoCaminho = textoNulo(arquivoCaminho)
	m.ArquivoMime = textoNulo(arquivoMime)
	m.Pasta = textoNulo(pasta)
	return &m, nil
}

func (s *repositorioMySQL) ObterParaAluno(ctx context.Context, alunoID, materialID string) (*Material, error) {
	var m Material
	var descricao, url, texto, editalID, arquivoNome, arquivoCaminho, arquivoMime, pasta sql.NullString
	consulta := `SELECT m.id,m.titulo,m.descricao,m.tipo,m.url,m.texto,m.escopo,m.edital_id,m.arquivo_nome,m.arquivo_caminho,m.arquivo_mime,m.pasta FROM materiais_apoio m WHERE m.id=? AND m.ativo=TRUE AND (m.escopo='global'`
	args := []any{materialID}
	consulta += ` OR (m.escopo='edital' AND m.edital_id IN (SELECT edital_id FROM aluno_editais ae WHERE ae.aluno_id=? AND ae.ativo=TRUE)))`
	args = append(args, alunoID)
	if err := s.db.QueryRowContext(ctx, consulta, args...).Scan(&m.ID, &m.Titulo, &descricao, &m.Tipo, &url, &texto, &m.Escopo, &editalID, &arquivoNome, &arquivoCaminho, &arquivoMime, &pasta); err != nil {
		return nil, err
	}
	m.Descricao = textoNulo(descricao)
	m.Url = textoNulo(url)
	m.Texto = textoNulo(texto)
	m.EditalID = textoNulo(editalID)
	m.ArquivoNome = textoNulo(arquivoNome)
	m.ArquivoCaminho = textoNulo(arquivoCaminho)
	m.ArquivoMime = textoNulo(arquivoMime)
	m.Pasta = textoNulo(pasta)
	return &m, nil
}

func (s *repositorioMySQL) Alterar(ctx context.Context, mentorID, materialID string, e EntradaMaterial) (bool, error) {
	if e.Tipo == "arquivo" {
		resultado, err := s.db.ExecContext(ctx, `UPDATE materiais_apoio SET titulo=?,descricao=?,tipo=?,url=?,texto=?,escopo=?,edital_id=?,pasta=? WHERE id=? AND mentor_id=? AND ativo=TRUE`, e.Titulo, textoNuloRaw(e.Descricao), e.Tipo, nil, nil, e.Escopo, e.EditalID, textoNuloRaw(e.Pasta), materialID, mentorID)
		if err != nil {
			return false, err
		}
		n, err := resultado.RowsAffected()
		return n > 0, err
	}
	resultado, err := s.db.ExecContext(ctx, `UPDATE materiais_apoio SET titulo=?,descricao=?,tipo=?,url=?,texto=?,escopo=?,edital_id=?,arquivo_nome=NULL,arquivo_caminho=NULL,arquivo_mime=NULL,pasta=? WHERE id=? AND mentor_id=? AND ativo=TRUE`, e.Titulo, textoNuloRaw(e.Descricao), e.Tipo, textoNuloRaw(e.Url), textoNuloRaw(e.Texto), e.Escopo, e.EditalID, textoNuloRaw(e.Pasta), materialID, mentorID)
	if err != nil {
		return false, err
	}
	n, err := resultado.RowsAffected()
	return n > 0, err
}

func (s *repositorioMySQL) Desativar(ctx context.Context, mentorID, materialID string) (bool, error) {
	res, err := s.db.ExecContext(ctx, `UPDATE materiais_apoio SET ativo=FALSE WHERE id=? AND mentor_id=? AND ativo=TRUE`, materialID, mentorID)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func textoNulo(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}

func textoNuloRaw(v *string) interface{} {
	if v == nil {
		return nil
	}
	return *v
}
