package server

import (
	"encoding/json"
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
	jsonResp, _ := json.Marshal(s.db.Health(r.Context()))
	_, _ = w.Write(jsonResp)
}

func (s *Server) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DocumentNumber string `json:"document_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.db.CreateAccount(r.Context(), req.DocumentNumber); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing account id", http.StatusBadRequest)
		return
	}
	numid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	acc, err := s.db.GetAccount(r.Context(), numid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResp, _ := json.Marshal(acc)
	_, _ = w.Write(jsonResp)
}

func (s *Server) createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.Transaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.EventDate = s.cl.Now()
	if err := req.Validate(s.opTypes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.db.CreateTransaction(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
