package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/exp/slog"

	"github.com/giornetta/microshop/config"
	"github.com/giornetta/microshop/customers"
	"github.com/giornetta/microshop/customers/pg"
	"github.com/giornetta/microshop/events"
	"github.com/giornetta/microshop/kafka"
	"github.com/giornetta/microshop/log"
	"github.com/giornetta/microshop/postgres"
	"github.com/giornetta/microshop/server"
)

func main() {
	defer os.Exit(1)
	logger := slog.New(slog.NewTextHandler(os.Stderr))

	cfg, err := config.FromYaml("./config.yml")
	if err != nil {
		logger.Error("could not load yaml config", slog.String("err", err.Error()))
		runtime.Goexit()
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Setup Kafka
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Kafka.BrokerAddrs...),
		kgo.ConsumerGroup(cfg.Kafka.ConsumerGroup),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		logger.Error("could not create kafka client", slog.String("err", err.Error()))
		runtime.Goexit()
	}
	defer client.Close()

	producer := kafka.NewEventPublisher(client)
	listener := kafka.NewListener(client)

	// Setup Postgres
	pgPool, err := postgres.Connect(ctx, cfg.Postgres.ConnectionString())
	if err != nil {
		logger.Error("could not connect to postgres", slog.String("err", err.Error()))
		runtime.Goexit()
	}
	defer pgPool.Close()

	customerRepository := pg.NewCustomerRepository(pgPool)

	customersHandler := log.NewEventHandler(
		logger.With("svc", "CustomerHandler"),
		customers.NewCustomerHandler(customerRepository),
	)
	listener.Handle(events.CustomerTopic, customersHandler)

	customerService := customers.NewLoggingService(
		logger.With("svc", "Service"),
		customers.NewService(customerRepository, producer),
	)

	s := server.New(customers.NewRouter(customerService), &server.Options{
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
			logger.Error("could not listen to events", slog.String("err", err.Error()))
			signals <- os.Interrupt
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		logger.Info("Server started", slog.Int("port", cfg.Server.Port))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("could not run server", slog.String("err", err.Error()))
			signals <- os.Interrupt
		}

		wg.Done()
	}()

	<-signals
	logger.Info("Shutting down server")
	cancel()
	s.Shutdown(ctx)

	wg.Wait()
}
