package bancoquestoes

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestListarQuestoes(t *testing.T) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/chega_junto_concurseiro?parseTime=true")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := novoRepositorioMySQL(db)
	ctx := context.Background()

	// Teste com AlunoID vazio
	q1, err := repo.ListarQuestoes(ctx, FiltroQuestao{})
	if err != nil {
		t.Fatalf("Erro ao listar questoes sem aluno: %v", err)
	}
	t.Logf("Questoes sem aluno: %d", len(q1))

	// Teste com AlunoID
	q2, err := repo.ListarQuestoes(ctx, FiltroQuestao{AlunoID: "ec9c1eea-9459-4fec-9982-6ef77359d381"})
	if err != nil {
		t.Fatalf("Erro ao listar questoes com aluno: %v", err)
	}
	t.Logf("Questoes com aluno: %d", len(q2))

	// Teste com Filtro completo
	q3, err := repo.ListarQuestoes(ctx, FiltroQuestao{
		AlunoID:    "ec9c1eea-9459-4fec-9982-6ef77359d381",
		Disciplina: "Direito Constitucional",
	})
	if err != nil {
		t.Fatalf("Erro ao listar questoes com filtro: %v", err)
	}
	t.Logf("Questoes com filtro: %d", len(q3))

	// Teste Responder
	if len(q1) > 0 {
		res, err := repo.Responder(ctx, "ec9c1eea-9459-4fec-9982-6ef77359d381", RespostaQuestaoEntrada{
			QuestaoID:     q1[0].ID,
			RespostaAluno: q1[0].RespostaCorreta,
		})
		if err != nil {
			t.Fatalf("Erro ao responder questao: %v", err)
		}
		if !res.Correto {
			t.Errorf("Esperava resposta correta, obteve incorreta")
		}

		stats, err := repo.ObterEstatisticas(ctx, "ec9c1eea-9459-4fec-9982-6ef77359d381", "")
		if err != nil {
			t.Fatalf("Erro ao obter estatisticas: %v", err)
		}
		t.Logf("Stats: Total=%d, Acertos=%d, Erros=%d, Taxa=%.1f%%", stats.TotalRespondidas, stats.TotalAcertos, stats.TotalErros, stats.TaxaAcerto)
	}
}

func TestVerificarResposta(t *testing.T) {
	alts := []string{"Motivo", "Objeto", "Competência", "Conveniência"}

	// Multipla escolha: aluno envia C, gabarito é C
	if !verificarResposta("multipla_escolha", alts, "C", "C") {
		t.Errorf("Esperava true para C == C")
	}
	// Multipla escolha: aluno envia C, gabarito é o texto da alternativa
	if !verificarResposta("multipla_escolha", alts, "C", "Competência") {
		t.Errorf("Esperava true para C == Competência")
	}
	// Multipla escolha: aluno envia texto da alternativa, gabarito é C
	if !verificarResposta("multipla_escolha", alts, "Competência", "C") {
		t.Errorf("Esperava true para Competência == C")
	}
	// Multipla escolha: aluno envia C) Competência, gabarito é C
	if !verificarResposta("multipla_escolha", alts, "C) Competência", "C") {
		t.Errorf("Esperava true para C) Competência == C")
	}
	// Multipla escolha: aluno envia A, gabarito é C
	if verificarResposta("multipla_escolha", alts, "A", "C") {
		t.Errorf("Esperava false para A != C")
	}

	// Certo / Errado
	if !verificarResposta("certo_errado", nil, "Certo", "certo") {
		t.Errorf("Esperava true para Certo == certo")
	}
	if !verificarResposta("certo_errado", nil, "C", "certo") {
		t.Errorf("Esperava true para C == certo")
	}
	if !verificarResposta("certo_errado", nil, "Certo", "C") {
		t.Errorf("Esperava true para Certo == C")
	}
	if !verificarResposta("certo_errado", nil, "Errado", "E") {
		t.Errorf("Esperava true para Errado == E")
	}
	if verificarResposta("certo_errado", nil, "Certo", "Errado") {
		t.Errorf("Esperava false para Certo != Errado")
	}
}
