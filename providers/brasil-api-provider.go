// providers/brasil_api_provider.go
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rafael1abrao/goexpert/multithreading/domain"
)

type brasilApiRawResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type BrasilAPIProvider struct {
	client *http.Client
}

func NewBrasilAPIProvider(client *http.Client) *BrasilAPIProvider {
	return &BrasilAPIProvider{client: client}
}

func (b *BrasilAPIProvider) FetchCep(ctx context.Context, cep string) (*domain.CepResponse, error) {

	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BrasilAPI retornou status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw brasilApiRawResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	cepResponse := &domain.CepResponse{
		Cep:        raw.Cep,
		Logradouro: raw.Street,
		Bairro:     raw.Neighborhood,
		Localidade: raw.City,
		UF:         raw.State,
		Source:     "BrasilAPI",
	}

	return cepResponse, nil
}
