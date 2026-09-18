package mentoria

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"chega-junto-concurseiro-web/internal/dominio"
	"chega-junto-concurseiro-web/internal/identificador"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(banco *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: banco} }

func (r *repositorioMySQL) ResolverAluno(ctx context.Context, mentorID, alunoID string) (bool, error) {
	var quantidade int
	err := r.banco.QueryRowContext(ctx, `SELECT COUNT(*) FROM mentor_alunos WHERE mentor_id=? AND aluno_id=? AND ativo=TRUE`, mentorID, alunoID).Scan(&quantidade)
	return quantidade > 0, err
}

func (r *repositorioMySQL) DefinirMentor(ctx context.Context, alunoID, mentorID, executor string) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := validarPapel(ctx, tx, mentorID, "mentor"); err != nil {
		return err
	}
	if err := validarPapel(ctx, tx, alunoID, "aluno"); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE mentor_alunos SET ativo=FALSE,desativado_em=UTC_TIMESTAMP(),desativado_por=? WHERE aluno_id=? AND ativo=TRUE`, executor, alunoID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO mentor_alunos (id,mentor_id,aluno_id) VALUES (?,?,?)`, identificador.UUID(), mentorID, alunoID); err != nil {
		return err
	}
	return tx.Commit()
}

func validarPapel(ctx context.Context, tx *sql.Tx, id, papelEsperado string) error {
	var papel string
	err := tx.QueryRowContext(ctx, `SELECT papel FROM usuarios WHERE id=? AND ativo=TRUE`, id).Scan(&papel)
	if errors.Is(err, sql.ErrNoRows) || papel != papelEsperado {
		return dominio.ErrEntradaInvalida
	}
	return err
}

