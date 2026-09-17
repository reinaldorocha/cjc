package transporte

import (
	"encoding/json"
	"net/http"
)

type corpoErro struct {
	Codigo   string `json:"codigo"`
	Mensagem string `json:"mensagem"`
}

func responder(w http.ResponseWriter, status int, dados any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]any{"dados": dados}); err != nil {
		return
	}
}

func responderErro(w http.ResponseWriter, status int, codigo, mensagem string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]corpoErro{"erro": {Codigo: codigo, Mensagem: mensagem}}); err != nil {
		return
	}
}

func decodificar(w http.ResponseWriter, r *http.Request, destino any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(destino); err != nil {
		responderErro(w, http.StatusBadRequest, "ENTRADA_INVALIDA", "Corpo da requisição inválido.")
		return false
	}
	return true
}
