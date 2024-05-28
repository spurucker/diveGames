package handlerDTO

import (
	"diveGames/domain"
)

type TradePriceMapper struct{}

func NewTradePriceMapper() TradePriceMapper {
	return TradePriceMapper{}
}

func (t TradePriceMapper) MapTradeToLastTradesPrice(trades []*domain.Trade) (*LastTradePrices, error) {
	arrayTradePrice := make([]TradePrice, len(trades))
	for i, trade := range trades {
		tp, _ := t.mapTradeToTradePrice(trade)
		arrayTradePrice[i] = *tp
	}
	return &LastTradePrices{Trades: arrayTradePrice}, nil
}

func (t TradePriceMapper) mapTradeToTradePrice(trade *domain.Trade) (*TradePrice, error) {
	if trade == nil {
		return &TradePrice{}, nil
	}
	return &TradePrice{
		Pair:   trade.Pair,
		Amount: trade.Price,
	}, nil
}
