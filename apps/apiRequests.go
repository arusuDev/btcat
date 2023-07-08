package apps

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"arusu.info/btcat/conf"
	"github.com/gorilla/websocket"
)

// WebSocketのエンドポイント
// https://bf-lightning-api.readme.io/docs/endpoint-json-rpc

func RealtimeTicker(priceDataChan chan<- PriceData) {
	// 外部からのSignalを受け取るためのチャネルを作成
	// バッファサイズ1 →ひとつのみ値を保持し、吐き出されるまで入力をブロックする
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(conf.BitFlyerRealtimeEndpoint, nil)
	if err != nil {
		log.Fatalf("WebSocket connection error: %v", err)
	}
	defer c.Close()

	// リクエストデータの作成
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"params": map[string]string{
			"channel": "lightning_ticker_BTC_JPY",
		},
	}

	// WebSocketへリクエストデータを送信
	if err := c.WriteJSON(request); err != nil {
		log.Fatalf("request send error: %v", err)
	}

	done := make(chan struct{})

	//  goroutineを呼び出し　チャネルからメッセージを読み取りログに出力する
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("read message error : %s\n", err)
				done <- struct{}{}
				return
			}

			// 受信したメッセージをPriceData型に変換
			var priceData PriceData
			var data Data
			if err := json.Unmarshal(message, &data); err != nil {
				log.Printf("Error unmarshalling params: %v", err)
				continue
			}

			priceData = data.Params.Message
			// チャネルにpriceDataを送信
			priceDataChan <- priceData
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("detected intterupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("Websocket disconnection error:%s", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
