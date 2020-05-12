# Sending Simple Transactions

The `undcli` CMD can be used to generate, sign and broadcast new transactions
to the network. It can also be used to query transactions, accounts and
a variety of other network information.

::: warning IMPORTANT
Whenever you use `undcli` to send Txs or query the chain ensure you pass the correct data to the `--chain-id` and if necessary `--node=` flags so that you connect to the correct network!
:::

#### Contents

[[toc]]

## Sending a Transaction

In this example, we'll generate and sign a simple `send` transaction, which will
send 1 FUND. IF you have followed the documentation so far, you should already
have the software installed, be running a full node, and have an account
with funds.

The `send` command is as follows:

```bash
undcli tx send [from_key_or_address] [to_address] [amount] --chain-id [chain_id] --node=tcp://[ip]:[port] --gas=auto --gas-adjustment=1.5 --gas-prices=0.25nund --trust-node false
```

- `[from_key_or_address]` - this can be either your account identifier, or your `bech32` address
- `[to_address]` - the `bech32` address of the account you are sending FUND to
- `[amount]` - the amount, in `nund`
- `[chain_id]` - the ID of the chain to run the transaction on
- `[ip]:[port]` - the IP and Port of the RPC node to broadcast the Tx

::: tip
If you are running your own full node, you can set the `--trust-node` flag to `true`, which will tell `undcli` not to verify the proofs form the response.
:::

For example, we are running on DevNet, and would like to send 1 FUND from
our account `und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed`, to our friend's
account `und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy`. We would
therefore run:

```bash
undcli tx send und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy 1000000000nund --chain-id FUND-Mainchain-DevNet --node=tcp://172.25.0.3:26661 --gas=auto --gas-adjustment=1.5 --gas-prices=0.25nund --trust-node=false
```

You will be prompted for confirmation, along with your password for the account.

If all goes well, the transaction will be broadcast and you should see a result
similar to the following:

```json
{
  "height": "0",
  "txhash": "6FC93147D467E27C104BD68DADAC0CFD6AA130E37E8B29F6652570A891E38F71",
  "raw_log": "[]"
}

```

::: tip
you can set the `--broadcast-mode` flag in the command to `block`. This will tell `undcli` to wait for the transaction to be processed in a block before returning the result. This will take up to 5-6 seconds to complete, but the Tx result will be included in the output.
:::

## Query a Transaction

You can then query the transaction's progress and final result by running:

```bash
undcli query tx 6FC93147D467E27C104BD68DADAC0CFD6AA130E37E8B29F6652570A891E38F71 --chain-id FUND-Mainchain-DevNet
```

The output should be similar to:

```json
{
  "height": "7",
  "txhash": "6FC93147D467E27C104BD68DADAC0CFD6AA130E37E8B29F6652570A891E38F71",
  "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy\"},{\"key\":\"amount\",\"value\":\"100000000000nund\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "send"
            },
            {
              "key": "sender",
              "value": "und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed"
            },
            {
              "key": "module",
              "value": "bank"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy"
            },
            {
              "key": "amount",
              "value": "100000000000nund"
            }
          ]
        }
      ]
    }
  ],
  "gas_wanted": "75420",
  "gas_used": "63558",
  "tx": {
    "type": "cosmos-sdk/StdTx",
    "value": {
      "msg": [
        {
          "type": "cosmos-sdk/MsgSend",
          "value": {
            "from_address": "und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed",
            "to_address": "und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy",
            "amount": [
              {
                "denom": "nund",
                "amount": "100000000000"
              }
            ]
          }
        }
      ],
      "fee": {
        "amount": [
          {
            "denom": "nund",
            "amount": "1886"
          }
        ],
        "gas": "75420"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A1qL4KCBiGgrE/PYIrUtpN08HxA7+Up+Q7eh3XNbCdSD"
          },
          "signature": "VLldGBkI0C3xcqwGShR2ImIc76btDGtW7QlEVfeDHuZtONIHDR5Ckf87wROazxqVw3rM35RvPgTyoj8VkVFV4w=="
        }
      ],
      "memo": ""
    }
  },
  "timestamp": "2020-04-01T14:28:51Z"
}

```

## Query an Account

Finally, to check that the funds have been sent and received, we can query the
account:

```bash
undcli query account und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy --chain-id FUND-Mainchain-DevNet
```

Which will output a result similar to:

```json
{
  "account": {
    "type": "cosmos-sdk/Account",
    "value": {
      "address": "und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy",
      "coins": [
        {
          "denom": "nund",
          "amount": "100400000000000"
        }
      ],
      "public_key": "",
      "account_number": 3,
      "sequence": 0
    }
  },
  "enterprise": {
    "locked": {
      "denom": "nund",
      "amount": "0"
    },
    "available_for_wrkchain": [
      {
        "denom": "nund",
        "amount": "100400000000000"
      }
    ]
  }
}
```

#### Next

Example transactions for [registering a WRKChain and submitting hashes](wrkchain.md)
