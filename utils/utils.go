package utils

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
	"os"
)

func ValidateAccountAddress(val string) (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32(val)
}

func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(key, "=", value)
		return value
	} else {
		log.Fatal("Error loading environment variable: ", key)
		return ""
	}
}
