package utils

import sdk "github.com/cosmos/cosmos-sdk/types"

func SetupConfig(bech32Prefix string) {
	// set the address prefixes
	config := sdk.GetConfig()
	var (
		// bech32PrefixAccAddr defines the bech32 prefix of an account's address
		bech32PrefixAccAddr = bech32Prefix
		// bech32PrefixAccPub defines the bech32 prefix of an account's public key
		bech32PrefixAccPub = bech32Prefix + sdk.PrefixPublic
		// bech32PrefixValAddr defines the bech32 prefix of a validator's operator address
		bech32PrefixValAddr = bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
		// bech32PrefixValPub defines the bech32 prefix of a validator's operator public key
		bech32PrefixValPub = bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
		// bech32PrefixConsAddr defines the bech32 prefix of a consensus node address
		bech32PrefixConsAddr = bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
		// bech32PrefixConsPub defines the bech32 prefix of a consensus node public key
		bech32PrefixConsPub = bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
	)

	config.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
	config.Seal()
}
