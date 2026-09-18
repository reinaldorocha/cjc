package autenticacao

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"chega-junto-concurseiro-web/internal/configuracao"
	"chega-junto-concurseiro-web/internal/identificador"
	"chega-junto-concurseiro-web/internal/modelos"

	"golang.org/x/crypto/bcrypt"
)

var ErrCredenciais = errors.New("credenciais inválidas")
var ErrSessao = errors.New("sessão inválida")
var ErrPlanoExpirado = errors.New("plano expirado")
var ErrSemVinculo = errors.New("aluno sem vínculo ativo")

type Servico struct {
	banco *sql.DB
	cfg   configuracao.Configuracao
}

func Novo(banco *sql.DB, cfg configuracao.Configuracao) *Servico {
	return &Servico{banco: banco, cfg: cfg}
}

func (s *Servico) SemearMestre() error {
	if s.cfg.MestreEmail == "" || s.cfg.MestreSenha == "" {
		return nil
	}
	var quantidade int
	if err := s.banco.QueryRow(`SELECT COUNT(*) FROM usuarios WHERE papel='mestre'`).Scan(&quantidade); err != nil {
		return err
	}
	if quantidade > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(s.cfg.MestreSenha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.banco.Exec(`INSERT INTO usuarios (id,nome,email,senha_hash,papel) VALUES (?,?,?,?, 'mestre')`, identificador.UUID(), "Administrador", normalizarEmail(s.cfg.MestreEmail), string(hash))
	return err
}

func (s *Servico) Entrar(ctx context.Context, email, senha string) (modelos.Sessao, error) {
	usuario, hash, err := s.buscarPorEmail(ctx, normalizarEmail(email))
	if err != nil || !usuario.Ativo || bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha)) != nil {
		return modelos.Sessao{}, ErrCredenciais
	}
	if usuario.Papel == "aluno" {
		if err = s.aplicarPlano(ctx, &usuario); err != nil {
			return modelos.Sessao{}, err
		}
	}
	token, csrf := identificador.Segredo(), identificador.Segredo()
	expira := time.Now().UTC().Add(24 * time.Hour)
	_, err = s.banco.ExecContext(ctx, `INSERT INTO sessoes_autenticacao (id,usuario_id,token_hash,csrf_hash,expira_em) VALUES (?,?,?,?,?)`, identificador.UUID(), usuario.ID, identificador.Hash(token), identificador.Hash(csrf), expira)
	if err != nil {
		return modelos.Sessao{}, err
	}
	if _, err = s.banco.ExecContext(ctx, `UPDATE usuarios SET ultimo_acesso_em=UTC_TIMESTAMP() WHERE id=?`, usuario.ID); err != nil {
		return modelos.Sessao{}, err
	}
	return modelos.Sessao{Usuario: usuario, Token: token, TokenCSRF: csrf}, nil
}

func (s *Servico) Autenticar(ctx context.Context, token string) (modelos.ContextoAutenticado, error) {
	if token == "" {
		return modelos.ContextoAutenticado{}, ErrSessao
	}
	var u modelos.Usuario
	var expira time.Time
	var csrf string
	var permite sql.NullBool
	var data sql.NullTime
	err := s.banco.QueryRowContext(ctx, `SELECT u.id,u.nome,u.email,u.papel,u.ativo,sa.expira_em,sa.csrf_hash,ma.permite_cronograma_inteligente,ma.data_expiracao_plano FROM sessoes_autenticacao sa JOIN usuarios u ON u.id=sa.usuario_id LEFT JOIN mentor_alunos ma ON ma.aluno_id=u.id AND ma.ativo=TRUE WHERE sa.token_hash=? AND sa.revogado_em IS NULL`, identificador.Hash(token)).Scan(&u.ID, &u.Nome, &u.Email, &u.Papel, &u.Ativo, &expira, &csrf, &permite, &data)
	if err != nil || !u.Ativo || time.Now().UTC().After(expira) {
		return modelos.ContextoAutenticado{}, ErrSessao
	}
	u.PermiteCronogramaInteligente = permite.Valid && permite.Bool
	if data.Valid {
		u.DataExpiracaoPlano = dataTexto(data.Time)
	}
	if u.Papel == "aluno" {
		if err = s.aplicarPlano(ctx, &u); err != nil {
			return modelos.ContextoAutenticado{}, err
		}
	}
	if _, err = s.banco.ExecContext(ctx, `UPDATE sessoes_autenticacao SET ultimo_uso_em=UTC_TIMESTAMP() WHERE token_hash=?`, identificador.Hash(token)); err != nil {
		return modelos.ContextoAutenticado{}, err
	}
	return modelos.ContextoAutenticado{Usuario: u, CSRFHash: csrf}, nil
}

