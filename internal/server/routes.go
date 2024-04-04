package server

import (
	"encoding/json"
	"log"
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
	msg := s.db.Health(r.Context())
	_, _ = w.Write(fmtResponse(msg))
}

func (s *Server) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DocumentNumber string `json:"document_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(fmtResponse(err.Error()))
		return
	}

	if err := s.db.CreateAccount(r.Context(), req.DocumentNumber); err != nil {
		log.Printf("error creating account: %s", err)
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

	acc, err := s.db.GetAccount(r.Context(), numid)
	if err != nil {
		log.Printf("error getting account %d: %s", numid, err)
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

	req.EventDate = s.cl.Now()
	if err := req.Validate(s.opTypes); err != nil {
		log.Printf("error validating transaction: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(fmtResponse(err.Error()))
		return
	}

	if err := s.db.CreateTransaction(r.Context(), req); err != nil {
		log.Printf("error creating transaction: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(fmtResponse("failed to create transaction"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
