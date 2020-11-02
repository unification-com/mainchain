package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagHeight     = "height"
	flagDumpWhat   = "dump-what"
	flagDumpId     = "dump-id"
	flagTraceStore = "trace-store"
)

// DumpDataCmd dumps app state to JSON.
func DumpDataCmd(ctx *server.Context, cdc *codec.Codec,
	dumper func(log.Logger, dbm.DB, io.Writer, int64, string, uint64) (dataDump json.RawMessage, err error),
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump-data",
		Short: "Dump WRKChain or BEACON data to JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			traceWriterFile := viper.GetString(flagTraceStore)

			dataDir := filepath.Join(config.RootDir, "data")
			db, err := sdk.NewLevelDB("application", dataDir)
			if err != nil {
				return err
			}

			if isEmptyState(db) || dumper == nil {
				if _, err := fmt.Fprintln(os.Stderr, "WARNING: State is not initialized. Returning empty result."); err != nil {
					return err
				}

				fmt.Println("[]")
				return nil
			}

			traceWriter, err := openTraceWriter(traceWriterFile)

			if err != nil {
				return err
			}

			height := viper.GetInt64(flagHeight)
			dumpWhat := viper.GetString(flagDumpWhat)
			dumpId := viper.GetUint64(flagDumpId)

			switch dumpWhat {
			case "beacon":
			case "wrkchain":
				break
			default:
				if _, err := fmt.Fprintln(os.Stderr, "WARNING: --dump-what must be either beacon or wrkchain."); err != nil {
					return err
				}

				fmt.Println("[]")
				return nil
			}

			dataDump, err := dumper(ctx.Logger, db, traceWriter, height, dumpWhat, dumpId)
			if err != nil {
				return fmt.Errorf("error dumping data: %v", err)
			}

			encoded, err := codec.MarshalJSONIndent(cdc, dataDump)
			if err != nil {
				return err
			}

			fmt.Println(string(sdk.MustSortJSON(encoded)))
			return nil
		},
	}

	cmd.Flags().Int64(flagHeight, -1, "Export state from a particular height (-1 means latest height)")
	cmd.Flags().String(flagDumpWhat, "", "What to dump (beacon | wrkchain")
	cmd.Flags().Uint64(flagDumpId, 0, "ID of entity to dump")
	return cmd
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile != "" {
		w, err = os.OpenFile(
			traceWriterFile,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0666,
		)
		return
	}
	return
}

func isEmptyState(db dbm.DB) bool {
	return db.Stats()["leveldb.sstables"] == ""
}
