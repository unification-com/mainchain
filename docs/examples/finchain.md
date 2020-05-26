# WRKChain: Finchains

Finchain is a live WRKChain. The full source code is available on [Github](https://github.com/unification-com/finchain).

[[toc]]

## Introduction

Finchains is WRKChain built using the Go-Ethereum codebase, and utilises smart contracts to analyse Crypto, Stocks, and other data. Data is written to the smart contracts from several API sources - each source API has an oracle periodically querying the APIs and submitting the data to the smart contracts. The smart contracts analyses input sent from the API oracles, and emit events for both submitted price updates and when a discrepancy is found between APIs' submissions. Discrepancies are detected when price differences exceed a configurable threshold value.  

## Finchains Public UI

Finchains' front-end can be viewed at [https://finchains.io](https://finchains.io)

The Finchains WRKChain writes its block hashes to the [Mainchain MainNet](https://explorer.unification.io/).

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
