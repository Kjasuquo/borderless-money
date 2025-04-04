package postgres

import (
	"context"
	"fmt"
	"gorm.io/gorm"

	"operation-borderless/internal/domain/model"
	"operation-borderless/pkg/config"
)

func (d *database) createUser(ctx context.Context, tx *gorm.DB, user *model.User) (string, error) {
	err := tx.WithContext(ctx).Create(user).Error
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (d *database) GetUserByID(ctx context.Context, id string) (user model.User, err error) {
	err = d.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Wallets").
		First(&user).Error

	return
}

func (d *database) CreateUserWallet(ctx context.Context, email string) (string, error) {

	tx := d.db.Begin()
	if tx.Error != nil {
		return "", fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := &model.User{
		Email: email,
	}

	userID, err := d.createUser(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	var wallets []model.Wallet
	for _, currency := range config.SupportedCurrencyList {
		wallets = append(wallets, model.Wallet{
			UserID:   userID,
			Currency: string(currency),
			Balance:  0,
		})

	}

	err = d.createWallets(ctx, tx, wallets)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create wallets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}
