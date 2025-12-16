package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"client-server-api/pkg/models"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Erro ao obter cotação:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Erro: servidor retornou status %d\n", response.StatusCode)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		return
	}

	var bidResponse models.BidResponse
	err = json.Unmarshal(body, &bidResponse)
	if err != nil {
		fmt.Println("Erro ao deserializar resposta:", err)
		fmt.Println("Body recebido:", string(body))
		return
	}

	select {
	case <-ctx.Done():
		fmt.Println("Timeout ao obter cotação (300ms expirado):", ctx.Err())
		return
	default:
		fmt.Printf("Bid recebido: %s\n", bidResponse.Bid)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Println("Erro ao criar arquivo:", err)
		return
	}
	defer file.Close()

	conteudo := fmt.Sprintf("Dólar: %s", bidResponse.Bid)

	_, err = file.WriteString(conteudo)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}

	fmt.Println("Cotacao salva com sucesso no arquivo cotacao.txt")
}
