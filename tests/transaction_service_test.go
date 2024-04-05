package tests

import (
	"context"
	"errors"
	"testing"
	"time"
	"transaction-routine/internal/entity"
	"transaction-routine/internal/service"
	"transaction-routine/tests/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type transactionSvcTestSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	ctx     context.Context
	repo    *mocks.MockRepository
	cl      *mocks.MockClock
	opTypes entity.OperationType
	txSvc   service.TransactionService
}

func TestTransactionSvcSuite(t *testing.T) {
	suite.Run(t, new(transactionSvcTestSuite))
}

func (s *transactionSvcTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockRepository(s.ctrl)
	s.cl = mocks.NewMockClock(s.ctrl)
	s.opTypes = entity.OperationType{
		1: &entity.Operation{Description: "COMPRA A VISTA", PositiveAmount: false},
		2: &entity.Operation{Description: "PAGAMENTO", PositiveAmount: true},
	}
	s.txSvc = service.NewTransactionService(s.cl, s.repo, s.opTypes)
}

func (s *transactionSvcTestSuite) TestCreateTransaction() {
	now := time.Now()
	s.T().Run("success", func(t *testing.T) {
		tx := entity.Transaction{AccountID: 1, OperationTypeID: 2, Amount: decimal.NewFromInt(100), EventDate: now}
		s.cl.EXPECT().Now().Return(now)
		s.repo.EXPECT().CreateTransaction(gomock.Any(), tx).Return(nil)
		err := s.txSvc.CreateTransaction(s.ctx, tx)
		s.NoError(err)
	})

	s.T().Run("invalid transaction", func(t *testing.T) {
		tx := entity.Transaction{AccountID: 1, OperationTypeID: 5, Amount: decimal.NewFromInt(100), EventDate: now}
		s.cl.EXPECT().Now().Return(tx.EventDate)
		err := s.txSvc.CreateTransaction(s.ctx, tx)
		s.Error(err)
	})

	s.T().Run("repo error", func(t *testing.T) {
		tx := entity.Transaction{AccountID: 1, OperationTypeID: 1, Amount: decimal.NewFromInt(-100), EventDate: now}
		s.cl.EXPECT().Now().Return(tx.EventDate)
		s.repo.EXPECT().CreateTransaction(gomock.Any(), tx).Return(errors.New("error"))
		err := s.txSvc.CreateTransaction(s.ctx, tx)
		s.Error(err)
	})
}

func (s *transactionSvcTestSuite) TestUpdateTransaction() {
	now := time.Now()
	s.T().Run("success", func(t *testing.T) {
		tx := entity.Transaction{ID: 1, AccountID: 1, OperationTypeID: 2, Amount: decimal.NewFromInt(100), EventDate: now}
		s.repo.EXPECT().FindTransactions(gomock.Any(), entity.TransactionFilter{ID: &tx.ID}).Return([]entity.Transaction{tx}, nil)
		s.repo.EXPECT().UpdateTransaction(gomock.Any(), tx).Return(nil)
		err := s.txSvc.UpdateTransaction(s.ctx, tx)
		s.NoError(err)
	})

	s.T().Run("invalid transaction", func(t *testing.T) {
		tx := entity.Transaction{ID: 1, AccountID: 1, OperationTypeID: 5, Amount: decimal.NewFromInt(100), EventDate: now}
		err := s.txSvc.UpdateTransaction(s.ctx, tx)
		s.Error(err)
	})

	s.T().Run("not found", func(t *testing.T) {
		tx := entity.Transaction{ID: 1, AccountID: 1, OperationTypeID: 2, Amount: decimal.NewFromInt(100), EventDate: now}
		s.repo.EXPECT().FindTransactions(gomock.Any(), entity.TransactionFilter{ID: &tx.ID}).Return(nil, nil)
		err := s.txSvc.UpdateTransaction(s.ctx, tx)
		s.Error(err)
	})

	s.T().Run("repo error", func(t *testing.T) {
		tx := entity.Transaction{ID: 1, AccountID: 1, OperationTypeID: 2, Amount: decimal.NewFromInt(100), EventDate: now}
		s.repo.EXPECT().FindTransactions(gomock.Any(), entity.TransactionFilter{ID: &tx.ID}).Return([]entity.Transaction{tx}, nil)
		s.repo.EXPECT().UpdateTransaction(gomock.Any(), tx).Return(errors.New("error"))
		err := s.txSvc.UpdateTransaction(s.ctx, tx)
		s.Error(err)
	})
}
