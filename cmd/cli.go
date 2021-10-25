package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	faucetconfig "github.com/gsk967/cosmos-faucet/config"
	"github.com/gsk967/cosmos-faucet/utils"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

const (
	FlagTo = "to"
)

func FromCli(cfg *faucetconfig.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cli",
		Short: "Transfer tokens from cli command of ethermint.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			to, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}

			toAddress, err := utils.ValidateAccountAddress(to)
			if err != nil {
				return err
			}

			return utils.GetTokens(clientCtx, cfg, cmd.Flags(), toAddress)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagTo, "", "The destination address.")
	err := cmd.MarkFlagRequired(FlagTo)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.MarkFlagRequired(flags.FlagFrom)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}
