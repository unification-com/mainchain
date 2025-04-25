package cmd

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	appparams "github.com/unification-com/mainchain/app/params"
)

const (
	FundPow = 1e9  // multiplier for converting from und to (nano) nund
	NundPow = 1e-9 // multiplier for converting from (nano) nund to und
)

var (
	flagBech32Prefix = "prefix"
)

func ConvertUndDenomination(amount string, from string, to string) (string, error) {

	if from == to {
		return amount + from, nil
	}

	switch from {
	case appparams.HumanCoinUnit: // from und to nund
		fromAmt, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return "", err
		}
		fromAmtBf := new(big.Float).SetFloat64(fromAmt)
		res := fromAmtBf.Mul(fromAmtBf, big.NewFloat(FundPow))
		result := new(big.Int)
		res.Int(result)
		return result.String() + to, nil
	case appparams.BaseCoinUnit: // from nund to fund
		fromAmt, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return "", err
		}
		fromAmtBf := new(big.Float).SetFloat64(fromAmt)
		res := fromAmtBf.Mul(fromAmtBf, big.NewFloat(NundPow))
		return res.Text('f', 9) + to, nil
	}

	return "", nil
}

// ConvertBech32Prefix convert bech32 address to specified prefix.
func ConvertBech32Prefix(address, prefix string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", fmt.Errorf("cannot decode %s address: %s", address, err)
	}

	convertedAddress, err := bech32.ConvertAndEncode(prefix, bz)
	if err != nil {
		return "", fmt.Errorf("cannot convert %s address: %s", address, err)
	}

	return convertedAddress, nil
}

func GetDenomConversionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "convert [amount] [from_denom] [to_denom]",
		Short: "convert between FUND denominations",
		Long: strings.TrimSpace(
			fmt.Sprintf(`convert between FUND denominations'
Example:
$ %s convert 24 fund nund
`,
				version.AppName,
			),
		),

		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, err := ConvertUndDenomination(args[0], args[1], args[2])

			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%s%s = %s\n", args[0], args[1], res))
		},
	}
}

// AddBech32ConvertCommand returns bech32-convert cobra Command.
func AddBech32ConvertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bech32-convert [address]",
		Short: "Convert any bech32 string to the und prefix",
		Long: `Convert any bech32 string to the und prefix

Example:
	gaiad debug bech32-convert akash1a6zlyvpnksx8wr6wz8wemur2xe8zyh0ytz6d88

	gaiad debug bech32-convert stride1673f0t8p893rqyqe420mgwwz92ac4qv6synvx2 --prefix osmo
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bech32prefix, err := cmd.Flags().GetString(flagBech32Prefix)
			if err != nil {
				return err
			}

			address := args[0]
			convertedAddress, err := ConvertBech32Prefix(address, bech32prefix)
			if err != nil {
				return fmt.Errorf("conversation failed: %s", err)
			}

			cmd.Println(convertedAddress)

			return nil
		},
	}

	cmd.Flags().StringP(flagBech32Prefix, "p", "und", "Bech32 Prefix to encode to")

	return cmd
}

// addDebugCommands injects custom debug commands into another command as children.
func addDebugCommands(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(AddBech32ConvertCommand())
	return cmd
}
