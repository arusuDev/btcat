package apps

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 全てのオリジンを許可する
	},
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

func PriceHandler(priceDataChan <-chan PriceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// WebSocket接続を開始
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("failed to upgrade:", err)
			return
		}
		defer conn.Close()

		// 接続が閉じられたことを示すためのチャンネル
		closedChan := make(chan struct{})

		// 別のgoroutineでエラーチェックを行う
		go func() {
			for {
				if _, _, err := conn.NextReader(); err != nil {
					// WebSocket接続が閉じられている場合、そのことをclosedChanに通知する
					closedChan <- struct{}{}
					break
				}
			}
		}()

		for {
			select {
			case priceData := <-priceDataChan:
				// 価格データを取得
				fmt.Println(priceData.BestAsk)

				if err := conn.WriteJSON(priceData); err != nil {
					log.Println(err)
					return
				}
			case <-closedChan:
				// WebSocket接続が閉じられた場合、ループを抜ける
				return
			}
		}
	}
}
