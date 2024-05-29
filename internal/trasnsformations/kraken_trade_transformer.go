package trasnsformations

import (
	"bytes"
	models2 "diveGames/internal/models"
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

func TransformKrakenTrade(body io.ReadCloser, pair string) (*models2.TradePrice, error) {
	var tradesResponse models2.KrakenTrade
	if err := json.NewDecoder(body).Decode(&tradesResponse); err != nil {
		return nil, err
	}
	if tradesResponse.Error != nil && len(tradesResponse.Error) > 0 {
		return nil, errors.New(mergeErrorMessages(tradesResponse.Error))
	}
	return transformKrakenBody(tradesResponse.Result, pair)
}

func mergeErrorMessages(errors []string) string {
	var errorBuffer bytes.Buffer
	errorBuffer.WriteString("We received the following errors from Kraken services \n")
	for _, errMessage := range errors {
		errorBuffer.WriteString(errMessage + "\n")
	}
	return errorBuffer.String()
}
func transformKrakenBody(json map[string]interface{}, pair string) (*models2.TradePrice, error) {
	for key, value := range json {
		if key == pair {
			trades, ok := value.([]interface{})
			if !ok {
				return nil, errors.New("trade could not be mapped")
			}
			trade, ok := trades[len(trades)-1].([]interface{})
			if !ok {
				return nil, errors.New("trade could not be mapped")
			}
			priceString, ok := trade[0].(string)
			if !ok {
				return nil, errors.New("price could not be mapped")
			}
			price, err := strconv.ParseFloat(priceString, 64)
			if err != nil {
				return nil, errors.New("price could not be read as float64")
			}
			result := &models2.TradePrice{
				Pair:   MapKrakenKeysToPairs[pair],
				Amount: price,
			}
			return result, nil
		}
	}
	return nil, errors.New("no such key")
}
