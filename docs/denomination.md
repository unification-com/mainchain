# Native Coin Denomination `nund`

The native coin denomination on Mainchain (MainNet, TestNet and DevNet) is **`nund`**, or "Nano UND", such that 1,000,000,000 nund == 1 UND.

All transactions, fees and stakes are paid for in `nund`. For example, if you need to send 1 UND to your friend, you will need to set your Tx to send `1000000000nund`.

The `undcli` CMD has a simple conversion utility to help convert any fees
and UND transactions into `nund`, and vice-versa:

```bash
undcli convert 1000000000 nund und
```

will result in:

```bash
1000000000nund = 1.000000000und
```

Likewise,

```bash
undcli convert 10 und nund
```

will result in:

```bash
10und = 10000000000nund
```
