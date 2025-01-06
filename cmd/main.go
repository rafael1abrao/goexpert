// cmd/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rafael1abrao/goexpert/multithreading/domain"
	"github.com/rafael1abrao/goexpert/multithreading/handlers"
	"github.com/rafael1abrao/goexpert/multithreading/providers"
)

func main() {
	client := &http.Client{Timeout: 2 * time.Second}

	providerBrasil := providers.NewBrasilAPIProvider(client)
	providerViaCep := providers.NewViaCepProvider(client)

	providersList := []domain.CepProvider{
		providerBrasil,
		providerViaCep,
	}

	mux := http.NewServeMux()
	cepHandler := handlers.NewCepHandler(providersList)

	mux.HandleFunc("/ceps/", cepHandler.GetCep)

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	go func() {
		log.Println("Servidor iniciando na porta 8000...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Servidor encerrando...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao encerrar servidor: %v", err)
	}

	log.Println("Servidor encerrado com sucesso.")
}
