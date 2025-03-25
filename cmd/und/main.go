package main

import (
	"fmt"
	"github.com/unification-com/mainchain/app"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	appparams "github.com/unification-com/mainchain/app/params"
	"github.com/unification-com/mainchain/cmd/und/cmd"
)

func main() {
	// Set config for address prefixes to "und", hd path to 5555 etc.
	appparams.SetAddressPrefixes()

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
