package apps

// EMA(Exponential Moving Average)
// 時系列データから滑らかな線を生成するための一般的な手法
// 最新の価格動向に重きを置くため、急速な価格変動に素早く反応するのが特徴

// EMA = (Close_price * Alpha) + (EMA_before * ( 1 - Alpha ))
func EMA(data []float64, period int) []float64 {
	smoothing := 2.0
	ema := make([]float64, len(data))
	multiplier := smoothing / (1.0 + float64(period))
	// fmt.Printf("price len: %d\n", len(data))
	for i, price := range data {
		if i == 0 {
			ema[i] = price
			continue
		}
		ema[i] = price*multiplier + ema[i-1]*(1-multiplier)
	}

	return ema
}
