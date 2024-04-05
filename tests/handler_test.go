package tests

import (
	"context"
	"io"
	"net/http"
	"testing"
	"transaction-routine/internal/config"
	"transaction-routine/internal/entity"
	"transaction-routine/internal/server"
	"transaction-routine/tests/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type handlersTestSuite struct {
	suite.Suite
	ctrl   *gomock.Controller
	ctx    context.Context
	cfg    *config.Config
	accSvc *mocks.MockAccountService
	opSvc  *mocks.MockOpTypeService
	txSvc  *mocks.MockTransactionService
	srv    *http.Server
	url    string
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}

func (s *handlersTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.cfg = &config.Config{
		Port: 8081,
	}
	s.url = "http://localhost:8081"
	s.opSvc = mocks.NewMockOpTypeService(s.ctrl)
	s.accSvc = mocks.NewMockAccountService(s.ctrl)
	s.txSvc = mocks.NewMockTransactionService(s.ctrl)
	s.srv = server.NewServer(s.ctx, s.cfg, nil, s.accSvc, s.opSvc, s.txSvc)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.T().Fail()
		}
	}()
}

func (s *handlersTestSuite) TearDownSuite() {
	if err := s.srv.Shutdown(s.ctx); err != nil {
		s.T().Fail()
	}
}

func (s *handlersTestSuite) TestHandlers() {
	s.T().Run("getAccountHandler success", func(t *testing.T) {
		acc := entity.Account{ID: 1, DocumentNumber: "123456"}
		s.accSvc.EXPECT().GetAccountByID(gomock.Any(), acc.ID).Return(&acc, nil)
		resp, err := http.Get(s.url + "/accounts/1")
		if err != nil {
			t.Errorf("getAccountHandler request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("getAccountHandler status code: %d", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("getAccountHandler read body: %v", err)
		}
		if string(body) != `{"id":1,"document_number":"123456"}` {
			t.Errorf("getAccountHandler body: %s", body)
		}
	})
	s.T().Run("getAccountHandler invalid id", func(t *testing.T) {
		resp, err := http.Get(s.url + "/accounts/abcd")
		if err != nil {
			t.Errorf("getAccountHandler request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("getAccountHandler status code: %d", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("getAccountHandler read body: %v", err)
		}
		if string(body) != `{"message":"invalid account id"}` {
			t.Errorf("getAccountHandler body: %s", body)
		}
	})
	s.T().Run("getAccountHandler account not found", func(t *testing.T) {
		s.accSvc.EXPECT().GetAccountByID(gomock.Any(), 1).Return(nil, nil)
		resp, err := http.Get(s.url + "/accounts/1")
		if err != nil {
			t.Errorf("getAccountHandler request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("getAccountHandler status code: %d", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("getAccountHandler read body: %v", err)
		}
		if string(body) != `{"message":"account not found"}` {
			t.Errorf("getAccountHandler body: %s", body)
		}
	})
}
