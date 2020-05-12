# Total Supply

The command

```bash
undcli query supply
```

Will return the complete supply information.

The equivalent REST query is at the enpoint `/supply/total` - for example on the public TestNet REST server [https://rest-testnet.unification.io/supply/total](https://rest-testnet.unification.io/supply/total)


Three quantity values are returned, all representing `nund`:

1. **amount**: Liquid FUND in active circulation, and the actual circulating total supply which is available and can be used for FUND transfers, staking, Tx fees etc. It is the **locked** amount subtracted from **total**. _This is the important value when processing any calculations dependent on FUND circulation/total supply of FUND etc._
2. **locked**: Total FUND locked through Enterprise purchases. This FUND is only available specifically to pay WRKChain/BEACON fees and **cannot** be used for transfers, staking/delegation or any other transactions. _Locked FUND only enters the active circulation supply once it has been used to pay for WRKChain/BEACON fees. Until then, it is considered "dormant", and not part of the circulating total supply_
3. **total**: The total amount of FUND currently known on the chain, including any Enterprise **locked** FUND. This is for informational purposes only and should not be used for any "circulating/total supply" calculations.

The **amount** value is the important value regarding total supply _currently in active circulation_, and is the information that should be used to represent any "total supply/circulation" values for example in block explorers, wallets, exchanges etc.

Consider the following `undcli query supply` result:

```json
{
  "denom": "nund",
  "amount": "120010263000000000",
  "locked": "89737000000000",
  "total": "120100000000000000"
}
```

Or, the equivalent REST query result:

```json
{
  "height":"84695",
  "result":{
    "denom":"nund",
    "amount":"120010263000000000",
    "locked":"89737000000000",
    "total":"120100000000000000"
    }
  }
```

In the above example, the active circulating supply - usable for transfers and standard transactions etc. - is currently 120,010,263 FUND. 89,737 FUND is currently locked, and can only be used for paying for WRKChain/BEACON fees - it is "dormant" and _cannot be used for any other purpose until it has been used to pay for WRKChain/BEACON fees_. Finally, the total amount of FUND known on the chain is 120,100,000 FUND, and is the equivalent of 120,010,263 + 89,737.

::: tip Note
The REST endpoint `/supply/total/nund` will return only the appropriate **amount** value, for example https://rest-testnet.unification.io/supply/total/nund would just return

```json
{
  "height":"84746",
  "result":"120010263000000000"
}
```
:::

## Converting to FUND

In much the same way that Ethereum uses `wei` and Cosmos uses `uatom` as the smallest on-chain denomination, all results for Unification Mainchain return the native on-chain coin denomination values in `nund`. **1,000,000,000 nund == 1 FUND**. As such, simply dividing the result by 1000000000 will yield the FUND value.

See [denomination](denomination.md) for further information.

See [undcli-query-supply](../software/undcli-commands.md#undcli-query-supply) for more details on command flags and parameters.
