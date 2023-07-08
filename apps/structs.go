package apps

type ChartData struct {
	Price     string `json:"price_data"`
	EMA       string `json:"ema"`
	Timestamp string `json:"timestamp"`
}

type PriceData struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	State           string  `json:"state"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	LTP             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

type Message struct {
	Message PriceData `json:"message"`
}

type Params struct {
	Channel string    `json:"channel"`
	Message PriceData `json:"message"`
}

type Data struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}
