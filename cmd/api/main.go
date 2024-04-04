package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transaction-routine/internal/clock"
	"transaction-routine/internal/config"
	"transaction-routine/internal/server"
)

func main() {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	cl := clock.New()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}
	srv := server.NewServer(serverCtx, cl, cfg)

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	log.Printf("Server running on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server terminated with error: %s", err)
	}

	<-serverCtx.Done()
}
