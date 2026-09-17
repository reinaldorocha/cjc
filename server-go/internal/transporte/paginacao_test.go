package transporte

import "testing"

func TestPaginarMantemMetadadosELimites(t *testing.T) {
	itens, dados := paginar([]int{1, 2, 3, 4, 5}, 2, 2)
	if len(itens) != 2 || itens[0] != 3 || dados.Total != 5 || dados.TotalPaginas != 3 || dados.Pagina != 2 {
		t.Fatalf("resultado inesperado: %#v %#v", itens, dados)
	}
	ultimos, dados := paginar([]int{1, 2, 3}, 9, 2)
	if len(ultimos) != 1 || ultimos[0] != 3 || dados.Pagina != 2 {
		t.Fatalf("última página inesperada: %#v %#v", ultimos, dados)
	}
}
