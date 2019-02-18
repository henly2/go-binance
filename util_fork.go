package binance

import (
	"fmt"
	"strconv"
	"github.com/pkg/errors"
	"math/big"
	"strings"
)

func FloatFromStringByAutoDecimals(raw interface{}) (float64, error) {
	str, ok := raw.(string)
	if !ok {
		return 0, errors.New(fmt.Sprintf("unable to parse, value not string: %T", raw))
	}

	return StringToFloat64ByDecimals(str, GetDecimalsByStr(str))
}

func GetDecimalsByStr(str string) int {
	index := strings.Index(str, ".")
	if index == -1 {
		return 0
	}

	return len(str) - index
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