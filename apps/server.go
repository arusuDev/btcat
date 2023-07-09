package apps

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 全てのオリジンを許可する
	},
}

var macd_Flag bool
var macd_sum float64
var buy float64
var sel float64

func PriceHandler(priceDataChan <-chan PriceData) http.HandlerFunc {
	macd_Flag = false
	return func(w http.ResponseWriter, r *http.Request) {
		var prices []float64
		var macdValues []float64
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
				var tech Technical
				// 価格データを追加
				tech, prices, macdValues = CalcTechnical(prices, macdValues, priceData.BestAsk)
				// fmt.Println("\npriceData.BestAsk", priceData.BestAsk)
				chartData := ChartData{
					Price:     fmt.Sprintf("%.10f", priceData.BestAsk),
					Technical: tech,
					Timestamp: priceData.Timestamp,
				}
				if err := conn.WriteJSON(chartData); err != nil {
					log.Println(err)
					return
				}

				// 確認用
				HistogramFloat, err := strconv.ParseFloat(tech.Histogram, 64)
				if err != nil {
					// エラーハンドリング
					log.Printf("Failed to parse Histogram value: %v\n", err)
					return
				}
				if HistogramFloat > 0 && !macd_Flag {
					fmt.Printf("%s  -- MACD購入サイン：%f\n", priceData.Timestamp, priceData.BestAsk)
					buy = priceData.BestAsk
					macd_Flag = true
				} else if HistogramFloat < 0 && macd_Flag {
					fmt.Printf("%s  -- MACD売却サイン：%f\n", priceData.Timestamp, priceData.BestAsk)
					sel = priceData.BestAsk
					fmt.Printf("今回の購入価格：%f 売却価格：%f 差額:%f\n", buy, sel, sel-buy)
					macd_sum += sel - buy
					fmt.Printf("MACD差額：%f\n", macd_sum)
					macd_Flag = false
				}

			case <-closedChan:
				// WebSocket接続が閉じられた場合、ループを抜ける
				return
			}
		}
	}
}
