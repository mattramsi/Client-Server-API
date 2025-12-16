package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"client-server-api/pkg/errors"
	"client-server-api/pkg/models"
)

type CotacaoClient struct {
	baseURL string
	client  *http.Client
	timeout time.Duration
}

func NewCotacaoClient(serverURL string, timeout time.Duration) *CotacaoClient {
	return &CotacaoClient{
		baseURL: serverURL,
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

func (c *CotacaoClient) GetBid(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return "", errors.ErroInterno(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if err == context.DeadlineExceeded {
			return "", errors.ErroTimeoutContext("chamada ao servidor", err)
		}
		return "", errors.ErroInterno(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.ErroInterno(fmt.Errorf("status %d: %s", resp.StatusCode, resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.ErroInterno(err)
	}

	var bidResponse models.BidResponse
	if err := json.Unmarshal(body, &bidResponse); err != nil {
		return "", errors.ErroInterno(fmt.Errorf("erro ao fazer parse do JSON: %w", err))
	}

	if bidResponse.Bid == "" {
		return "", errors.ErroValidacao("bid não pode estar vazio")
	}

	return bidResponse.Bid, nil
}



