package apps

import (
	"errors"

	"arusu.info/btcat/conf"
)

// MACDを計算する
func MACD(macdValues, prices []float64) ([]float64, float64, float64, float64, float64, float64, error) {
	short, long, macdLine, err := macd_Line(prices)
	if err != nil {
		return macdValues, 0, 0, 0, 0, 0, err
	}

	// Signal Lineの計算
	macdValues = append(macdValues, macdLine)
	if len(macdValues) > conf.SignalEMAPeriod {
		macdValues = macdValues[1:]
	}

	signalLine := EMA(macdValues, conf.SignalEMAPeriod)

	if len(signalLine) < 1 {
		return macdValues, 0, 0, 0, 0, 0, errors.New("not enough data")
	}

	histogram := macdLine - signalLine[len(signalLine)-1]

	return macdValues, short, long, macdLine, signalLine[len(signalLine)-1], histogram, nil
}

func macd_Line(prices []float64) (float64, float64, float64, error) {
	short_EMA := EMA(prices, conf.ShortEMAPeriod)
	long_EMA := EMA(prices, conf.LongEMAPeriod)
	if len(short_EMA) < 1 || len(long_EMA) < 1 {
		return 0, 0, 0, errors.New("not enough data to calculate MACD Line")
	}
	return short_EMA[len(short_EMA)-1], long_EMA[len(long_EMA)-1], short_EMA[len(short_EMA)-1] - long_EMA[len(long_EMA)-1], nil
}
