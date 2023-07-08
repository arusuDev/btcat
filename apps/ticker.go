package apps

import (
	"fmt"

	"arusu.info/btcat/conf"
)

func addPriceData(prices []float64, newPrice float64, period int) []float64 {
	prices = append(prices, newPrice)
	if len(prices) > period {
		// 0番目の要素を削除
		prices = prices[1:]
	}
	return prices
}

// 価格を保持してテクニカル分析をする
func CalcTechnical(prices []float64, newPrice float64) []float64 {
	prices = addPriceData(prices, newPrice, conf.Period_tick)
	EMAs := EMA(prices, conf.TechnicalPeriod)
	if len(EMAs) > 1 {
		fmt.Printf("EMA:%f", EMAs[len(EMAs)-1])
	}
	return prices

}
