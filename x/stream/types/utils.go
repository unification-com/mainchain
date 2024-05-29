package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
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
		baseDuration = 60 * 60
	case StreamPeriodDay:
		baseDuration = 24 * 60 * 60
	case StreamPeriodWeek:
		baseDuration = 7 * 24 * 60 * 60
	case StreamPeriodMonth:
		baseDuration = (365 / 12) * 24 * 60 * 60
	case StreamPeriodYear:
		baseDuration = 365 * 24 * 60 * 60
	}

	totalDuration = baseDuration * duration

	// flow rate calculation from deposit and duration
	decCoin := sdk.NewDecCoinFromCoin(coin)
	decDuration := sdk.NewDecFromInt(sdk.NewIntFromUint64(totalDuration))

	flowRate := decCoin.Amount.QuoTruncateMut(decDuration)

	return totalDuration, flowRate, flowRate.TruncateInt64()
}

func CalculateDuration(deposit sdk.Coin, flowRate int64) int64 {
	// calculate duration in seconds
	decFlowRate := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(flowRate)))
	decDeposit := sdk.NewDecCoinFromCoin(deposit)
	decDuration := decDeposit.Amount.QuoTruncateMut(decFlowRate)
	return decDuration.TruncateInt64()
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

	if nowTime.After(depositZeroTime) {
		// now > deposit_zero_time, use all remaining deposit
		amountToClaim = deposit
		remainingDepositValue = sdk.NewCoin(deposit.Denom, sdk.NewInt(0))
	} else {
		// calculate based on flow rate and remaining deposit
		timeSinceLast := nowTime.Sub(lastOutflowTime)
		secondsSinceLast := int64(timeSinceLast.Seconds())
		numCoins := secondsSinceLast * flowRate
		amountToClaim = sdk.NewCoin(deposit.Denom, sdk.NewIntFromUint64(uint64(numCoins)))
		// ToDo - use SafeSub
		remainingDepositValue = deposit.Sub(amountToClaim)
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
		// ToDo - use SafeSub
		finalClaimCoin = amountToClaim.Sub(valFeeCoin)
	} else {
		valFeeCoin = sdk.NewCoin(amountToClaim.Denom, sdk.NewIntFromUint64(0))
		finalClaimCoin = amountToClaim
	}

	return finalClaimCoin, valFeeCoin
}
