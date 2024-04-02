package main

import (
	"context"
	"log"
	"transaction-routine/internal/clock"
	"transaction-routine/internal/server"
)

func main() {
	ctx := context.Background()
	cl := clock.New()
	server := server.NewServer(ctx, cl)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("cannot start server: %s", err)
	}
}
