package simulation

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
)

var (
	bPk1   = ed25519.GenPrivKey().PubKey()
	bAddr1 = sdk.AccAddress(bPk1.Address())
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()

	beacon := types.NewBeacon()
	beacon.BeaconID = 1
	beacon.Moniker = "beacon1"
	beacon.Name = "Test BEACON 1"
	beacon.LastTimestampID = 1
	beacon.Owner = bAddr1

	beaconTs := types.NewBeaconTimestamp()
	beaconTs.BeaconID = 1
	beaconTs.TimestampID = 1
	beaconTs.Owner = bAddr1
	beaconTs.Hash = "arbitraryblockhashvalue"
	beaconTs.SubmitTime = uint64(time.Now().Unix())

	kvPairs := cmn.KVPairs{
		cmn.KVPair{Key: types.BeaconKey(1), Value: cdc.MustMarshalBinaryLengthPrefixed(beacon)},
		cmn.KVPair{Key: types.BeaconTimestampKey(1, 1), Value: cdc.MustMarshalBinaryLengthPrefixed(beaconTs)},
		cmn.KVPair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"beacon", fmt.Sprintf("%v\n%v", beacon, beacon)},
		{"beacon timestamp", fmt.Sprintf("%v\n%v", beaconTs, beaconTs)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
