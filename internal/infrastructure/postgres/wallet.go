package postgres

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"operation-borderless/internal/domain/model"
)

func (d *database) CreateWallet(ctx context.Context, wallet *model.Wallet) (string, error) {
	err := d.db.WithContext(ctx).Create(wallet).Error
	if err != nil {
		return "", err
	}
	return wallet.ID, nil
}

func (d *database) updateUserBalanceByCurrency(ctx context.Context, tx *gorm.DB, userID, currency string, balance float64) error {
	err := tx.WithContext(ctx).Model(&model.Wallet{}).
		Where("user_id = ? AND currency = ?", userID, currency).
		Update("balance", balance).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *database) createWallets(ctx context.Context, tx *gorm.DB, wallets []model.Wallet) error {

	err := tx.WithContext(ctx).Create(&wallets).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *database) GetUserWalletByCurrency(ctx context.Context, userID, currency string) (wallet model.Wallet, err error) {
	err = d.db.WithContext(ctx).
		Where("user_id = ?", userID).Where("currency = ?", currency).First(&wallet).Error

	return
}

func (d *database) Deposit(ctx context.Context, trx *model.Transaction) (string, error) {

	tx := d.db.Begin()
	if tx.Error != nil {
		return "", fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	wallet, err := d.GetUserWalletByCurrency(ctx, trx.ReceiverID, trx.ToCurrency)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed get wallet: %v", err)
	}

	balance := wallet.Balance + trx.ToAmount

	err = d.updateUserBalanceByCurrency(ctx, tx, trx.ReceiverID, trx.ToCurrency, balance)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed update wallet balance: %v", err)
	}

	trxID, err := d.createTransaction(ctx, tx, trx)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed create transaction: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	return trxID, nil
}

func (d *database) Transfer(ctx context.Context, trx *model.Transaction) (string, error) {
	tx := d.db.Begin()
	if tx.Error != nil {
		return "", fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	baseCurrencyWallet, err := d.GetUserWalletByCurrency(ctx, trx.SenderID, trx.FromCurrency)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed get wallet: %v", err)
	}

	wallet, err := d.GetUserWalletByCurrency(ctx, trx.ReceiverID, trx.ToCurrency)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed get wallet: %v", err)
	}

	if baseCurrencyWallet.Balance < trx.FromAmount {
		tx.Rollback()
		return "", fmt.Errorf("insufficient balance: %v", err)
	}

	baseCurrencyBalance := baseCurrencyWallet.Balance - trx.FromAmount
	balance := wallet.Balance + trx.ToAmount

	err = d.updateUserBalanceByCurrency(ctx, tx, trx.SenderID, trx.FromCurrency, baseCurrencyBalance)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed update wallet balance: %v", err)
	}

	err = d.updateUserBalanceByCurrency(ctx, tx, trx.ReceiverID, trx.ToCurrency, balance)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed update wallet balance: %v", err)
	}

	trxID, err := d.createTransaction(ctx, tx, trx)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed create transaction: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	return trxID, nil
}
