package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gsk967/cosmos-faucet/cmd"
	faucetconfig "github.com/gsk967/cosmos-faucet/config"
	"github.com/gsk967/cosmos-faucet/utils"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/tharsis/ethermint/app"
	"github.com/tharsis/ethermint/encoding"
	"os"
)

var cfg *faucetconfig.Config

func init() {
	// read the configuration
	fileConfig, err := faucetconfig.ReadConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}
	cfg = fileConfig
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "ethermint-faucet",
		Short: "ethermint faucet",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// setting up sdk bech32 config
			utils.SetupConfig(cfg.Faucet.AccountPrefix)

			// register all module interfaces for encoding/decoding
			encodingConfig := encoding.MakeConfig(app.ModuleBasics)

			// client context
			initClientCtx := client.Context{}.
				WithCodec(encodingConfig.Marshaler).
				WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
				WithTxConfig(encodingConfig.TxConfig).
				WithLegacyAmino(encodingConfig.Amino).
				WithAccountRetriever(types.AccountRetriever{}).
				WithBroadcastMode(flags.BroadcastBlock)

			if cfg.Faucet.EnvPrefix == "" {
				initClientCtx = initClientCtx.WithViper("")
			} else {
				initClientCtx = initClientCtx.WithViper(cfg.Faucet.EnvPrefix)
			}

			initClientCtx = client.ReadHomeFlag(initClientCtx, cmd)
			initClientCtx, err := config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				log.Fatal(err)
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	// add sub-commands
	rootCmd.AddCommand(
		// ui interface
		cmd.StartServer(cfg),
		// cli interface for transferring the test tokens
		cmd.FromCli(cfg),
	)

	// add flags
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")
	err := rootCmd.MarkPersistentFlagRequired(flags.FlagChainID)
	if err != nil {
		log.Fatal(err)
	}

	if err := svrcmd.Execute(rootCmd, ""); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
