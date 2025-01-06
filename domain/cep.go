// domain/cep.go
package domain

import "context"

type CepResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Source     string
}

type CepProvider interface {
	FetchCep(ctx context.Context, cep string) (*CepResponse, error)
}
