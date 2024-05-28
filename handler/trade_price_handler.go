package handler

import (
	"diveGames/domain"
	"diveGames/handler/handlerDTO"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TradePriceHandler struct {
	tradeFetcherService TradeFetcherService
	tradeHandlerMapper  TradeHandlerMapper
}

type TradeFetcherService interface {
	GetLastTradePrices() ([]*domain.Trade, error)
}

type TradeHandlerMapper interface {
	MapTradeToLastTradesPrice(trades []*domain.Trade) (*handlerDTO.LastTradePrices, error)
}

func NewTradePriceHandler(r *gin.Engine, tradeFetcherService TradeFetcherService, thm TradeHandlerMapper) {
	handler := &TradePriceHandler{tradeFetcherService: tradeFetcherService, tradeHandlerMapper: thm}
	r.GET("/api/v1/ltp", handler.FetchTradePrice)
}

func (h *TradePriceHandler) FetchTradePrice(c *gin.Context) {
	trades, err := h.tradeFetcherService.GetLastTradePrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ltp, err := h.tradeHandlerMapper.MapTradeToLastTradesPrice(trades)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error mapping response": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ltp)
}
