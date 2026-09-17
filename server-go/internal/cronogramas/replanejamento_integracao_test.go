package cronogramas

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"track-concursos-web/internal/banco"
	"track-concursos-web/internal/configuracao"
	"track-concursos-web/internal/identificador"
)

func TestIntegracaoReplanejamento(t *testing.T) {
	if os.Getenv("CRONOGRAMA_TEST_MYSQL") != "1" {
		t.Skip("defina CRONOGRAMA_TEST_MYSQL=1")
	}
	t.Chdir("../..")
	cfg, err := configuracao.Carregar()
	if err != nil {
		t.Fatal(err)
	}
	cfg.DBNome = fmt.Sprintf("test_replanejar_%d", time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	db, err := banco.Abrir(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer func() {
		if _, err := db.Exec("DROP DATABASE `" + cfg.DBNome + "`"); err != nil {
			t.Error(err)
		}
	}()
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.ExecContext(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
	for _, nome := range []string{"001_estrutura_inicial.sql", "002_dominio_estudos.sql", "003_cronogramas_versionados.sql", "005_rotinas_estudo.sql"} {
		dados, err := os.ReadFile("migracoes/sql/" + nome)
		if err != nil {
			t.Fatal(err)
		}
		exec(string(dados))
	}
	aluno, outro, concurso, c2 := identificador.UUID(), identificador.UUID(), identificador.UUID(), identificador.UUID()
	for _, id := range []string{aluno, outro} {
		exec(`INSERT INTO usuarios (id,nome,email,senha_hash,papel) VALUES (?,?,?,'teste','aluno')`, id, id, id+"@teste.local")
	}
	for _, id := range []string{concurso, c2} {
		exec(`INSERT INTO concursos (id,nome,banca,criado_por) VALUES (?,'Teste','Banca',?)`, id, aluno)
	}
	agenda, outraAgenda := identificador.UUID(), identificador.UUID()
	for i, id := range []string{agenda, outraAgenda} {
		co := concurso
		if i == 1 {
			co = c2
		}
		exec(`INSERT INTO cronogramas (id,aluno_id,concurso_id,tipo,criado_por,configuracao) VALUES (?,?,?,'agendado',?,?)`, id, aluno, co, aluno, `{"horas":{"seg":2,"ter":2,"qua":2},"minutosTopico":60}`)
	}
	a, b, feito, hoje, inalterado := identificador.UUID(), identificador.UUID(), identificador.UUID(), identificador.UUID(), identificador.UUID()
	for _, id := range []string{a, b} {
		exec(`INSERT INTO cronograma_itens (id,cronograma_id,data_planejada,duracao_minutos,ordem) VALUES (?,?,'2026-09-14',60,?)`, id, agenda, map[string]int{a: 1, b: 2}[id])
	}
	exec(`INSERT INTO cronograma_itens (id,cronograma_id,data_planejada,duracao_minutos,situacao,concluido_em) VALUES (?,?,'2026-09-14',30,'concluido','2026-09-14 08:00:00')`, feito, agenda)
	exec(`INSERT INTO cronograma_itens (id,cronograma_id,data_planejada,duracao_minutos) VALUES (?,?,'2026-09-13',60)`, hoje, agenda)
	exec(`INSERT INTO cronograma_itens (id,cronograma_id,data_planejada,duracao_minutos) VALUES (?,?,'2026-09-14',60)`, inalterado, outraAgenda)
	revisao, revisaoOutro := identificador.UUID(), identificador.UUID()
	exec(`INSERT INTO revisoes_programadas (id,aluno_id,concurso_id,proxima_data) VALUES (?,?,?,'2026-09-14')`, revisao, aluno, concurso)
	exec(`INSERT INTO revisoes_programadas (id,aluno_id,concurso_id,proxima_data) VALUES (?,?,?,'2026-09-14')`, revisaoOutro, outro, concurso)
	r := novoRepositorioMySQL(db)
	for tentativa := 0; tentativa < 2; tentativa++ {
		if err = r.ReplanejarCalendario(ctx, aluno, concurso, "2026-09-14", "2026-09-13"); err != nil {
			t.Fatal(err)
		}
		for id, esperado := range map[string]string{a: "2026-09-14", b: "2026-09-15", feito: "2026-09-14", hoje: "2026-09-13", inalterado: "2026-09-14"} {
			var data string
			if err = db.QueryRow(`SELECT DATE_FORMAT(data_planejada,'%Y-%m-%d') FROM cronograma_itens WHERE id=?`, id).Scan(&data); err != nil || data != esperado {
				t.Fatal("data incorreta", id, data, esperado, err)
			}
		}
	}
	var total int
	if err = db.QueryRow(`SELECT COUNT(*) FROM revisoes_programadas`).Scan(&total); err != nil || total != 2 {
		t.Fatal("duplicou revisões", total, err)
	}
	var conclusao string
	if err = db.QueryRow(`SELECT DATE_FORMAT(concluido_em,'%Y-%m-%d %H:%i:%s') FROM cronograma_itens WHERE id=?`, feito).Scan(&conclusao); err != nil || conclusao != "2026-09-14 08:00:00" {
		t.Fatal("alterou histórico", err)
	}
	exec(`UPDATE cronogramas SET configuracao='{"horas":{"seg":0.25}}' WHERE id=?`, agenda)
	if err = r.ReplanejarCalendario(ctx, aluno, concurso, "2026-09-14", "2026-09-13"); err == nil {
		t.Fatal("aceitou sobrecarga")
	}
	var data string
	if err = db.QueryRow(`SELECT DATE_FORMAT(data_planejada,'%Y-%m-%d') FROM cronograma_itens WHERE id=?`, b).Scan(&data); err != nil || data != "2026-09-15" {
		t.Fatal("não preservou plano após erro", err)
	}
}
