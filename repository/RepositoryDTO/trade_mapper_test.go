package RepositoryDTO

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

const jsonTestData = `{
        "error": [],
        "result": {
            "XXBTZUSD": [
                [
                    "30243.40000",
                    "0.34507674",
                    1688669597.8277369,
                    "b",
                    "m",
                    "",
                    61044952
                ]
            ],
            "last": "1688671969993150842"
        }
    }`

func TestMapTradeResponse(t *testing.T) {
	tm := TradeMapper{}
	body := io.NopCloser(bytes.NewBufferString(jsonTestData))

	trade, err := tm.MapResponseToTrade(body, "XXBTZUSD")

	assert.Nil(t, err)
	assert.NotNil(t, trade)
	assert.Equal(t, trade.Pair, "BTC/USD")
	assert.Equal(t, trade.Price, "30243.40000")
	assert.Equal(t, trade.Volume, "0.34507674")
	assert.Equal(t, trade.Time, tm.timestampToTime(1688669597.8277369))
	assert.Equal(t, trade.BuySell, "b")
	assert.Equal(t, trade.MarketLimit, "m")
	assert.Empty(t, trade.Misc)
	assert.Equal(t, trade.TradeID, 61044952)
}
