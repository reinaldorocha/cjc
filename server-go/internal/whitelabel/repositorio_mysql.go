package whitelabel

import (
	"context"
	"database/sql"
)

type repositorioMySQL struct{ banco *sql.DB }

func novoRepositorioMySQL(banco *sql.DB) *repositorioMySQL { return &repositorioMySQL{banco: banco} }
func (r *repositorioMySQL) ObterPorMentor(ctx context.Context, id string) (Config, bool, error) {
	var c Config
	var logo, banner, msg sql.NullString
	var atualizado sql.NullTime
	c.MentorID = id
	err := r.banco.QueryRowContext(ctx, `SELECT nome_plataforma,logo_url,banner_url,cor_primaria,cor_secundaria,mensagem_boas_vindas,atualizado_em FROM mentor_whitelabel WHERE mentor_id=?`, id).Scan(&c.NomePlataforma, &logo, &banner, &c.CorPrimaria, &c.CorSecundaria, &msg, &atualizado)
	if err == sql.ErrNoRows {
		return Config{}, false, nil
	}
	if err != nil {
		return Config{}, false, err
	}
	if logo.Valid {
		c.LogoURL = &logo.String
	}
	if banner.Valid {
		c.BannerURL = &banner.String
	}
	if msg.Valid {
		c.MensagemBoasVindas = &msg.String
	}
	if atualizado.Valid {
		c.AtualizadoEm = atualizado.Time.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return c, true, nil
}
func (r *repositorioMySQL) MentorDoAluno(ctx context.Context, aluno string) (string, bool, error) {
	var mentor string
	err := r.banco.QueryRowContext(ctx, `SELECT mentor_id FROM mentor_alunos WHERE aluno_id=? AND ativo=TRUE LIMIT 1`, aluno).Scan(&mentor)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	return mentor, err == nil, err
}
func (r *repositorioMySQL) Salvar(ctx context.Context, c Config) error {
	_, err := r.banco.ExecContext(ctx, `INSERT INTO mentor_whitelabel (mentor_id,nome_plataforma,logo_url,banner_url,cor_primaria,cor_secundaria,mensagem_boas_vindas) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE nome_plataforma=VALUES(nome_plataforma),logo_url=VALUES(logo_url),banner_url=VALUES(banner_url),cor_primaria=VALUES(cor_primaria),cor_secundaria=VALUES(cor_secundaria),mensagem_boas_vindas=VALUES(mensagem_boas_vindas)`, c.MentorID, c.NomePlataforma, c.LogoURL, c.BannerURL, c.CorPrimaria, c.CorSecundaria, c.MensagemBoasVindas)
	return err
}
