package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"client-server-api/internal/client"
)

func main() {
	// Configuração (pode vir de env vars no futuro)
	serverURL := "http://localhost:8080/cotacao"
	timeout := 300 * time.Millisecond
	filename := "cotacao.txt"

	// Criar cliente HTTP
	cotacaoClient := client.NewCotacaoClient(serverURL, timeout)

	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Buscar cotação do servidor
	bid, err := cotacaoClient.GetBid(ctx)
	if err != nil {
		log.Fatal("Erro ao obter cotação:", err)
	}

	// Escrever em arquivo
	if err := client.WriteCotacaoToFile(filename, bid); err != nil {
		log.Fatal("Erro ao escrever arquivo:", err)
	}

	fmt.Printf("Cotação salva com sucesso no arquivo %s: %s\n", filename, bid)
}


