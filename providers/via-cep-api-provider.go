// providers/via_cep_provider.go
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rafael1abrao/goexpert/multithreading/domain"
)

type viaCepRawResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

type ViaCepProvider struct {
	client *http.Client
}

func NewViaCepProvider(client *http.Client) *ViaCepProvider {
	return &ViaCepProvider{client: client}
}

func (v *ViaCepProvider) FetchCep(ctx context.Context, cep string) (*domain.CepResponse, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP retornou status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw viaCepRawResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	cepResponse := &domain.CepResponse{
		Cep:        raw.Cep,
		Logradouro: raw.Logradouro,
		Bairro:     raw.Bairro,
		Localidade: raw.Localidade,
		UF:         raw.UF,
		Source:     "ViaCEP",
	}

	return cepResponse, nil
}
