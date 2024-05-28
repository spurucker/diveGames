package RepositoryDTO

import (
	"bytes"
	"diveGames/domain"
	domainErrors "diveGames/error"
	"encoding/json"
	"io"
	"time"
)

type TradeMapper struct{}

func NewTradeMapper() TradeMapper {
	return TradeMapper{}
}

func (tm TradeMapper) MapResponseToTrade(body io.ReadCloser, pair string) (*domain.Trade, error) {
	var tradesResponse TradesResponse
	if err := json.NewDecoder(body).Decode(&tradesResponse); err != nil {
		return nil, err
	}
	if tradesResponse.Error != nil && len(tradesResponse.Error) > 0 {
		return nil, domainErrors.NewDependencyError(tm.mergeErrorMessages(tradesResponse.Error))
	}
	return tm.mapBodyToTrade(tradesResponse.Result, pair)
}

func (tm *TradeMapper) mergeErrorMessages(errors []string) string {
	var errorBuffer bytes.Buffer
	errorBuffer.WriteString("We received the following errors from Kraken services \n")
	for _, errMessage := range errors {
		errorBuffer.WriteString(errMessage + "\n")
	}
	return errorBuffer.String()
}
func (tm TradeMapper) mapBodyToTrade(json map[string]interface{}, pair string) (*domain.Trade, error) {
	for key, value := range json {
		if key == pair {
			trades, ok := value.([]interface{})
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Trade could not be mapped")
			}
			trade, ok := trades[len(trades)-1].([]interface{})
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Trade could not be mapped")
			}

			price, ok := trade[0].(string)
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Price could not be mapped")
			}
			vol, ok := trade[1].(string)
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Volume could not be mapped")
			}
			timestamp, ok := trade[2].(float64)
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Timestamp could not be mapped")
			}
			buySell, ok := trade[3].(string)
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Buy/Sell could not be mapped")
			}
			marketLimit, ok := trade[4].(string)
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Market/Limit could not be mapped")
			}
			misc, ok := trade[5].(string)
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Miscellaneus could not be mapped")
			}
			tradeId, ok := trade[6].(float64)
			if !ok {
				return nil, domainErrors.NewInvalidTradeError("Trade ID could not be mapped")
			}
			result := &domain.Trade{
				Pair:        MapKrakenKeysToPairs[pair],
				Price:       price,
				Volume:      vol,
				Time:        tm.timestampToTime(timestamp),
				BuySell:     buySell,
				MarketLimit: marketLimit,
				Misc:        misc,
				TradeID:     int(tradeId),
			}
			return result, nil
		}
	}
	return nil, domainErrors.NewInvalidTradeError("No such key")
}

func (TradeMapper) timestampToTime(timestamp float64) time.Time {
	seconds := int64(timestamp)
	nanoseconds := int64((timestamp - float64(seconds)) * 1e9)
	return time.Unix(seconds, nanoseconds)
}
