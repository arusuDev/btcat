package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	LastTradePrice  float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}
type Allowance struct {
	Cost      float64 `json:"cost"`
	Remaining float64 `json:"remaining"`
}

type OHLCResponse struct {
	Result    map[string][][]float64 `json:"result"`
	Allowance Allowance              `json:"allowance"`
}

type Candlestick struct {
	Time   time.Time `json:"time"`
	Open   float64   `json:"open"`
	Close  float64   `json:"close"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Volume float64   `json:"volume"`
}

var clients = make(map[*websocket.Conn]bool)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade http to websocket: %v", err)
	}
	clients[conn] = true
}

func broadcastPrice(price float64) {
	for client := range clients {
		err := client.WriteJSON(price)
		if err != nil {
			log.Printf("Error occurred while writing message to client: %v", err)
			client.Close()
			delete(clients, client)
		} else {
			log.Printf("Broadcasted price: %f", price)
		}
	}
}
func broadcastCandlestick(candlestick Candlestick) {
	for client := range clients {
		err := client.WriteJSON(candlestick)
		if err != nil {
			log.Printf("Error occurred while writing message to client: %v", err)
			client.Close()
			delete(clients, client)
		} else {
			log.Printf("Broadcasted candlestick: %+v", candlestick)
		}
	}
}
func getBitcoinPriceAndBroadcast() {
	resp, err := http.Get("https://api.bitflyer.com/v1/ticker?product_code=BTC_JPY")
	if err != nil {
		log.Fatalf("Failed to get BTC price: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var ticker Ticker
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	log.Printf("BTC Price: %f", ticker.LastTradePrice)
	broadcastPrice(ticker.LastTradePrice)
}
func getBitcoinCandlestickAndBroadcast() {
	resp, err := http.Get("https://api.cryptowat.ch/markets/bitflyer/btcjpy/ohlc?periods=60")
	if err != nil {
		log.Fatalf("Failed to get BTC candlestick: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var ohlcRes OHLCResponse
	err = json.Unmarshal(body, &ohlcRes)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	candlestickData := ohlcRes.Result["60"][0] // Get the latest minute candlestick data

	candlestick := Candlestick{
		Time:   time.Unix(int64(candlestickData[0]), 0),
		Open:   candlestickData[1],
		Close:  candlestickData[4],
		High:   candlestickData[2],
		Low:    candlestickData[3],
		Volume: candlestickData[5],
	}

	log.Printf("BTC Candlestick: %+v", candlestick)
	broadcastCandlestick(candlestick)
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	go func() {
		for range time.Tick(30 * time.Second) {
			// getBitcoinPriceAndBroadcast()
			getBitcoinCandlestickAndBroadcast()
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