func (s *Servico) Revogar(ctx context.Context, token string) error {
	_, err := s.banco.ExecContext(ctx, `UPDATE sessoes_autenticacao SET revogado_em=UTC_TIMESTAMP() WHERE token_hash=?`, identificador.Hash(token))
	return err
}
func (s *Servico) RevogarTodas(ctx context.Context, usuarioID string) error {
	_, err := s.banco.ExecContext(ctx, `UPDATE sessoes_autenticacao SET revogado_em=UTC_TIMESTAMP() WHERE usuario_id=? AND revogado_em IS NULL`, usuarioID)
	return err
}
func (s *Servico) RenovarCSRF(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", ErrSessao
	}
	novo := identificador.Segredo()
	resultado, err := s.banco.ExecContext(ctx, `UPDATE sessoes_autenticacao SET csrf_hash=?,ultimo_uso_em=UTC_TIMESTAMP() WHERE token_hash=? AND revogado_em IS NULL AND expira_em>UTC_TIMESTAMP()`, identificador.Hash(novo), identificador.Hash(token))
	if err != nil {
		return "", err
	}
	quantidade, err := resultado.RowsAffected()
	if err != nil || quantidade == 0 {
		return "", ErrSessao
	}
	return novo, nil
}
func (s *Servico) AlterarSenha(ctx context.Context, usuario modelos.Usuario, atual, nova string) error {
	if len(nova) < 8 {
		return errors.New("senha curta")
	}
	_, hash, err := s.buscarPorEmail(ctx, usuario.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(atual)) != nil {
		return ErrCredenciais
	}
	novo, err := bcrypt.GenerateFromPassword([]byte(nova), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.banco.ExecContext(ctx, `UPDATE usuarios SET senha_hash=? WHERE id=?`, string(novo), usuario.ID)
	return err
}
func (s *Servico) HashCSRF(valor string) string { return identificador.Hash(valor) }

func (s *Servico) aplicarPlano(ctx context.Context, u *modelos.Usuario) error {
	var data sql.NullTime
	var permite bool
	err := s.banco.QueryRowContext(ctx, `SELECT data_expiracao_plano,permite_cronograma_inteligente FROM mentor_alunos WHERE aluno_id=? AND ativo=TRUE`, u.ID).Scan(&data, &permite)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrSemVinculo
	}
	if err != nil {
		return err
	}
	u.PermiteCronogramaInteligente = permite
	if data.Valid {
		u.DataExpiracaoPlano = dataTexto(data.Time)
		hoje := time.Now().In(s.cfg.FusoHorario).Format("2006-01-02")
		if hoje >= data.Time.In(s.cfg.FusoHorario).Format("2006-01-02") {
			u.Expirado = true
			return ErrPlanoExpirado
		}
	}
	return nil
}
func (s *Servico) buscarPorEmail(ctx context.Context, email string) (modelos.Usuario, string, error) {
	var u modelos.Usuario
	var hash string
	err := s.banco.QueryRowContext(ctx, `SELECT id,nome,email,senha_hash,papel,ativo FROM usuarios WHERE email=?`, email).Scan(&u.ID, &u.Nome, &u.Email, &hash, &u.Papel, &u.Ativo)
	return u, hash, err
}
func normalizarEmail(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
func dataTexto(v time.Time) *string   { x := v.Format("2006-01-02"); return &x }
