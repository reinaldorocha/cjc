package cursos

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Progresso struct {
	AulaID       string  `json:"aulaId"`
	VideoID      string  `json:"-"`
	Posicao      float64 `json:"posicao"`
	Duracao      float64 `json:"duracao"`
	Concluida    bool    `json:"concluida"`
	AtualizadoEm string  `json:"atualizadoEm"`
}
type Resumo struct {
	Total           int    `json:"total"`
	Concluidas      int    `json:"concluidas"`
	Percentual      int    `json:"percentual"`
	RetomarAulaID   string `json:"retomarAulaId"`
	UltimaAtividade string `json:"ultimaAtividade"`
}

func idLegado(valor string) string {
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(valor)))
	return h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

// Identidades determinísticas para o conteúdo anterior à migração 021.
func normalizarAulas(c *Curso) error {
	if c.Aulas == nil {
		c.Aulas = []Aula{}
	}
	aulas := map[string]bool{}
	modulos := map[string]Aula{}
	nomes := map[string]string{}
	fechados := map[string]bool{}
	anterior := ""
	for i := range c.Aulas {
		a := &c.Aulas[i]
		if a.Modulo == "" {
			a.Modulo = "Aulas"
		}
		if a.ModuloID == "" {
			a.ModuloID = idLegado(c.ID + "/modulo/" + a.Modulo)
		}
		if a.ID == "" {
			a.ID = idLegado(fmt.Sprintf("%s/aula/%d/%s", c.ID, i, a.VideoID))
		}
		if !uuidValido.MatchString(a.ID) || !uuidValido.MatchString(a.ModuloID) || aulas[a.ID] {
			return ErrEntrada
		}
		aulas[a.ID] = true
		if a.Questoes == nil {
			a.Questoes = []string{}
		}
		if nomeID, ok := nomes[a.Modulo]; ok && nomeID != a.ModuloID {
			return ErrEntrada
		}
		nomes[a.Modulo] = a.ModuloID
		if ref, ok := modulos[a.ModuloID]; ok {
			if ref.Modulo != a.Modulo || ref.LiberarEm != a.LiberarEm || ref.ExigeAnterior != a.ExigeAnterior || ref.ModuloCapaURL != a.ModuloCapaURL || ref.ModuloDescricao != a.ModuloDescricao {
				return ErrEntrada
			}
		} else {
			modulos[a.ModuloID] = *a
		}
		if anterior != "" && anterior != a.ModuloID {
			fechados[anterior] = true
		}
		if fechados[a.ModuloID] {
			return ErrEntrada
		}
		anterior = a.ModuloID
	}
	return nil
}

func aplicarAprendizagem(c *Curso, progresso map[string]Progresso, agora time.Time) {
	c.Resumo = Resumo{Total: len(c.Aulas)}
	anteriorConcluido := true
	moduloConcluido := true
	moduloID := ""
	bloqueio := ""
	ultimaRetomada := ""
	for i := range c.Aulas {
		a := &c.Aulas[i]
		if moduloID != a.ModuloID {
			if moduloID != "" {
				anteriorConcluido = anteriorConcluido && moduloConcluido
			}
			moduloConcluido = true
			moduloID = a.ModuloID
			bloqueio = ""
			if a.LiberarEm != "" && agora.Format("2006-01-02") < a.LiberarEm {
				bloqueio = "Disponível em " + a.LiberarEm
			}
			if a.ExigeAnterior && !anteriorConcluido {
				bloqueio = "Conclua os módulos anteriores para liberar"
			}
		}
		a.Progresso = nil
		if p, ok := progresso[a.ID]; ok && p.VideoID == a.VideoID {
			a.Progresso = &p
			if p.Concluida {
				c.Resumo.Concluidas++
			} else {
				moduloConcluido = false
			}
			if p.AtualizadoEm > c.Resumo.UltimaAtividade {
				c.Resumo.UltimaAtividade = p.AtualizadoEm
			}
			if !p.Concluida && bloqueio == "" && p.AtualizadoEm > ultimaRetomada {
				c.Resumo.RetomarAulaID = a.ID
				ultimaRetomada = p.AtualizadoEm
			}
		} else {
			moduloConcluido = false
		}
		a.Bloqueada = bloqueio != ""
		a.MotivoBloqueio = bloqueio
		if a.Bloqueada {
			a.VideoID = ""
			a.PDFID = ""
			a.PDFNome = ""
			a.Questoes = []string{}
		}
	}
	if c.Resumo.Total > 0 {
		c.Resumo.Percentual = c.Resumo.Concluidas * 100 / c.Resumo.Total
	}
	if c.Resumo.RetomarAulaID == "" {
		for _, a := range c.Aulas {
			if !a.Bloqueada && (a.Progresso == nil || !a.Progresso.Concluida) {
				c.Resumo.RetomarAulaID = a.ID
				break
			}
		}
	}
	if c.Resumo.RetomarAulaID == "" {
		for _, a := range c.Aulas {
			if !a.Bloqueada {
				c.Resumo.RetomarAulaID = a.ID
				break
			}
		}
	}
	c.Destinatarios = []string{}
}
