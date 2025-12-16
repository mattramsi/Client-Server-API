package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"client-server-api/internal/server/config"
	"client-server-api/pkg/errors"
	"client-server-api/pkg/models"
)

type AwesomeAPIClient struct {
	baseURL string
	client  *http.Client
	timeout time.Duration
}

func NewAwesomeAPIClient(cfg config.APIConfig) *AwesomeAPIClient {
	return &AwesomeAPIClient{
		baseURL: cfg.BaseURL,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		timeout: cfg.Timeout,
	}
}

func (c *AwesomeAPIClient) FetchUSD(ctx context.Context) (*models.Cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return nil, errors.ErroAPI(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, errors.ErroTimeoutContext("chamada à API", err)
		}
		return nil, errors.ErroAPI(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.ErroAPI(fmt.Errorf("status %d: %s", resp.StatusCode, resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErroAPI(err)
	}

	var dolarResponse models.DolarResponse
	if err := json.Unmarshal(body, &dolarResponse); err != nil {
		return nil, errors.ErroAPI(fmt.Errorf("erro ao fazer parse do JSON: %w", err))
	}

	cotacao := &models.Cotacao{
		Code:       dolarResponse.USDBRL.Code,
		Codein:     dolarResponse.USDBRL.Codein,
		Name:       dolarResponse.USDBRL.Name,
		High:       dolarResponse.USDBRL.High,
		Low:        dolarResponse.USDBRL.Low,
		VarBid:     dolarResponse.USDBRL.VarBid,
		PctChange:  dolarResponse.USDBRL.PctChange,
		Bid:        dolarResponse.USDBRL.Bid,
		Ask:        dolarResponse.USDBRL.Ask,
		Timestamp:  dolarResponse.USDBRL.Timestamp,
		CreateDate: dolarResponse.USDBRL.CreateDate,
		CreatedAt:  time.Now(),
	}

	if cotacao.Bid == "" {
		return nil, errors.ErroValidacao("bid não pode estar vazio")
	}

	return cotacao, nil
}

