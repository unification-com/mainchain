# What is Mainchain?

Mainchain is the backbone of the Unification Network. It is a Tendermint based chain, and is where WRKChains and BEACONs submit their hashes, and FUND transactions take place.

>Mainchain is a public proof-of-stake chain. The core coin used for staking,
>rewards, and network fees is FUND. The on-chain denomination is `nund`
> (Nano Unification Denomination) which is 10^-9 FUND.

## Networks

There are currently two live/planned public Mainchain networks, and one private development network:

#### Testnet
Launched in Q4 2019, [UND-Mainchain-TestNet](https://github.com/unification-com/testnet) serves as the official Unification public Test network.  

#### MainNet
MainNet is the upcoming official Unification `MainNet` - this will be the live public Unification Main network.

#### DevNet
Additionally, the [Mainchain repository](https://github.com/unification-com/mainchain) comes with a full private [DevNet](local-devnet.md) - a completely self-contained network for development and testing Mainchain features.

## Software

The Mainchain suite comes with two binaries: `und` and `undcli`

- `und` - "Unification Daemon": the Mainchain server-side daemon, used to run a full node for Mainchain. Validators run this service to produce blocks. See [full command reference](und-commands.md)  
- `undcli` - "Unification Daemon Client": the command line interface for interacting with Mainchain nodes. It can also be used to run a light-client RPC node service. See [full command reference](undcli-commands.md)

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

Unification have also developed the following modules for Mainchain:

- `x/beacon`: BEACON hash timestamp submission logic
- `x/enterprise`: Handles purchasing, locking and unlocking of Enterprise FUND
- `x/mint`: A modified core Cosmos SDK module, handling inflation and Block Reward logic.
- `x/wrkchain`: WRKChain block hash submission handling

#### Next

[Installing](../software/installation.md) the Mainchain binaries
