package transporte

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"time"
	"track-concursos-web/internal/cursos"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) enviarPDFCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	r.Body = http.MaxBytesReader(w, r.Body, cursos.LimitePDF+(1<<20))
	defer r.Body.Close()
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		responderErro(w, 400, "PDF_INVALIDO", "Envie um PDF de até 10 MB.")
		return
	}
	defer r.MultipartForm.RemoveAll()
	file, header, err := r.FormFile("arquivo")
	if err != nil {
		responderErro(w, 400, "PDF_INVALIDO", "Selecione o PDF da aula.")
		return
	}
	defer file.Close()
	conteudo, err := io.ReadAll(io.LimitReader(file, cursos.LimitePDF+1))
	if err != nil {
		responderErro(w, 400, "PDF_INVALIDO", "Não foi possível ler o PDF.")
		return
	}
	id, nome, err := s.cursos.SalvarPDF(r.Context(), c.Usuario.ID, header.Filename, conteudo)
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	responder(w, 201, map[string]string{"id": id, "nome": nome})
}

func (s *Servidor) baixarPDFCurso(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	mentor := r.PathValue("alunoId") == ""
	usuario := c.Usuario.ID
	if !mentor {
		var ok bool
		usuario, ok = s.resolverAluno(r.Context(), w, c, r.PathValue("alunoId"))
		if !ok {
			return
		}
	}
	nome, conteudo, err := s.cursos.ObterPDF(r.Context(), usuario, mentor, r.PathValue("cursoId"), r.PathValue("pdfId"))
	if err != nil {
		s.erroCurso(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": nome}))
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeContent(w, r, nome, time.Time{}, bytes.NewReader(conteudo))
}
