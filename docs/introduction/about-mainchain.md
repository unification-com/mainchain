# What is Mainchain?

Mainchain is the backbone of the Unification Network. It is a Tendermint based chain, and is where WRKChains and BEACONs submit their hashes, and FUND transactions take place.

>Mainchain is a public proof-of-stake chain. The core coin used for staking,
>rewards, and network fees is FUND. The on-chain denomination is `nund`
> (Nano Unification Denomination) which is 10^-9 FUND.

## Networks

There are currently two live public Mainchain networks, and one private development network:

#### MainNet
Launched on 14/05/2020, [FUND-Mainchain-MainNet](https://github.com/unification-com/mainnet) is the live public Unification **Main Network**.

#### Testnet
Launched in Q4 2019, [FUND-Mainchain-TestNet](https://github.com/unification-com/testnet) serves as the official Unification public **Test network**, where developers can test WRKChains/BEACONs or third party application ssuch as wallets and block explorers before deploying on MainNet. It also serves as a test platform for Mainchain developers to test updates and new features to the Mainchain code in a live environment.  

#### DevNet
Additionally, the [Mainchain repository](https://github.com/unification-com/mainchain) comes with a full private [DevNet](local-devnet.md) - a completely self-contained network for development and testing Mainchain features.

## Software

The Mainchain suite comes with a unified binary: `und`. This binary is used for running the server-side daemon (e.g.
for Validator nodes, seeds, sentries, RPCs etc.), and as a client to interact with the network (e.g. for sending queries
and generating and broadcasting transactions).

Mainchain has been built with the following core Cosmos SDK modules:

- `x/auth`: Accounts and signatures.
- `x/bank`: Token transfers.
- `x/staking`: Staking logic.
- `x/distribution`: Fee distribution logic.
- `x/slashing`: Slashing logic.
- `x/gov`: on-chain governance logic.
- `x/supply`: FUND Coin supply logic
- `x/params`: Handles module-level parameters, which can be modified via governance.
- `x/crisis`: Handles potential network errors during the early days of deployment
- `x/upgrade`: Handles processing software upgrades and associated state migrations

Unification have also developed the following modules for Mainchain:

- `x/beacon`: BEACON hash timestamp submission logic
- `x/enterprise`: Handles purchasing, locking and unlocking of Enterprise FUND
- `x/wrkchain`: WRKChain block hash submission handling

#### Next

[Installing](../software/installation.md) the Mainchain binaries
