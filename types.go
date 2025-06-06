package main

type TradeData struct {
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Timestamp int64  `json:"t"`
}

type PublishTradeData struct {
	Symbol    string  `json:"s"`
	Price     float64 `json:"p"`
	Timestamp int64   `json:"t"`
}

type WebSocketMsg struct {
	Stream string    `json:"stream"`
	Data   TradeData `json:"data"`
}
