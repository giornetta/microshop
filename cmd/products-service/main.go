package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/giornetta/microshop/config"
	"github.com/giornetta/microshop/events"
	"github.com/giornetta/microshop/kafka"
	"github.com/giornetta/microshop/postgres"
	"github.com/giornetta/microshop/products"
	"github.com/giornetta/microshop/server"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/giornetta/microshop/products/http"
	"github.com/giornetta/microshop/products/pg"
)

func main() {
	cfg, err := config.FromYaml("./config.yml")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Setup Kafka
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Kafka.BrokerAddrs...),
		kgo.ConsumerGroup(cfg.Kafka.ConsumerGroup),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		log.Fatalf("could not create kafka client: %v", err)
	}
	defer client.Close()

	producer := kafka.NewEventPublisher(client)
	listener := kafka.NewListener(client)

	// Setup Postgres
	pgPool, err := postgres.Connect(ctx, cfg.Postgres.ConnectionString())
	if err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}
	defer pgPool.Close()

	productRepository := pg.NewProductRepository(pgPool)

	productHandler := products.NewProductHandler(productRepository)
	listener.Handle(events.ProductTopic, productHandler)

	productService := products.NewService(productRepository, producer)

	s := server.New(http.Router(productService), &server.Options{
		Port:         cfg.Server.Port,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 120,
	})
	defer s.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if err := listener.Listen(ctx); err != nil {
			log.Println(err)
			signals <- os.Interrupt
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Printf("Listening on port %d\n", cfg.Server.Port)
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}

		wg.Done()
	}()

	<-signals
	cancel()
	s.Shutdown(ctx)

	wg.Wait()
}
