package transporte

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"chega-junto-concurseiro-web/internal/arquivos"
	"chega-junto-concurseiro-web/internal/modelos"
	"chega-junto-concurseiro-web/internal/whitelabel"
)

func (s *Servidor) obterWhiteLabelMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	config, err := s.whitelabel.ObterPorMentor(r.Context(), c.Usuario.ID)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "ERRO_INTERNO", err.Error())
		return
	}
	responder(w, http.StatusOK, map[string]any{"configuracao": whitelabel.URLsProtegidas(config, "/api/v1/mentor/whitelabel")})
}

func (s *Servidor) obterWhiteLabelMentorPorID(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	mentorID := r.PathValue("mentorId")
	if mentorID == "" {
		responderErro(w, http.StatusBadRequest, "DADOS_INVALIDOS", "mentorId é obrigatório")
		return
	}
	config, err := s.whitelabel.ObterPorMentor(r.Context(), mentorID)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "ERRO_INTERNO", err.Error())
		return
	}
	if !s.podeVerWhiteLabel(r, c, mentorID) {
		responderErro(w, http.StatusForbidden, "ACESSO_NEGADO", "Configuração não disponível.")
		return
	}
	responder(w, http.StatusOK, map[string]any{"configuracao": whitelabel.URLsProtegidas(config, "/api/v1/mentores/"+mentorID+"/whitelabel")})
}

func (s *Servidor) salvarWhiteLabelMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var req whitelabel.Config
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderErro(w, http.StatusBadRequest, "DADOS_INVALIDOS", "payload inválido")
		return
	}
	req.MentorID = c.Usuario.ID
	salvo, err := s.whitelabel.Salvar(r.Context(), req)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "ERRO_INTERNO", err.Error())
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, c.Usuario.ID, "salvar", "whitelabel", c.Usuario.ID, salvo)
	responder(w, http.StatusOK, map[string]any{"configuracao": whitelabel.URLsProtegidas(salvo, "/api/v1/mentor/whitelabel")})
}

func (s *Servidor) uploadLogoWhiteLabel(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	r.Body = http.MaxBytesReader(w, r.Body, arquivos.LimiteImagem+(1<<20))
	defer r.Body.Close()
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		responderErro(w, http.StatusBadRequest, "ENTRADA_INVALIDA", "Arquivo muito grande ou formulário inválido.")
		return
	}
	file, header, err := r.FormFile("logo")
	if err != nil {
		responderErro(w, http.StatusBadRequest, "DADOS_INVALIDOS", "Arquivo de logo não enviado.")
		return
	}
	defer file.Close()

	_, err = s.whitelabel.SalvarLogo(r.Context(), c.Usuario.ID, header.Filename, file, s.cfg.UploadsPath)
	if err != nil {
		s.log.Error("falha ao salvar logo", "erro", err, "mentor", c.Usuario.ID)
		if errors.Is(err, arquivos.ErrImagemInvalida) {
			responderErro(w, http.StatusBadRequest, "IMAGEM_INVALIDA", arquivos.ErrImagemInvalida.Error())
			return
		}
		responderErro(w, http.StatusInternalServerError, "ERRO_INTERNO", "Não foi possível salvar o logo.")
		return
	}
	responder(w, http.StatusOK, map[string]string{"url": "/api/v1/mentor/whitelabel/logo"})
}

func (s *Servidor) uploadBannerWhiteLabel(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	r.Body = http.MaxBytesReader(w, r.Body, arquivos.LimiteImagem+(1<<20))
	defer r.Body.Close()
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		responderErro(w, http.StatusBadRequest, "ENTRADA_INVALIDA", "Arquivo muito grande ou formulário inválido.")
		return
	}
	file, header, err := r.FormFile("banner")
	if err != nil {
		responderErro(w, http.StatusBadRequest, "DADOS_INVALIDOS", "Imagem do banner não enviada.")
		return
	}
	defer file.Close()

	_, err = s.whitelabel.SalvarBanner(r.Context(), c.Usuario.ID, header.Filename, file, s.cfg.UploadsPath)
	if err != nil {
		s.log.Error("falha ao salvar banner", "erro", err, "mentor", c.Usuario.ID)
		if errors.Is(err, arquivos.ErrImagemInvalida) {
			responderErro(w, http.StatusBadRequest, "IMAGEM_INVALIDA", arquivos.ErrImagemInvalida.Error())
			return
		}
		responderErro(w, http.StatusInternalServerError, "ERRO_INTERNO", "Não foi possível salvar o banner.")
		return
	}
	responder(w, http.StatusOK, map[string]string{"url": "/api/v1/mentor/whitelabel/banner"})
}

