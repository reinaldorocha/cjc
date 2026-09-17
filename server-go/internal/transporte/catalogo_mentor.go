package transporte

import (
	"net/http"
	"track-concursos-web/internal/cartoes"
	"track-concursos-web/internal/concursos"
	"track-concursos-web/internal/editais"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) catalogarConcursosMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	lista, err := s.concursos.Catalogar(r.Context(), c.Usuario.ID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os concursos.")
		return
	}
	responder(w, 200, map[string]any{"concursos": lista})
}
func (s *Servidor) criarConcursoMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e concursos.Atribuicao
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.concursos.CriarCatalogo(r.Context(), c.Usuario.ID, e)
	if err != nil {
		responderErro(w, 400, "CONCURSO_INVALIDO", "Não foi possível cadastrar o concurso.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, "", "criar", "concurso_catalogo", id, nil)
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) catalogarEditaisMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	lista, err := s.editais.ListarCatalogo(r.Context(), c.Usuario.ID, r.URL.Query().Get("concursoId"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os editais.")
		return
	}
	responder(w, 200, map[string]any{"editais": lista})
}
func (s *Servidor) criarEditalMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e editais.Entrada
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.editais.CriarCatalogo(r.Context(), c.Usuario.ID, e)
	if err != nil {
		responderErro(w, 400, "EDITAL_INVALIDO", "Não foi possível cadastrar o edital.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, "", "criar", "edital_catalogo", id, nil)
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) listarBaralhosMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	lista, err := s.cartoes.ListarBaralhosMentor(r.Context(), c.Usuario.ID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os baralhos.")
		return
	}
	responder(w, 200, map[string]any{"baralhos": lista})
}
func (s *Servidor) criarBaralhoMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e cartoes.EntradaBaralho
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.cartoes.CriarBaralhoMentor(r.Context(), c.Usuario.ID, e)
	if err != nil {
		responderErro(w, 400, "BARALHO_INVALIDO", "Não foi possível criar o baralho.")
		return
	}
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarBaralhoMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e cartoes.AlteracaoBaralho
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.cartoes.AlterarBaralho(r.Context(), "", c.Usuario.ID, "mentor", r.PathValue("baralhoId"), e); err != nil {
		responderErro(w, 403, "BARALHO_NAO_EDITAVEL", "O baralho não pode ser alterado.")
		return
	}
	responder(w, 200, map[string]any{"atualizado": true})
}
func (s *Servidor) desativarBaralhoMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	if err := s.cartoes.DesativarBaralho(r.Context(), "", c.Usuario.ID, "mentor", r.PathValue("baralhoId")); err != nil {
		responderErro(w, 403, "BARALHO_NAO_EDITAVEL", "O baralho não pode ser desativado.")
		return
	}
	responder(w, 200, map[string]any{"desativado": true})
}
func (s *Servidor) criarCartaoMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e cartoes.EntradaCartao
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.cartoes.CriarCartao(r.Context(), "", c.Usuario.ID, "mentor", r.PathValue("baralhoId"), e)
	if err != nil {
		responderErro(w, 400, "CARTAO_INVALIDO", "Não foi possível criar o cartão.")
		return
	}
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarCartaoMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e cartoes.AlteracaoCartao
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.cartoes.AlterarCartao(r.Context(), "", c.Usuario.ID, "mentor", r.PathValue("cartaoId"), e); err != nil {
		responderErro(w, 403, "CARTAO_NAO_EDITAVEL", "O cartão não pode ser alterado.")
		return
	}
	responder(w, 200, map[string]any{"atualizado": true})
}
func (s *Servidor) desativarCartaoMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	if err := s.cartoes.DesativarCartao(r.Context(), "", c.Usuario.ID, "mentor", r.PathValue("cartaoId")); err != nil {
		responderErro(w, 403, "CARTAO_NAO_EDITAVEL", "O cartão não pode ser desativado.")
		return
	}
	responder(w, 200, map[string]any{"desativado": true})
}
