package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gsk967/cosmos-faucet/utils"
	"log"
	"os"
)

func main() {
	// register all module interfaces for encoding/decoding
	encodingConfig := simapp.MakeTestEncodingConfig()

	//codec.RegisterCrypto(encodingConfig.Amino)
	codec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	utils.SetupConfig("cosmos")

	// client context
	_ = client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).WithViper("")

	kr, err := keyring.New("sim", "os", "/home/gsk967/.simapp", os.Stdin, func(options *keyring.Options) {
		options.SupportedAlgos = keyring.SigningAlgoList{
			hd.Secp256k1,
		}
	})
	if err != nil {
		log.Fatal("Error init keyring ", err.Error())
	}
	//
	accountAddress, err := sdk.AccAddressFromBech32("cosmos1kzkwx43pgyrr3tk4x8dqkaea5svv0tl9vpu0ty")
	if err != nil {
		log.Fatal("Error account keyring ", err.Error())
	}

	fmt.Println("acc address ", accountAddress)
	////
	info, err := kr.KeyByAddress(accountAddress)
	if err != nil {
		log.Fatal("Error keyring info ", err.Error())
	}

	fmt.Println("name ", info.GetName())
	fmt.Println("address ", info.GetAddress())

	list, err := kr.List()
	for _, in := range list {
		fmt.Println("name ", in.GetName())
		fmt.Println("address ", in.GetAddress())
	}

	info, err = kr.Key("root")
	if err != nil {
		log.Fatal("Error keyring info ", err.Error())
	}

	fmt.Println("1 name ", info.GetName())
	fmt.Println("1 address ", info.GetAddress())

}