func (s *Servidor) obterWhiteLabelAluno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	alunoID := r.PathValue("alunoId")
	aluno, ok := s.resolverAluno(r.Context(), w, c, alunoID)
	if !ok {
		return
	}
	config, err := s.whitelabel.ObterPorAluno(r.Context(), aluno)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "ERRO_INTERNO", err.Error())
		return
	}
	responder(w, http.StatusOK, map[string]any{"configuracao": whitelabel.URLsProtegidas(config, "/api/v1/alunos/"+aluno+"/whitelabel")})
}

func (s *Servidor) podeVerWhiteLabel(r *http.Request, c modelos.ContextoAutenticado, mentorID string) bool {
	if c.Usuario.Papel == "mentor" {
		return c.Usuario.ID == mentorID
	}
	config, err := s.whitelabel.ObterPorAluno(r.Context(), c.Usuario.ID)
	return err == nil && config.MentorID == mentorID
}

func (s *Servidor) exibirImagemWhiteLabel(w http.ResponseWriter, r *http.Request, mentorID, tipo string) {
	mime, conteudo, err := s.whitelabel.ObterImagem(r.Context(), mentorID, tipo, s.cfg.UploadsPath)
	if err != nil {
		responderErro(w, http.StatusNotFound, "ARQUIVO_NAO_ENCONTRADO", "Imagem não encontrada.")
		return
	}
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeContent(w, r, "imagem", time.Time{}, bytes.NewReader(conteudo))
}

func (s *Servidor) exibirLogoWhiteLabelMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	s.exibirImagemWhiteLabel(w, r, c.Usuario.ID, "logo")
}
func (s *Servidor) exibirBannerWhiteLabelMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	s.exibirImagemWhiteLabel(w, r, c.Usuario.ID, "banner")
}
func (s *Servidor) exibirLogoWhiteLabelMentorPorID(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	mentorID := r.PathValue("mentorId")
	if !s.podeVerWhiteLabel(r, c, mentorID) {
		responderErro(w, http.StatusForbidden, "ACESSO_NEGADO", "Imagem não disponível.")
		return
	}
	s.exibirImagemWhiteLabel(w, r, mentorID, "logo")
}
func (s *Servidor) exibirBannerWhiteLabelMentorPorID(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	mentorID := r.PathValue("mentorId")
	if !s.podeVerWhiteLabel(r, c, mentorID) {
		responderErro(w, http.StatusForbidden, "ACESSO_NEGADO", "Imagem não disponível.")
		return
	}
	s.exibirImagemWhiteLabel(w, r, mentorID, "banner")
}
func (s *Servidor) exibirLogoWhiteLabelAluno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	config, err := s.whitelabel.ObterPorAluno(r.Context(), aluno)
	if err != nil {
		responderErro(w, 404, "ARQUIVO_NAO_ENCONTRADO", "Imagem não encontrada.")
		return
	}
	s.exibirImagemWhiteLabel(w, r, config.MentorID, "logo")
}
func (s *Servidor) exibirBannerWhiteLabelAluno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, ok := s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
	if !ok {
		return
	}
	config, err := s.whitelabel.ObterPorAluno(r.Context(), aluno)
	if err != nil {
		responderErro(w, 404, "ARQUIVO_NAO_ENCONTRADO", "Imagem não encontrada.")
		return
	}
	s.exibirImagemWhiteLabel(w, r, config.MentorID, "banner")
}
