package v2_test

import (
	"testing"
	"time"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/app/params"
	"github.com/unification-com/mainchain/x/stream"
	v2 "github.com/unification-com/mainchain/x/stream/migrations/v2"
	"github.com/unification-com/mainchain/x/stream/types"
)

func mkAddr(b byte) sdk.AccAddress {
	out := make([]byte, 20)
	for i := range out {
		out[i] = b
	}
	return sdk.AccAddress(out)
}

// v1 key shape: 0x11 | len(receiver) | receiver | len(sender) | sender
func v1Key(receiver, sender sdk.AccAddress) []byte {
	out := append([]byte{}, types.StreamKeyPrefix...)
	out = append(out, address.MustLengthPrefix(receiver)...)
	out = append(out, address.MustLengthPrefix(sender)...)
	return out
}

func TestMigrateStore(t *testing.T) {
	params.SetAddressPrefixes()

	encCfg := moduletestutil.MakeTestEncodingConfig(stream.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	store := ctx.KVStore(storeKey)

	now := time.Unix(1700000000, 0).UTC()

	type seed struct {
		receiver sdk.AccAddress
		sender   sdk.AccAddress
		stream   types.Stream
	}

	seeds := []seed{
		{
			receiver: mkAddr(0xA1), sender: mkAddr(0xB1),
			stream: types.Stream{
				Deposit:         sdk.NewInt64Coin("nund", 1000),
				FlowRate:        10,
				LastOutflowTime: now,
				DepositZeroTime: now.Add(time.Minute),
				Cancellable:     true,
			},
		},
		{
			receiver: mkAddr(0xA1), sender: mkAddr(0xB2),
			stream: types.Stream{
				Deposit:         sdk.NewInt64Coin("ibc/ABC", 2000),
				FlowRate:        20,
				LastOutflowTime: now,
				DepositZeroTime: now.Add(2 * time.Minute),
				Cancellable:     false,
			},
		},
		{
			receiver: mkAddr(0xA2), sender: mkAddr(0xB1),
			stream: types.Stream{
				Deposit:         sdk.NewInt64Coin("stake", 12345),
				FlowRate:        1,
				LastOutflowTime: now,
				DepositZeroTime: now.Add(time.Hour),
				Cancellable:     true,
			},
		},
	}

	// Seed under v1 key shape (no denom).
	for _, s := range seeds {
		store.Set(v1Key(s.receiver, s.sender), cdc.MustMarshal(&s.stream))
	}

	// Migrate.
	require.NoError(t, v2.MigrateStore(ctx, storeKey, cdc))

	for _, s := range seeds {
		// Old key gone.
		require.Nil(t, store.Get(v1Key(s.receiver, s.sender)), "v1 key must be deleted")

		// New key resolves.
		v2KeyBytes := types.GetStreamKey(s.receiver, s.sender, s.stream.Deposit.Denom)
		raw := store.Get(v2KeyBytes)
		require.NotNil(t, raw, "v2 key must be set for denom=%s", s.stream.Deposit.Denom)

		var got types.Stream
		require.NoError(t, cdc.Unmarshal(raw, &got))
		require.Equal(t, s.stream.Deposit, got.Deposit)
		require.Equal(t, s.stream.FlowRate, got.FlowRate)
		require.Equal(t, s.stream.LastOutflowTime, got.LastOutflowTime)
		require.Equal(t, s.stream.DepositZeroTime, got.DepositZeroTime)
		require.Equal(t, s.stream.Cancellable, got.Cancellable)
	}

	// Re-key-parse round trip: every key parses to (receiver, sender, denom) we set.
	expected := map[string]string{}
	for _, s := range seeds {
		expected[s.receiver.String()+"|"+s.sender.String()] = s.stream.Deposit.Denom
	}
	iter := storetypes.KVStorePrefixIterator(store, types.StreamKeyPrefix)
	defer iter.Close()
	saw := map[string]string{}
	for ; iter.Valid(); iter.Next() {
		r, snd, d := types.AddressesFromStreamKey(iter.Key())
		saw[r.String()+"|"+snd.String()] = d
	}
	require.Equal(t, expected, saw)
}
