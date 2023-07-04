package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"arusu.info/btcat/apps"
)

func main() {
	// リアルタイムで価格を取得するためのチャネルを作成する
	priceDataChan := make(chan apps.PriceData)
	// Goroutineで価格情報を取得する
	go apps.RealtimeTicker(priceDataChan)

	// サーバの起動
	http.HandleFunc("/price", apps.PriceHandler(priceDataChan))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
