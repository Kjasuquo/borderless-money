package postgres

import (
	"context"
	"gorm.io/gorm"

	"operation-borderless/internal/domain/model"
)

func (d *database) createTransaction(ctx context.Context, tx *gorm.DB, transaction *model.Transaction) (string, error) {
	err := tx.WithContext(ctx).Create(transaction).Error
	if err != nil {
		return "", err
	}
	return transaction.ID, nil
}

func (d *database) GetUserTransactions(ctx context.Context, userID string) (transaction []model.Transaction, err error) {
	err = d.db.WithContext(ctx).
		Where("sender_id = ? OR receiver_id = ? ", userID, userID).
		Find(&transaction).Error

	return
}
