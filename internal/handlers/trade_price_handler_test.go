package handlers

import (
	models2 "diveGames/internal/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) GetLastTradePriceByPair(pair string) (*models2.TradePrice, error) {
	args := m.Called(pair)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models2.TradePrice), args.Error(1)
}

func TestMapTradeResponseOk(t *testing.T) {
	serviceMock := &ServiceMock{}

	handler := TradePriceHandler{tradeFetcherService: serviceMock}

	usdTrade := &models2.TradePrice{Pair: "BTC/USD", Amount: "1"}
	serviceMock.On("GetLastTradePriceByPair", "BTC/USD").Return(usdTrade, nil)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request, _ = http.NewRequest("GET", "/api/v1/ltp?pairs=BTC/USD", nil)

	handler.FetchTradePriceByPairs(context)

	var result models2.LastTradePrices
	json.NewDecoder(recorder.Body).Decode(&result)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, models2.LastTradePrices{[]*models2.TradePrice{usdTrade}}, result)
}

func TestMapTradeResponseErrorInService(t *testing.T) {
	serviceMock := &ServiceMock{}

	handler := TradePriceHandler{tradeFetcherService: serviceMock}
	serviceMock.On("GetLastTradePriceByPair", "BTC/USD").Return(nil, errors.New("some error"))
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request, _ = http.NewRequest("GET", "/api/v1/ltp?pairs=BTC/USD", nil)
	handler.FetchTradePriceByPairs(context)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "{\"error\":\"some error\",\"status\":500}", recorder.Body.String())
}

func TestMapTradeResponseNoneQueryParams(t *testing.T) {
	serviceMock := &ServiceMock{}

	handler := TradePriceHandler{tradeFetcherService: serviceMock}

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request, _ = http.NewRequest("GET", "/api/v1/ltp", nil)
	handler.FetchTradePriceByPairs(context)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "{\"error\":\"'pairs' param is required\",\"status\":400}", recorder.Body.String())
}
