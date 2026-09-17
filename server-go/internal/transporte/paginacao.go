package transporte

import (
	"net/http"
	"strconv"
)

type DadosPaginacao struct {
	Pagina       int `json:"pagina"`
	Limite       int `json:"limite"`
	Total        int `json:"total"`
	TotalPaginas int `json:"totalPaginas"`
}

func paginacaoSolicitada(r *http.Request) (pagina, limite int, solicitada bool) {
	consulta := r.URL.Query()
	if consulta.Get("pagina") == "" && consulta.Get("limite") == "" {
		return 0, 0, false
	}
	pagina, _ = strconv.Atoi(consulta.Get("pagina"))
	limite, _ = strconv.Atoi(consulta.Get("limite"))
	if pagina < 1 {
		pagina = 1
	}
	if limite < 1 {
		limite = 12
	}
	if limite > 100 {
		limite = 100
	}
	return pagina, limite, true
}

func paginar[T any](lista []T, pagina, limite int) ([]T, DadosPaginacao) {
	total := len(lista)
	totalPaginas := (total + limite - 1) / limite
	if totalPaginas == 0 {
		totalPaginas = 1
	}
	if pagina > totalPaginas {
		pagina = totalPaginas
	}
	inicio := (pagina - 1) * limite
	fim := inicio + limite
	if inicio > total {
		inicio = total
	}
	if fim > total {
		fim = total
	}
	return lista[inicio:fim], DadosPaginacao{Pagina: pagina, Limite: limite, Total: total, TotalPaginas: totalPaginas}
}
