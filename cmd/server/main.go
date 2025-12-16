package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"client-server-api/internal/external"
	"client-server-api/internal/repository"
	"client-server-api/internal/server/config"
	"client-server-api/internal/server/handler"
	"client-server-api/internal/server/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Erro ao carregar configuração:", err)
	}

	apiClient := external.NewAwesomeAPIClient(cfg.API)

	repo, err := repository.NewSQLiteRepository(cfg.Database)
	if err != nil {
		log.Fatal("Erro ao criar repositório:", err)
	}
	defer repo.Close()

	cotacaoService := service.NewCotacaoService(apiClient, repo)

	cotacaoHandler := handler.NewCotacaoHandler(cotacaoService)

	http.HandleFunc("/cotacao", cotacaoHandler.GetCotacao)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: nil,
	}

	go func() {
		log.Printf("Servidor iniciado na porta %s\n", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Erro ao iniciar servidor:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Encerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Erro ao encerrar servidor:", err)
	}

	log.Println("Servidor encerrado com sucesso")
}
