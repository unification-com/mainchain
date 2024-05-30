package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/stream/types"
	"testing"
	"time"
)

func TestPeriodEnumFromString(t *testing.T) {
	testCases := []struct {
		name     string
		period   string
		expected types.StreamPeriod
	}{
		{"parse Second", "Second", types.StreamPeriodSecond},
		{"parse second", "second", types.StreamPeriodSecond},
		{"parse sec", "sec", types.StreamPeriodSecond},
		{"parse Minute", "Minute", types.StreamPeriodMinute},
		{"parse minute", "minute", types.StreamPeriodMinute},
		{"parse min", "min", types.StreamPeriodMinute},
		{"parse Hour", "Hour", types.StreamPeriodHour},
		{"parse hour", "hour", types.StreamPeriodHour},
		{"parse Day", "Day", types.StreamPeriodDay},
		{"parse day", "day", types.StreamPeriodDay},
		{"parse Week", "Week", types.StreamPeriodWeek},
		{"parse week", "week", types.StreamPeriodWeek},
		{"parse Month", "Month", types.StreamPeriodMonth},
		{"parse month", "month", types.StreamPeriodMonth},
		{"parse mon", "mon", types.StreamPeriodMonth},
		{"parse Year", "Year", types.StreamPeriodYear},
		{"parse year", "year", types.StreamPeriodYear},
		{"parse kwbefi", "kwbefi", types.StreamPeriodUnspecified},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			res := types.PeriodEnumFromString(tc.period)
			require.Equal(t, tc.expected, res)
		})
	}
}

func TestCalculateFlowRateForCoin(t *testing.T) {
	testCases := []struct {
		name             string
		coin             sdk.Coin
		period           types.StreamPeriod
		duration         uint64
		expectedDuration uint64
		expectedFlowRate int64
	}{
		{"1", sdk.NewInt64Coin("testdenom", 1000), types.StreamPeriodSecond, 1, 1, 1000},
		{"2", sdk.NewInt64Coin("testdenom", 1000), types.StreamPeriodMinute, 1, 60, 16},        // 1000 / 60 = 16.666666667
		{"3", sdk.NewInt64Coin("testdenom", 23423423), types.StreamPeriodMonth, 1, 2628000, 8}, // 23423423 / 2628000 = 8.913022451
		{"4", sdk.NewInt64Coin("testdenom", 23467645081223423), types.StreamPeriodMonth, 2, 5256000, 4464924863},
		{"5", sdk.NewInt64Coin("testdenom", 23467645081223423), types.StreamPeriodYear, 1, 31536000, 744154143},
		{"6", sdk.NewInt64Coin("testdenom", 77000000000), types.StreamPeriodMonth, 1, 2628000, 29299},
		{"7", sdk.NewInt64Coin("testdenom", 46000000000), types.StreamPeriodMonth, 1, 2628000, 17503},
		{"8", sdk.NewInt64Coin("testdenom", 459000000000), types.StreamPeriodMonth, 1, 2628000, 174657},
		{"9", sdk.NewInt64Coin("testdenom", 4584000000000), types.StreamPeriodMonth, 1, 2628000, 1744292},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			duration, _, flowRate := types.CalculateFlowRateForCoin(tc.coin, tc.period, tc.duration)
			require.Equal(t, tc.expectedDuration, duration, "duration")
			require.Equal(t, tc.expectedFlowRate, flowRate, "flowRate")
		})
	}
}

func TestCalculateDuration(t *testing.T) {
	testCases := []struct {
		name             string
		coin             sdk.Coin
		flowRate         int64
		expectedDuration int64
	}{
		{"1", sdk.NewInt64Coin("testdenom", 1000), 1000, 1},
		{"2", sdk.NewInt64Coin("testdenom", 1000), 16, 62},
		{"3", sdk.NewInt64Coin("testdenom", 23423423), 8, 2927927},
		{"4", sdk.NewInt64Coin("testdenom", 23467645081223423), 4464924863, 5256000},
		{"5", sdk.NewInt64Coin("testdenom", 23467645081223423), 744154143, 31536000},
		{"6", sdk.NewInt64Coin("testdenom", 77000000000), 29299, 2628076},
		{"7", sdk.NewInt64Coin("testdenom", 46000000000), 17503, 2628120},
		{"8", sdk.NewInt64Coin("testdenom", 459000000000), 174657, 2628008},
		{"9", sdk.NewInt64Coin("testdenom", 4584000000000), 1744292, 2628000},
		{"10", sdk.NewInt64Coin("testdenom", 0), 123456789, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			duration := types.CalculateDuration(tc.coin, tc.flowRate)
			require.Equal(t, tc.expectedDuration, duration)
		})
	}
}

