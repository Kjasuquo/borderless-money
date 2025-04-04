package repositories

import (
	"context"

	"operation-borderless/internal/domain/model"
)

type Repository interface {
	CreateUserWallet(ctx context.Context, email string) (string, error)
	GetUserByID(ctx context.Context, id string) (user model.User, err error)
	GetUserWalletByCurrency(ctx context.Context, userID, currency string) (wallet model.Wallet, err error)
	Deposit(ctx context.Context, trx *model.Transaction) (string, error)
	Transfer(ctx context.Context, trx *model.Transaction) (string, error)
	GetUserTransactions(ctx context.Context, userID string) (transaction []model.Transaction, err error)
	CreateAuditLogs(ctx context.Context, auditLog model.AuditLog) error
	GetUserByEmail(ctx context.Context, email string) (user model.User, err error)
}
