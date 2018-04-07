package decimal

import (
	"github.com/shopspring/decimal"
)

type PriceDecimal decimal.Decimal

func (p PriceDecimal) Fen() int64 {
	return decimal.Decimal(p).Round(0).IntPart()
}

func (p PriceDecimal) Yuan() float64 {
	res, _ := decimal.Decimal(p).Round(2).Float64()
	return res
}

func (p PriceDecimal) Decimal() decimal.Decimal {
	return decimal.Decimal(p)
}

func ToDecimal(num interface{}) decimal.Decimal {
	switch t := num.(type) {
	case decimal.Decimal:
		return t
	case PriceDecimal:
		return decimal.Decimal(t)
	case int:
		return decimal.New(int64(t), 0)
	case int8:
		return decimal.New(int64(t), 0)
	case int16:
		return decimal.New(int64(t), 0)
	case int32:
		return decimal.New(int64(t), 0)
	case int64:
		return decimal.New(t, 0)
	case uint:
		return decimal.New(int64(t), 0)
	case uint8:
		return decimal.New(int64(t), 0)
	case uint16:
		return decimal.New(int64(t), 0)
	case uint32:
		return decimal.New(int64(t), 0)
	case uint64:
		return decimal.New(int64(t), 0)
	case float32:
		return decimal.NewFromFloat(float64(t))
	case float64:
		return decimal.NewFromFloat(t)
	case string:
		n, err := decimal.NewFromString(t)
		if err != nil {
			panic(err)
		}

		return n
	default:
		panic("not support type")
	}
}

func Div(nums ...interface{}) PriceDecimal {
	if len(nums) <= 0 {
		return PriceDecimal(ToDecimal(0))
	}

	var res decimal.Decimal

	for i, num := range nums {
		if i == 0 {
			res = ToDecimal(num)
		} else {
			res = res.Div(ToDecimal(num))
		}
	}

	return PriceDecimal(res)
}

func Mul(nums ...interface{}) PriceDecimal {
	if len(nums) <= 0 {
		return PriceDecimal(ToDecimal(0))
	}

	var res decimal.Decimal

	for i, num := range nums {
		if i == 0 {
			res = ToDecimal(num)
		} else {
			res = res.Mul(ToDecimal(num))
		}
	}

	return PriceDecimal(res)
}

func Add(nums ...interface{}) PriceDecimal {
	var res decimal.Decimal

	for _, num := range nums {
		res = res.Add(ToDecimal(num))
	}

	return PriceDecimal(res)
}

func Sub(nums ...interface{}) PriceDecimal {
	if len(nums) <= 0 {
		return PriceDecimal(ToDecimal(0))
	}

	var res decimal.Decimal

	for i, num := range nums {
		if i == 0 {
			res = ToDecimal(num)
		} else {
			res = res.Sub(ToDecimal(num))
		}
	}

	return PriceDecimal(res)
}
