package RepositoryDTO

type TradesResponse struct {
	Error  []string               `json:"error"`
	Result map[string]interface{} `json:"result"`
}
