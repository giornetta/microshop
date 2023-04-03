package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/giornetta/microshop/events"
	"github.com/giornetta/microshop/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	client, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
	)
	if err != nil {
		log.Fatalf("could not connect to kafka: %v", err)
	}
	defer client.Close()

	listener := kafka.NewListener(client)

	ctx, cancel := context.WithCancel(context.Background())

	listener.Handle(events.ProductTopic, &h{})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if err := listener.Listen(ctx); err != nil {
			log.Println(err)
		}

		wg.Done()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	<-signals
	cancel()

	wg.Wait()
}

type h struct {
}

func (h *h) Handle(e events.Event, ctx context.Context) error {
	fmt.Printf("got: %v, %v\n", e.Type(), e.Key())
	return nil
}
