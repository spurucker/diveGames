package handler

import (
	"diveGames/domain"
	"diveGames/handler/handlerDTO"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TradePriceHandler struct {
	tradeFetcherService TradeFetcherService
	tradeHandlerMapper  TradeHandlerMapper
}

type TradeFetcherService interface {
	GetLastTradePricesByPairs(pairs []string) ([]*domain.Trade, error)
}

type TradeHandlerMapper interface {
	MapTradeToLastTradesPrice(trades []*domain.Trade) (*handlerDTO.LastTradePrices, error)
}

func NewTradePriceHandler(r *gin.Engine, tradeFetcherService TradeFetcherService, thm TradeHandlerMapper) {
	handler := &TradePriceHandler{tradeFetcherService: tradeFetcherService, tradeHandlerMapper: thm}
	r.GET("/api/v1/ltp", handler.FetchTradePriceByPairs)
}

func (h *TradePriceHandler) FetchTradePriceByPairs(c *gin.Context) {
	err := h.validateFetchTradePriceByPairsParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		return
	}
	var pairs []string
	if c.Query("pair") != "" {
		pairs = []string{c.Query("pair")}
	} else {
		pairs = c.QueryArray("pairs")
	}
	trades, err := h.tradeFetcherService.GetLastTradePricesByPairs(pairs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	ltp, err := h.tradeHandlerMapper.MapTradeToLastTradesPrice(trades)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error mapping response": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, ltp)
}

func (h *TradePriceHandler) validateFetchTradePriceByPairsParams(c *gin.Context) error {
	pairParam := c.Query("pair")
	pairsParam := c.QueryArray("pairs")
	if pairParam == "" && len(pairsParam) == 0 {
		return errors.New("missing pair parameter. One (and only one) 'pair' or 'pairs' query parameter is required")
	}
	if pairParam != "" && len(pairsParam) > 0 {
		return errors.New("one (and only one) 'pair' or 'pairs' query parameter is required")
	}
	return nil
}
