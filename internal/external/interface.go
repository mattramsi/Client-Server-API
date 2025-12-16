package external

import (
	"context"
	"client-server-api/pkg/models"
)

// ExchangeRateClient define o contrato para clientes de API de câmbio
type ExchangeRateClient interface {
	FetchUSD(ctx context.Context) (*models.Cotacao, error)
}

