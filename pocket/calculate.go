package pocket

import (
	"github.com/shopspring/decimal"
)

func Sub(x float64, y float64) float64 {
	a := decimal.NewFromFloat(x)
	b := decimal.NewFromFloat(y)
	sub := a.Sub(b)
	result, _ := sub.Float64()
	return result
}

func Add(x float64, y float64) float64 {
	a := decimal.NewFromFloat(x)
	b := decimal.NewFromFloat(y)
	sum := a.Add(b)
	result, _ := sum.Float64()
	return result
}
