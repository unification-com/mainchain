![Unification Logo](./.vuepress/public/assets/img/unification_logoblack.png)

# Unification Mainchain Documentation

Welcome to the documentation for Unification's Mainchain. These docs
cover how to build and install the two main applications `und` and `undcli`,
how to run a node on MainNet, TestNet and DevNet, and how to interact with the Mainchain network.

## 1. About Mainchain

- [What is Mainchain?](introduction/about-mainchain.md)
- [Native Coin Denomination: nund](introduction/denomination.md)
- [Total Supply Queries and Conversions](introduction/total-supply.md)
- [Fees and Gas](introduction/fees-and-gas.md)
- [Introduction to Genesis Params](introduction/genesis-settings.md)
- [Introduction to Delegators and Staking](introduction/delegators.md)
- [Introduction to Validators](introduction/validators.md)
- [FAQs](introduction/faqs.md)

## 2. Install & Use the Software

- [Installation](software/installation.md)
- [Accounts and Wallets](software/accounts-wallets.md)
- [Run a full node & join a Network](networks/join-network.md)
- [Run `und` as a daemon](software/run-und-as-service.md)

### Light Client & REST

- [Running a Light Client/REST server](software/light-client-rpc.md)

### Full Command References

- [`und` - the und server command reference](und_cmd/und.md)

### Full Config File References

- [`.und_mainchain/config/config.toml` Reference](software/und-mainchain-config-ref.md)
- [`.und_mainchain/config/app.toml` Reference](software/und-mainchain-app-config-ref.md)
- [`.und_mainchain/config/client.toml` Reference](software/und-mainchain-client-config-ref.md)

## 3. Mainchain Networks

- [Join a Network](networks/join-network.md)
- [Becoming a Validator](networks/become-validator.md)

### Private DevNet

- [Play with DevNet](networks/local-devnet.md)

## 4. Guides & Examples

### Tx & Query Examples
- [Sending Simple Transactions](examples/transactions.md)
- [Example WRKChain Transactions and Queries](examples/wrkchain.md)
- [Example BEACON Transactions and Queries](examples/beacon.md)
- [Example Enterprise FUND Transactions and Queries](examples/enterprise-fund.md)
- [WRKChain: Finchains](examples/finchain.md)

### In-depth guides

- [AWS 101: Introduction to installing `und` on AWS EC2 instances](guides/cloud/install-aws.md)
- [Google Cloud 101: Introduction to installing `und` on Google Cloud VMs](guides/cloud/install-gc.md)

## 5. Developer guides

- [Third Party Tool Development](developers/third-party.md)

### Disclaimer
Please note that this software is still in development. In these early days, we can expect to have issues, updates, and bugs. The existing `und` and `undcli` CLI tools require advanced technical skills and may involve risks which are outside of the control of the Unification Foundation and/or the Unification dev team. Any use of this open source [Apache 2.0 licensed](https://github.com/unification-com/mainchain/blob/master/LICENSE) software is done at your own risk and on a "AS IS" basis, without warranties or conditions of any kind, and any and all liability of the Unification Foundation and/or the Unification dev team for damages arising in connection to the software is excluded. Please exercise extreme caution!
