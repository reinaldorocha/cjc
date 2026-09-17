package transporte

import (
	"net/http"
	"track-concursos-web/internal/editais"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) catalogarEditais(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	concursoID := r.URL.Query().Get("concursoId")
	if concursoID == "" {
		responderErro(w, 400, "CONCURSO_OBRIGATORIO", "Informe o concurso para consultar o catálogo.")
		return
	}
	lista, err := s.editais.Catalogar(r.Context(), c.Usuario.ID, concursoID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar o catálogo de editais.")
		return
	}
	responder(w, 200, map[string]any{"editais": lista})
}

func (s *Servidor) listarEditais(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	lista, err := s.editais.Listar(r.Context(), aluno, r.URL.Query().Get("concursoId"))
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar o edital.")
		return
	}
	responder(w, 200, map[string]any{"editais": lista})
}
func (s *Servidor) atribuirEdital(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e editais.Entrada
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.editais.Atribuir(r.Context(), aluno, c.Usuario.ID, e)
	if err != nil {
		responderErro(w, 400, "EDITAL_INVALIDO", "Não foi possível atribuir o edital.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "atribuir", "edital", id, nil)
	responder(w, 200, map[string]any{"id": id})
}
func (s *Servidor) criarItemEdital(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e editais.NovoItem
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.editais.CriarItem(r.Context(), aluno, e)
	if err != nil {
		responderErro(w, 400, "ITEM_INVALIDO", "Não foi possível criar o item do edital.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "criar", "item_edital", id, map[string]any{"tipo": e.Tipo})
	responder(w, 201, map[string]any{"id": id})
}
func (s *Servidor) alterarItemEdital(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e editais.AlteracaoItem
	if !decodificar(w, r, &e) {
		return
	}
	id := r.PathValue("itemId")
	if err := s.editais.AlterarItem(r.Context(), aluno, id, e); err != nil {
		responderErro(w, 400, "ITEM_INVALIDO", "Não foi possível alterar o item do edital.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "alterar", "item_edital", id, map[string]any{"tipo": e.Tipo})
	responder(w, 200, map[string]any{"atualizado": true})
}
func (s *Servidor) atualizarProgressoEdital(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var e editais.Progresso
	if !decodificar(w, r, &e) {
		return
	}
	id := r.PathValue("itemId")
	if err := s.editais.AtualizarProgresso(r.Context(), aluno, id, e); err != nil {
		responderErro(w, 400, "PROGRESSO_INVALIDO", "Não foi possível atualizar o progresso.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "progresso", "item_edital", id, map[string]any{"estudado": e.Estudado})
	resultado := s.replanejarAposRevisoes(r, aluno, true)
	resultado["atualizado"] = true
	responder(w, 200, resultado)
}

func (s *Servidor) reordenarItensEdital(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var entrada struct {
		Itens []editais.ItemOrdem `json:"itens"`
	}
	if !decodificar(w, r, &entrada) {
		return
	}
	if err := s.editais.Reordenar(r.Context(), aluno, entrada.Itens); err != nil {
		responderErro(w, 400, "ORDEM_INVALIDA", "Não foi possível reordenar os itens do edital.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "reordenar", "itens_edital", "", map[string]any{"quantidade": len(entrada.Itens)})
	responder(w, 200, map[string]any{"atualizado": true})
}

func (s *Servidor) atualizarProgressoMaterial(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	var entrada editais.EntradaProgressoMaterial
	if !decodificar(w, r, &entrada) {
		return
	}
	materialID := r.PathValue("materialId")
	if err := s.editais.AtualizarProgressoMaterial(r.Context(), aluno, materialID, entrada); err != nil {
		responderErro(w, 400, "PROGRESSO_INVALIDO", "Não foi possível atualizar o material.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, aluno, "progresso", "material_edital", materialID, map[string]any{"concluido": entrada.Concluido})
	responder(w, 200, map[string]any{"atualizado": true})
}
