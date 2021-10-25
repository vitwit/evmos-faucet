package utils

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/pflag"
)

func ReadTxCommandFlags(clientCtx client.Context, flagSet *pflag.FlagSet) (client.Context, error) {
	clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, flagSet)
	if err != nil {
		return clientCtx, err
	}

	if !clientCtx.GenerateOnly || flagSet.Changed(flags.FlagGenerateOnly) {
		genOnly, _ := flagSet.GetBool(flags.FlagGenerateOnly)
		clientCtx = clientCtx.WithGenerateOnly(genOnly)
	}

	if !clientCtx.Offline || flagSet.Changed(flags.FlagOffline) {
		offline, _ := flagSet.GetBool(flags.FlagOffline)
		clientCtx = clientCtx.WithOffline(offline)
	}

	if !clientCtx.UseLedger || flagSet.Changed(flags.FlagUseLedger) {
		useLedger, _ := flagSet.GetBool(flags.FlagUseLedger)
		clientCtx = clientCtx.WithUseLedger(useLedger)
	}

	if clientCtx.BroadcastMode == "" || flagSet.Changed(flags.FlagBroadcastMode) {
		bMode, _ := flagSet.GetString(flags.FlagBroadcastMode)
		clientCtx = clientCtx.WithBroadcastMode(bMode)
	}

	if !clientCtx.SkipConfirm || flagSet.Changed(flags.FlagSkipConfirmation) {
		skipConfirm, _ := flagSet.GetBool(flags.FlagSkipConfirmation)
		clientCtx = clientCtx.WithSkipConfirmation(skipConfirm)
	}

	if clientCtx.SignModeStr == "" || flagSet.Changed(flags.FlagSignMode) {
		signModeStr, _ := flagSet.GetString(flags.FlagSignMode)
		clientCtx = clientCtx.WithSignModeStr(signModeStr)
	}

	if clientCtx.FeeGranter == nil || flagSet.Changed(flags.FlagFeeAccount) {
		granter, _ := flagSet.GetString(flags.FlagFeeAccount)

		if granter != "" {
			granterAcc, err := sdk.AccAddressFromBech32(granter)
			if err != nil {
				return clientCtx, err
			}

			clientCtx = clientCtx.WithFeeGranterAddress(granterAcc)
		}
	}

	if clientCtx.From == "" || flagSet.Changed(flags.FlagFrom) {
		from, _ := flagSet.GetString(flags.FlagFrom)
		fromAddr, fromName, keyType, err := client.GetFromFields(clientCtx.Keyring, from, clientCtx.GenerateOnly)
		if err != nil {
			return clientCtx, err
		}

		clientCtx = clientCtx.WithFrom(from).WithFromAddress(fromAddr).WithFromName(fromName)

		// If the `from` signer account is a ledger key, we need to use
		// SIGN_MODE_AMINO_JSON, because ledger doesn't support proto yet.
		// ref: https://github.com/cosmos/cosmos-sdk/issues/8109
		if keyType == keyring.TypeLedger && clientCtx.SignModeStr != flags.SignModeLegacyAminoJSON {
			fmt.Println("Default sign-mode 'direct' not supported by Ledger, using sign-mode 'amino-json'.")
			clientCtx = clientCtx.WithSignModeStr(flags.SignModeLegacyAminoJSON)
		}
	}

	return clientCtx, nil
}
