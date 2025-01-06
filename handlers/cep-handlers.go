package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rafael1abrao/goexpert/multithreading/domain"
)

type CepHandler struct {
	providers []domain.CepProvider
}

type HandlerInterface interface {
	GetCep(w http.ResponseWriter, r *http.Request)
}

func NewCepHandler(providers []domain.CepProvider) HandlerInterface {
	return &CepHandler{
		providers: providers,
	}
}

func (ch *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cep := strings.TrimPrefix(r.URL.Path, "/cep/")
	cep = strings.TrimSpace(cep)

	if cep == "" {
		http.Error(w, `{"error": "CEP obrigatório na URL"}`, http.StatusBadRequest)
		return
	}

	result, err := domain.FetchFasterCep(context.Background(), cep, ch.providers)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Erro ao buscar CEP: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Erro ao converter resposta em JSON: %v"}`, err), http.StatusInternalServerError)
		return
	}
}
