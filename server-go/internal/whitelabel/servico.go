package whitelabel

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"strings"

	"chega-junto-concurseiro-web/internal/arquivos"
)

type Config struct {
	MentorID           string  `json:"mentorId"`
	NomePlataforma     string  `json:"nomePlataforma"`
	LogoURL            *string `json:"logoUrl"`
	BannerURL          *string `json:"bannerUrl"`
	CorPrimaria        string  `json:"corPrimaria"`
	CorSecundaria      string  `json:"corSecundaria"`
	MensagemBoasVindas *string `json:"mensagemBoasVindas"`
	AtualizadoEm       string  `json:"atualizadoEm,omitempty"`
}
type Servico struct{ repositorio repositorio }

func NovoServico(b *sql.DB) *Servico            { return NovoComRepositorio(novoRepositorioMySQL(b)) }
func NovoComRepositorio(r repositorio) *Servico { return &Servico{repositorio: r} }
func padrao(id string) Config {
	m := "Análise completa da sua preparação"
	return Config{MentorID: id, NomePlataforma: "Chega Junto Concurseiro", CorPrimaria: "#4f8ef7", CorSecundaria: "#7c5cfc", MensagemBoasVindas: &m}
}
func (s *Servico) ObterPorMentor(ctx context.Context, id string) (Config, error) {
	c, encontrado, err := s.repositorio.ObterPorMentor(ctx, id)
	if err != nil {
		return Config{}, err
	}
	if !encontrado {
		return padrao(id), nil
	}
	return c, nil
}
func (s *Servico) ObterPorAluno(ctx context.Context, aluno string) (Config, error) {
	mentor, encontrado, err := s.repositorio.MentorDoAluno(ctx, aluno)
	if err != nil {
		return Config{}, err
	}
	if !encontrado {
		return padrao(""), nil
	}
	return s.ObterPorMentor(ctx, mentor)
}
func (s *Servico) Salvar(ctx context.Context, c Config) (Config, error) {
	atual, err := s.ObterPorMentor(ctx, c.MentorID)
	if err != nil {
		return Config{}, err
	}
	c.NomePlataforma = valor(c.NomePlataforma, "Chega Junto Concurseiro")
	c.CorPrimaria = valor(c.CorPrimaria, "#4f8ef7")
	c.CorSecundaria = valor(c.CorSecundaria, "#7c5cfc")
	c.LogoURL = limpar(c.LogoURL)
	c.BannerURL = limpar(c.BannerURL)
	if c.LogoURL != nil {
		c.LogoURL = atual.LogoURL
	}
	if c.BannerURL != nil {
		c.BannerURL = atual.BannerURL
	}
	c.MensagemBoasVindas = limpar(c.MensagemBoasVindas)
	if err := s.repositorio.Salvar(ctx, c); err != nil {
		return Config{}, err
	}
	return s.ObterPorMentor(ctx, c.MentorID)
}

func URLsProtegidas(c Config, base string) Config {
	if c.LogoURL != nil {
		url := base + "/logo"
		c.LogoURL = &url
	}
	if c.BannerURL != nil {
		url := base + "/banner"
		c.BannerURL = &url
	}
	return c
}

func (s *Servico) SalvarLogo(ctx context.Context, mentorID, nomeArquivo string, conteudo io.Reader, raiz string) (string, error) {
	return s.salvarImagem(ctx, mentorID, conteudo, raiz, "logos", true)
}

func (s *Servico) SalvarBanner(ctx context.Context, mentorID, nomeArquivo string, conteudo io.Reader, raiz string) (string, error) {
	return s.salvarImagem(ctx, mentorID, conteudo, raiz, "banners", false)
}

func (s *Servico) salvarImagem(ctx context.Context, mentorID string, conteudo io.Reader, raiz, pasta string, logo bool) (string, error) {
	imagem, err := arquivos.LerImagem(conteudo)
	if err != nil {
		return "", err
	}
	caminho, err := arquivos.SalvarImagem(raiz, pasta, imagem)
	if err != nil {
		return "", err
	}
	url := "/uploads/" + pasta + "/" + filepath.Base(caminho)
	config, err := s.ObterPorMentor(ctx, mentorID)
	if err != nil {
		_ = os.Remove(caminho)
		return "", err
	}
	antiga := config.BannerURL
	if logo {
		antiga = config.LogoURL
	}
	if logo {
		config.LogoURL = &url
	} else {
		config.BannerURL = &url
	}
	if err = s.repositorio.Salvar(ctx, config); err != nil {
		_ = os.Remove(caminho)
		return "", err
	}
	if antiga != nil && strings.HasPrefix(*antiga, "/uploads/"+pasta+"/") {
		_ = os.Remove(filepath.Join(raiz, pasta, filepath.Base(*antiga)))
	}
	return url, nil
}

func (s *Servico) ObterImagem(ctx context.Context, mentorID, tipo, raiz string) (string, []byte, error) {
	config, err := s.ObterPorMentor(ctx, mentorID)
	if err != nil {
		return "", nil, err
	}
	url := config.LogoURL
	pasta := "logos"
	if tipo == "banner" {
		url, pasta = config.BannerURL, "banners"
	}
	if url == nil || !strings.HasPrefix(*url, "/uploads/"+pasta+"/") {
		return "", nil, os.ErrNotExist
	}
	nome := filepath.Base(*url)
	if nome != filepath.Base(filepath.Clean(nome)) {
		return "", nil, os.ErrNotExist
	}
	conteudo, err := os.ReadFile(filepath.Join(raiz, pasta, nome))
	if err != nil {
		return "", nil, err
	}
	imagem, err := arquivos.LerImagem(bytes.NewReader(conteudo))
	if err != nil {
		return "", nil, err
	}
	return imagem.Mime, conteudo, nil
}
func valor(v, p string) string {
	if x := strings.TrimSpace(v); x != "" {
		return x
	}
	return p
}
func limpar(v *string) *string {
	if v == nil {
		return nil
	}
	x := strings.TrimSpace(*v)
	if x == "" {
		return nil
	}
	return &x
}
