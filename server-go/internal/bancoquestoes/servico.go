package bancoquestoes

import (
	"context"
	"database/sql"
)

type Servico struct{ repositorio repositorio }

func Novo(banco *sql.DB) *Servico { return &Servico{repositorio: novoRepositorioMySQL(banco)} }

func (s *Servico) ListarQuestoes(ctx context.Context, f FiltroQuestao) ([]Questao, error) {
	return s.repositorio.ListarQuestoes(ctx, f)
}
func (s *Servico) CriarQuestao(ctx context.Context, e EntradaQuestao, mentorID string) (*Questao, error) {
	return s.repositorio.CriarQuestao(ctx, e, mentorID)
}
func (s *Servico) ImportarQuestoes(ctx context.Context, lista []EntradaQuestao, mentorID string) (int, error) {
	return s.repositorio.ImportarQuestoes(ctx, lista, mentorID)
}
func (s *Servico) AlterarQuestao(ctx context.Context, id string, e EntradaQuestao) error {
	return s.repositorio.AlterarQuestao(ctx, id, e)
}
func (s *Servico) DesativarQuestao(ctx context.Context, id string) error {
	return s.repositorio.DesativarQuestao(ctx, id)
}
func (s *Servico) Responder(ctx context.Context, alunoID string, e RespostaQuestaoEntrada) (*ResultadoResposta, error) {
	return s.repositorio.Responder(ctx, alunoID, e)
}
func (s *Servico) HistoricoQuestao(ctx context.Context, alunoID, questaoID, concursoID string) ([]RegistroHistoricoQuestoes, error) {
	return s.repositorio.HistoricoQuestao(ctx, alunoID, questaoID, concursoID)
}
func (s *Servico) ObterEstatisticas(ctx context.Context, alunoID, concursoID string) (*EstatisticasBancoQuestoes, error) {
	return s.repositorio.ObterEstatisticas(ctx, alunoID, concursoID)
}
