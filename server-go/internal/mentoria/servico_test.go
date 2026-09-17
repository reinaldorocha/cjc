package mentoria

import "testing"

func TestClassificarRiscoSemEstudo(t *testing.T) {
	nivel, motivo := classificarRisco(ItemRadarAluno{DiasSemEstudar: 999})
	if nivel != "vermelho" || motivo != "Nenhum estudo registrado" {
		t.Fatalf("classificação inesperada: %s, %s", nivel, motivo)
	}
}

func TestClassificarRiscoPorRendimento(t *testing.T) {
	nivel, _ := classificarRisco(ItemRadarAluno{QuestoesResolvidas: 10, QuestoesAcertos: 4, TaxaAcerto: 40})
	if nivel != "vermelho" {
		t.Fatalf("nível = %s, esperado vermelho", nivel)
	}
}

func TestClassificarRiscoEmDia(t *testing.T) {
	nivel, motivo := classificarRisco(ItemRadarAluno{})
	if nivel != "verde" || motivo != "Ritmo constante e em dia" {
		t.Fatalf("classificação inesperada: %s, %s", nivel, motivo)
	}
}
