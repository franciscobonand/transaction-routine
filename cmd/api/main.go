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
	"transaction-routine/internal/database"
	"transaction-routine/internal/server"
	"transaction-routine/internal/service"
)

func main() {
	appCtx, appStopCtx := context.WithCancel(context.Background())
	cl := clock.New()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}
	db, err := database.New(appCtx, cfg)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}
	opTypes, err := db.FindOperationType(appCtx)
	if err != nil {
		log.Fatalf("cannot get operation types to initialize application: %s", err)
	}

	healthSvc := service.NewHealthService(db)
	accsvc := service.NewAccountService(db)
	opsvc := service.NewOpTypeService(db, opTypes)
	txsvc := service.NewTransactionService(cl, db, opTypes)
	srv := server.NewServer(appCtx, cfg, healthSvc, accsvc, opsvc, txsvc)

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, cancel := context.WithTimeout(appCtx, 30*time.Second)
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
		appStopCtx()
	}()

	log.Printf("Server running on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server terminated with error: %s", err)
	}

	<-appCtx.Done()
}
