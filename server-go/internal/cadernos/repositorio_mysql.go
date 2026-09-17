package cadernos

import (
	"context"
	"database/sql"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(banco *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: banco} }

func (r *repositorioMySQL) Criar(ctx context.Context, id, alunoID string, e EntradaCadernoNormalizada) error {
	_, err := r.banco.ExecContext(ctx, `INSERT INTO cadernos (id, aluno_id, titulo, pasta, conteudo, edital_id, materia_id, topico_id, cor) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, alunoID, e.Titulo, e.Pasta, e.Conteudo, e.EditalID, e.MateriaID, e.TopicoID, e.Cor)
	return err
}
func (r *repositorioMySQL) Listar(ctx context.Context, alunoID, editalID string) ([]Caderno, error) {
	consulta := `SELECT id, aluno_id, titulo, pasta, conteudo, edital_id, materia_id, topico_id, cor, ativo, criado_em, atualizado_em FROM cadernos WHERE aluno_id = ? AND ativo = TRUE`
	args := []any{alunoID}
	if editalID != "" {
		consulta += ` AND (edital_id IS NULL OR edital_id = ?)`
		args = append(args, editalID)
	}
	linhas, err := r.banco.QueryContext(ctx, consulta+` ORDER BY atualizado_em DESC, criado_em DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := []Caderno{}
	for linhas.Next() {
		var c Caderno
		var edital, materia, topico sql.NullString
		if err := linhas.Scan(&c.ID, &c.AlunoID, &c.Titulo, &c.Pasta, &c.Conteudo, &edital, &materia, &topico, &c.Cor, &c.Ativo, &c.CriadoEm, &c.AtualizadoEm); err != nil {
			return nil, err
		}
		c.EditalID, c.MateriaID, c.TopicoID = textoNulo(edital), textoNulo(materia), textoNulo(topico)
		lista = append(lista, c)
	}
	return lista, linhas.Err()
}
func (r *repositorioMySQL) Obter(ctx context.Context, alunoID, id string) (*Caderno, error) {
	var c Caderno
	var edital, materia, topico sql.NullString
	err := r.banco.QueryRowContext(ctx, `SELECT id, aluno_id, titulo, pasta, conteudo, edital_id, materia_id, topico_id, cor, ativo, criado_em, atualizado_em FROM cadernos WHERE id = ? AND aluno_id = ? AND ativo = TRUE`, id, alunoID).Scan(&c.ID, &c.AlunoID, &c.Titulo, &c.Pasta, &c.Conteudo, &edital, &materia, &topico, &c.Cor, &c.Ativo, &c.CriadoEm, &c.AtualizadoEm)
	if err != nil {
		return nil, err
	}
	c.EditalID, c.MateriaID, c.TopicoID = textoNulo(edital), textoNulo(materia), textoNulo(topico)
	return &c, nil
}
func (r *repositorioMySQL) Alterar(ctx context.Context, alunoID, id string, e EntradaCadernoNormalizada) error {
	_, err := r.banco.ExecContext(ctx, `UPDATE cadernos SET titulo = ?, pasta = ?, conteudo = ?, edital_id = ?, materia_id = ?, topico_id = ?, cor = ? WHERE id = ? AND aluno_id = ? AND ativo = TRUE`, e.Titulo, e.Pasta, e.Conteudo, e.EditalID, e.MateriaID, e.TopicoID, e.Cor, id, alunoID)
	return err
}
func (r *repositorioMySQL) Desativar(ctx context.Context, alunoID, id string) error {
	_, err := r.banco.ExecContext(ctx, `UPDATE cadernos SET ativo = FALSE WHERE id = ? AND aluno_id = ?`, id, alunoID)
	return err
}
