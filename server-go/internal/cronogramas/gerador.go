package cronogramas

import (
	"database/sql"
	"time"
)

func gerarCiclo(materias []materiaInfo) []Item {
	itens := []Item{}
	for i, m := range materias {
		duracao := 30
		if m.prioridade >= 70 {
			duracao = 90
		} else if m.prioridade >= 55 {
			duracao = 70
		} else if m.prioridade >= 40 {
			duracao = 50
		}
		pos, id := i+1, m.id
		itens = append(itens, Item{PosicaoCiclo: &pos, MateriaID: &id, Nome: m.nome, DuracaoMinutos: duracao, Prioridade: m.prioridade, Ordem: pos, Situacao: "pendente"})
	}
	return itens
}

func gerarAgenda(materias []materiaInfo, unidades map[string][]unidade, cfg Configuracao, fuso *time.Location) []Item {
	if len(materias) == 0 {
		return []Item{}
	}
	minutos := cfg.MinutosTopico
	if minutos <= 0 {
		minutos = 60
	}
	horas := cfg.Horas
	if horas == nil {
		horas = map[string]float64{"seg": 2, "ter": 2, "qua": 2, "qui": 2, "sex": 2}
	}
	chaves := []string{"dom", "seg", "ter", "qua", "qui", "sex", "sab"}
	hoje := time.Now().In(fuso)
	hoje = time.Date(hoje.Year(), hoje.Month(), hoje.Day(), 0, 0, 0, 0, fuso)
	repeticoes := cfg.RepeticoesEdital
	if repeticoes <= 0 {
		repeticoes = 1
	}
	ponteiros := map[string]int{}
	indiceMateria := 0
	itens := []Item{}
	for dia := 0; dia < 365*3; dia++ {
		data := hoje.AddDate(0, 0, dia)
		capacidade := int(horas[chaves[int(data.Weekday())]]*60) / minutos
		if horas[chaves[int(data.Weekday())]] > 0 && capacidade < 1 {
			capacidade = 1
		}
		if cfg.MaxTopicosDia > 0 && capacidade > cfg.MaxTopicosDia {
			capacidade = cfg.MaxTopicosDia
		}
		if capacidade <= 0 {
			continue
		}
		usadas := map[string]bool{}
		for ordemDia := 1; ordemDia <= capacidade; ordemDia++ {
			escolhida := -1
			for tentativa := 0; tentativa < len(materias); tentativa++ {
				candidato := (indiceMateria + tentativa) % len(materias)
				m := materias[candidato]
				if ponteiros[m.id] < len(unidades[m.id])*repeticoes && !usadas[m.id] {
					escolhida = candidato
					break
				}
			}
			if escolhida == -1 {
				for candidato, m := range materias {
					if ponteiros[m.id] < len(unidades[m.id])*repeticoes {
						escolhida = candidato
						break
					}
				}
			}
			if escolhida == -1 {
				return itens
			}
			m := materias[escolhida]
			usadas[m.id] = true
			indiceMateria = (escolhida + 1) % len(materias)
			u := unidades[m.id][ponteiros[m.id]%len(unidades[m.id])]
			ponteiros[m.id]++
			dataTexto := data.Format("2006-01-02")
			mat := u.materiaID
			item := Item{DataPlanejada: &dataTexto, MateriaID: &mat, Nome: u.nome, DuracaoMinutos: minutos, Prioridade: u.prioridade, Ordem: ordemDia, Situacao: "pendente"}
			if u.topicoID != "" {
				v := u.topicoID
				item.TopicoID = &v
			}
			if u.subtopicoID != "" {
				v := u.subtopicoID
				item.SubtopicoID = &v
			}
			itens = append(itens, item)
		}
	}
	return itens
}
func valor(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
func ordem(v, indice int) int {
	if v > 0 {
		return v
	}
	return indice + 1
}
func dataPtr(v sql.NullTime) *string {
	if !v.Valid {
		return nil
	}
	x := v.Time.Format("2006-01-02")
	return &x
}
func intPtr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	x := int(v.Int64)
	return &x
}
func stringPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}
