package client

import (
	"fmt"
	"os"
)

func WriteCotacaoToFile(filename string, bid string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	content := fmt.Sprintf("Dólar: %s\n", bid)
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo: %w", err)
	}

	return nil
}



