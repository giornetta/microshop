package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"golang.org/x/exp/slog"

	"github.com/giornetta/microshop/auth"
	"github.com/giornetta/microshop/config"
	"github.com/giornetta/microshop/identity"
	"github.com/giornetta/microshop/identity/pg"
	"github.com/giornetta/microshop/postgres"
	"github.com/giornetta/microshop/server"
)

func main() {
	defer os.Exit(1)
	logger := slog.New(slog.NewTextHandler(os.Stderr))

	cfg, err := config.FromYaml("./config.yml")
	if err != nil {
		logger.Error("could not load yaml config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Setup Postgres
	pgPool, err := postgres.Connect(ctx, cfg.Postgres.ConnectionString())
	if err != nil {
		logger.Error("could not connect to postgres", slog.String("err", err.Error()))
		runtime.Goexit()
	}
	defer pgPool.Close()

	identityRepository := pg.NewIdentityRepository(pgPool)

	// Setup Issuer
	authIssuer, err := auth.NewJWTIssuer(
		[]byte(cfg.Auth.Key),
		auth.WithPEM(false),
		auth.WithName(cfg.Auth.Issuer),
		auth.WithTokenDuration(cfg.Auth.TokenDuration),
	)
	if err != nil {
		logger.Error("could not create issuer", slog.String("err", err.Error()))
		runtime.Goexit()
	}

	// Setup Verifier
	authVerifier, err := auth.NewJWTVerifier(
		[]byte(cfg.Auth.Key),
		auth.WithAllowedIssuer(cfg.Auth.Issuer),
	)
	if err != nil {
		logger.Error("could not create verifier", slog.String("err", err.Error()))
		runtime.Goexit()
	}

	identityService := identity.NewService(authIssuer, identityRepository)

	s := server.New(identity.NewRouter(identityService, authVerifier), &server.Options{
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
