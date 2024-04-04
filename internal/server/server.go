package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"transaction-routine/internal/clock"
	"transaction-routine/internal/config"
	"transaction-routine/internal/database"
	"transaction-routine/internal/entity"
)

type Server struct {
	port    int
	db      database.Service
	cl      clock.Clock
	cfg     *config.Config
	opTypes entity.OperationType
}

func NewServer(ctx context.Context, cl clock.Clock, cfg *config.Config) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db, err := database.New(ctx, cfg)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}
	ops, err := db.GetOperationTypes(ctx)
	if err != nil {
		log.Fatalf("cannot get operation types: %s", err)
	}

	NewServer := &Server{
		port:    port,
		db:      db,
		cl:      cl,
		cfg:     cfg,
		opTypes: ops,
	}
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}

func fmtResponse(msg string) []byte {
	resp, _ := json.Marshal(map[string]string{
		"message": msg,
	})
	return resp
}
