package types

import (
	"math/big"
	"strconv"
)

const (
	UndDenom  = "und"  // 1 (base unit)
	PundDenom = "pund" // 10^-12 (pico)

	UndPow  = 1e12  // multiplier for converting from und to pund
	PundPow = 1e-12 // multiplier for converting from pund to und
)

func ConvertUndDenomination(amount string, from string, to string) (string, error) {

	if from == to {
		return amount + from, nil
	}

	switch from {
	case UndDenom: // from und to pund
		fromAmt, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return "", err
		}
		fromAmtBf := new(big.Float).SetFloat64(fromAmt)
		res := fromAmtBf.Mul(fromAmtBf, big.NewFloat(UndPow))
		result := new(big.Int)
		res.Int(result)
		return result.String() + to, nil
	case PundDenom: // from pund to und
		fromAmt, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return "", err
		}
		fromAmtBf := new(big.Float).SetFloat64(fromAmt)
		res := fromAmtBf.Mul(fromAmtBf, big.NewFloat(PundPow))
		return res.Text('f', 12) + to, nil
	}

	return "", nil
}