func (r *repositorioMySQL) ListarAlunos(ctx context.Context, mentorID string) ([]Aluno, error) {
	linhas, err := r.banco.QueryContext(ctx, `SELECT u.id,u.nome,u.email,u.telefone,ma.permite_cronograma_inteligente,ma.data_expiracao_plano FROM mentor_alunos ma JOIN usuarios u ON u.id=ma.aluno_id WHERE ma.mentor_id=? AND ma.ativo=TRUE ORDER BY u.nome`, mentorID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := make([]Aluno, 0)
	for linhas.Next() {
		aluno, err := lerAluno(linhas)
		if err != nil {
			return nil, err
		}
		lista = append(lista, aluno)
	}
	return lista, linhas.Err()
}

func (r *repositorioMySQL) ObterAluno(ctx context.Context, mentorID, alunoID string) (Aluno, error) {
	linha := r.banco.QueryRowContext(ctx, `SELECT u.id,u.nome,u.email,u.telefone,ma.permite_cronograma_inteligente,ma.data_expiracao_plano FROM mentor_alunos ma JOIN usuarios u ON u.id=ma.aluno_id WHERE ma.mentor_id=? AND ma.aluno_id=? AND ma.ativo=TRUE`, mentorID, alunoID)
	aluno, err := lerAluno(linha)
	if errors.Is(err, sql.ErrNoRows) {
		return Aluno{}, dominio.ErrNaoEncontrado
	}
	return aluno, err
}

type scanner interface{ Scan(...any) error }

func lerAluno(linha scanner) (Aluno, error) {
	var aluno Aluno
	var telefone sql.NullString
	var data sql.NullTime
	if err := linha.Scan(&aluno.ID, &aluno.Nome, &aluno.Email, &telefone, &aluno.PermiteCronogramaInteligente, &data); err != nil {
		return Aluno{}, err
	}
	if telefone.Valid && strings.TrimSpace(telefone.String) != "" {
		valor := strings.TrimSpace(telefone.String)
		aluno.Telefone = &valor
	}
	if data.Valid {
		valor := data.Time.Format("2006-01-02")
		aluno.DataExpiracaoPlano = &valor
	}
	return aluno, nil
}

func (r *repositorioMySQL) VisaoGeral(ctx context.Context, alunoID string) (map[string]int64, error) {
	const consulta = `SELECT
	(SELECT COUNT(*) FROM aluno_concursos WHERE aluno_id=? AND ativo=TRUE),
	(SELECT COUNT(*) FROM aluno_editais WHERE aluno_id=? AND ativo=TRUE),
	(SELECT COUNT(*) FROM sessoes_estudo WHERE aluno_id=?),
	(SELECT COALESCE(SUM(segundos),0) FROM sessoes_estudo WHERE aluno_id=?),
	(SELECT COALESCE(SUM(resolvidas),0) FROM registros_questoes WHERE aluno_id=?),
	(SELECT COALESCE(SUM(acertos),0) FROM registros_questoes WHERE aluno_id=?),
	(SELECT COUNT(*) FROM revisoes_programadas WHERE aluno_id=? AND concluida=FALSE)`
	var concursos, editais, sessoes, segundos, questoes, acertos, revisoes int64
	err := r.banco.QueryRowContext(ctx, consulta, alunoID, alunoID, alunoID, alunoID, alunoID, alunoID, alunoID).Scan(&concursos, &editais, &sessoes, &segundos, &questoes, &acertos, &revisoes)
	if err != nil {
		return nil, err
	}
	return map[string]int64{"concursos": concursos, "editais": editais, "sessoes": sessoes, "segundosEstudados": segundos, "questoesResolvidas": questoes, "acertos": acertos, "revisoesPendentes": revisoes}, nil
}

func (r *repositorioMySQL) Configurar(ctx context.Context, mentorID, alunoID string, permite *bool, atualizarData bool, data *time.Time, telefone *string, nome *string, senhaHash *string) error {
	tx, err := r.banco.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var quantidade int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM mentor_alunos WHERE mentor_id=? AND aluno_id=? AND ativo=TRUE`, mentorID, alunoID).Scan(&quantidade); err != nil {
		return err
	}
	if quantidade == 0 {
		return dominio.ErrNaoEncontrado
	}
	if permite != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE mentor_alunos SET permite_cronograma_inteligente=? WHERE mentor_id=? AND aluno_id=? AND ativo=TRUE`, *permite, mentorID, alunoID); err != nil {
			return err
		}
	}
	if telefone != nil {
		if _, err = tx.ExecContext(ctx, `UPDATE usuarios SET telefone=NULLIF(?, '') WHERE id=?`, strings.TrimSpace(*telefone), alunoID); err != nil {
			return err
		}
	}
	if nome != nil && strings.TrimSpace(*nome) != "" {
		if _, err = tx.ExecContext(ctx, `UPDATE usuarios SET nome=? WHERE id=?`, strings.TrimSpace(*nome), alunoID); err != nil {
			return err
		}
	}
	if senhaHash != nil && *senhaHash != "" {
		if _, err = tx.ExecContext(ctx, `UPDATE usuarios SET senha_hash=? WHERE id=?`, *senhaHash, alunoID); err != nil {
			return err
		}
	}
	if atualizarData {
		if _, err = tx.ExecContext(ctx, `UPDATE mentor_alunos SET data_expiracao_plano=? WHERE mentor_id=? AND aluno_id=? AND ativo=TRUE`, data, mentorID, alunoID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *repositorioMySQL) Radar(ctx context.Context, mentorID string) ([]dadosRadar, error) {
	const consulta = `SELECT u.id,u.nome,u.email,NULLIF(TRIM(u.telefone),''),
	COALESCE((SELECT c.nome FROM aluno_concursos ac JOIN concursos c ON c.id=ac.concurso_id WHERE ac.aluno_id=u.id AND ac.ativo=TRUE AND c.ativo=TRUE ORDER BY ac.ordem,c.nome LIMIT 1),'Sem concurso definido'),
	(SELECT MAX(data_estudo) FROM (SELECT MAX(estudado_em) data_estudo FROM sessoes_estudo WHERE aluno_id=u.id AND ativo=TRUE UNION ALL SELECT MAX(registrado_em) FROM registros_questoes WHERE aluno_id=u.id AND ativo=TRUE) estudos),
	(SELECT COALESCE(SUM(segundos),0) FROM sessoes_estudo WHERE aluno_id=u.id AND ativo=TRUE AND estudado_em>=UTC_TIMESTAMP()-INTERVAL 7 DAY),
	(SELECT COALESCE(SUM(segundos),0) FROM sessoes_estudo WHERE aluno_id=u.id AND ativo=TRUE AND estudado_em>=UTC_TIMESTAMP()-INTERVAL 30 DAY),
	(SELECT COALESCE(SUM(resolvidas),0) FROM registros_questoes WHERE aluno_id=u.id AND ativo=TRUE AND registrado_em>=UTC_TIMESTAMP()-INTERVAL 30 DAY),
	(SELECT COALESCE(SUM(acertos),0) FROM registros_questoes WHERE aluno_id=u.id AND ativo=TRUE AND registrado_em>=UTC_TIMESTAMP()-INTERVAL 30 DAY),
	(SELECT COUNT(*) FROM revisoes_programadas WHERE aluno_id=u.id AND concluida=FALSE AND ativo=TRUE AND proxima_data<=CURRENT_DATE)
	FROM mentor_alunos ma JOIN usuarios u ON u.id=ma.aluno_id WHERE ma.mentor_id=? AND ma.ativo=TRUE ORDER BY u.nome`
	linhas, err := r.banco.QueryContext(ctx, consulta, mentorID)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()
	lista := make([]dadosRadar, 0)
	for linhas.Next() {
		var item dadosRadar
		var telefone sql.NullString
		var ultimo sql.NullTime
		if err := linhas.Scan(&item.AlunoID, &item.AlunoNome, &item.AlunoEmail, &telefone, &item.ConcursoNome, &ultimo, &item.SegundosSemana, &item.SegundosMes, &item.QuestoesResolvidas, &item.QuestoesAcertos, &item.RevisoesPendentes); err != nil {
			return nil, err
		}
		if telefone.Valid {
			valor := telefone.String
			item.Telefone = &valor
		}
		if ultimo.Valid {
			valor := ultimo.Time
			item.UltimoEstudoEm = &valor
		}
		lista = append(lista, item)
	}
	return lista, linhas.Err()
}
