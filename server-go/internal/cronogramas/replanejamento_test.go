package cronogramas

import (
	"reflect"
	"testing"
)

func TestRevisoesOcupamTempoDosAssuntos(t *testing.T) {
	cfg := Configuracao{Horas: map[string]float64{"seg": 2, "ter": 2}, MinutosTopico: 60}
	tarefas := []tarefaPlanejada{{ID: "a", Minutos: 60}, {ID: "b", Minutos: 60}, {ID: "r", Data: "2026-09-14", Revisao: true}}
	plano, err := distribuirComRevisoes(tarefas, nil, cfg, "2026-09-14")
	if err != nil {
		t.Fatal(err)
	}
	datas := map[string]string{}
	carga := map[string]int{}
	for _, p := range plano {
		datas[p.ID] = p.Data
		carga[p.Data] += p.Minutos
	}
	if datas["r"] != "2026-09-14" || datas["a"] != "2026-09-14" || datas["b"] != "2026-09-15" {
		t.Fatal(datas)
	}
	if carga["2026-09-14"] != 90 {
		t.Fatal(carga)
	}
	repetido, err := distribuirComRevisoes(plano, nil, cfg, "2026-09-14")
	if err != nil || !reflect.DeepEqual(plano, repetido) {
		t.Fatal("replanejamento repetido alterou o plano", err)
	}
}

func TestMuitasRevisoesRespeitamDiasELimite(t *testing.T) {
	cfg := Configuracao{Horas: map[string]float64{"seg": 2, "ter": 2, "qua": 2}, MaxTopicosDia: 2}
	tarefas := []tarefaPlanejada{{ID: "r1", Data: "2026-09-12", Grupo: "a", Revisao: true}, {ID: "r2", Data: "2026-09-12", Grupo: "b", Revisao: true}, {ID: "r3", Data: "2026-09-13", Grupo: "a", Revisao: true}, {ID: "t", Minutos: 60}}
	plano, err := distribuirComRevisoes(tarefas, nil, cfg, "2026-09-12")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range plano {
		if p.Data < "2026-09-14" {
			t.Fatal("usou fim de semana", p)
		}
	}
	if plano[0].Data != "2026-09-14" || plano[1].Data != "2026-09-14" || plano[2].Data != "2026-09-15" || plano[3].Data != "2026-09-15" {
		t.Fatal(plano)
	}
}

func TestPreservaCargaConcluidaEDuracaoReal(t *testing.T) {
	cfg := Configuracao{Horas: map[string]float64{"seg": 2, "ter": 2}}
	fixos := map[string]ocupacaoDia{"2026-09-14": {Minutos: 60, Quantidade: 1}}
	plano, err := distribuirComRevisoes([]tarefaPlanejada{{ID: "t", Minutos: 90}, {ID: "r", Data: "2026-09-14", Revisao: true}}, fixos, cfg, "2026-09-14")
	if err != nil {
		t.Fatal(err)
	}
	if plano[0].Data != "2026-09-14" || plano[1].Data != "2026-09-15" {
		t.Fatal(plano)
	}
	if fixos["2026-09-14"].Minutos != 60 {
		t.Fatal("alterou carga original")
	}
}

func TestSemCapacidadeNaoRetornaPlanoParcial(t *testing.T) {
	for _, cfg := range []Configuracao{{Horas: map[string]float64{}}, {Horas: map[string]float64{"seg": .25}}} {
		plano, err := distribuirComRevisoes([]tarefaPlanejada{{ID: "r", Data: "2026-09-14", Revisao: true}}, nil, cfg, "2026-09-14")
		if err == nil || plano != nil {
			t.Fatal("aceitou sobrecarga", plano, err)
		}
	}
}

func TestNaoAntecipaRevisao(t *testing.T) {
	plano, err := distribuirComRevisoes([]tarefaPlanejada{{ID: "r", Data: "2026-09-21", Revisao: true}}, nil, Configuracao{Horas: map[string]float64{"seg": 1}}, "2026-09-14")
	if err != nil || plano[0].Data != "2026-09-21" {
		t.Fatal(plano, err)
	}
}
