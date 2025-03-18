package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/cmd/und/cmd"
)

func main() {
	// Set config for address prefixes to "und", hd path to 5555 etc.
	app.SetConfig()

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
