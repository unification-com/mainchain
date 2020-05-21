package helpers

import (

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	SimAppChainID = "UND-Simulation-App" // SimAppChainID hardcoded chainID for simulation
)

// GenTx generates a signed mock transaction.
func GenTx(msgs []sdk.Msg, feeAmt sdk.Coins, chainID string, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) auth.StdTx {
	fee := auth.StdFee{
		Amount: feeAmt,
		Gas:    2500000, // TODO: this should be a param
	}

	sigs := make([]auth.StdSignature, len(priv))

	for i, p := range priv {
		sig, err := p.Sign(auth.StdSignBytes(chainID, accnums[i], seq[i], fee, msgs, ""))
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{
			PubKey:    p.PubKey(),
			Signature: sig,
		}
	}

	return auth.NewStdTx(msgs, fee, sigs, "")
}
