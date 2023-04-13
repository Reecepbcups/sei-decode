package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/spf13/cobra"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/CosmosContracts/juno/v13/app"
	"github.com/CosmosContracts/juno/v13/app/params"
)

// juno-decode tx decode <tx-b64>

func main() {
	rootCmd, _ := NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	cfg.Seal()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		// WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		// WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   version.AppName,
		Short: "Juno Network Lightweight Decoder",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	// Keeping it as close to the original as possible
	rootCmd.AddCommand(
		txCommand(),
	)
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetDecodeCommand(),
		GetFileDecodeCommand(),
	)

	return cmd
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// create a struct which is {"key":"value"} with both being strings

// GetDecodeCommand returns the decode command to take serialized bytes and turn
// it into a JSON-encoded transaction.
func GetFileDecodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode-file [file]",
		Short: "Decode a bunch of amino bytre strings in 1 file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			var txBytes []byte

			dat, err := os.ReadFile("amino.json")
			check(err)
			// fmt.Print(string(dat))

			values := []string{}
			err = json.Unmarshal(dat, &values)
			check(err)

			// fmt.Println(values)

			new_values := []string{}
			for _, value := range values {
				txBytes, err = base64.StdEncoding.DecodeString(value)
				if err != nil {
					return err
				}

				tx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
				if err != nil {
					return err
				}

				json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
				if err != nil {
					return err
				}

				new_values = append(new_values, string(json))
			}

			// return new_values (or do we save to a file maybe?)
			return clientCtx.PrintBytes([]byte(strings.Join(new_values, ";;;")))
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
