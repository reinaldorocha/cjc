package usuarios

import (
	"context"
	"testing"
	"time"

	"chega-junto-concurseiro-web/internal/dominio"
)

type repositorioFalso struct {
	mentorAtivo bool
	criado      registroCriacao
	vinculado   bool
}

func (r *repositorioFalso) Listar(context.Context) ([]Resumo, error) { return nil, nil }
func (r *repositorioFalso) Criar(_ context.Context, registro registroCriacao) error {
	r.criado = registro
	return nil
}
func (r *repositorioFalso) CriarAlunoVinculado(_ context.Context, registro registroCriacao, _ string, _ bool, _ *time.Time) error {
	r.criado = registro
	r.vinculado = true
	return nil
}
func (r *repositorioFalso) MentorAtivo(context.Context, string) (bool, error) {
	return r.mentorAtivo, nil
}
func (r *repositorioFalso) Alterar(context.Context, string, *string, *bool, string) error {
	return nil
}
func (r *repositorioFalso) Reativar(context.Context, string) error { return nil }
func (r *repositorioFalso) RedefinirSenha(context.Context, string, string) error {
	return nil
}
func (r *repositorioFalso) AlterarNomeProprio(context.Context, string, string) error {
	return nil
}

func TestCriarNormalizaDadosDoUsuario(t *testing.T) {
	repositorio := &repositorioFalso{}
	servico := NovoComRepositorio(repositorio)

	id, err := servico.Criar(context.Background(), Criacao{Nome: "  Ana  ", Email: "ANA@EXEMPLO.COM ", Senha: "senha-segura", Papel: "mentor"})

	if err != nil {
		t.Fatalf("criar retornou erro: %v", err)
	}
	if id == "" || repositorio.criado.ID != id {
		t.Fatal("identificador não foi persistido")
	}
	if repositorio.criado.Nome != "Ana" || repositorio.criado.Email != "ana@exemplo.com" {
		t.Fatalf("dados não normalizados: %+v", repositorio.criado)
	}
	if repositorio.criado.SenhaHash == "senha-segura" {
		t.Fatal("senha não foi protegida antes de persistir")
	}
}

func TestCriarAlunoExigeMentorAtivo(t *testing.T) {
	servico := NovoComRepositorio(&repositorioFalso{})
	_, err := servico.CriarAlunoVinculado(context.Background(), Criacao{Nome: "Aluno", Email: "aluno@exemplo.com", Senha: "senha-segura"}, "mentor-inativo", false, nil)
	if err != dominio.ErrNaoEncontrado {
		t.Fatalf("erro = %v, esperado %v", err, dominio.ErrNaoEncontrado)
	}
}

func TestCriarAlunoValidaDataDeExpiracao(t *testing.T) {
	servico := NovoComRepositorio(&repositorioFalso{mentorAtivo: true})
	data := "2026-99-99"
	_, err := servico.CriarAlunoVinculado(context.Background(), Criacao{Nome: "Aluno", Email: "aluno@exemplo.com", Senha: "senha-segura"}, "mentor", false, &data)
	if err != dominio.ErrEntradaInvalida {
		t.Fatalf("erro = %v, esperado %v", err, dominio.ErrEntradaInvalida)
	}
}
