package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"transaction-routine/internal/entity"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", s.healthHandler)

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/{id}", s.getAccountHandler)
		r.Post("/", s.createAccountHandler)
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", s.createTransactionHandler)
	})
	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if isHealthy := s.healthsvc.HealthCheck(r.Context()); isHealthy {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fmtResponse("service is healthy"))
	}
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(fmtResponse("service is unhealthy"))
}

func (s *Server) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(fmtResponse(err.Error()))
		return
	}

	if err := s.accsvc.CreateAccount(r.Context(), req); err != nil {
		if errors.Is(err, entity.ErrMissingDocumentNumber) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(fmtResponse(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(fmtResponse("failed to create account"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(fmtResponse("missing account id"))
		return
	}
	numid, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(fmtResponse("invalid account id"))
		return
	}

	acc, err := s.accsvc.GetAccountByID(r.Context(), numid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(fmtResponse("failed to get account"))
		return
	}
	if acc == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(fmtResponse("account not found"))
		return
	}

	jsonResp, _ := json.Marshal(acc)
	_, _ = w.Write(jsonResp)
}

func (s *Server) createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.Transaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(fmtResponse(err.Error()))
		return
	}

	if err := s.txsvc.CreateTransaction(r.Context(), req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("failed to create transaction: %s", err.Error())
		_, _ = w.Write(fmtResponse(msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
