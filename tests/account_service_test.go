package tests

import (
	"context"
	"errors"
	"testing"
	"transaction-routine/internal/entity"
	"transaction-routine/internal/service"
	"transaction-routine/tests/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type accountSvcTestSuite struct {
	suite.Suite
	ctrl   *gomock.Controller
	ctx    context.Context
	repo   *mocks.MockRepository
	accSvc service.AccountService
}

func TestAccountSvcSuite(t *testing.T) {
	suite.Run(t, new(accountSvcTestSuite))
}

func (s *accountSvcTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockRepository(s.ctrl)
	s.accSvc = service.NewAccountService(s.repo)
}

func (s *accountSvcTestSuite) TestGetAccountByID() {
	s.T().Run("success", func(t *testing.T) {
		acc := entity.Account{ID: 1, DocumentNumber: "123456"}
		s.repo.EXPECT().FindAccounts(gomock.Any(), entity.AccountFilter{ID: &acc.ID}).Return([]entity.Account{acc}, nil)
		res, err := s.accSvc.GetAccountByID(s.ctx, acc.ID)
		s.NoError(err)
		s.Equal(&acc, res)
	})

	s.T().Run("repo error", func(t *testing.T) {
		acc := entity.Account{ID: 1, DocumentNumber: "123456"}
		s.repo.EXPECT().FindAccounts(gomock.Any(), entity.AccountFilter{ID: &acc.ID}).Return(nil, errors.New("error"))
		res, err := s.accSvc.GetAccountByID(s.ctx, acc.ID)
		s.Error(err)
		s.Nil(res)
	})

	s.T().Run("not found", func(t *testing.T) {
		acc := entity.Account{ID: 1, DocumentNumber: "123456"}
		s.repo.EXPECT().FindAccounts(gomock.Any(), entity.AccountFilter{ID: &acc.ID}).Return(nil, nil)
		res, err := s.accSvc.GetAccountByID(s.ctx, acc.ID)
		s.NoError(err)
		s.Nil(res)
	})
}

func (s *accountSvcTestSuite) TestCreateAccount() {
	s.T().Run("success", func(t *testing.T) {
		acc := entity.Account{DocumentNumber: "123456"}
		s.repo.EXPECT().CreateAccount(gomock.Any(), acc).Return(nil)
		err := s.accSvc.CreateAccount(s.ctx, acc)
		s.NoError(err)
	})

	s.T().Run("repo error", func(t *testing.T) {
		acc := entity.Account{DocumentNumber: "123456"}
		s.repo.EXPECT().CreateAccount(gomock.Any(), acc).Return(errors.New("error"))
		err := s.accSvc.CreateAccount(s.ctx, acc)
		s.Error(err)
	})

	s.T().Run("missing document number", func(t *testing.T) {
		acc := entity.Account{}
		err := s.accSvc.CreateAccount(s.ctx, acc)
		s.Error(err)
		s.True(errors.Is(err, entity.ErrMissingDocumentNumber))
	})
}

func (s *accountSvcTestSuite) TestGetAccountBalance() {
	s.T().Run("success", func(t *testing.T) {
		id := 1
		txs := []entity.Transaction{
			{AccountID: id, Amount: decimal.NewFromInt(10)},
			{AccountID: id, Amount: decimal.NewFromInt(20)},
			{AccountID: id, Amount: decimal.NewFromInt(-5)},
		}
		s.repo.EXPECT().FindTransactions(gomock.Any(), entity.TransactionFilter{AccountID: &id}).Return(txs, nil)
		res, err := s.accSvc.GetAccountBalance(s.ctx, id)
		s.NoError(err)
		s.Equal("25", res.String())
	})

	s.T().Run("repo error", func(t *testing.T) {
		id := 1
		s.repo.EXPECT().FindTransactions(gomock.Any(), entity.TransactionFilter{AccountID: &id}).Return(nil, errors.New("error"))
		res, err := s.accSvc.GetAccountBalance(s.ctx, id)
		s.Error(err)
		s.Equal("0", res.String())
	})
}
