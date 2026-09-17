package transporte

import (
	"net/http"
	"time"
)

// A ação principal já foi persistida. Uma falha no planejamento deve ser
// informada sem induzir o cliente a repetir a conclusão ou o cadastro.
func (s *Servidor) replanejarAposRevisoes(r *http.Request, aluno string, amanha bool) map[string]any {
	data := time.Now().In(s.cfg.FusoHorario)
	if amanha {
		data = data.AddDate(0, 0, 1)
	}
	if err := s.cronogramas.ReprogramarPendentes(r.Context(), aluno, data.Format("2006-01-02")); err != nil {
		s.log.Error("replanejamento após revisão", "aluno", aluno, "erro", err)
		return map[string]any{"replanejado": false, "avisoReplanejamento": "Alteração salva, mas o cronograma não foi reorganizado: " + err.Error() + ". Ajuste a carga e use Reprogramar Pendentes."}
	}
	return map[string]any{"replanejado": true}
}
