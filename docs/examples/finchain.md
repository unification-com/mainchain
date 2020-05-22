# WRKChain: Finchains

Finchain is a live WRKChain. The full source code is available on [Github](https://github.com/unification-com/finchain).

Finchain is an Ethereum based WRKChain that utilises a smart contract to analyse Stock data. Stock data is written to the smart contract from several API sources - each source API has an oracle periodically querying the APIs and submitting the data to the smart contract. The smart contract analyses input sent from the API oracles, and emits events for both submitted stock price updates and when a discrepancy is found between APIs' submitted prices. Discrepancies are detected when price differences exceed a configurable threshold value.  

[[toc]]

## Public Finchains

The TestNet Finchains can be viewed here: [https://finchain.unification.io](https://finchain.unification.io)

The TestNet Finchains writes its block hashes to the [Mainchain TestNet](https://explorer-testnet.unification.io/).

## Running Finchains Locally

Finchains can also be run locally as a completely self-contained environment, to allow developers to play with different configurations, and see how the internals of a WRKChain work.

Docker and Docker Compose are required to run the localised, self-contained
Finchains.

Copy `example.env` to `.env` and make any required changes. API keys are required
for the composition to work - see `example.env` for details on where to obtain the
necessary API keys.

Run the composition using:

```bash
make
```

### WRKChain: Localised Finchains Docker Composition Info

Finchains's WRKChain is a `geth` (Ethereum) based WRKChain.

Network ID: `2339117895`  

Block Explorer: [http://localhost:8081](http://localhost:8081)

JSON RPC: [http://localhost:8547](http://localhost:8547)

WRKChain Block Validation UI: [http://localhost:4040](http://localhost:4040)


### Local UND Mainchain DevNet

The local Finchains composition contains a completely self-contained pre-configured local Mainchain DevNet:

Chain ID: `FUND-Mainchain-DevNet  `

Block Explorer: [http://localhost:3000](http://localhost:3000)

RPC: [http://localhost:26661](http://localhost:26661)
REST: [https://localhost:1318](https://localhost:1318)