func TestCalculateAmountToClaim(t *testing.T) {
	nowTime := time.Now()
	testCases := []struct {
		name                          string
		nowTime                       time.Time
		depositZeroTime               time.Time
		lastOutflowTime               time.Time
		deposit                       sdk.Coin
		flowRate                      int64
		expectedAmountToClaim         sdk.Coin
		expectedRemainingDepositValue sdk.Coin
	}{
		{
			"1",
			nowTime,
			nowTime,
			nowTime,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			1000,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(0)),
		},
		{
			"2",
			nowTime,
			nowTime.Add(time.Second * time.Duration(1000)),
			time.Unix(nowTime.Unix()-1000, 0),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(2000)),
			1,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
		},
		{
			"3",
			nowTime,
			nowTime.Add(time.Second * time.Duration(1)),
			time.Unix(nowTime.Unix()-999, 0),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			1,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(999)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1)),
		},
		{
			"4",
			nowTime,
			nowTime.Add(time.Second * time.Duration(940)),
			time.Unix(nowTime.Unix()-60, 0),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			1,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(60)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(940)),
		},
		{
			"5",
			nowTime,
			nowTime,
			time.Unix(nowTime.Unix()-234276, 0),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1494667526000)),
			6379943,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1494667526000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(0)),
		},
		{
			"6",
			nowTime,
			nowTime.Add(time.Second * time.Duration(8626)),
			time.Unix(nowTime.Unix()-23427, 0),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(204496312979)),
			6379943,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(149462924661)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(55033388318)),
		},
		{
			"7",
			nowTime,
			nowTime.Add(time.Second * time.Duration(1)),
			time.Unix(nowTime.Unix()-2627999, 0),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(4584000003123)),
			1744292,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(4583997631708)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(2371415)),
		},
		{
			"8",
			nowTime,
			nowTime.Add(time.Second * time.Duration(10)),
			time.Unix(nowTime.Unix()-2627999, 0),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(450000000000)), // deposit is less than 1744292 * 2627999
			1744292,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(450000000000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(0)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			amountToClaim, remainingDeposit := types.CalculateAmountToClaim(tc.nowTime, tc.depositZeroTime, tc.lastOutflowTime, tc.deposit, tc.flowRate)
			require.Equal(t, tc.expectedAmountToClaim, amountToClaim, "amountToClaim")
			require.Equal(t, tc.expectedRemainingDepositValue, remainingDeposit, "remainingDeposit")
		})
	}
}

func TestCalculateValidatorFee(t *testing.T) {

	zeroPerc, _ := sdk.NewDecFromStr("0.0")
	onePerc, _ := sdk.NewDecFromStr("0.01")
	fivePerc, _ := sdk.NewDecFromStr("0.05")
	tenPerc, _ := sdk.NewDecFromStr("0.1")
	twentyFourPerc, _ := sdk.NewDecFromStr("0.24")
	ninetyNinePerc, _ := sdk.NewDecFromStr("0.99")
	hundredPerc, _ := sdk.NewDecFromStr("1.0")

	testCases := []struct {
		name                   string
		valFee                 sdk.Dec
		amountToClaim          sdk.Coin
		expectedFinalClaimCoin sdk.Coin
		expectedValFeeCoin     sdk.Coin
	}{
		{
			"1",
			zeroPerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(0)),
		},
		{
			"2",
			onePerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(990)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(10)),
		},
		{
			"3",
			tenPerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(900)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(100)),
		},
		{
			"4",
			fivePerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(950)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(50)),
		},
		{
			"5",
			twentyFourPerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(760)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(240)),
		},
		{
			"6",
			ninetyNinePerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(10)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(990)),
		},
		{
			"7",
			hundredPerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(0)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(1000)),
		},
		{
			"8",
			onePerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(8723642874687)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(8636406445941)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(87236428746)),
		},
		{
			"9",
			twentyFourPerc,
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(912742861395)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(693684574661)),
			sdk.NewCoin("testdenom", sdk.NewIntFromUint64(219058286734)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			finalClaimCoin, valFeeCoin := types.CalculateValidatorFee(tc.valFee, tc.amountToClaim)
			if tc.expectedFinalClaimCoin.Amount.IsZero() {
				require.True(t, finalClaimCoin.IsZero())
			} else {
				require.Equal(t, tc.expectedFinalClaimCoin, finalClaimCoin, "finalClaimCoin ")
			}
			require.Equal(t, tc.expectedValFeeCoin, valFeeCoin, "valFeeCoin")
			require.Equal(t, tc.amountToClaim, finalClaimCoin.Add(valFeeCoin), "total")
		})
	}
}
