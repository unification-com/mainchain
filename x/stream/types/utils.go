package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PeriodEnumFromString(period string) StreamPeriod {
	switch period {
	case "Second", "second", "sec":
		return StreamPeriodSecond
	case "Minute", "minute", "min":
		return StreamPeriodMinute
	case "Hour", "hour":
		return StreamPeriodHour
	case "Day", "day":
		return StreamPeriodDay
	case "Week", "week":
		return StreamPeriodWeek
	case "Month", "month", "mon":
		return StreamPeriodMonth
	case "Year", "year":
		return StreamPeriodYear
	}
	return StreamPeriodUnspecified
}

func CalculateFlowRateForCoin(coin sdk.Coin, period StreamPeriod, duration uint64) (uint64, sdk.Dec, int64) {
	baseDuration := uint64(1)
	totalDuration := uint64(1)

	switch period {
	case StreamPeriodUnspecified:
		baseDuration = 1
	default:
		baseDuration = 1
	case StreamPeriodSecond:
		baseDuration = 1
	case StreamPeriodMinute:
		baseDuration = 60
	case StreamPeriodHour:
		baseDuration = 3600
	case StreamPeriodDay:
		baseDuration = 86400
	case StreamPeriodWeek:
		baseDuration = 604800
	case StreamPeriodMonth:
		baseDuration = 2628000 // (365 / 12) * 24 * 60 * 60 = 30.416666667 * 24 * 60 * 60
	case StreamPeriodYear:
		baseDuration = 31536000
	}

	totalDuration = baseDuration * duration

	if coin.IsNil() || coin.IsNegative() || coin.IsZero() || totalDuration == 0 {
		return totalDuration, sdk.NewDecWithPrec(0, 0), 0
	}

	// flow rate calculation from deposit and duration
	decCoin := sdk.NewDecCoinFromCoin(coin)
	decDuration := sdk.NewDecFromInt(sdk.NewIntFromUint64(totalDuration))

	flowRate := decCoin.Amount.QuoTruncateMut(decDuration)

	// note: decimal values are rounded down, e.g. 8.9 to just 8.
	return totalDuration, flowRate, flowRate.TruncateInt64()
}

func CalculateDuration(deposit sdk.Coin, flowRate int64) int64 {
	// no point if flowRate is <= 0
	if flowRate <= 0 {
		return 0
	}
	// no point if the deposit value is zero - e.g. if re-calculating from a new flow rate
	// of an existing stream
	if deposit.Amount.GT(sdk.NewIntFromUint64(0)) {
		// calculate duration in seconds
		decFlowRate := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(flowRate)))
		decDeposit := sdk.NewDecCoinFromCoin(deposit)
		decDuration := decDeposit.Amount.QuoTruncateMut(decFlowRate)
		// note: decimal values are rounded down, e.g. 2628008.9 to just 2628008.
		return decDuration.TruncateInt64()
	}

	return 0
}

func CalculateAmountToClaim(
	nowTime,
	depositZeroTime,
	lastOutflowTime time.Time,
	deposit sdk.Coin,
	flowRate int64,
) (sdk.Coin, sdk.Coin) {
	var amountToClaim sdk.Coin
	var remainingDepositValue sdk.Coin

	if nowTime.After(depositZeroTime) || nowTime.Equal(depositZeroTime) {
		// now > deposit_zero_time, use all remaining deposit
		amountToClaim = deposit
		remainingDepositValue = sdk.NewCoin(deposit.Denom, sdk.NewInt(0))
	} else {
		// calculate based on flow rate and remaining deposit
		timeSinceLast := nowTime.Sub(lastOutflowTime)
		secondsSinceLast := int64(timeSinceLast.Seconds())
		numCoins := secondsSinceLast * flowRate
		amountToClaim = sdk.NewCoin(deposit.Denom, sdk.NewIntFromUint64(uint64(numCoins)))
		if deposit.Amount.GT(amountToClaim.Amount) {
			remainingDepositValue = deposit.Sub(amountToClaim)
		} else {
			// just in case
			amountToClaim = deposit
			remainingDepositValue = sdk.NewCoin(deposit.Denom, sdk.NewInt(0))
		}
	}

	return amountToClaim, remainingDepositValue
}

func CalculateValidatorFee(valFee sdk.Dec, amountToClaim sdk.Coin) (sdk.Coin, sdk.Coin) {
	var valFeeCoin sdk.Coin
	var finalClaimCoin sdk.Coin

	if valFee.GT(sdk.NewDecFromInt(sdk.NewIntFromUint64(0))) {
		decCoin := sdk.NewDecCoinFromCoin(amountToClaim)
		valFeeAmount := decCoin.Amount.Mul(valFee).TruncateInt64()
		valFeeCoin = sdk.NewCoin(amountToClaim.Denom, sdk.NewIntFromUint64(uint64(valFeeAmount)))
		finalClaimCoin = amountToClaim.Sub(valFeeCoin)
	} else {
		valFeeCoin = sdk.NewCoin(amountToClaim.Denom, sdk.NewIntFromUint64(0))
		finalClaimCoin = amountToClaim
	}

	return finalClaimCoin, valFeeCoin
}
