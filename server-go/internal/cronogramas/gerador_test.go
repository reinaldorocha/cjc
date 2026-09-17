package cronogramas

import (
	"testing"
	"time"
)

func TestGerarCicloOrdenaPorPrioridade(t *testing.T) {
	materias := []materiaInfo{{id: "a", nome: "A", prioridade: 72}, {id: "b", nome: "B", prioridade: 40}}
	itens := gerarCiclo(materias)
	if len(itens) != 2 {
		t.Fatalf("itens = %d, esperado 2", len(itens))
	}
	if itens[0].DuracaoMinutos != 90 || itens[1].DuracaoMinutos != 50 {
		t.Fatalf("durações inesperadas: %d, %d", itens[0].DuracaoMinutos, itens[1].DuracaoMinutos)
	}
	if itens[0].PosicaoCiclo == nil || *itens[0].PosicaoCiclo != 1 {
		t.Fatal("primeira posição de ciclo não foi preservada")
	}
}

func TestGerarAgendaRespeitaCapacidadeDiaria(t *testing.T) {
	fuso := time.FixedZone("teste", 0)
	materias := []materiaInfo{{id: "a", nome: "A", prioridade: 60}, {id: "b", nome: "B", prioridade: 50}}
	unidades := map[string][]unidade{
		"a": {{materiaID: "a", topicoID: "ta", nome: "Tópico A", prioridade: 60}},
		"b": {{materiaID: "b", topicoID: "tb", nome: "Tópico B", prioridade: 50}},
	}
	cfg := Configuracao{MinutosTopico: 60, MaxTopicosDia: 1, Horas: map[string]float64{"dom": 1, "seg": 1, "ter": 1, "qua": 1, "qui": 1, "sex": 1, "sab": 1}}
	itens := gerarAgenda(materias, unidades, cfg, fuso)
	if len(itens) != 2 {
		t.Fatalf("itens = %d, esperado 2", len(itens))
	}
	if itens[0].DataPlanejada == nil || itens[1].DataPlanejada == nil || *itens[0].DataPlanejada == *itens[1].DataPlanejada {
		t.Fatal("capacidade diária não foi respeitada")
	}
}
