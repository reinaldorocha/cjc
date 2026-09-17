package cursos

import (
	"testing"
	"time"
)

func TestLiberacaoPorDataEPrerequisito(t *testing.T) {
	fuso := time.FixedZone("Fortaleza", -3*3600)
	base := Curso{ID: idLegado("curso"), Aulas: []Aula{
		{Titulo: "Inicial", Modulo: "1", VideoID: "abcdefghijk"},
		{Titulo: "Agendada", Modulo: "2", VideoID: "abcdefghijk", LiberarEm: "2026-09-11", ExigeAnterior: true, PDFID: idLegado("pdf"), Questoes: []string{idLegado("q")}},
	}}
	if err := normalizarAulas(&base); err != nil {
		t.Fatal(err)
	}
	p := map[string]Progresso{base.Aulas[0].ID: {VideoID: "abcdefghijk", Concluida: true}}
	for _, caso := range []struct {
		instante  string
		bloqueada bool
	}{{"2026-09-11T02:59:59Z", true}, {"2026-09-11T03:00:00Z", false}} {
		agora, _ := time.Parse(time.RFC3339, caso.instante)
		c := base
		c.Aulas = append([]Aula{}, base.Aulas...)
		aplicarAprendizagem(&c, p, agora.In(fuso))
		if c.Aulas[1].Bloqueada != caso.bloqueada {
			t.Fatal("data deve usar fuso configurado", caso)
		}
		if caso.bloqueada && (c.Aulas[1].VideoID != "" || c.Aulas[1].PDFID != "" || len(c.Aulas[1].Questoes) != 0) {
			t.Fatal("mídia bloqueada exposta")
		}
	}
	c := base
	c.Aulas = append([]Aula{}, base.Aulas...)
	aplicarAprendizagem(&c, nil, time.Date(2026, 9, 12, 0, 0, 0, 0, fuso))
	if !c.Aulas[1].Bloqueada {
		t.Fatal("pré-requisito ignorado")
	}
}

func TestIdentidadePreservada(t *testing.T) {
	c := Curso{ID: idLegado("curso"), Aulas: []Aula{{Modulo: "A", VideoID: "abcdefghijk"}, {Modulo: "B", VideoID: "abcdefghijk"}}}
	if err := normalizarAulas(&c); err != nil {
		t.Fatal(err)
	}
	id := c.Aulas[0].ID
	c.Aulas[0], c.Aulas[1] = c.Aulas[1], c.Aulas[0]
	c.Aulas[1].Modulo = "Renomeado"
	if err := normalizarAulas(&c); err != nil {
		t.Fatal(err)
	}
	if c.Aulas[1].ID != id {
		t.Fatal("identidade alterada ao reordenar")
	}
	aplicarAprendizagem(&c, map[string]Progresso{id: {VideoID: "abcdefghijk", Posicao: 37, AtualizadoEm: "2026-09-10T12:00:00Z"}}, time.Now())
	if c.Resumo.RetomarAulaID != id || c.Aulas[1].Progresso.Posicao != 37 {
		t.Fatal("retomada perdida")
	}
}
