package handlers

import (
	models2 "diveGames/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TradePriceHandler struct {
	tradeFetcherService TradeFetcherService
}

type TradeFetcherService interface {
	GetLastTradePriceByPair(pair string) (*models2.TradePrice, error)
}

func NewTradePriceHandler(tradeFetcherService TradeFetcherService) TradePriceHandler {
	return TradePriceHandler{tradeFetcherService: tradeFetcherService}
}

func (h *TradePriceHandler) FetchTradePriceByPairs(c *gin.Context) {
	err := h.validateFetchTradePriceByPairsParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		return
	}
	ltp := make([]*models2.TradePrice, len(c.QueryArray("pairs")))
	for i, pair := range c.QueryArray("pairs") {
		ltp[i], err = h.tradeFetcherService.GetLastTradePriceByPair(pair)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
			return
		}
	}
	c.JSON(http.StatusOK, models2.LastTradePrices{Trades: ltp})
}

func (h *TradePriceHandler) validateFetchTradePriceByPairsParams(c *gin.Context) error {
	pairsParam := c.QueryArray("pairs")
	if len(pairsParam) == 0 {
		return errors.New("'pairs' param is required")
	}
	return nil
}
