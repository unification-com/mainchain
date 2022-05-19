# Native Coin Denomination `nund`

The currency used on Mainchain is FUND. However, the native on-chain coin denomination (on MainNet, TestNet and DevNet)
is **`nund`**, or "Nano Unification Denomination", such that **1,000,000,000 nund == 1 FUND**.

All transactions, fees and stakes are defined and paid for in `nund`. For example, if you need to send **1 FUND** to 
your friend, you will need to set your Tx to send `1000000000nund`.

The `undcli` CMD has a simple conversion utility to help convert any fees
and FUND transactions into `nund`, and vice-versa:

```bash
und convert 1000000000 nund fund
```

will result in:

```bash
1000000000nund = 1.000000000fund
```

Likewise,

```bash
und convert 10 fund nund
```

will result in:

```bash
10fund = 10000000000nund
```

## HD Wallet Path

The BIP-0044 Path for our HD Wallets is as follows:

`44'/5555'/0'/0`   

SLIP-0044 Coin ID is `5555`
