package models

type LastTradePrices struct {
	Trades []*TradePrice `json:"ltp"`
}
