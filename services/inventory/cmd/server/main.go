package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/giornetta/microshop/events"
	"github.com/giornetta/microshop/kafka"
	"github.com/giornetta/microshop/server"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/giornetta/microshop/services/inventory/http"
	"github.com/giornetta/microshop/services/inventory/inmem"
	"github.com/giornetta/microshop/services/inventory/service"
)

const port int = 3000

func main() {
	brokerAddr := []string{"localhost:9092"}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokerAddr...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		log.Fatalf("could not create kafka client: %v", err)
	}

	producer := kafka.NewEventPublisher(client)
	listener := kafka.NewListener(client)

	productRepository := inmem.NewRepository()

	productHandler := service.NewProductHandler(productRepository)
	listener.Handle(events.ProductTopic, productHandler)

	productService := service.New(productRepository, producer)

	s := server.New(http.Router(productService), &server.Options{
		Port:         port,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 10,
	})
	defer s.Close()

	ctx, cancel := context.WithCancel(context.Background())

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
		log.Printf("Listening on port %d\n", port)
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
