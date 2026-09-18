package transporte

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (s *Servidor) servirFrontend(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	relativo := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
	arquivo := filepath.Join(s.cfg.FrontendPath, filepath.FromSlash(relativo))
	if info, err := os.Stat(arquivo); err == nil && !info.IsDir() {
		http.ServeFile(w, r, arquivo)
		return
	}

	http.ServeFile(w, r, filepath.Join(s.cfg.FrontendPath, "index.html"))
}
