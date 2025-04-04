package model

import (
	"operation-borderless/internal/domain/dto"
)

type Transaction struct {
	Models
	SenderID       string `gorm:"index"`
	ReceiverID     string `gorm:"index"`
	FromAmount     float64
	ToAmount       float64
	FromCurrency   string
	ToCurrency     string
	Type           string
	ConversionRate float64

	// Foreign key relationships
	Sender   User `gorm:"foreignKey:SenderID;references:ID"`
	Receiver User `gorm:"foreignKey:ReceiverID;references:ID"`
}

func (t *Transaction) ToTransactionDTO() dto.Transaction {

	out := dto.Transaction{
		Models: dto.Models{
			ID:        t.ID,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
		SenderID:       t.SenderID,
		ReceiverID:     t.ReceiverID,
		FromAmount:     t.FromAmount,
		ToAmount:       t.ToAmount,
		FromCurrency:   t.FromCurrency,
		ToCurrency:     t.ToCurrency,
		Type:           t.Type,
		ConversionRate: t.ConversionRate,
	}

	return out
}

func FromTransactionDTO(t dto.Transaction) Transaction {
	response := Transaction{
		Models: Models{
			ID:        t.ID,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
		SenderID:       t.SenderID,
		ReceiverID:     t.ReceiverID,
		FromAmount:     t.FromAmount,
		ToAmount:       t.ToAmount,
		FromCurrency:   t.FromCurrency,
		ToCurrency:     t.ToCurrency,
		Type:           t.Type,
		ConversionRate: t.ConversionRate,
	}

	return response
}
