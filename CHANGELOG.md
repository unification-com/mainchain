# Changelog

## v1.2.0

- Updated Tendermint Core to v0.33.0 [PR 98](https://github.com/unification-com/mainchain/pull/98)
- Updated Cosmos SDK to v0.38.0 [PR 98](https://github.com/unification-com/mainchain/pull/98)
- Updated Keyring [PR 98](https://github.com/unification-com/mainchain/pull/98)  
**Note**: Keyring update affects how keys are stored
by `undcli`. Keys created/imported before v1.2.0 will need to be migrated. 
Additionally, the `undcli keys` command has a new flag 
`--keyring-backend os|file|test`, which defaults to `os` to use 
the operating system's keyring.  
**IMPORTANT**: Since the `os` backend generally calls a GUI prompt, 
headless systems will need to pass `file` to `--keyring-backend` in order
to use the command line password prompt for all `undcli` commands, including
re-importing keys.  
Migration can be run via: `undcli keys migrate`
- Update DevNet config & genesis: [PR 98](https://github.com/unification-com/mainchain/pull/98)
modified `genesis.json` and `config.toml` for DevNet EV nodes

## v1.1.0 - 23/01/2020

- Beacon module: remove debug output msg during `CheckTx`
[PR 88](https://github.com/unification-com/mainchain/pull/88)
- Docs: Update to include TestNet info and WRKChain examples 
[PR 89](https://github.com/unification-com/mainchain/pull/89)
- All modules: temporarily disable bulk front-end queries 
[PR 90](https://github.com/unification-com/mainchain/pull/90)
- WRKChain module: track number of blocks submitted & fix
genesis import logger
[PR 91](https://github.com/unification-com/mainchain/pull/91)
- DevNet Docker: add reverse proxy to Docker composition for REST
server to allow CORS policy [PR 92](https://github.com/unification-com/mainchain/pull/92)
- Mint module: set DevNet inflation to 1% 
[PR 93](https://github.com/unification-com/mainchain/pull/93)
- DevNet Genesis: set unbond time to 3 minutes for DevNet 
[PR 96](https://github.com/unification-com/mainchain/pull/96)

## v1.0 - 20/12/2019

- Enterprise Module: [PR 85](https://github.com/unification-com/mainchain/pull/85)
 modify `ent_signers` param to accept multiple addresses as an array
- Docs: initial documentation [PR 74](https://github.com/unification-com/mainchain/pull/74)
- DevNet Genesis: [PR 86](https://github.com/unification-com/mainchain/pull/85)
`ent_signers` accounts need funds
- DevNet Genesis: [PR 87](https://github.com/unification-com/mainchain/pull/87)
add more test accounts

## v0.9 - 13/12/2019

- Initial release PRs 1 - 73