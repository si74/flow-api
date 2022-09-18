package main

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/si74/flow-api/internal/flowd"
	"github.com/sirupsen/logrus"
)

func main() {
	// TODO(sneha): make this configurable via a flag)
	addr := ":8080"

	// TODO(sneha): make debug logging configurable

	// Initiate prometheus
	// Initiate logrus logger
	ll := logrus.New()
	ll.SetFormatter(&logrus.TextFormatter{})
	ll.SetLevel(logrus.DebugLevel)

	// Create a cancellable context that may be cancelled in response
	// to keyboard interrupts and other system signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	srv, err := flowd.NewServer(addr, ll, reg)
	if err != nil {
		log.Fatalf("unable to create new flow server: %v", err)
	}

	ll.Info("starting flowd server...")
	if err := srv.Serve(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("failed to run flowd server: %v", err)
	}
}
