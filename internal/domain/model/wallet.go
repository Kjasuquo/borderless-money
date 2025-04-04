package model

import "operation-borderless/internal/domain/dto"

type Wallet struct {
	Models
	UserID   string `gorm:"index:user_currency,unique"`
	Currency string `gorm:"index:user_currency,unique"`
	Balance  float64
}

func (w *Wallet) ToWalletDTO() dto.Wallet {
	out := dto.Wallet{
		Models: dto.Models{
			ID:        w.ID,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		},
		UserID:   w.UserID,
		Currency: w.Currency,
		Balance:  w.Balance,
	}

	return out
}

func FromWalletDTO(w dto.Wallet) Wallet {
	response := Wallet{
		Models: Models{
			ID:        w.ID,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		},
		UserID:   w.UserID,
		Currency: w.Currency,
		Balance:  w.Balance,
	}

	return response
}
