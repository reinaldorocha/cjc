package transporte

import (
	"bytes"
	"io"
	"net/http"
	"time"
	"track-concursos-web/internal/cursos"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) enviarCapaCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	r.Body = http.MaxBytesReader(w, r.Body, cursos.LimiteCapa+(1<<20))
	defer r.Body.Close()
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		responderErro(w, 400, "CAPA_INVALIDA", cursos.ErrCapa.Error())
		return
	}
	defer r.MultipartForm.RemoveAll()
	file, _, err := r.FormFile("arquivo")
	if err != nil {
		responderErro(w, 400, "CAPA_INVALIDA", cursos.ErrCapa.Error())
		return
	}
	defer file.Close()
	conteudo, err := io.ReadAll(io.LimitReader(file, cursos.LimiteCapa+1))
	if err != nil {
		responderErro(w, 400, "CAPA_INVALIDA", cursos.ErrCapa.Error())
		return
	}
	url, err := s.cursos.SalvarCapa(r.Context(), c.Usuario.ID, conteudo)
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	responder(w, 201, map[string]string{"url": url})
}

func (s *Servidor) exibirCapaCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	tipo, conteudo, err := s.cursos.ObterCapa(r.Context(), c.Usuario.ID, c.Usuario.Papel == "mentor", r.PathValue("capaId"))
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	w.Header().Set("Content-Type", tipo)
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeContent(w, r, "capa", time.Time{}, bytes.NewReader(conteudo))
}
