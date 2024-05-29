package test

import (
	"diveGames/internal/config"
	"diveGames/internal/external"
	"diveGames/internal/handlers"
	models2 "diveGames/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	jsonTestDataWithErrors = `{
        "error": ["Some error"],
        "result": {}
    }`
	jsonTestDataUSD = `{
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
	jsonTestDataEUR = `{
        "error": [],
        "result": {
            "XXBTZEUR": [
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
	jsonTestDataCHF = `{
        "error": [],
        "result": {
            "XBTCHF": [
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
	expectedResponse = "{\n  \"ltp\": [\n    {\n      \"pair\": \"BTC/USD\",\n      \"amount\": \"30243.40000\"\n    },\n    {\n      \"pair\": \"BTC/CHF\",\n      \"amount\": \"30243.40000\"\n    },\n    {\n      \"pair\": \"BTC/EUR\",\n      \"amount\": \"30243.40000\"\n    }\n  ]\n}"
)

func TestGetLastTradePricesOk(t *testing.T) {
	r := gin.Default()
	server := httptest.NewServer(r)
	defer server.Close()

	setTradePriceHandlerToServer(r, server.URL)
	setKrakenMockServer(r)

	resp, err := http.Get(server.URL + "/api/v1/ltp?pairs=BTC/USD&pairs=BTC/CHF&pairs=BTC/EUR")
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, unmarshalLastTradePrices([]byte(expectedResponse)), unmarshalLastTradePrices(body))
}

func TestGetLastTradePricesKrakenInternalError(t *testing.T) {
	r := gin.Default()
	server := httptest.NewServer(r)
	defer server.Close()

	setTradePriceHandlerToServer(r, server.URL)
	setKrakenMockServerWithInternalErrors(r)

	resp, err := http.Get(server.URL + "/api/v1/ltp?pairs=BTC/USD&pairs=BTC/EUR&pairs=BTC/CHF")
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"error\":\"We received the following errors from Kraken services \\nSome error\\n\",\"status\":500}", string(body))
}

func TestGetLastTradePricesKrakenNotFoundError(t *testing.T) {
	r := gin.Default()
	server := httptest.NewServer(r)
	defer server.Close()

	setTradePriceHandlerToServer(r, server.URL)
	setKrakenMockServerNotFoundError(r)

	resp, err := http.Get(server.URL + "/api/v1/ltp?pairs=BTC/USD&pairs=BTC/EUR&pairs=BTC/CHF")
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"error\":\"Kraken endpoint returned status code 404\",\"status\":500}", string(body))
}

func setTradePriceHandlerToServer(r *gin.Engine, url string) {
	httpClient, _ := config.NewHTTPClient("../config/config.yml")
	krakenService := external.NewKrakenServiceImpl(url+"/0/public/Trades?pair=%s&count=%d", httpClient)
	tph := handlers.NewTradePriceHandler(krakenService)
	r.GET("/api/v1/ltp", tph.FetchTradePriceByPairs)
}

func setKrakenMockServer(r *gin.Engine) {
	r.GET("/0/public/Trades", func(c *gin.Context) {
		pair := c.Query("pair")
		if pair == "XXBTZUSD" {
			c.JSON(http.StatusOK, unmarshalMockResponse(jsonTestDataUSD))
		} else if pair == "XXBTZEUR" {
			c.JSON(http.StatusOK, unmarshalMockResponse(jsonTestDataEUR))
		} else {
			c.JSON(http.StatusOK, unmarshalMockResponse(jsonTestDataCHF))
		}
	})
}

func setKrakenMockServerNotFoundError(r *gin.Engine) {
	r.GET("/0/public/Trades", func(c *gin.Context) {
		pair := c.Query("pair")
		if pair == "XXBTZUSD" {
			c.JSON(http.StatusNotFound, nil)
		} else if pair == "XXBTZEUR" {
			c.JSON(http.StatusOK, unmarshalMockResponse(jsonTestDataEUR))
		} else {
			c.JSON(http.StatusOK, unmarshalMockResponse(jsonTestDataCHF))
		}
	})
}

func setKrakenMockServerWithInternalErrors(r *gin.Engine) {
	r.GET("/0/public/Trades", func(c *gin.Context) {
		pair := c.Query("pair")
		if pair == "XXBTZUSD" {
			c.JSON(http.StatusOK, unmarshalMockResponse(jsonTestDataWithErrors))
		} else if pair == "XXBTZEUR" {
			c.JSON(http.StatusOK, unmarshalMockResponse(jsonTestDataEUR))
		} else {
			c.JSON(http.StatusOK, unmarshalMockResponse(jsonTestDataCHF))
		}
	})
}

func unmarshalMockResponse(s string) models2.KrakenTrade {
	var tradesResponse models2.KrakenTrade
	if err := json.Unmarshal([]byte(s), &tradesResponse); err != nil {
		panic(err)
	}
	return tradesResponse
}

func unmarshalLastTradePrices(b []byte) models2.LastTradePrices {
	var lastTradePrices models2.LastTradePrices
	if err := json.Unmarshal(b, &lastTradePrices); err != nil {
		panic(err)
	}
	return lastTradePrices
}
