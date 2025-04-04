package config

type Currency string
type TransactionType string

const (
	NGN Currency = "cNGN"
	XAF Currency = "cXAF"
	USD Currency = "cUSD"
	EUR Currency = "cEUR"

	Deposit  TransactionType = "DEPOSIT"
	Swap     TransactionType = "SWAP"
	Transfer TransactionType = "TRANSFER"
)

var SupportedCurrencyList = []Currency{
	NGN, XAF, USD, EUR,
}
