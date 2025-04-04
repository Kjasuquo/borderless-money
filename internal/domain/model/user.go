package model

import (
	"operation-borderless/internal/domain/dto"
)

type User struct {
	Models
	Email   string   `gorm:"not null;unique"`
	Wallets []Wallet `gorm:"foreignKey:UserID;references:ID"`
}

func (u *User) ToUserDTO() dto.User {
	var wallets []dto.Wallet
	for _, v := range u.Wallets {
		wallets = append(wallets, v.ToWalletDTO())
	}
	out := dto.User{
		Models: dto.Models{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Email:   u.Email,
		Wallets: wallets,
	}

	return out
}

func FromUserDTO(u dto.User) User {
	var wallets []Wallet
	for _, v := range u.Wallets {
		wallets = append(wallets, FromWalletDTO(v))
	}
	response := User{
		Models: Models{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Email:   u.Email,
		Wallets: wallets,
	}

	return response
}
