package cursos

import "testing"

func TestVideoID(t *testing.T) {
	for _, v := range []string{"dQw4w9WgXcQ", "https://youtu.be/dQw4w9WgXcQ?t=20", "https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=abc", "https://youtube.com/shorts/dQw4w9WgXcQ", "https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ", "https://youtube.com/live/dQw4w9WgXcQ"} {
		if id, ok := videoID(v); !ok || id != "dQw4w9WgXcQ" {
			t.Errorf("link válido rejeitado: %s", v)
		}
	}
	for _, v := range []string{"https://youtube.com.evil.test/watch?v=dQw4w9WgXcQ", "javascript:alert(1)", "https://evil.test/dQw4w9WgXcQ", "https://youtube.com/watch?v=curto", "https://youtube.com@evil.test/watch?v=dQw4w9WgXcQ", "https://youtube.com/playlist?list=abc"} {
		if _, ok := videoID(v); ok {
			t.Errorf("link inválido aceito: %s", v)
		}
	}
}
func TestValidarCurso(t *testing.T) {
	base := func() Curso {
		return Curso{Titulo: " Curso ", Escopo: "global", Aulas: []Aula{{Titulo: " Aula ", VideoID: "https://youtu.be/dQw4w9WgXcQ"}}}
	}
	c := base()
	if err := validar(&c); err != nil {
		t.Fatal(err)
	}
	if c.Titulo != "Curso" || c.Aulas[0].VideoID != "dQw4w9WgXcQ" || c.Categoria != "Geral" {
		t.Fatal("normalização incorreta")
	}
	for _, alterar := range []func(*Curso){func(c *Curso) { c.Titulo = " " }, func(c *Curso) { c.Escopo = "publico" }, func(c *Curso) { c.CapaURL = "javascript:alert(1)" }, func(c *Curso) { c.Aulas[0].VideoID = "invalido" }} {
		c := base()
		alterar(&c)
		if validar(&c) == nil {
			t.Fatal("curso inválido aceito")
		}
	}
	for _, alterar := range []func(*Curso){func(c *Curso) { c.Escopo = "alunos" }, func(c *Curso) { c.Escopo = "concursos" }, func(c *Curso) { c.Aulas = nil }} {
		c := base()
		alterar(&c)
		if err := validar(&c); err != nil {
			t.Fatal("rascunho rejeitado", err)
		}
		if validarPublicacao(c) != ErrEntrada {
			t.Fatal("publicação incompleta aceita")
		}
	}
}
