package transporte

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"chega-junto-concurseiro-web/internal/materiaisapoio"
	"chega-junto-concurseiro-web/internal/modelos"
)

func (s *Servidor) listarMateriaisApoio(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	editalID := r.URL.Query().Get("editalId")
	materiais, err := s.materiaisApoio.ListarAluno(r.Context(), aluno, editalID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os materiais de apoio.")
		return
	}
	if pagina, limite, ok := paginacaoSolicitada(r); ok {
		itens, dados := paginar(materiais, pagina, limite)
		responder(w, 200, map[string]any{"materiais": itens, "paginacao": dados})
		return
	}
	responder(w, 200, map[string]any{"materiais": materiais})
}

func (s *Servidor) listarMateriaisApoioMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	editalID := r.URL.Query().Get("editalId")
	materiais, err := s.materiaisApoio.ListarMentor(r.Context(), c.Usuario.ID, editalID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar os materiais de apoio do mentor.")
		return
	}
	if pagina, limite, ok := paginacaoSolicitada(r); ok {
		itens, dados := paginar(materiais, pagina, limite)
		responder(w, 200, map[string]any{"materiais": itens, "paginacao": dados})
		return
	}
	responder(w, 200, map[string]any{"materiais": materiais})
}

func (s *Servidor) criarMaterialApoioMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		salvarMaterialArquivo(w, r, c, s)
		return
	}
	var e materiaisapoio.EntradaMaterial
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.materiaisApoio.Criar(r.Context(), c.Usuario.ID, e)
	if err != nil {
		responderErro(w, 400, "MATERIAL_INVALIDO", "Não foi possível criar o material de apoio.")
		return
	}
	responder(w, 201, map[string]any{"id": id})
}

func (s *Servidor) alterarMaterialApoioMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		salvarMaterialArquivo(w, r, c, s)
		return
	}
	var e materiaisapoio.EntradaMaterial
	if !decodificar(w, r, &e) {
		return
	}
	materialID := r.PathValue("materialId")
	if err := s.materiaisApoio.Alterar(r.Context(), c.Usuario.ID, materialID, e); err != nil {
		responderErro(w, 400, "MATERIAL_INVALIDO", "Não foi possível alterar o material de apoio.")
		return
	}
	responder(w, 200, map[string]any{"atualizado": true})
}

func (s *Servidor) baixarArquivoMaterialApoio(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	materialID := r.PathValue("materialId")
	var material *materiaisapoio.Material
	var err error
	if c.Usuario.Papel == "mentor" {
		material, err = s.materiaisApoio.ObterParaMentor(r.Context(), c.Usuario.ID, materialID)
	} else {
		material, err = s.materiaisApoio.ObterParaAluno(r.Context(), c.Usuario.ID, materialID)
	}
	if err != nil || material == nil || material.ArquivoCaminho == nil || material.ArquivoNome == nil || material.ArquivoMime == nil {
		responderErro(w, 404, "MATERIAL_NAO_ENCONTRADO", "Arquivo de material não encontrado.")
		return
	}
	if _, err := os.Stat(*material.ArquivoCaminho); err != nil {
		responderErro(w, 404, "MATERIAL_NAO_ENCONTRADO", "Arquivo de material não encontrado.")
		return
	}
	w.Header().Set("Content-Type", *material.ArquivoMime)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filepath.Base(*material.ArquivoNome)}))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "private, no-store")
	http.ServeFile(w, r, *material.ArquivoCaminho)
}

func (s *Servidor) desativarMaterialApoioMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	materialID := r.PathValue("materialId")
	if err := s.materiaisApoio.Desativar(r.Context(), c.Usuario.ID, materialID); err != nil {
		responderErro(w, 400, "MATERIAL_INVALIDO", "Não foi possível desativar o material de apoio.")
		return
	}
	responder(w, 200, map[string]any{"desativado": true})
}
