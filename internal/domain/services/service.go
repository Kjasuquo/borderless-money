package services

import (
	"context"

	"operation-borderless/internal/domain/dto"
)

type Services interface {
	CreateWallet(ctx context.Context, email string) (string, error)
	DepositFunds(ctx context.Context, transaction dto.Transaction) (string, error)
	Transfer(ctx context.Context, transaction dto.Transaction) (string, error)
	GetUserWalletByCurrency(ctx context.Context, userID, currency string) (wallet dto.Wallet, err error)
	GetUserByID(ctx context.Context, id string) (user dto.User, err error)
	GetUserTransactions(ctx context.Context, userID string) (transactions []dto.Transaction, err error)
	CreateAuditLogs(ctx context.Context, auditLog dto.AuditLog) error
	GetUserByEmail(ctx context.Context, email string) (user dto.User, err error)
}
