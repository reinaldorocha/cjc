package mentoria

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"chega-junto-concurseiro-web/internal/dominio"

	"golang.org/x/crypto/bcrypt"
)

type Servico struct {
	repositorio repositorio
	agora       func() time.Time
}

type Aluno struct {
	ID                           string  `json:"id"`
	Nome                         string  `json:"nome"`
	Email                        string  `json:"email"`
	Telefone                     *string `json:"telefone,omitempty"`
	PermiteCronogramaInteligente bool    `json:"permiteCronogramaInteligente"`
	DataExpiracaoPlano           *string `json:"dataExpiracaoPlano,omitempty"`
}

type ItemRadarAluno struct {
	AlunoID            string  `json:"alunoId"`
	AlunoNome          string  `json:"alunoNome"`
	AlunoEmail         string  `json:"alunoEmail"`
	Telefone           *string `json:"telefone,omitempty"`
	ConcursoNome       string  `json:"concursoNome"`
	UltimoEstudoEm     *string `json:"ultimoEstudoEm"`
	DiasSemEstudar     int     `json:"diasSemEstudar"`
	SegundosSemana     int64   `json:"segundosSemana"`
	SegundosMes        int64   `json:"segundosMes"`
	QuestoesResolvidas int64   `json:"questoesResolvidas"`
	QuestoesAcertos    int64   `json:"questoesAcertos"`
	TaxaAcerto         float64 `json:"taxaAcerto"`
	RevisoesPendentes  int64   `json:"revisoesPendentes"`
	NivelRisco         string  `json:"nivelRisco"`
	MotivoRisco        string  `json:"motivoRisco"`
}

func Novo(banco *sql.DB) *Servico { return NovoComRepositorio(novoRepositorioMySQL(banco)) }

func NovoComRepositorio(repositorio repositorio) *Servico {
	return &Servico{repositorio: repositorio, agora: time.Now}
}

func (s *Servico) ResolverAluno(ctx context.Context, usuarioID, papel, solicitado string) (string, error) {
	if solicitado == "eu" || solicitado == usuarioID {
		return usuarioID, nil
	}
	if papel != "mentor" || solicitado == "" {
		return "", dominio.ErrNaoAutorizado
	}
	vinculado, err := s.repositorio.ResolverAluno(ctx, usuarioID, solicitado)
	if err != nil {
		return "", err
	}
	if !vinculado {
		return "", dominio.ErrNaoEncontrado
	}
	return solicitado, nil
}

func (s *Servico) DefinirMentor(ctx context.Context, alunoID, mentorID, executor string) error {
	if strings.TrimSpace(alunoID) == "" || strings.TrimSpace(mentorID) == "" {
		return dominio.ErrEntradaInvalida
	}
	return s.repositorio.DefinirMentor(ctx, alunoID, mentorID, executor)
}

func (s *Servico) ListarAlunos(ctx context.Context, mentorID string) ([]Aluno, error) {
	return s.repositorio.ListarAlunos(ctx, mentorID)
}

func (s *Servico) ObterAluno(ctx context.Context, mentorID, alunoID string) (Aluno, error) {
	return s.repositorio.ObterAluno(ctx, mentorID, alunoID)
}

func (s *Servico) VisaoGeral(ctx context.Context, mentorID, alunoID string) (map[string]any, error) {
	if _, err := s.ObterAluno(ctx, mentorID, alunoID); err != nil {
		return nil, err
	}
	resumo, err := s.repositorio.VisaoGeral(ctx, alunoID)
	if err != nil {
		return nil, err
	}
	resultado := make(map[string]any, len(resumo))
	for chave, valor := range resumo {
		resultado[chave] = valor
	}
	return resultado, nil
}

func (s *Servico) Configurar(ctx context.Context, mentorID, alunoID string, permite *bool, data *string, telefone *string, nome *string, senha *string) error {
	if permite == nil && data == nil && telefone == nil && nome == nil && senha == nil {
		return dominio.ErrEntradaInvalida
	}
	var dataExpiracao *time.Time
	if data != nil && strings.TrimSpace(*data) != "" {
		valor, err := time.Parse("2006-01-02", *data)
		if err != nil {
			return dominio.ErrEntradaInvalida
		}
		dataExpiracao = &valor
	}
	var senhaHash *string
	if senha != nil && strings.TrimSpace(*senha) != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(*senha)), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		strHash := string(hash)
		senhaHash = &strHash
	}
	return s.repositorio.Configurar(ctx, mentorID, alunoID, permite, data != nil, dataExpiracao, telefone, nome, senhaHash)
}

func (s *Servico) RadarAlunos(ctx context.Context, mentorID string) ([]ItemRadarAluno, error) {
	dados, err := s.repositorio.Radar(ctx, mentorID)
	if err != nil {
		return nil, err
	}
	agora := s.agora()
	lista := make([]ItemRadarAluno, 0, len(dados))
	for _, dado := range dados {
		item := ItemRadarAluno{AlunoID: dado.AlunoID, AlunoNome: dado.AlunoNome, AlunoEmail: dado.AlunoEmail, Telefone: dado.Telefone, ConcursoNome: dado.ConcursoNome, SegundosSemana: dado.SegundosSemana, SegundosMes: dado.SegundosMes, QuestoesResolvidas: dado.QuestoesResolvidas, QuestoesAcertos: dado.QuestoesAcertos, RevisoesPendentes: dado.RevisoesPendentes}
		if dado.UltimoEstudoEm != nil {
			data := dado.UltimoEstudoEm.Format("2006-01-02")
			item.UltimoEstudoEm = &data
			item.DiasSemEstudar = diasSemEstudar(agora, *dado.UltimoEstudoEm)
		} else {
			item.DiasSemEstudar = 999
		}
		if item.QuestoesResolvidas > 0 {
			item.TaxaAcerto = float64(item.QuestoesAcertos) / float64(item.QuestoesResolvidas) * 100
		}
		item.NivelRisco, item.MotivoRisco = classificarRisco(item)
		lista = append(lista, item)
	}
	return lista, nil
}

func diasSemEstudar(agora, ultimo time.Time) int {
	dias := int(agora.Sub(ultimo).Hours() / 24)
	if dias < 0 {
		return 0
	}
	return dias
}

func classificarRisco(item ItemRadarAluno) (string, string) {
	if item.DiasSemEstudar >= 2 || (item.QuestoesResolvidas >= 10 && item.TaxaAcerto < 50) {
		if item.DiasSemEstudar >= 999 {
			return "vermelho", "Nenhum estudo registrado"
		}
		if item.DiasSemEstudar >= 2 {
			return "vermelho", fmt.Sprintf("Sem estudar há %d dias", item.DiasSemEstudar)
		}
		return "vermelho", fmt.Sprintf("Baixo rendimento em questões (%.0f%%)", item.TaxaAcerto)
	}
	if item.DiasSemEstudar == 1 || (item.QuestoesResolvidas >= 10 && item.TaxaAcerto < 70) || item.RevisoesPendentes >= 5 {
		if item.DiasSemEstudar == 1 {
			return "amarelo", "Sem estudar ontem"
		}
		if item.RevisoesPendentes >= 5 {
			return "amarelo", fmt.Sprintf("%d revisões atrasadas", item.RevisoesPendentes)
		}
		return "amarelo", fmt.Sprintf("Atenção ao rendimento (%.0f%%)", item.TaxaAcerto)
	}
	return "verde", "Ritmo constante e em dia"
}
