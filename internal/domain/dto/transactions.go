package dto

type Transaction struct {
	Models
	SenderID       string  `json:"sender_id"`
	ReceiverID     string  `json:"receiver_id"`
	FromAmount     float64 `json:"from_amount"`
	ToAmount       float64 `json:"to_amount"`
	FromCurrency   string  `json:"from_currency"`
	ToCurrency     string  `json:"to_currency"`
	Type           string  `json:"type"`
	ConversionRate float64 `json:"conversion_rate"`
}
