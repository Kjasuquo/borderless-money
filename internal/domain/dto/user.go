package dto

import (
	"errors"
	"strings"
)

type User struct {
	Models
	Email           string `json:"email" binding:"required"`
	Wallets         []Wallet
	TotalUSDBalance float64 `json:"totalUSDBalance"`
}

func (u *User) ConfirmEmailFormat() error {
	parts := strings.Split(u.Email, "@")
	if len(parts) != 2 {
		return errors.New("wrong email format")
	}
	return nil
}
