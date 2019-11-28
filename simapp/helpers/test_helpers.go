package helpers

import (
	"math/rand"
	"time"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

const (
	SimAppChainID = "UND-Simulation-App" // SimAppChainID hardcoded chainID for simulation
)

// GenTx generates a signed mock transaction.
func GenTx(msgs []sdk.Msg, feeAmt sdk.Coins, chainID string, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) auth.StdTx {
	fee := auth.StdFee{
		Amount: feeAmt,
		Gas:    1500000, // TODO: this should be a param
	}

	sigs := make([]auth.StdSignature, len(priv))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memo := simulation.RandStringOfLength(r, 100)

	// fuzzy adjust gas for memo size
	fee.Gas = fee.Gas + (auth.DefaultTxSizeCostPerByte * uint64(len([]byte(memo)))) + 500

	for i, p := range priv {
		sig, err := p.Sign(auth.StdSignBytes(chainID, accnums[i], seq[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{
			PubKey:    p.PubKey(),
			Signature: sig,
		}
	}

	return auth.NewStdTx(msgs, fee, sigs, memo)
}
