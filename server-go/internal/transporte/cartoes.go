package transporte

import (
	"net/http"
	"chega-junto-concurseiro-web/internal/cartoes"
	"chega-junto-concurseiro-web/internal/modelos"
)

func (s *Servidor) listarBaralhos(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.cartoes.ListarBaralhos(r.Context(), aluno, c.Usuario.ID, c.Usuario.Papel, r.URL.Query().Get("concursoId"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os baralhos.")
		return
	}
	responder(w, 200, map[string]any{"baralhos": lista})
}

func (s *Servidor) criarBaralho(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e cartoes.EntradaBaralho
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.cartoes.CriarBaralho(r.Context(), aluno, c.Usuario.ID, c.Usuario.Papel, e)
	if err != nil {
		responderErro(w, 400, "BARALHO_INVALIDO", "Não foi possível criar o baralho.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "criar", "baralho_cartoes", id, e)
	responder(w, 201, map[string]any{"id": id})
}

func (s *Servidor) alterarBaralho(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e cartoes.AlteracaoBaralho
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.cartoes.AlterarBaralho(r.Context(), aluno, c.Usuario.ID, c.Usuario.Papel, r.PathValue("baralhoId"), e); err != nil {
		responderErro(w, 403, "BARALHO_NAO_EDITAVEL", "O baralho não pode ser alterado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "alterar", "baralho_cartoes", r.PathValue("baralhoId"), e)
	responder(w, 200, map[string]any{"atualizado": true})
}

func (s *Servidor) desativarBaralho(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	if err := s.cartoes.DesativarBaralho(r.Context(), aluno, c.Usuario.ID, c.Usuario.Papel, r.PathValue("baralhoId")); err != nil {
		responderErro(w, 403, "BARALHO_NAO_EDITAVEL", "O baralho não pode ser desativado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "desativar", "baralho_cartoes", r.PathValue("baralhoId"), nil)
	responder(w, 200, map[string]any{"desativado": true})
}

func (s *Servidor) criarCartao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e cartoes.EntradaCartao
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.cartoes.CriarCartao(r.Context(), aluno, c.Usuario.ID, c.Usuario.Papel, r.PathValue("baralhoId"), e)
	if err != nil {
		responderErro(w, 403, "CARTAO_NAO_EDITAVEL", "Não foi possível criar o cartão neste baralho.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "criar", "cartao_estudo", id, e)
	responder(w, 201, map[string]any{"id": id})
}

func (s *Servidor) alterarCartao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e cartoes.AlteracaoCartao
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.cartoes.AlterarCartao(r.Context(), aluno, c.Usuario.ID, c.Usuario.Papel, r.PathValue("cartaoId"), e); err != nil {
		responderErro(w, 403, "CARTAO_NAO_EDITAVEL", "O cartão não pode ser alterado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "alterar", "cartao_estudo", r.PathValue("cartaoId"), e)
	responder(w, 200, map[string]any{"atualizado": true})
}

func (s *Servidor) desativarCartao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	if err := s.cartoes.DesativarCartao(r.Context(), aluno, c.Usuario.ID, c.Usuario.Papel, r.PathValue("cartaoId")); err != nil {
		responderErro(w, 403, "CARTAO_NAO_EDITAVEL", "O cartão não pode ser desativado.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "desativar", "cartao_estudo", r.PathValue("cartaoId"), nil)
	responder(w, 200, map[string]any{"desativado": true})
}

func (s *Servidor) revisarCartao(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	if c.Usuario.Papel != "aluno" {
		responderErro(w, 403, "REVISAO_NAO_PERMITIDA", "Somente o aluno pode registrar sua revisão.")
		return
	}
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e struct {
		Qualidade  int    `json:"qualidade"`
		ConcursoID string `json:"concursoId"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	revisao, err := s.cartoes.Revisar(r.Context(), aluno, r.PathValue("cartaoId"), e.Qualidade, e.ConcursoID)
	if err != nil {
		responderErro(w, 400, "REVISAO_INVALIDA", "Não foi possível registrar a revisão.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "revisar", "cartao_estudo", r.PathValue("cartaoId"), e)
	responder(w, 201, revisao)
}

func (s *Servidor) listarCartoesPendentes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.cartoes.ListarPendentes(r.Context(), aluno)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os cartões pendentes.")
		return
	}
	responder(w, 200, map[string]any{"cartoes": lista})
}

func (s *Servidor) listarHistoricoCartoes(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.cartoes.ListarHistorico(r.Context(), aluno, r.URL.Query().Get("concursoId"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar o histórico dos cartões.")
		return
	}
	responder(w, 200, map[string]any{"revisoes": lista})
}
