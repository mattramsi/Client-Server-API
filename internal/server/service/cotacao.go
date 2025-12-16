package service

import (
	"context"

	"client-server-api/internal/external"
	"client-server-api/internal/repository"
)

type CotacaoService struct {
	apiClient  external.ExchangeRateClient
	repository repository.CotacaoRepository
}

func NewCotacaoService(
	apiClient external.ExchangeRateClient,
	repo repository.CotacaoRepository,
) *CotacaoService {
	return &CotacaoService{
		apiClient:  apiClient,
		repository: repo,
	}
}

func (s *CotacaoService) GetBid(ctx context.Context) (string, error) {
	cotacao, err := s.apiClient.FetchUSD(ctx)
	if err != nil {
		return "", err
	}

	if err := s.repository.Save(ctx, cotacao); err != nil {
		return "", err
	}

	return cotacao.Bid, nil
}


