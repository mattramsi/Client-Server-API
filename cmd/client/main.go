package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"client-server-api/internal/client"
)

func main() {
	serverURL := "http://localhost:8080/cotacao"
	timeout := 300 * time.Millisecond
	filename := "cotacao.txt"

	cotacaoClient := client.NewCotacaoClient(serverURL, timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	bid, err := cotacaoClient.GetBid(ctx)
	if err != nil {
		log.Fatal("Erro ao obter cotação:", err)
	}

	if err := client.WriteCotacaoToFile(filename, bid); err != nil {
		log.Fatal("Erro ao escrever arquivo:", err)
	}

	fmt.Printf("Cotação salva com sucesso no arquivo %s: %s\n", filename, bid)
}



