## und query wrkchain search

Query all WRKChains with optional filters

### Synopsis

Query for all paginated WRKChains that match optional filters:

Example:
$ und query wrkchain search --moniker wrkchain1
$ und query wrkchain search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
$ und query wrkchain search --page=2 --limit=100

```
und query wrkchain search [flags]
```

### Options

```
      --height int       Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help             help for search
      --moniker string   (optional) filter wrkchains by name
      --node string      <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string    Output format (text|json) (default "text")
      --owner string     (optional) filter wrkchains by owner address
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID
```

### SEE ALSO

* [und query wrkchain](und_query_wrkchain.md)	 - Querying commands for the wrkchain module

###### Auto generated by spf13/cobra on 28-Feb-2022