package handler

import (
	"diveGames/domain"
	"diveGames/handler/handlerDTO"
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

func (m *ServiceMock) GetLastTradePrices() ([]*domain.Trade, error) {
	args := m.Called()
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Trade), args.Error(1)
}

type MapperMock struct {
	mock.Mock
}

func (m *MapperMock) MapTradeToLastTradesPrice(trades []*domain.Trade) (*handlerDTO.LastTradePrices, error) {
	args := m.Called(trades)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*handlerDTO.LastTradePrices), args.Error(1)
}

func TestMapTradeResponseOk(t *testing.T) {
	serviceMock := &ServiceMock{}
	mapperMock := &MapperMock{}
	handler := TradePriceHandler{tradeFetcherService: serviceMock, tradeHandlerMapper: mapperMock}

	trades := make([]*domain.Trade, 1)
	mockResponse := handlerDTO.LastTradePrices{
		Trades: make([]handlerDTO.TradePrice, 0),
	}
	serviceMock.On("GetLastTradePrices").Return(trades, nil)
	mapperMock.On("MapTradeToLastTradesPrice", trades).Return(&mockResponse, nil)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	handler.FetchTradePrice(context)

	var result handlerDTO.LastTradePrices
	json.NewDecoder(recorder.Body).Decode(&result)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, mockResponse, result)
}

func TestMapTradeResponseErrorInService(t *testing.T) {
	serviceMock := &ServiceMock{}
	mapperMock := &MapperMock{}
	handler := TradePriceHandler{tradeFetcherService: serviceMock, tradeHandlerMapper: mapperMock}

	serviceMock.On("GetLastTradePrices").Return(nil, errors.New("some error"))
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	handler.FetchTradePrice(context)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "{\"error\":\"some error\"}", recorder.Body.String())
}

func TestMapTradeResponseErrorInMapper(t *testing.T) {
	serviceMock := &ServiceMock{}
	mapperMock := &MapperMock{}
	handler := TradePriceHandler{tradeFetcherService: serviceMock, tradeHandlerMapper: mapperMock}

	trades := make([]*domain.Trade, 1)
	serviceMock.On("GetLastTradePrices").Return(trades, nil)
	mapperMock.On("MapTradeToLastTradesPrice", trades).Return(nil, errors.New("some error"))
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	handler.FetchTradePrice(context)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "{\"error mapping response\":\"some error\"}", recorder.Body.String())
}
