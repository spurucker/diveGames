package external

import (
	"diveGames/internal/models"
	trasnsformations2 "diveGames/internal/trasnsformations"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const TradeSearchSize = 1

type KrakenServiceImpl struct {
	apiUrl     string
	httpClient Client
}

type Client interface {
	Get(url string) (*http.Response, error)
}

func NewKrakenServiceImpl(apiUrl string, httpClient Client) *KrakenServiceImpl {
	return &KrakenServiceImpl{apiUrl: apiUrl, httpClient: httpClient}
}

func (ksi *KrakenServiceImpl) GetLastTradePriceByPair(pair string) (*models.TradePrice, error) {
	krakenPair := trasnsformations2.MapPairsToKrakenKeys[pair]
	response, err := ksi.httpClient.Get(fmt.Sprintf(ksi.apiUrl, krakenPair, TradeSearchSize))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Kraken endpoint returned status code %d", response.StatusCode))
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("Error closing response body: %v\n", err)
		}
	}(response.Body)

	return trasnsformations2.TransformKrakenTrade(response.Body, krakenPair)
}
