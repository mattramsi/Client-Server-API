package external

import (
	"context"
	"client-server-api/pkg/models"
)

type ExchangeRateClient interface {
	FetchUSD(ctx context.Context) (*models.Cotacao, error)
}

