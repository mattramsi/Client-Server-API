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
	// 1. Carregar configuração
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Erro ao carregar configuração:", err)
	}

	// 2. Inicializar dependências
	// API Client (não tem dependências)
	apiClient := external.NewAwesomeAPIClient(cfg.API)

	// Repository (usa config)
	repo, err := repository.NewSQLiteRepository(cfg.Database)
	if err != nil {
		log.Fatal("Erro ao criar repositório:", err)
	}
	defer repo.Close() // Fechar pool ao encerrar

	// 3. Criar service (usa API client e repository)
	cotacaoService := service.NewCotacaoService(apiClient, repo)

	// 4. Criar handler (usa service)
	cotacaoHandler := handler.NewCotacaoHandler(cotacaoService)

	// 5. Configurar rotas
	http.HandleFunc("/cotacao", cotacaoHandler.GetCotacao)

	// 6. Criar servidor HTTP
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: nil,
	}

	// 7. Iniciar servidor em goroutine para permitir graceful shutdown
	go func() {
		log.Printf("Servidor iniciado na porta %s\n", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Erro ao iniciar servidor:", err)
		}
	}()

	// 8. Graceful shutdown - aguardar sinal de encerramento
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Encerrando servidor...")

	// Criar contexto com timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Encerrar servidor graciosamente
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Erro ao encerrar servidor:", err)
	}

	log.Println("Servidor encerrado com sucesso")
}
