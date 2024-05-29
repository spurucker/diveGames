package usercase

import (
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

func (u *FetchTradePriceUserCase) GetLastTradePricesByPairs(pairs []string) ([]*domain.Trade, error) {
	var err error
	trades := make([]*domain.Trade, len(pairs))
	for i, pair := range pairs {
		trades[i], err = u.tradeFetcher.GetLastTradePriceByPair(pair)
		if err != nil {
			return nil, err
		}
	}
	return trades, nil
}
