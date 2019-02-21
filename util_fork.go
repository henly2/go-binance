package binance

import (
	"fmt"
	"strconv"
	"math/big"
)

func FloatFromStringByAutoDecimals(str string) (float64, error) {
	return StringToFloat64ByDecimals(str, GetDecimalsByStr(str))
}

func GetDecimalsByStr(str string) int {
	/*
	2.345600
	-> cntZero = 2
	-> dotIndex = 1
	=> decimal = 8 - 1 - 1 - 2 = 4

	2.
	-> cntZero = 0
	-> dotIndex = 1
	=> decimal = 2 -1 - 1 - 0
	 */

	var (
		needZero = true

		cntZero = 0
		dotIndex = -1
	)

	ss := []byte(str)
	for i := len(ss) - 1; i >= 0; i-- {
		if ss[i] == '0' {
			if needZero {
				cntZero++
			}
		} else {
			needZero = false
			if ss[i] == '.' {
				dotIndex = i
				break
			}
		}
	}

	if dotIndex != -1 {
		return len(str) - 1 - dotIndex - cntZero
	}

	return 0
}

func StringToFloat64ByDecimals(am string, decimals int) (float64, error) {
	r := &big.Rat{}
	if _, ok := r.SetString(am); !ok {
		return 0, fmt.Errorf("cannot parse amount: %s", am)
	}

	f, err := strconv.ParseFloat(r.FloatString(decimals), 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

func StringFromFloatByAutoDecimals(v float64) (string) {
	tmp := Float64ToStringByDecimals(v, 10)

	return Float64ToStringByDecimals(v, GetDecimalsByStr(tmp))
}

func Float64ToStringByDecimals(v float64, decimals int) (string) {
	r := &big.Rat{}
	r.SetFloat64(v)

	return r.FloatString(decimals)
}