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

func PriceHandler(priceDataChan <-chan PriceData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var prices []float64
		// WebSocket接続を開始
		// http.ResponseWriterとhttp.Requestを受け取って、
		// WebSocket通信にアップグレードしてくれる。
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("failed to upgrade:", err)
			return
		}
		// WebSocket通信は、開いたら繋ぎっぱなしになるため、終了することを保証
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
				var ema float64
				// 価格データを追加
				ema, prices = CalcTechnical(prices, priceData.BestAsk)
				fmt.Println("\npriceData.BestAsk", priceData.BestAsk)
				chartData := ChartData{
					Price:     fmt.Sprintf("%.10f", priceData.BestAsk),
					EMA:       fmt.Sprintf("%.10f", ema),
					Timestamp: priceData.Timestamp,
				}
				if err := conn.WriteJSON(chartData); err != nil {
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
