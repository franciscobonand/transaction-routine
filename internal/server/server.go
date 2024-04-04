package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"transaction-routine/internal/config"
	"transaction-routine/internal/service"
)

type Server struct {
	port      int
	cfg       *config.Config
	healthsvc service.HealthService
	accsvc    service.AccountService
	opsvc     service.OpTypeService
	txsvc     service.TransactionService
}

func NewServer(
	ctx context.Context,
	cfg *config.Config,
	healthSvc service.HealthService,
	accSvc service.AccountService,
	opSvc service.OpTypeService,
	tSvc service.TransactionService,
) *http.Server {
	NewServer := &Server{
		port:      cfg.Port,
		cfg:       cfg,
		healthsvc: healthSvc,
		accsvc:    accSvc,
		opsvc:     opSvc,
		txsvc:     tSvc,
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
