package repository

import (
	"client-server-api/pkg/models"
	"context"
)

type CotacaoRepository interface {
	Save(ctx context.Context, cotacao *models.Cotacao) error
	FindByID(ctx context.Context, id int64) (*models.Cotacao, error)
}


