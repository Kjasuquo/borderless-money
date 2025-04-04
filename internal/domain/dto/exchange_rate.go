package dto

type ExchangeRate struct {
	Result             string          `json:"result"`
	TimeLastUpdateUnix int             `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string          `json:"time_last_update_utc"`
	TimeNextUpdateUnix int             `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string          `json:"time_next_update_utc"`
	BaseCode           string          `json:"base_code"`
	TargetCode         string          `json:"target_code"`
	ConversionRate     float64         `json:"conversion_rate"`
	ConversionRates    ConversionRates `json:"conversion_rates"`
}

type ConversionRates struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
	NGN float64 `json:"NGN"`
	XAF float64 `json:"XAF"`
}
