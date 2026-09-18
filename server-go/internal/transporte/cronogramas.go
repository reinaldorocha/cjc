package transporte

import (
	"net/http"
	"time"
	"chega-junto-concurseiro-web/internal/cronogramas"
	"chega-junto-concurseiro-web/internal/modelos"
)

func (s *Servidor) obterCronograma(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	concursoID := r.URL.Query().Get("concursoId")
	dados, err := s.cronogramas.ObterAtivo(r.Context(), aluno, concursoID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar o cronograma.")
		return
	}
	responder(w, 200, map[string]any{"cronograma": dados})
}
func (s *Servidor) salvarCronograma(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e cronogramas.Entrada
	if !decodificar(w, r, &e) {
		return
	}
	var id string
	var err error
	if c.Usuario.Papel == "aluno" {
		if !c.Usuario.PermiteCronogramaInteligente {
			responderErro(w, 403, "CRONOGRAMA_NAO_PERMITIDO", "O mentor não habilitou o cronograma inteligente.")
			return
		}
		id, err = s.cronogramas.Gerar(r.Context(), aluno, c.Usuario.ID, e)
	} else {
		id, err = s.cronogramas.Salvar(r.Context(), aluno, c.Usuario.ID, e)
	}
	if err != nil {
		responderErro(w, 400, "CRONOGRAMA_INVALIDO", "Não foi possível salvar o cronograma.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "substituir", "cronograma", id, nil)
	responder(w, 200, map[string]any{"id": id})
}
func (s *Servidor) gerarCronograma(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	if c.Usuario.Papel == "aluno" && !c.Usuario.PermiteCronogramaInteligente {
		responderErro(w, 403, "CRONOGRAMA_NAO_PERMITIDO", "O mentor não habilitou o cronograma inteligente.")
		return
	}
	var e cronogramas.Entrada
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.cronogramas.Gerar(r.Context(), aluno, c.Usuario.ID, e)
	if err != nil {
		responderErro(w, 400, "CRONOGRAMA_INVALIDO", "Não foi possível gerar o cronograma.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "gerar", "cronograma", id, nil)
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarItemCronograma(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e cronogramas.AlteracaoItem
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.cronogramas.AlterarItem(r.Context(), aluno, r.PathValue("itemId"), e, c.Usuario.Papel == "mentor"); err != nil {
		responderErro(w, 400, "ITEM_INVALIDO", "Não foi possível alterar o item.")
		return
	}
	responder(w, 200, map[string]any{"atualizado": true})
}
func (s *Servidor) criarItemCronograma(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var item cronogramas.Item
	if !decodificar(w, r, &item) {
		return
	}
	if err := s.cronogramas.CriarItemNoCronograma(r.Context(), aluno, item); err != nil {
		responderErro(w, 400, "ITEM_INVALIDO", "Não foi possível criar o item no cronograma.")
		return
	}
	responder(w, 201, map[string]any{"criado": true})
}
func (s *Servidor) reprogramarCronograma(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var req struct {
		APartirDe  string `json:"aPartirDe"`
		ConcursoID string `json:"concursoId"`
	}
	if r.Body != nil {
		if !decodificar(w, r, &req) {
			return
		}
	}
	if req.APartirDe == "" {
		req.APartirDe = time.Now().In(s.cfg.FusoHorario).Format("2006-01-02")
	}
	if err := s.cronogramas.ReprogramarPendentes(r.Context(), aluno, req.APartirDe, req.ConcursoID); err != nil {
		responderErro(w, 400, "ERRO_REPROGRAMAR", err.Error())
		return
	}
	responder(w, 200, map[string]any{"reprogramado": true})
}
func (s *Servidor) calendarioCronograma(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	concursoID := r.URL.Query().Get("concursoId")
	itens, err := s.cronogramas.Calendario(r.Context(), aluno, concursoID, r.URL.Query().Get("inicio"), r.URL.Query().Get("fim"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar o calendário.")
		return
	}
	responder(w, 200, map[string]any{"itens": itens})
}
