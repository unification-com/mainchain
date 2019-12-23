# Sending Simple Transactions

The `undcli` CMD can be used to generate, sign and broadcast new transactions
to the network. It can also be used to query transactions, accounts and 
a variety of other network information.

## Sending a Transaction

In this example, we'll generate and sign a simple `send` transaction, which will 
send 1UND. IF you have followed the documentation so far, you should already
have the software installed, be running a full node, and have an account
with funds.

The `send` command is as follows:

```bash
undcli tx send [from_key_or_address] [to_address] [amount] --chain-id [chain_id] --gas=auto --gas-adjustment=1.5 --gas-prices=0.025nund
```

- `[from_key_or_address]` - this can be either your account identifier, or your `bech32` address
- `[to_address]` - the `bech32` address of the account you are sending UND to
- `[amount]` - the amount, in `nund`
- `[chain_id]` - the ID of the chain to run the transaction on

If you do not have a local full node running, the `--node` flag can also
be passed to the `undcli` command to use a public node. For example, on
testnet, `--node tcp://3.136.43.0:26660` can be passed to send Txs
to our public TestNet node.

**Note**: For [TestNet](join-testnet.md), use `UND-Mainchain-TestNet` as the Chain ID.  
For [DevNet](local-devnet.md), use `UND-Mainchain-DevNet` as the Chain ID.

For example, we are running on TestNet, and would like to send 1 UND from
our account `und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed`, to our friend's
account `und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy`. We would
therefore run:

```bash
undcli tx send und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy 1000000000nund --chain-id UND-Mainchain-TestNet --gas=auto --gas-adjustment=1.5 --gas-prices=0.025nund
```

You will be prompted for confirmation, along with your password for the account.

If all goes well, the transaction will be broadcast and you should see a result
similar to the following:

```json
{
  "height": "0",
  "txhash": "0E0F6B2EFD2F0DF7593789F506E512B61B052515AEC3E26DD021B6020A8AF562",
  "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "success": true,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "send"
            }
          ]
        }
      ]
    }
  ]
}
```

## Query a Transaction

You can then query the transaction's progress and final result by running:

```bash
undcli query tx 0E0F6B2EFD2F0DF7593789F506E512B61B052515AEC3E26DD021B6020A8AF562 --chain-id UND-Mainchain-TestNet
```

The output should be similar to:

```json
{
  "height": "137",
  "txhash": "0E0F6B2EFD2F0DF7593789F506E512B61B052515AEC3E26DD021B6020A8AF562",
  "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"sender\",\"value\":\"und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed\"},{\"key\":\"module\",\"value\":\"bank\"},{\"key\":\"action\",\"value\":\"send\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy\"},{\"key\":\"amount\",\"value\":\"1000000000nund\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "success": true,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "sender",
              "value": "und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed"
            },
            {
              "key": "module",
              "value": "bank"
            },
            {
              "key": "action",
              "value": "send"
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
              "value": "1000000000nund"
            }
          ]
        }
      ]
    }
  ],
  "gas_wanted": "65000",
  "gas_used": "63607",
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
                "amount": "1000000000"
              }
            ]
          }
        }
      ],
      "fee": {
        "amount": [
          {
            "denom": "nund",
            "amount": "1950"
          }
        ],
        "gas": "65000"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A1qL4KCBiGgrE/PYIrUtpN08HxA7+Up+Q7eh3XNbCdSD"
          },
          "signature": "3NLXxU+tf1OkXE6jot70hepWk/DxPCoOMvlXC4BdQ61BG1XTnQho/WLUDURyhQ2IRaRxajMUh1GmZD35IKe7Bw=="
        }
      ],
      "memo": ""
    }
  },
  "timestamp": "2019-12-17T12:17:33Z",
  "events": [
    {
      "type": "message",
      "attributes": [
        {
          "key": "sender",
          "value": "und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed"
        },
        {
          "key": "module",
          "value": "bank"
        },
        {
          "key": "action",
          "value": "send"
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
          "value": "1000000000nund"
        }
      ]
    }
  ]
}
```

## Query an Account

Finally, to check that the funds have been sent and received, we can query the
account:

```bash
undcli query account und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy --chain-id UND-Mainchain-TestNet
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
          "amount": "1000000000"
        }
      ],
      "public_key": null,
      "account_number": "3",
      "sequence": "0"
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
        "amount": "1000000000"
      }
    ]
  }
}
```

## Next

Example transactions for [registering a WRKChain and submitting hashes](wrkchain.md)
