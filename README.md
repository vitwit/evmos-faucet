# Faucet for Ethermint

This wil help to get test tokens for ethermint

## Config

```yaml
[ ui ]
  port = 1314

  [ faucet ]
  # account_prefix for bech32 encoding
  account_prefix="ethm"
  # env_prefix for storing the all env variables
  env_prefix="ETHERMINT"
  # amount is tokens per request
  amount = 10000
  # maximum tokens allowed for an account
  max_tokens = 10000000
  # tendermint node address
  # <host>:<port> to Tendermint RPC interface for this chain
  node = "tcp://localhost:26657"
  # chain denom
  denom = "aphoton"
```

## Build

```shell
$ make clean && make  build 
```

## From Cli

```shell
$ ./build/faucet cli --from [from-account-name] --chain-id [chain-id] --to [from-address]
Ex: $ ./build/faucet cli --from root --chain-id ethermint_9000-1 --to ethm1wvy5qkdcwugq7dknxhv56yzj5yzkjejrkvcf4e

```

## UI Server

```shell
$ ./build/faucet server --from root --chain-id ethermint_9000-1
```