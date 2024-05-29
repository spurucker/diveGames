package app

import (
	"diveGames/internal/config"
	"diveGames/internal/external"
	handler "diveGames/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewApp(krakenURL string) *gin.Engine {
	r := gin.Default()
	httpClient, err := config.NewHTTPClient("config/config.yml")
	if err != nil {
		panic(err)
	}
	krakenService := external.NewKrakenServiceImpl(krakenURL, httpClient)
	tph := handler.NewTradePriceHandler(krakenService)

	r.GET("/api/v1/ltp", tph.FetchTradePriceByPairs)
	return r
}
