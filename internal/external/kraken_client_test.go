package external

import (
	"bytes"
	"diveGames/internal/models"
	"diveGames/internal/trasnsformations"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
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

type ClientMock struct {
	mock.Mock
}

func (m ClientMock) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetLastTradePriceByPairOk(t *testing.T) {
	clientMock := ClientMock{}
	url := "someURL.com?pair=%s&size=%d"

	body := io.NopCloser(bytes.NewBufferString(jsonTestData))
	httpResponse := &http.Response{StatusCode: 200, Body: body}
	tradeUSD := models.TradePrice{
		Pair:   "BTC/USD",
		Amount: 30243.40000,
	}
	clientMock.On("Get", fmt.Sprintf(url, trasnsformations.MapPairsToKrakenKeys["BTC/USD"], 1)).Return(httpResponse, nil)

	ksi := KrakenServiceImpl{url, clientMock}

	trade, err := ksi.GetLastTradePriceByPair("BTC/USD")

	assert.Nil(t, err)
	assert.NotNil(t, trade)
	assert.Equal(t, tradeUSD, *trade)
}

func TestGetLastTradePriceByPairHttpClientError(t *testing.T) {
	clientMock := ClientMock{}
	url := "someURL.com?pair=%s&size=%d"

	clientMock.On("Get", mock.Anything).Return(nil, errors.New("some error"))

	ksi := KrakenServiceImpl{url, clientMock}

	trade, err := ksi.GetLastTradePriceByPair("BTC/USD")

	assert.NotNil(t, err)
	assert.Nil(t, trade)
	assert.Equal(t, "some error", err.Error())
}
