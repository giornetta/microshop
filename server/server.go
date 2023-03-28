package server

import (
	"fmt"
	"net/http"
	"time"
)

func New(h http.Handler, opt *Options) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", opt.Port),
		Handler: h,
		// TODO Add TLS Support
		TLSConfig:    nil,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		IdleTimeout:  opt.IdleTimeout,
	}
}

type Options struct {
	Port int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}
