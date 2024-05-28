package repository

import (
	"diveGames/domain"
	domainErrors "diveGames/error"
	"diveGames/repository/RepositoryDTO"
	"fmt"
	"io"
	"net/http"
)

const TradeSearchSize = 1

type KrakenServiceImpl struct {
	apiUrl      string
	tradeMapper Mapper
	httpClient  Client
}

type Client interface {
	Get(url string) (*http.Response, error)
}

type Mapper interface {
	MapResponseToTrade(body io.ReadCloser, pair string) (*domain.Trade, error)
}

func NewKrakenServiceImpl(apiUrl string, tm Mapper, httpClient Client) *KrakenServiceImpl {
	return &KrakenServiceImpl{apiUrl: apiUrl, tradeMapper: tm, httpClient: httpClient}
}

func (ksi *KrakenServiceImpl) GetLastTradePriceByPair(pair string) (*domain.Trade, error) {
	krakenPair := RepositoryDTO.MapPairsToKrakenKeys[pair]
	response, err := ksi.httpClient.Get(fmt.Sprintf(ksi.apiUrl, krakenPair, TradeSearchSize))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, domainErrors.NewDependencyError(fmt.Sprintf("Kraken endpoint returned status code %d", response.StatusCode))
	}
	defer response.Body.Close()

	return ksi.tradeMapper.MapResponseToTrade(response.Body, krakenPair)
}
