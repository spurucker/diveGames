package main

import (
	"diveGames/handler"
	dto2 "diveGames/handler/handlerDTO"
	"diveGames/httpClient"
	"diveGames/repository"
	"diveGames/repository/RepositoryDTO"
	"diveGames/usercase"
	"github.com/gin-gonic/gin"
)

func InitializeServer(krakenURL string) *gin.Engine {
	r := gin.Default()
	httpConfig, err := httpClient.NewHTTPClientConfig("config/config.yml")
	if err != nil {
		panic(err)
	}
	client := httpClient.NewHTTPClient(httpConfig)
	krakenService := repository.NewKrakenServiceImpl(krakenURL, RepositoryDTO.NewTradeMapper(), client)
	fetchTradePriceUC := usercase.NewFetchTradePriceUserCase(krakenService)
	handler.NewTradePriceHandler(r, fetchTradePriceUC, dto2.NewTradePriceMapper())
	return r
}
