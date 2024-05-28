package diveGames

const (
	PairUSD = "BTC/USD"
	PairEUR = "BTC/EUR"
	PairCHF = "BTC/CHF"

	KrakenURL = "https://api.kraken.com/0/public/Trades?pair=%s&count%d"
)

var PairValues = [...]string{PairUSD, PairCHF, PairEUR}
