package service

import (
	"context"

	"client-server-api/internal/external"
	"client-server-api/internal/repository"
)

// CotacaoService orquestra a lógica de negócio para cotações
type CotacaoService struct {
	apiClient  external.ExchangeRateClient
	repository repository.CotacaoRepository
}

// NewCotacaoService cria uma nova instância do service de cotação
func NewCotacaoService(
	apiClient external.ExchangeRateClient,
	repo repository.CotacaoRepository,
) *CotacaoService {
	return &CotacaoService{
		apiClient:  apiClient,
		repository: repo,
	}
}

// GetBid busca a cotação do dólar, salva no banco e retorna o bid
func (s *CotacaoService) GetBid(ctx context.Context) (string, error) {
	// 1. Buscar cotação da API externa
	cotacao, err := s.apiClient.FetchUSD(ctx)
	if err != nil {
		return "", err // Erro já é AppError (timeout, API error, etc.)
	}

	// 2. Salvar no banco de dados
	if err := s.repository.Save(ctx, cotacao); err != nil {
		return "", err // Erro já é AppError (timeout, database error, etc.)
	}

	// 3. Retornar bid
	return cotacao.Bid, nil
}

