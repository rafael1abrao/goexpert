// domain/cep_service.go
package domain

import (
	"context"
	"errors"
	"regexp"
	"time"
)

func NormalizeAndValidateCep(cep string) (string, error) {
	re := regexp.MustCompile(`\D`)
	normalized := re.ReplaceAllString(cep, "")

	if len(normalized) != 8 {
		return "", errors.New("CEP inválido: deve conter 8 dígitos")
	}
	return normalized, nil
}

func FetchFasterCep(ctx context.Context, rawCep string, providers []CepProvider) (*CepResponse, error) {
	cep, err := NormalizeAndValidateCep(rawCep)
	if err != nil {
		return nil, err
	}

	// 2. Ajustamos o contexto principal com timeout de 1 segundo
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resultChan := make(chan *CepResponse, len(providers))
	errChan := make(chan error, len(providers))

	for _, provider := range providers {
		go func(p CepProvider) {
			resp, pErr := p.FetchCep(ctx, cep)
			if pErr != nil {
				errChan <- pErr
				return
			}
			resultChan <- resp
		}(provider)
	}
	for i := 0; i < len(providers); i++ {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout excedido ao buscar CEP")
		case err := <-errChan:
			_ = err
		case result := <-resultChan:
			return result, nil
		}
	}

	return nil, errors.New("nenhuma das requisições retornou com sucesso")
}
