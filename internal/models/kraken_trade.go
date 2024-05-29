package models

type KrakenTrade struct {
	Error  []string               `json:"error"`
	Result map[string]interface{} `json:"result"`
}
