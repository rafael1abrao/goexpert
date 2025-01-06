package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rafael1abrao/goexpert/multithreading/domain"
	"github.com/stretchr/testify/assert"
)

type mockCepProvider struct {
	result domain.CepResponse
	err    error
}

func (m *mockCepProvider) FetchCep(ctx context.Context, cep string) (*domain.CepResponse, error) {
	return &m.result, m.err
}

func TestGetCep(t *testing.T) {
	tests := []struct {
		name           string
		cep            string
		providers      []domain.CepProvider
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid cep",
			cep:            "12345678",
			providers:      []domain.CepProvider{&mockCepProvider{result: domain.CepResponse{Cep: "12345678"}}},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Source":"", "cep":"12345678", "logradouro":"", "bairro":"", "localidade":"", "uf":""}`,
		},
		{
			name:           "empty cep",
			cep:            "",
			providers:      []domain.CepProvider{&mockCepProvider{}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "CEP obrigatório na URL"}`,
		},
		{
			name:           "error fetching cep",
			cep:            "12345678",
			providers:      []domain.CepProvider{&mockCepProvider{err: errors.New("fetch error")}},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "Erro ao buscar CEP: nenhuma das requisições retornou com sucesso"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cep/"+tt.cep, nil)
			w := httptest.NewRecorder()

			handler := NewCepHandler(tt.providers)
			handler.GetCep(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var body strings.Builder
			_, err := io.Copy(&body, resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expectedBody, body.String())
		})
	}
}
