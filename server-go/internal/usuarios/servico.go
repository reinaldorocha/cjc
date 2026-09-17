package usuarios

import (
	"context"
	"database/sql"
	"net/mail"
	"strings"
	"time"

	"track-concursos-web/internal/dominio"
	"track-concursos-web/internal/identificador"

	"golang.org/x/crypto/bcrypt"
)

type Servico struct{ repositorio repositorio }

type Resumo struct {
	ID           string     `json:"id"`
	Nome         string     `json:"nome"`
	Email        string     `json:"email"`
	Papel        string     `json:"papel"`
	Ativo        bool       `json:"ativo"`
	DesativadoEm *time.Time `json:"desativadoEm,omitempty"`
}

type Criacao struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Senha    string `json:"senha"`
	Papel    string `json:"papel"`
	Telefone string `json:"telefone,omitempty"`
}

func Novo(banco *sql.DB) *Servico                         { return NovoComRepositorio(novoRepositorioMySQL(banco)) }
func NovoComRepositorio(repositorio repositorio) *Servico { return &Servico{repositorio: repositorio} }

func (s *Servico) Listar(ctx context.Context) ([]Resumo, error) { return s.repositorio.Listar(ctx) }

func (s *Servico) Criar(ctx context.Context, entrada Criacao) (string, error) {
	registro, err := prepararCriacao(entrada, false)
	if err != nil {
		return "", err
	}
	if err := s.repositorio.Criar(ctx, registro); err != nil {
		return "", err
	}
	return registro.ID, nil
}

func (s *Servico) CriarAlunoVinculado(ctx context.Context, entrada Criacao, mentorID string, permite bool, dataExpiracao *string) (string, error) {
	registro, err := prepararCriacao(entrada, true)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(mentorID) == "" {
		return "", dominio.ErrEntradaInvalida
	}
	mentorValido, err := s.repositorio.MentorAtivo(ctx, mentorID)
	if err != nil {
		return "", err
	}
	if !mentorValido {
		return "", dominio.ErrNaoEncontrado
	}
	data, err := dataOpcional(dataExpiracao)
	if err != nil {
		return "", err
	}
	if err := s.repositorio.CriarAlunoVinculado(ctx, registro, mentorID, permite, data); err != nil {
		return "", err
	}
	return registro.ID, nil
}

func (s *Servico) Alterar(ctx context.Context, id string, nome *string, ativo *bool, executor string) error {
	if nome != nil {
		valor := strings.TrimSpace(*nome)
		if valor == "" {
			return dominio.ErrEntradaInvalida
		}
		nome = &valor
	}
	if nome == nil && ativo == nil {
		return dominio.ErrEntradaInvalida
	}
	return s.repositorio.Alterar(ctx, id, nome, ativo, executor)
}

func (s *Servico) Reativar(ctx context.Context, id string) error {
	return s.repositorio.Reativar(ctx, id)
}

func (s *Servico) RedefinirSenha(ctx context.Context, id, nova string) error {
	if len(nova) < 8 {
		return dominio.ErrEntradaInvalida
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(nova), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repositorio.RedefinirSenha(ctx, id, string(hash))
}

func (s *Servico) AlterarNomeProprio(ctx context.Context, id, nome string) error {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return dominio.ErrEntradaInvalida
	}
	return s.repositorio.AlterarNomeProprio(ctx, id, nome)
}

func prepararCriacao(entrada Criacao, forcarAluno bool) (registroCriacao, error) {
	nome := strings.TrimSpace(entrada.Nome)
	email := strings.ToLower(strings.TrimSpace(entrada.Email))
	papel := entrada.Papel
	if forcarAluno {
		papel = "aluno"
	}
	if nome == "" || len(entrada.Senha) < 8 || (papel != "mentor" && papel != "aluno") {
		return registroCriacao{}, dominio.ErrEntradaInvalida
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return registroCriacao{}, dominio.ErrEntradaInvalida
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(entrada.Senha), bcrypt.DefaultCost)
	if err != nil {
		return registroCriacao{}, err
	}
	return registroCriacao{ID: identificador.UUID(), Nome: nome, Email: email, SenhaHash: string(hash), Papel: papel, Telefone: strings.TrimSpace(entrada.Telefone)}, nil
}

func dataOpcional(texto *string) (*time.Time, error) {
	if texto == nil || strings.TrimSpace(*texto) == "" {
		return nil, nil
	}
	data, err := time.Parse("2006-01-02", *texto)
	if err != nil {
		return nil, dominio.ErrEntradaInvalida
	}
	return &data, nil
}
