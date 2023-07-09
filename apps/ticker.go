package apps

import (
	"fmt"
	"log"

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
func CalcTechnical(prices, macdValues []float64, newPrice float64) (Technical, []float64, []float64) {
	prices = addPriceData(prices, newPrice, conf.Period_tick)

	macdValues, short, long, macd, signal, hist, err := MACD(macdValues, prices)
	if err != nil {
		log.Println(err)
	}
	technical := Technical{
		ShortEMA:  fmt.Sprintf("%.10f", short),
		LongEMA:   fmt.Sprintf("%.10f", long),
		MACD:      fmt.Sprintf("%.10f", macd),
		SignalEMA: fmt.Sprintf("%.10f", signal),
		Histogram: fmt.Sprintf("%.10f", hist),
	}
	return technical, prices, macdValues
}
