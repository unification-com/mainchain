# Third Party Tool Development

This page contains information that may be useful for developers of third party tools, such as wallets and block explorers

#### Contents

[[toc]]

## HD Wallet Path

The BIP-0044 Path for our HD Wallets is as follows:

`44'/5555'/0'/0`   

[SLIP-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md) Coin ID is `5555`

## REST Endpoints

The REST endpoints for API interaction (for example block explorers, wallets etc.), served by [light-clients](../software/light-client-rpc.md) via port **1337** can be found in [swagger.yaml](https://github.com/unification-com/mainchain/blob/master/client/lcd/swagger-ui/swagger.yaml)

Live examples can be found at [https://rest.unification.io/swagger-ui/](https://rest.unification.io/swagger-ui/).

## Tendermint RPC Endpoints

The Tendermint RPC endpoints, served by full-nodes via port **26657** can be found at [http://rpc1.unification.io:26657](http://rpc1.unification.io:26657).

The RPC specification is the same as [Tendermint](https://docs.tendermint.com/master/rpc/).
