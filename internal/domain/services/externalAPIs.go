package services

import (
	"context"

	"operation-borderless/internal/domain/dto"
)

type ExternalAPIs interface {
	GetPairExchangeRate(ctx context.Context, baseCurrency, targetCurrency string) (dto.ExchangeRate, error)
	GetUSDExchangeRate(ctx context.Context) (dto.ExchangeRate, error)
	GetUserCountry(ctx context.Context, ipAddress string) (string, error)
}
