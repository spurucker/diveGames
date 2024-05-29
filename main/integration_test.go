package main

import (
	"diveGames/handler"
	"diveGames/handler/handlerDTO"
	"diveGames/httpClient"
	"diveGames/repository"
	"diveGames/repository/RepositoryDTO"
	"diveGames/usercase"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	resp, err := http.Get(getUrl(server))
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

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

	resp, err := http.Get(getUrl(server))
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

	resp, err := http.Get(getUrl(server))
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.NotNil(t, body)

	assert.Equal(t, "{\"error\":\"Kraken endpoint returned status code 404\",\"status\":500}", string(body))
}

func setTradePriceHandlerToServer(r *gin.Engine, url string) {
	httpConfig, _ := httpClient.NewHTTPClientConfig("../config/config.yml")

	krakenService := repository.NewKrakenServiceImpl(url+"/0/public/Trades?pair=%s&count=%d", RepositoryDTO.NewTradeMapper(), httpClient.NewHTTPClient(httpConfig))
	fetchTradePriceUC := usercase.NewFetchTradePriceUserCase(krakenService)
	handler.NewTradePriceHandler(r, fetchTradePriceUC, handlerDTO.NewTradePriceMapper())
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

func unmarshalMockResponse(s string) RepositoryDTO.TradesResponse {
	var tradesResponse RepositoryDTO.TradesResponse
	if err := json.Unmarshal([]byte(s), &tradesResponse); err != nil {
		panic(err)
	}
	return tradesResponse
}

func unmarshalLastTradePrices(b []byte) handlerDTO.LastTradePrices {
	var lastTradePrices handlerDTO.LastTradePrices
	if err := json.Unmarshal(b, &lastTradePrices); err != nil {
		panic(err)
	}
	return lastTradePrices
}

func getUrl(server *httptest.Server) string {
	params := map[string][]string{
		"pairs": {"BTC/USD", "BTC/CHF", "BTC/EUR"},
	}
	u, err := url.Parse(server.URL + "/api/v1/ltp")
	q := u.Query()
	for key, values := range params {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	u.RawQuery = q.Encode()
	if err != nil {
		fmt.Println("Error creating url:", err)
		panic(err)
	}
	return u.String()
}
