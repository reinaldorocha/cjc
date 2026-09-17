package usuarios

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"track-concursos-web/internal/dominio"
	"track-concursos-web/internal/identificador"

	"github.com/go-sql-driver/mysql"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(banco *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: banco} }

func (r *repositorioMySQL) Listar(ctx context.Context) ([]Resumo, error) {
	linhas, err := r.banco.QueryContext(ctx, `SELECT id,nome,email,papel,ativo,desativado_em FROM usuarios ORDER BY criado_em DESC`)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := make([]Resumo, 0)
	for linhas.Next() {
		var item Resumo
		if err := linhas.Scan(&item.ID, &item.Nome, &item.Email, &item.Papel, &item.Ativo, &item.DesativadoEm); err != nil {
			return nil, err
		}
		lista = append(lista, item)
	}
	return lista, linhas.Err()
}

func (r *repositorioMySQL) Criar(ctx context.Context, registro registroCriacao) error {
	_, err := r.banco.ExecContext(ctx, `INSERT INTO usuarios (id,nome,email,senha_hash,papel,telefone) VALUES (?,?,?,?,?,NULLIF(?, ''))`, registro.ID, registro.Nome, registro.Email, registro.SenhaHash, registro.Papel, registro.Telefone)
	return traduzirErroDuplicidade(err)
}

func (r *repositorioMySQL) MentorAtivo(ctx context.Context, id string) (bool, error) {
	var quantidade int
	err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM usuarios WHERE id=? AND papel='mentor' AND ativo=TRUE`, id).Scan(&quantidade)
	return quantidade == 1, err
}

func (r *repositorioMySQL) CriarAlunoVinculado(ctx context.Context, registro registroCriacao, mentorID string, permite bool, dataExpiracao *time.Time) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `INSERT INTO usuarios (id,nome,email,senha_hash,papel,telefone) VALUES (?,?,?,?, 'aluno',NULLIF(?, ''))`, registro.ID, registro.Nome, registro.Email, registro.SenhaHash, registro.Telefone); err != nil {
		return traduzirErroDuplicidade(err)
	}
	var data any
	if dataExpiracao != nil {
		data = *dataExpiracao
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO mentor_alunos (id,mentor_id,aluno_id,permite_cronograma_inteligente,data_expiracao_plano) VALUES (?,?,?,?,?)`, identificador.UUID(), mentorID, registro.ID, permite, data); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *repositorioMySQL) Alterar(ctx context.Context, id string, nome *string, ativo *bool, executor string) error {
	if nome != nil {
		resultado, err := r.banco.ExecContext(ctx, `UPDATE usuarios SET nome=? WHERE id=?`, *nome, id)
		if err != nil {
			return err
		}
		if linhas, _ := resultado.RowsAffected(); linhas == 0 {
			return dominio.ErrNaoEncontrado
		}
	}
	if ativo == nil {
		return nil
	}
	var consulta string
	if *ativo {
		consulta = `UPDATE usuarios SET ativo=TRUE,desativado_em=NULL,desativado_por=NULL WHERE id=?`
		_, err := r.banco.ExecContext(ctx, consulta, id)
		return err
	}
	consulta = `UPDATE usuarios SET ativo=FALSE,desativado_em=UTC_TIMESTAMP(),desativado_por=? WHERE id=?`
	_, err := r.banco.ExecContext(ctx, consulta, executor, id)
	return err
}

func (r *repositorioMySQL) Reativar(ctx context.Context, id string) error {
	_, err := r.banco.ExecContext(ctx, `UPDATE usuarios SET ativo=TRUE,desativado_em=NULL,desativado_por=NULL WHERE id=?`, id)
	return err
}

func (r *repositorioMySQL) RedefinirSenha(ctx context.Context, id, hash string) error {
	_, err := r.banco.ExecContext(ctx, `UPDATE usuarios SET senha_hash=? WHERE id=?`, hash, id)
	return err
}

func (r *repositorioMySQL) AlterarNomeProprio(ctx context.Context, id, nome string) error {
	_, err := r.banco.ExecContext(ctx, `UPDATE usuarios SET nome=? WHERE id=?`, nome, id)
	return err
}

func traduzirErroDuplicidade(err error) error {
	if err == nil {
		return nil
	}
	var mysqlErro *mysql.MySQLError
	if errors.As(err, &mysqlErro) && mysqlErro.Number == 1062 {
		return dominio.ErrConflito
	}
	return err
}
