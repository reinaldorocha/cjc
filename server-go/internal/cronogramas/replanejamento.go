package cronogramas

import (
	"fmt"
	"math"
	"time"
)

// A revisão reserva 30 minutos no mesmo orçamento dos assuntos.
const MinutosRevisao = 30

type tarefaPlanejada struct {
	ID, Data, Grupo string
	Minutos         int
	Revisao         bool
}
type ocupacaoDia struct{ Minutos, Quantidade int }

func distribuirComRevisoes(tarefas []tarefaPlanejada, fixos map[string]ocupacaoDia, cfg Configuracao, inicio string) ([]tarefaPlanejada, error) {
	data, err := time.Parse("2006-01-02", inicio)
	if err != nil {
		return nil, err
	}
	horas := cfg.Horas
	if horas == nil {
		horas = map[string]float64{"seg": 2, "ter": 2, "qua": 2, "qui": 2, "sex": 2}
	}
	chaves := []string{"dom", "seg", "ter", "qua", "qui", "sex", "sab"}
	maior := 0
	for _, h := range horas {
		if math.IsNaN(h) || math.IsInf(h, 0) || h < 0 || h > 24 {
			return nil, fmt.Errorf("carga horária diária inválida")
		}
		if int(h*60) > maior {
			maior = int(h * 60)
		}
	}
	usados := map[string]ocupacaoDia{}
	for d, o := range fixos {
		usados[d] = o
	}
	resultado := make([]tarefaPlanejada, 0, len(tarefas))
	ultimaRevisao := map[string]string{}
	// Primeiro reservar revisões; o espaço restante recebe os assuntos na ordem existente.
	for _, revisao := range []bool{true, false} {
		cursor := data
		for _, t := range tarefas {
			if t.Revisao != revisao {
				continue
			}
			if t.Minutos <= 0 {
				t.Minutos = cfg.MinutosTopico
				if t.Minutos <= 0 {
					t.Minutos = 60
				}
			}
			if t.Revisao {
				t.Minutos = MinutosRevisao
			}
			if t.Minutos > maior {
				return nil, fmt.Errorf("a atividade de %d minutos não cabe na carga diária; ajuste a duração ou as horas disponíveis", t.Minutos)
			}
			inicioTarefa := cursor
			if revisao {
				inicioTarefa = data
				if t.Data > inicio {
					inicioTarefa, err = time.Parse("2006-01-02", t.Data)
					if err != nil {
						return nil, err
					}
				}
				if ultima := ultimaRevisao[t.Grupo]; t.Grupo != "" && ultima != "" && inicioTarefa.Format("2006-01-02") <= ultima {
					anterior, _ := time.Parse("2006-01-02", ultima)
					inicioTarefa = anterior.AddDate(0, 0, 1)
				}
			}
			alocada := false
			for dia := inicioTarefa; dia.Before(data.AddDate(10, 0, 0)); dia = dia.AddDate(0, 0, 1) {
				chave := dia.Format("2006-01-02")
				limite := int(horas[chaves[int(dia.Weekday())]] * 60)
				o := usados[chave]
				if limite <= 0 || o.Minutos+t.Minutos > limite || (cfg.MaxTopicosDia > 0 && o.Quantidade >= cfg.MaxTopicosDia) {
					continue
				}
				t.Data = chave
				resultado = append(resultado, t)
				usados[chave] = ocupacaoDia{o.Minutos + t.Minutos, o.Quantidade + 1}
				if revisao && t.Grupo != "" {
					ultimaRevisao[t.Grupo] = chave
				}
				if !revisao {
					cursor = dia
				}
				alocada = true
				break
			}
			if !alocada {
				return nil, fmt.Errorf("não foi possível distribuir todas as atividades; aumente a carga disponível")
			}
		}
	}
	return resultado, nil
}
