package main

import (
	"client-server-api/pkg/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {

	dolarApiAddress := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	_ = dolarApiAddress

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", dolarApiAddress, nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	request = request.WithContext(ctx)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Println("Erro ao chamar API:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		return
	}

	fmt.Println(string(body))

	var dolarResponse models.DolarResponse
	err = json.Unmarshal(body, &dolarResponse)
	if err != nil {
		fmt.Println("Erro ao deserializar resposta:", err)
		return
	}

	fmt.Println(dolarResponse.USDBRL.Bid)

	select {
	case <-ctx.Done():
		fmt.Println("Requisição cancelada ou expirou:", ctx.Err())
		return
	default:
		fmt.Println("Resposta da API recebida com status:", response.Status)
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	if err != nil {
		fmt.Println("Erro ao criar contexto para gravar no SQLite:", err)
		return
	}
	defer cancelDB()

	db, err := sql.Open("sqlite3", ".cotacoes.db")
	if err != nil {
		fmt.Println("Erro ao abrir banco de dados:", err)
		return
	}
	defer db.Close()

	createTableSQL := `
			CREATE TABLE IF NOT EXISTS cotacoes (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				code TEXT,
				codein TEXT,
				name TEXT,
				high TEXT,
				low TEXT,
				var_bid TEXT,
				pct_change TEXT,
				bid TEXT,
				ask TEXT,
				timestamp TEXT,
				create_date TEXT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Erro ao criar tabela:", err)
		return
	}

	insertCotacaoSQL := `
			INSERT INTO cotacoes (code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = db.ExecContext(ctxDB, insertCotacaoSQL, dolarResponse.USDBRL.Code, dolarResponse.USDBRL.Codein, dolarResponse.USDBRL.Name, dolarResponse.USDBRL.High, dolarResponse.USDBRL.Low, dolarResponse.USDBRL.VarBid, dolarResponse.USDBRL.PctChange, dolarResponse.USDBRL.Bid, dolarResponse.USDBRL.Ask, dolarResponse.USDBRL.Timestamp, dolarResponse.USDBRL.CreateDate)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Timeout ao salvar no banco (10ms expirado)")
		} else {
			fmt.Println("Erro ao salvar no banco:", err)
		}
		return
	}

	select {
	case <-ctxDB.Done():
		fmt.Println("Timeout ao salvar no banco (10ms expirado)")
		return
	default:
		fmt.Println("Cotação gravada com sucesso")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.BidResponse{Bid: dolarResponse.USDBRL.Bid})
}
