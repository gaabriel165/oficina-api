package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

var ErrCNPJNotFound = errors.New("CNPJ not found")

const brasilAPIBaseURL = "https://brasilapi.com.br/api"

type brasilAPICNPJResponse struct {
	CNPJ              string `json:"cnpj"`
	RazaoSocial       string `json:"razao_social"`
	NomeFantasia      string `json:"nome_fantasia"`
	SituacaoCadastral int    `json:"situacao_cadastral"`
}

type BrasilAPIClient struct {
	httpClient *http.Client
}

func NewBrasilAPIClient() *BrasilAPIClient {
	return &BrasilAPIClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *BrasilAPIClient) Fetch(cnpj string) (*repository.CNPJInfo, error) {
	url := fmt.Sprintf("%s/cnpj/v1/%s", brasilAPIBaseURL, cnpj)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logIntegrationError(cnpj, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrCNPJNotFound
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("brasil api returned status %d", resp.StatusCode)
		logIntegrationError(cnpj, err)
		return nil, err
	}

	var data brasilAPICNPJResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logIntegrationError(cnpj, err)
		return nil, err
	}

	return &repository.CNPJInfo{
		CNPJ:              data.CNPJ,
		CompanyName:       data.RazaoSocial,
		TradeName:         data.NomeFantasia,
		SituacaoCadastral: data.SituacaoCadastral,
		IsActive:          data.SituacaoCadastral == 2,
	}, nil
}

func logIntegrationError(cnpj string, err error) {
	slog.Error("integration.error",
		slog.String("event", "integration.error"),
		slog.String("integration", "brasilapi"),
		slog.String("operation", "fetch_cnpj"),
		slog.String("cnpj", cnpj),
		slog.String("error", err.Error()),
	)
}
