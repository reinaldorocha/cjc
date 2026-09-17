package materiaisapoio

type EntradaArquivo struct {
	Titulo      string  `json:"titulo"`
	Descricao   *string `json:"descricao"`
	Escopo      string  `json:"escopo"`
	EditalID    *string `json:"editalId"`
	NomeArquivo string  `json:"nomeArquivo"`
	Tipo        string  `json:"tipo"`
	Pasta       *string `json:"pasta"`
}
