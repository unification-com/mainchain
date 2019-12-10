# Deploying a Local Devnet

The repository contains a ready to deploy Docker composition for local
development and testing. The devnet can be run with the following command:

```bash
make devnet
```

To bring the local devnet down cleanly, use <kbd>Ctrl</kbd>+<kbd>C</kbd>, followed by:

```bash
make devnet-down
```

## Devnet Chain ID

**Important**: Devnet's Chain ID is `UND-Mainchain-DevNet`. Any `und` or `undcli` commands
intended for Devnet should use the flag `--chain-id UND-Mainchain-DevNet`

## Devnet Nodes

The devnet composition will spin up three full nodes, and one light client in the following
Docker containers:

- `node1` - Full validation node, on 172.25.0.3:26661
- `node2` - Full validation node, on 172.25.0.4:26662
- `node3` - Full validation node, on 172.25.0.5:26663
- `rest-server` - Light Client for RPC interaction on 172.25.0.4:1317

## Devnet wallets and keys

See [Docker README](../Docker/README.md) for the mnemonic phrases and keys used
by the above nodes, and for test accounts included in Devnet's genesis.

### Importing the Devnet keys

The Devnet accounts can be imported as follows. First, build the `und` and 
`undcli` binaries:

```bash
make clean && make build
```

Then, for each account run the following command:

```bash
./build/undcli keys add node1 --recover
```

You will be asked for a password, and to enter the mnemonic phrase itself.
Change `node1` to an appropriate moniker for each imported account.
