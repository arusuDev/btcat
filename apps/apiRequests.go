package apps

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketのエンドポイント
// https://bf-lightning-api.readme.io/docs/endpoint-json-rpc
const bitFlyerRealtimeEndpoint = "wss://ws.lightstream.bitflyer.com/json-rpc"

func RealtimeTicker() {
	// 外部からのSignalを受け取るためのチャネルを作成
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(bitFlyerRealtimeEndpoint, nil)
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
				log.Println("read message error : %s\n", err)
				return
			}
			log.Printf("message : %s\n", message)
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
				log.Println("Websocket disconnection error:%s", err)
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
