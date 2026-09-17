package configuracao

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Configuracao struct {
	Ambiente       string
	Porta          string
	FusoHorario    *time.Location
	DBHost         string
	DBPorta        int
	DBUsuario      string
	DBSenha        string
	DBNome         string
	DBCriarBanco   bool
	MestreEmail    string
	MestreSenha    string
	OrigemFrontend string
	UploadsPath    string

	// Rate limiting
	RateGlobalRPS   float64 // requisições por segundo (global, por IP)
	RateGlobalBurst int     // burst global
	RateLoginRPS    float64 // requisições por segundo (rota de login, por IP)
	RateLoginBurst  int     // burst de login
}

func texto(chave, padrao string) string {
	if valor := os.Getenv(chave); valor != "" {
		return valor
	}
	return padrao
}

func Carregar() (Configuracao, error) {
	carregarArquivoEnv(".env")
	portaBanco, err := strconv.Atoi(texto("DB_PORT", "3306"))
	if err != nil {
		return Configuracao{}, fmt.Errorf("DB_PORT inválida: %w", err)
	}
	fuso, err := time.LoadLocation(texto("FUSO_HORARIO", "America/Fortaleza"))
	if err != nil {
		return Configuracao{}, fmt.Errorf("FUSO_HORARIO inválido: %w", err)
	}
	rateGlobalRPS, err := strconv.ParseFloat(texto("RATE_GLOBAL_RPS", "100.0"), 64)
	if err != nil {
		return Configuracao{}, fmt.Errorf("RATE_GLOBAL_RPS inválido: %w", err)
	}
	rateGlobalBurst, err := strconv.Atoi(texto("RATE_GLOBAL_BURST", "500"))
	if err != nil {
		return Configuracao{}, fmt.Errorf("RATE_GLOBAL_BURST inválido: %w", err)
	}
	rateLoginRPS, err := strconv.ParseFloat(texto("RATE_LOGIN_RPS", "10.0"), 64)
	if err != nil {
		return Configuracao{}, fmt.Errorf("RATE_LOGIN_RPS inválido: %w", err)
	}
	rateLoginBurst, err := strconv.Atoi(texto("RATE_LOGIN_BURST", "50"))
	if err != nil {
		return Configuracao{}, fmt.Errorf("RATE_LOGIN_BURST inválido: %w", err)
	}
	dbCriarBanco, err := strconv.ParseBool(texto("DB_CRIAR_BANCO", "true"))
	if err != nil {
		return Configuracao{}, fmt.Errorf("DB_CRIAR_BANCO inválida: %w", err)
	}
	cfg := Configuracao{
		Ambiente: texto("AMBIENTE", "desenvolvimento"), Porta: texto("PORTA", "8080"), FusoHorario: fuso,
		DBHost: texto("DB_HOST", "127.0.0.1"), DBPorta: portaBanco, DBUsuario: texto("DB_USUARIO", "root"),
		DBSenha: os.Getenv("DB_SENHA"), DBNome: texto("DB_NOME", "track_concursos"),
		DBCriarBanco: dbCriarBanco,
		MestreEmail:  texto("MESTRE_EMAIL", ""), MestreSenha: texto("MESTRE_SENHA", ""),
		OrigemFrontend:  texto("ORIGEM_FRONTEND", "http://127.0.0.1:5173"),
		UploadsPath:     texto("UPLOADS_PATH", "./uploads"),
		RateGlobalRPS:   rateGlobalRPS,
		RateGlobalBurst: rateGlobalBurst,
		RateLoginRPS:    rateLoginRPS,
		RateLoginBurst:  rateLoginBurst,
	}
	return cfg, nil
}

func carregarArquivoEnv(caminho string) {
	conteudo, err := os.ReadFile(caminho)
	if err != nil {
		return
	}
	for _, linha := range strings.Split(string(conteudo), "\n") {
		linha = strings.TrimSpace(linha)
		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}
		partes := strings.SplitN(linha, "=", 2)
		if len(partes) != 2 {
			continue
		}
		chave, valor := strings.TrimSpace(partes[0]), strings.Trim(strings.TrimSpace(partes[1]), `"`)
		if chave != "" && os.Getenv(chave) == "" {
			_ = os.Setenv(chave, valor)
		}
	}
}
