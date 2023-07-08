package apps

import (
	"errors"

	"arusu.info/btcat/conf"
)

func MACD(prices []float64) {
	// TODO
	// short, long, macd, err := macd_Line(prices)
	// if err != nil {
	// 	panic(err)
	// }

	// signal_EMA := EMA(prices,conf.SignalEMA)

}

func macd_Line(prices []float64) (float64, float64, float64, error) {
	short_EMA := EMA(prices, conf.ShortEMAPeriod)
	long_EMA := EMA(prices, conf.LongEMAPeriod)
	if len(short_EMA) < 1 || len(long_EMA) < 1 {
		return 0, 0, 0, errors.New("not enough data to calculate MACD Line")
	}
	return short_EMA[len(short_EMA)-1], long_EMA[len(long_EMA)-1], short_EMA[len(short_EMA)-1] - long_EMA[len(long_EMA)-1], nil
}
