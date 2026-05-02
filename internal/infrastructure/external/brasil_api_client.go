package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

var ErrCNPJNotFound = errors.New("CNPJ not found")

const brasilAPIBaseURL = "https://brasilapi.com.br/api"

type brasilAPICNPJResponse struct {
	CNPJ                    string `json:"cnpj"`
	RazaoSocial             string `json:"razao_social"`
	NomeFantasia            string `json:"nome_fantasia"`
	SituacaoCadastral       int    `json:"situacao_cadastral"`
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrCNPJNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brasil api returned status %d", resp.StatusCode)
	}

	var data brasilAPICNPJResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
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
