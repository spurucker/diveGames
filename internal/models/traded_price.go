package models

type TradePrice struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}
