package trasnsformations

var MapPairsToKrakenKeys = map[string]string{
	"BTC/USD": "XXBTZUSD",
	"BTC/EUR": "XXBTZEUR",
	"BTC/CHF": "XBTCHF",
}

var MapKrakenKeysToPairs = map[string]string{
	"XXBTZUSD": "BTC/USD",
	"XXBTZEUR": "BTC/EUR",
	"XBTCHF":   "BTC/CHF",
}
