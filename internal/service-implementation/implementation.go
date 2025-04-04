package service_implementation

import (
	"context"
	"net/http"

	"operation-borderless/internal/domain/dto"
	"operation-borderless/internal/domain/model"
	"operation-borderless/internal/domain/repositories"
	"operation-borderless/pkg/config"
)

type ServiceClient struct {
	Http       http.Client
	Config     *config.Config
	repository repositories.Repository
}

func NewServiceClient(conf *config.Config, repository repositories.Repository) *ServiceClient {
	return &ServiceClient{Config: conf, repository: repository}
}

func (s *ServiceClient) CreateWallet(ctx context.Context, email string) (string, error) {

	userID, err := s.repository.CreateUserWallet(ctx, email)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *ServiceClient) DepositFunds(ctx context.Context, transaction dto.Transaction) (string, error) {

	t := model.FromTransactionDTO(transaction)
	trxID, err := s.repository.Deposit(ctx, &t)
	if err != nil {
		return "", err
	}

	return trxID, err
}

func (s *ServiceClient) Transfer(ctx context.Context, transaction dto.Transaction) (string, error) {

	t := model.FromTransactionDTO(transaction)
	trxID, err := s.repository.Transfer(ctx, &t)
	if err != nil {
		return "", err
	}

	return trxID, err
}

func (s *ServiceClient) GetUserWalletByCurrency(ctx context.Context, userID, currency string) (wallet dto.Wallet, err error) {
	modelWallet, err := s.repository.GetUserWalletByCurrency(ctx, userID, currency)
	if err != nil {
		return dto.Wallet{}, err
	}

	wallet = modelWallet.ToWalletDTO()

	return wallet, nil
}

func (s *ServiceClient) GetUserByID(ctx context.Context, id string) (user dto.User, err error) {
	modelUser, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		return dto.User{}, err
	}

	user = modelUser.ToUserDTO()

	return user, nil
}

func (s *ServiceClient) GetUserTransactions(ctx context.Context, userID string) (transactions []dto.Transaction, err error) {
	trans, err := s.repository.GetUserTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, tx := range trans {
		transactions = append(transactions, tx.ToTransactionDTO())
	}

	return transactions, nil
}

func (s *ServiceClient) CreateAuditLogs(ctx context.Context, auditLog dto.AuditLog) error {
	al := model.FromAuditLogDTO(auditLog)
	err := s.repository.CreateAuditLogs(ctx, al)
	if err != nil {
		return err
	}
	return nil
}
