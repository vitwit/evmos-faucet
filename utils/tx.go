package utils

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	faucetconfig "github.com/gsk967/cosmos-faucet/config"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

func GetTokens(ctx client.Context, cfg *faucetconfig.Config, flagSet *pflag.FlagSet, toAddress sdk.AccAddress) error {
	ctx, err := ReadTxCommandFlags(ctx, flagSet)
	if err != nil {
		fmt.Errorf("couldn't ReadPersistentCommandFlags: %v", err)
		return err
	}
	ctx = ctx.WithBroadcastMode(flags.BroadcastBlock)
	fromAddress := ctx.GetFromAddress()
	log.Info(fmt.Sprintf("Transfering the %d%s tokens from %s to %s", cfg.Faucet.Amount, cfg.Faucet.Denom, fromAddress, toAddress))
	msgSend := banktypes.NewMsgSend(fromAddress, toAddress, sdk.NewCoins(sdk.NewCoin(cfg.Faucet.Denom, sdk.NewInt(int64(cfg.Faucet.Amount)))))
	sdkResponse, err := submitTx(ctx, flagSet, msgSend)
	if err != nil {
		_ = fmt.Errorf("error at submit the tx %v", err)
		return err
	}
	if sdkResponse.Code != 0 {
		log.Error(fmt.Sprintf("Error while transnfering the %d%s tokens to %s \n Err : %s", cfg.Faucet.Amount, cfg.Faucet.Denom, toAddress, sdkResponse.RawLog))
		return errors.New(sdkResponse.RawLog)
	} else {
		log.Info(fmt.Sprintf("%d%s tokens are successfully transfered to %s", cfg.Faucet.Amount, cfg.Faucet.Denom, toAddress))
	}
	return nil
}

// submitTx will submit the signed sdk.Msg to tendermint node
func submitTx(clientCtx client.Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	// validate the messages
	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
	}
	txf := tx.NewFactoryCLI(clientCtx, flagSet)
	txf, err := prepareFactory(clientCtx, txf)
	if err != nil {
		return nil, err
	}

	txBuilder, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}

	txBuilder.SetFeeGranter(clientCtx.GetFeeGranterAddress())
	err = tx.Sign(txf, clientCtx.GetFromName(), txBuilder, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	// broadcast to a Tendermint node
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}
	return res, err
}

func prepareFactory(clientCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := clientCtx.GetFromAddress()
	//
	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}
