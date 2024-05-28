package domain

import (
	"time"
)

type Trade struct {
	Pair        string
	Price       string
	Volume      string
	Time        time.Time
	BuySell     string
	MarketLimit string
	Misc        string
	TradeID     int
}
