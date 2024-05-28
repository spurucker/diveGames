package repository

import (
	"diveGames"
	"diveGames/domain"
	"diveGames/repository/RepositoryDTO"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"testing"
)

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

type MapperMock struct {
	mock.Mock
}

func (m MapperMock) MapResponseToTrade(body io.ReadCloser, pair string) (*domain.Trade, error) {
	args := m.Called(body, pair)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Trade), args.Error(1)
}

type EmptyReadCloser struct{}

func (erc EmptyReadCloser) Read(p []byte) (int, error) {
	return 0, nil
}

func (erc EmptyReadCloser) Close() error {
	return nil
}

func TestGetLastTradePriceByPairOk(t *testing.T) {
	clientMock := ClientMock{}
	mapperMock := MapperMock{}
	url := "someURL.com?pair=%s&size=%d"

	body := EmptyReadCloser{}
	httpResponse := &http.Response{StatusCode: 200, Body: &body}
	tradeUSD := domain.Trade{
		Price: "1.00",
	}
	clientMock.On("Get", fmt.Sprintf(url, RepositoryDTO.MapPairsToKrakenKeys[diveGames.PairUSD], 1)).Return(httpResponse, nil)
	mapperMock.On("MapResponseToTrade", mock.Anything, RepositoryDTO.MapPairsToKrakenKeys[diveGames.PairUSD]).Return(&tradeUSD, nil)

	ksi := KrakenServiceImpl{url, mapperMock, clientMock}

	trade, err := ksi.GetLastTradePriceByPair(diveGames.PairUSD)

	assert.Nil(t, err)
	assert.NotNil(t, trade)
	assert.Equal(t, tradeUSD, *trade)
}

func TestGetLastTradePriceByPairHttpClientError(t *testing.T) {
	clientMock := ClientMock{}
	mapperMock := MapperMock{}
	url := "someURL.com?pair=%s&size=%d"

	clientMock.On("Get", mock.Anything).Return(nil, errors.New("some error"))

	ksi := KrakenServiceImpl{url, mapperMock, clientMock}

	trade, err := ksi.GetLastTradePriceByPair(diveGames.PairUSD)

	assert.NotNil(t, err)
	assert.Nil(t, trade)
	assert.Equal(t, "some error", err.Error())
}

func TestGetLastTradePriceByPairMapperError(t *testing.T) {
	clientMock := ClientMock{}
	mapperMock := MapperMock{}
	url := "someURL.com?pair=%s&size=%d"

	body := EmptyReadCloser{}
	httpResponse := &http.Response{StatusCode: 200, Body: &body}
	clientMock.On("Get", fmt.Sprintf(url, RepositoryDTO.MapPairsToKrakenKeys[diveGames.PairUSD], 1)).Return(httpResponse, nil)
	mapperMock.On("MapResponseToTrade", mock.Anything, RepositoryDTO.MapPairsToKrakenKeys[diveGames.PairUSD]).Return(nil, errors.New("some error"))

	ksi := KrakenServiceImpl{url, mapperMock, clientMock}

	trade, err := ksi.GetLastTradePriceByPair(diveGames.PairUSD)

	assert.NotNil(t, err)
	assert.Nil(t, trade)
	assert.Equal(t, "some error", err.Error())
}
