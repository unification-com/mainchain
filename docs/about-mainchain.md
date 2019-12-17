# What is Mainchain?

Mainchain is the backbone of the Unification Network, and has been built using
the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk).

_**Note**: Mainchain has been built using the Cosmos SDK, but is 
completely independent network, and not part of the Cosmos Hub (or IBC) 
ecosystem._

Mainchain is where WRKChains and BEACONs submit their hashes, and in future
iterations of the network, will host UVCID and HAIKUs.

The Mainchain suite comes with two binaries: `und` and `undcli`

- `und` - The Mainchain daemon, used to run a full node for Mainchain  
- `undcli` - The command line interface for interacting with Mainchain full nodes.  
It can also be used to run a light-client RPC node service.

>Mainchain is a public proof-of-stake chain. The core coin used for staking,
>rewards, and network fees is UND. The on-chain denomination is `nund`
> (nano UND) which is 10^-9 UND.

Mainchain has been built with the following core Cosmos SDK modules:

- `x/auth`: Accounts and signatures.
- `x/bank`: Token transfers.
- `x/staking`: Staking logic.
- `x/mint`: Inflation and Block Reward logic.
- `x/distribution`: Fee distribution logic.
- `x/slashing`: Slashing logic.
- `x/supply`: UND Coin supply logic
- `x/params`: Handles app-level parameters.
- `x/crisis`: Handles potential network errors during the early days of deployment

Unification have also developed the following modules for Mainchain:

- `x/beacon`: BEACON hash timestamp submission logic
- `x/enterprise`: Handles purchasing, locking and unlocking of Enterprise UND
- `x/wrkchain`: WRKChain block hash submission handling

## Next

[Installing](installation.md) the Mainchain binaries
