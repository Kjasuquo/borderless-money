package dto

import "operation-borderless/pkg/config"

type Wallet struct {
	Models
	UserID   string  `json:"user_id"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type DepositRequest struct {
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
}

func IsValidCurrency(currency string) bool {
	switch currency {
	case string(config.USD), string(config.EUR), string(config.XAF), string(config.NGN):
		return true
	default:
		return false
	}
}

type TransferRequest struct {
	FromCurrency  string  `json:"from_currency" binding:"required"`
	ToCurrency    string  `json:"to_currency" binding:"required"`
	FromAmount    float64 `json:"from_amount"`
	ToAmount      float64 `json:"to_amount"`
	ReceiverEmail string  `json:"receiver_email"`
}
