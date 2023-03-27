package main

import (
	"log"
	"net/http"

	"github.com/giornetta/microshop/kafka"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/giornetta/microshop/services/inventory"
	"github.com/giornetta/microshop/services/inventory/inmem"
	"github.com/giornetta/microshop/services/inventory/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	brokerAddr := []string{"localhost:9092"}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokerAddr...),
	)
	if err != nil {
		log.Fatalf("could not create kafka client: %v", err)
	}

	producer := kafka.NewEventPublisher(client)

	productRepository := inmem.NewRepository()
	productService := inventory.NewService(productRepository, producer)

	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
	)

	router.Mount("/api/inventory", server.Router(productService))

	http.ListenAndServe(":3000", router)
}
