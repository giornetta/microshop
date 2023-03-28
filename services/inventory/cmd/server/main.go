package main

import (
	"log"
	"time"

	"github.com/giornetta/microshop/kafka"
	"github.com/giornetta/microshop/server"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/giornetta/microshop/services/inventory"
	"github.com/giornetta/microshop/services/inventory/http"
	"github.com/giornetta/microshop/services/inventory/inmem"
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

	router.Mount("/api/inventory", http.Router(productService))

	opts := &server.Options{
		Port:         3000,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 10,
	}
	s := server.New(router, opts)
	defer s.Close()

	log.Printf("Listening on port %d\n", opts.Port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
