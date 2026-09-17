package cursos

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

// Opt-in: cria e remove somente um banco temporário exclusivo deste teste.
func TestIntegracaoAcessoCursos(t *testing.T) {
	if os.Getenv("CURSOS_TEST_MYSQL") != "1" {
		t.Skip("defina CURSOS_TEST_MYSQL=1 para validar com MySQL local")
	}
	t.Chdir("../..")
	cfg, err := configuracao.Carregar()
	if err != nil {
		t.Fatal(err)
	}
	cfg.DBNome = fmt.Sprintf("test_cursos_%d", time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	db, err := banco.Abrir(ctx, cfg)
	if err != nil {
		t.Fatal("MySQL de teste indisponível:", err)
	}
	defer db.Close()
	defer func() {
		if _, err := db.Exec("DROP DATABASE `" + cfg.DBNome + "`"); err != nil {
			t.Error("limpeza do banco temporário:", err)
		}
	}()
	// Apenas o esquema usado por cursos; a migração antiga 009 de outros
	// domínios tem um problema independente no executor de migrações.
	for _, nome := range []string{"001_estrutura_inicial.sql", "002_dominio_estudos.sql", "017_banco_questoes.sql", "019_cursos.sql", "020_cursos_pdfs.sql", "021_cursos_aprendizagem.sql", "022_cursos_capas.sql", "023_cursos_modo_exibicao.sql"} {
		conteudo, err := os.ReadFile("migracoes/sql/" + nome)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = db.ExecContext(ctx, string(conteudo)); err != nil {
			t.Fatal(nome, err)
		}
	}
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.ExecContext(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
	mentor, outro, a, b, intruso, concurso := identificador.UUID(), identificador.UUID(), identificador.UUID(), identificador.UUID(), identificador.UUID(), identificador.UUID()
	for i, id := range []string{mentor, outro, a, b, intruso} {
		papel := "aluno"
		if i < 2 {
			papel = "mentor"
		}
		exec(`INSERT INTO usuarios (id,nome,email,senha_hash,papel) VALUES (?,?,?,?,?)`, id, id, id+"@teste.local", "teste", papel)
	}
	for _, id := range []string{a, b} {
		exec(`INSERT INTO mentor_alunos (id,mentor_id,aluno_id) VALUES (?,?,?)`, identificador.UUID(), mentor, id)
	}
	exec(`INSERT INTO mentor_alunos (id,mentor_id,aluno_id) VALUES (?,?,?)`, identificador.UUID(), outro, intruso)
	exec(`INSERT INTO concursos (id,nome,banca,criado_por) VALUES (?,?,?,?)`, concurso, "Concurso", "Banca", mentor)
	exec(`INSERT INTO aluno_concursos (id,aluno_id,concurso_id,atribuido_por) VALUES (?,?,?,?)`, identificador.UUID(), a, concurso, mentor)
	s := Novo(db)
	publicar := func(id string) Curso {
		t.Helper()
		lista, err := s.Listar(ctx, mentor, true, id)
		if err != nil || len(lista) != 1 {
			t.Fatalf("rascunho: %v", err)
		}
		if err := s.Publicar(ctx, mentor, id, lista[0].Revisao, true); err != nil {
			t.Fatal(err)
		}
		lista, err = s.Listar(ctx, mentor, true, id)
		if err != nil || len(lista) != 1 {
			t.Fatalf("publicado: %v", err)
		}
		if !lista[0].Publicado || lista[0].AlteracoesPendentes {
			t.Fatal("revisão publicada inconsistente")
		}
		return lista[0]
	}
	criar := func(escopo string, ids ...string) string {
		t.Helper()
		id, err := s.Salvar(ctx, mentor, "", Curso{Titulo: escopo, Escopo: escopo, Destinatarios: ids, Aulas: []Aula{{Titulo: "Aula", VideoID: "dQw4w9WgXcQ"}}})
		if err != nil {
			t.Fatal(err)
		}
		publicar(id)
		return id
	}
	global, individual, porConcurso := criar("global"), criar("alunos", a), criar("concursos", concurso)
	verificar := func(usuario string, gestor bool, id string, n int) {
		t.Helper()
		lista, err := s.Listar(ctx, usuario, gestor, id)
		if err != nil || len(lista) != n {
			t.Fatalf("acesso: obtido %d, esperado %d; erro %v", len(lista), n, err)
		}
	}
	verificar(a, false, "", 3)
	verificar(b, false, "", 1)
	verificar(intruso, false, "", 0)
	verificar(outro, true, "", 0)
	verificar(intruso, false, global, 0)
	verificar(b, false, individual, 0)
	verificar(b, false, porConcurso, 0)
	for _, c := range []Curso{{Titulo: "x", Escopo: "alunos", Destinatarios: []string{intruso}, Aulas: []Aula{{Titulo: "a", VideoID: "dQw4w9WgXcQ"}}}, {Titulo: "x", Escopo: "concursos", Destinatarios: []string{identificador.UUID()}, Aulas: []Aula{{Titulo: "a", VideoID: "dQw4w9WgXcQ"}}}} {
		if _, err := s.Salvar(ctx, mentor, "", c); err != ErrEntrada {
			t.Fatalf("destinatário fora da mentoria: %v", err)
		}
	}
	lista, err := s.Listar(ctx, mentor, true, individual)
	if err != nil {
		t.Fatal(err)
	}
	editado := lista[0]
	editado.ModoExibicao = "modulos"
	capaModulo, err := s.SalvarCapa(ctx, mentor, imagemTeste(t))
	if err != nil {
		t.Fatal(err)
	}
	editado.Aulas[0].ModuloCapaURL = capaModulo
	capa, err := s.SalvarCapa(ctx, mentor, imagemTeste(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = s.ObterCapa(ctx, mentor, true, idCapa(capa)); err != nil {
		t.Fatal("prévia da capa", err)
	}
	if _, _, err = s.ObterCapa(ctx, a, false, idCapa(capa)); err != ErrNaoEncontrado {
		t.Fatal("capa não publicada exposta", err)
	}
	if _, _, err = s.ObterCapa(ctx, outro, true, idCapa(capa)); err != ErrNaoEncontrado {
		t.Fatal("capa de outro mentor exposta", err)
	}
	editado.CapaURL = capa
	pdfID, _, err := s.SalvarPDF(ctx, mentor, "aula.pdf", []byte("%PDF-1.4\n%%EOF"))
	if err != nil {
		t.Fatal(err)
	}
	editado.Aulas[0].PDFID = pdfID
	editado.Aulas[0].PDFNome = "nome-injetado.pdf"
	if _, err = s.Salvar(ctx, mentor, individual, editado); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.ObterPDF(ctx, a, false, individual, pdfID); err != ErrNaoEncontrado {
		t.Fatal("PDF do rascunho exposto", err)
	}
	antes, err := s.Listar(ctx, a, false, individual)
	if err != nil || antes[0].ModoExibicao != "curso" {
		t.Fatal("modo do rascunho exposto", err)
	}
	if _, _, err = s.ObterCapa(ctx, a, false, idCapa(capaModulo)); err != ErrNaoEncontrado {
		t.Fatal("capa de módulo do rascunho exposta", err)
	}
	editado = publicar(individual)
	depois, err := s.Listar(ctx, a, false, individual)
	if err != nil || depois[0].ModoExibicao != "modulos" {
		t.Fatal("modo não publicado", err)
	}
	if _, _, err = s.ObterCapa(ctx, a, false, idCapa(capaModulo)); err != nil {
		t.Fatal("capa do módulo inacessível", err)
	}
	if tipo, _, err := s.ObterCapa(ctx, a, false, idCapa(capa)); err != nil || tipo != "image/png" {
		t.Fatal("capa publicada", tipo, err)
	}
	if _, _, err = s.ObterCapa(ctx, b, false, idCapa(capa)); err != ErrNaoEncontrado {
		t.Fatal("capa fora do acesso", err)
	}
	capaOutro, err := s.SalvarCapa(ctx, outro, imagemTeste(t))
	if err != nil {
		t.Fatal(err)
	}
	editado.CapaURL = capaOutro
	editado.Aulas[0].ModuloCapaURL = capaOutro
	if _, err = s.Salvar(ctx, mentor, individual, editado); err != ErrEntrada {
		t.Fatal("vinculou capa alheia", err)
	}
	editado.CapaURL = capa
	if _, err = s.Salvar(ctx, mentor, individual, editado); err != ErrEntrada {
		t.Fatal("capa de módulo alheia aceita", err)
	}
	editado.Aulas[0].ModuloCapaURL = capaModulo
	if nome, _, err := s.ObterPDF(ctx, a, false, individual, pdfID); err != nil || nome != "aula.pdf" {
		t.Fatal("PDF autorizado", nome, err)
	}
	for _, usuario := range []string{b, intruso} {
		if _, _, err := s.ObterPDF(ctx, usuario, false, individual, pdfID); err != ErrNaoEncontrado {
			t.Fatal("PDF fora do acesso", err)
		}
	}
	if _, _, err := s.ObterPDF(ctx, mentor, true, global, pdfID); err != ErrNaoEncontrado {
		t.Fatal("PDF sem vínculo ao curso", err)
	}
	pdfOutro, _, err := s.SalvarPDF(ctx, outro, "privado.pdf", []byte("%PDF-1.4\n%%EOF"))
	if err != nil {
		t.Fatal(err)
	}
	editado.Aulas[0].PDFID = pdfOutro
	if _, err = s.Salvar(ctx, mentor, individual, editado); err != ErrEntrada {
		t.Fatal("anexo de outro mentor aceito", err)
	}
	editado.Aulas[0].PDFID = pdfID
	editado.Destinatarios = []string{b}
	if _, err = s.Salvar(ctx, outro, individual, editado); err != ErrNaoEncontrado {
		t.Fatal("outro mentor alterou curso", err)
	}
	if err = s.Desativar(ctx, outro, global); err != ErrNaoEncontrado {
		t.Fatal("outro mentor desativou curso", err)
	}
	if _, err = s.Salvar(ctx, mentor, individual, editado); err != nil {
		t.Fatal(err)
	}
	verificar(a, false, individual, 1)
	editado = publicar(individual)
	verificar(a, false, individual, 0)
	if _, _, err = s.ObterCapa(ctx, a, false, idCapa(capa)); err != ErrNaoEncontrado {
		t.Fatal("capa acessível após revogação", err)
	}
	if _, _, err := s.ObterPDF(ctx, a, false, individual, pdfID); err != ErrNaoEncontrado {
		t.Fatal("PDF acessível após revogar acesso", err)
	}
	verificar(b, false, individual, 1)
	editado.Aulas[0].PDFID = ""
	if _, err = s.Salvar(ctx, mentor, individual, editado); err != nil {
		t.Fatal(err)
	}
	editado = publicar(individual)
	if _, _, err := s.ObterPDF(ctx, b, false, individual, pdfID); err != ErrNaoEncontrado {
		t.Fatal("PDF removido ainda acessível", err)
	}
	// Publicação, progresso, bloqueios e exercícios usam o mesmo controle de acesso.
	qID := identificador.UUID()
	exec(`INSERT INTO banco_questoes (id,disciplina,assunto,enunciado,alternativas,resposta_correta,explicacao,criado_por) VALUES (?,'Direito','Teste','Pergunta',JSON_ARRAY('Sim','Não'),'A','Explicação',?)`, qID, mentor)
	id, err := s.Salvar(ctx, mentor, "", Curso{Titulo: "Aprendizagem", Escopo: "global"})
	if err != nil {
		t.Fatal(err)
	}
	verificar(a, false, id, 0)
	if err := s.Publicar(ctx, mentor, id, 1, true); err != ErrEntrada {
		t.Fatal("publicou curso vazio", err)
	}
	cursos, err := s.Listar(ctx, mentor, true, id)
	if err != nil {
		t.Fatal(err)
	}
	c := cursos[0]
	c.Aulas = []Aula{
		{Titulo: "Introdução", Modulo: "Módulo 1", VideoID: "dQw4w9WgXcQ", Questoes: []string{qID}},
		{Titulo: "Avançado", Modulo: "Módulo 2", VideoID: "dQw4w9WgXcQ", ExigeAnterior: true, PDFID: pdfID, Questoes: []string{qID}},
	}
	if _, err = s.Salvar(ctx, mentor, id, c); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Salvar(ctx, mentor, id, c); err != ErrConflito {
		t.Fatal("edição desatualizada aceita", err)
	}
	c = publicar(id)
	a1, a2 := c.Aulas[0].ID, c.Aulas[1].ID
	cursos, err = s.Listar(ctx, a, false, id)
	if err != nil {
		t.Fatal(err)
	}
	locked := cursos[0].Aulas[1]
	if !locked.Bloqueada || locked.VideoID != "" || locked.PDFID != "" || len(locked.Questoes) != 0 {
		t.Fatal("conteúdo bloqueado exposto")
	}
	if _, err = s.SalvarProgresso(ctx, a, id, a2, EntradaProgresso{}); err != ErrBloqueado {
		t.Fatal("progresso bloqueado", err)
	}
	if _, _, err = s.ObterPDF(ctx, a, false, id, pdfID); err != ErrNaoEncontrado {
		t.Fatal("PDF bloqueado", err)
	}
	if _, err = s.QuestoesAula(ctx, a, false, id, a2); err != ErrBloqueado {
		t.Fatal("questões bloqueadas", err)
	}
	if _, err = s.SalvarProgresso(ctx, intruso, id, a1, EntradaProgresso{}); err != ErrNaoEncontrado {
		t.Fatal("progresso de intruso", err)
	}
	qs, err := s.QuestoesAula(ctx, a, false, id, a1)
	if err != nil || len(qs) != 1 {
		t.Fatal("questões vinculadas", qs, err)
	}
	res, err := s.ResponderQuestao(ctx, a, false, id, a1, qID, "B")
	if err != nil || res.Correta || res.RespostaCorreta != "A" {
		t.Fatal("correção incorreta", res, err)
	}
	res, err = s.ResponderQuestao(ctx, a, false, id, a1, qID, "A")
	if err != nil || !res.Correta {
		t.Fatal("acerto incorreto", res, err)
	}
	if _, err = s.ResponderQuestao(ctx, a, false, id, a1, qID, "Z"); err != ErrEntrada {
		t.Fatal("alternativa inválida", err)
	}
	if _, err = s.ResponderQuestao(ctx, mentor, true, id, a1, qID, "A"); err != nil {
		t.Fatal("resposta na prévia", err)
	}
	var n int
	if err = db.QueryRow(`SELECT COUNT(*) FROM cursos_respostas WHERE curso_id=?`, id).Scan(&n); err != nil || n != 1 {
		t.Fatal("prévia gravou resposta", n, err)
	}
	sim := true
	progresso, err := s.SalvarProgresso(ctx, a, id, a1, EntradaProgresso{Posicao: 30, Duracao: 120, Concluida: &sim})
	if err != nil || progresso.Resumo.Percentual != 50 || progresso.Aulas[1].Bloqueada {
		t.Fatal("conclusão/liberação", progresso.Resumo, err)
	}
	progresso, err = s.SalvarProgresso(ctx, a, id, a1, EntradaProgresso{Posicao: 45, Duracao: 120})
	if err != nil || !progresso.Aulas[0].Progresso.Concluida {
		t.Fatal("heartbeat desmarcou conclusão", err)
	}
	previa, err := s.Previa(ctx, mentor, id)
	if err != nil || previa.Resumo.Concluidas != 0 || !previa.Aulas[1].Bloqueada {
		t.Fatal("prévia incorreta", err)
	}
	relatorio, err := s.Acompanhar(ctx, mentor, id)
	if err != nil || len(relatorio) != 2 {
		t.Fatal("acompanhamento", err, len(relatorio))
	}
	for _, al := range relatorio {
		if al.ID == a && (al.Resumo.Percentual != 50 || al.Acertos != 1 || al.Respostas != 1) {
			t.Fatal("métricas incorretas", al)
		}
	}
	c.Aulas[0].Titulo = "Renomeada"
	if _, err = s.Salvar(ctx, mentor, id, c); err != nil {
		t.Fatal(err)
	}
	cursos, err = s.Listar(ctx, a, false, id)
	if err != nil || cursos[0].Aulas[0].Titulo == "Renomeada" {
		t.Fatal("rascunho vazou", err)
	}
	c = publicar(id)
	cursos, err = s.Listar(ctx, a, false, id)
	if err != nil || cursos[0].Resumo.Percentual != 50 || cursos[0].Aulas[0].Progresso.Posicao != 45 {
		t.Fatal("retomada perdida", err)
	}
	c.Aulas[0].VideoID = "abcdefghijk"
	if _, err = s.Salvar(ctx, mentor, id, c); err != nil {
		t.Fatal(err)
	}
	c = publicar(id)
	cursos, err = s.Listar(ctx, a, false, id)
	if err != nil || cursos[0].Resumo.Concluidas != 0 || !cursos[0].Aulas[1].Bloqueada {
		t.Fatal("vídeo novo herdou conclusão", err)
	}
	if err = s.Publicar(ctx, mentor, id, c.Revisao, false); err != nil {
		t.Fatal(err)
	}
	verificar(a, false, id, 0)
	publicar(id)
	verificar(a, false, id, 1)
	exec(`UPDATE mentor_alunos SET data_expiracao_plano=DATE_SUB(CURDATE(),INTERVAL 1 DAY) WHERE aluno_id=?`, a)
	verificar(a, false, id, 0)
	exec(`UPDATE mentor_alunos SET data_expiracao_plano=NULL WHERE aluno_id=?`, a)
	exec(`UPDATE aluno_concursos SET ativo=FALSE WHERE aluno_id=?`, a)
	verificar(a, false, porConcurso, 0)
	exec(`UPDATE mentor_alunos SET ativo=FALSE WHERE aluno_id=?`, b)
	verificar(b, false, "", 0)
	if err = s.Desativar(ctx, mentor, global); err != nil {
		t.Fatal(err)
	}
	verificar(a, false, global, 0)
}
