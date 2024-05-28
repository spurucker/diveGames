package usercase

import (
	"diveGames"
	"diveGames/domain"
)

type TradeFetcher interface {
	GetLastTradePriceByPair(pair string) (*domain.Trade, error)
}

type FetchTradePriceUserCase struct {
	tradeFetcher TradeFetcher
}

func NewFetchTradePriceUserCase(tradeFetcher TradeFetcher) *FetchTradePriceUserCase {
	return &FetchTradePriceUserCase{tradeFetcher: tradeFetcher}
}

func (u *FetchTradePriceUserCase) GetLastTradePrices() ([]*domain.Trade, error) {
	var err error
	trades := make([]*domain.Trade, len(diveGames.PairValues))
	for i, pair := range diveGames.PairValues {
		trades[i], err = u.tradeFetcher.GetLastTradePriceByPair(pair)
		if err != nil {
			return nil, err
		}
	}
	return trades, nil
}
