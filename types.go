package main

type TradeData struct {
	Symbol    string  `json:"s"`
	Price     float64 `json:"p"`
	Timestamp int64   `json:"t"`
}

type WebSocketMsg struct {
	Type string      `json:"type"`
	Data []TradeData `json:"data"`
}
