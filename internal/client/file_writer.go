package client

import (
	"fmt"
	"os"
)

// WriteCotacaoToFile escreve a cotação em um arquivo
func WriteCotacaoToFile(filename string, bid string) error {
	// Criar/abrir arquivo
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	// Escrever conteúdo (adiciona \n no final)
	content := fmt.Sprintf("Dólar: %s\n", bid)
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo: %w", err)
	}

	return nil
}


