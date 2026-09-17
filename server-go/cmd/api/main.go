package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"track-concursos-web/internal/banco"
	"track-concursos-web/internal/configuracao"
	"track-concursos-web/internal/transporte"
	"track-concursos-web/migracoes"
)

func main() {
	cfg, err := configuracao.Carregar()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelar := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelar()
	db, err := banco.Abrir(ctx, cfg)
	if err != nil {
		log.Fatal("falha ao conectar ao MySQL: ", err)
	}
	defer db.Close()
	if err := migracoes.Aplicar(ctx, db); err != nil {
		log.Fatal("falha ao aplicar migrações: ", err)
	}
	api := transporte.Novo(db, cfg, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err := api.SemearMestre(); err != nil {
		log.Fatal("falha ao criar mestre inicial: ", err)
	}
	servidor := &http.Server{Addr: ":" + cfg.Porta, Handler: api.Rotas(), ReadHeaderTimeout: 10 * time.Second, ReadTimeout: 20 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	erros := make(chan error, 1)
	go func() { erros <- servidor.ListenAndServe() }()
	log.Printf("API Go em http://127.0.0.1:%s", cfg.Porta)
	sinais := make(chan os.Signal, 1)
	signal.Notify(sinais, os.Interrupt, syscall.SIGTERM)
	select {
	case sinal := <-sinais:
		log.Printf("encerrando API após sinal %s", sinal)
		encerrar, cancelarEncerramento := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelarEncerramento()
		if err := servidor.Shutdown(encerrar); err != nil {
			log.Printf("falha no encerramento gracioso: %v", err)
		}
	case err := <-erros:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
