package main

import (
	"context"
	"errors"
	"log"

	"github.com/si74/flow-api/internal/flowd"
)

func main() {
	// TODO(make this configurable via a flag)
	addr := ":8080"

	// Initiate prometheus
	// Initiate logrus logger

	ctx := context.Background()

	srv, err := flowd.NewServer(addr)
	if err != nil {
		log.Fatalf("unable to create new flow server: %v", err)
	}

	if err := srv.Serve(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("failed to run flowd server: %v", err)
	}
}
