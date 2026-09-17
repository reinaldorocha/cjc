package transporte

import (
	"context"
	"net/http"
	"track-concursos-web/internal/modelos"
)

func (s *Servidor) resolverAluno(ctx context.Context, w http.ResponseWriter, c modelos.ContextoAutenticado, solicitado string) (string, bool) {
	id, err := s.mentoria.ResolverAluno(ctx, c.Usuario.ID, c.Usuario.Papel, solicitado)
	if err != nil {
		responderErro(w, 404, "ALUNO_NAO_ENCONTRADO", "Aluno não encontrado no seu escopo de acesso.")
		return "", false
	}
	return id, true
}
