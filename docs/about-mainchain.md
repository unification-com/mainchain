# What is Mainchain?

Mainchain is the backbone of the Unification Network. It is a Tendermint based chain, and is where WRKChains and BEACONs submit their hashes, and UND transactions take place.

>Mainchain is a public proof-of-stake chain. The core coin used for staking,
>rewards, and network fees is UND. The on-chain denomination is `nund`
> (nano UND) which is 10^-9 UND.

There are currently two live/planned public Mainchain networks:

**TestNet**: Launched in Q1 2020, `UND-Mainchain-TestNet-v3` serves as the official Unification public Test network.
**MainNet**: MainNet is the upcoming official Unification `MainNet` - this will be the live public Unification Main network.

Additionally, this repository comes with a full private [DevNet](local-devnet.md) - a completely self-contained network for development and testing Mainchain features.

The Mainchain suite comes with two binaries: `und` and `undcli`

- `und` - The Mainchain server-side daemon, used to run a full node for Mainchain. Validators run this service to produce blocks.  
- `undcli` - The command line interface for interacting with Mainchain nodes. It can also be used to run a light-client RPC node service.

Mainchain has been built with the following core Cosmos SDK modules:

- `x/auth`: Accounts and signatures.
- `x/bank`: Token transfers.
- `x/staking`: Staking logic.
- `x/distribution`: Fee distribution logic.
- `x/slashing`: Slashing logic.
- `x/gov`: on-chain governance logic.
- `x/supply`: UND Coin supply logic
- `x/params`: Handles module-level parameters, which can be modified via governance.
- `x/crisis`: Handles potential network errors during the early days of deployment

Unification have also developed the following modules for Mainchain:

- `x/beacon`: BEACON hash timestamp submission logic
- `x/enterprise`: Handles purchasing, locking and unlocking of Enterprise UND
- `x/mint`: A modified core Cosmos SDK module, handling inflation and Block Reward logic.
- `x/wrkchain`: WRKChain block hash submission handling

#### Next

[Installing](installation.md) the Mainchain binaries
