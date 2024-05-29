package usercase

import (
	"diveGames/domain"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type TradeFetcherMock struct {
	mock.Mock
}

func (m *TradeFetcherMock) GetLastTradePriceByPair(pair string) (*domain.Trade, error) {
	args := m.Called(pair)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Trade), args.Error(1)
}

func TestGetLastTradePricesOk(t *testing.T) {
	mock := TradeFetcherMock{}
	userCase := FetchTradePriceUserCase{tradeFetcher: &mock}

	tradeUSD := domain.Trade{
		Price: "1.00",
	}
	tradeEUR := domain.Trade{
		Price: "2.00",
	}
	tradeCHF := domain.Trade{
		Price: "3.00",
	}
	mock.On("GetLastTradePriceByPair", "BTC/USD").Return(&tradeUSD, nil)
	mock.On("GetLastTradePriceByPair", "BTC/EUR").Return(&tradeEUR, nil)
	mock.On("GetLastTradePriceByPair", "BTC/CHF").Return(&tradeCHF, nil)

	result, err := userCase.GetLastTradePricesByPairs([]string{"BTC/USD", "BTC/EUR", "BTC/CHF"})

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, tradeUSD, *result[0])
	assert.Equal(t, tradeEUR, *result[1])
	assert.Equal(t, tradeCHF, *result[2])
}

func TestGetLastTradePricesError(t *testing.T) {
	mock := TradeFetcherMock{}
	userCase := FetchTradePriceUserCase{tradeFetcher: &mock}
	mock.On("GetLastTradePriceByPair", "BTC/USD").Return(nil, errors.New("Some error"))

	trades, err := userCase.GetLastTradePricesByPairs([]string{"BTC/USD", "BTC/EUR", "BTC/CHF"})

	assert.NotNil(t, err)
	assert.Nil(t, trades)
	assert.Equal(t, "Some error", err.Error())
}
