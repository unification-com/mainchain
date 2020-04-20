# Demo WRKChain: Finchain

[Finchain](https://finchain-testnet.unification.io/) is a live WRKChain example currently running on TestNet. The full source code is available on [Github](https://github.com/unification-com/finchain-demo), and can be run as a completely self-contained localised demo.

[[toc]]

## Public Finchain Demo

The public Finchain Demo can be viewed here: [https://finchain.unification.io](https://finchain.unification.io)  
The public Finchain Block explorer can be found here: [https://finchain.unification.io/explorer](https://finchain.unification.io/explorer)

The public Finchain demo writes its block hashes to the [UND Mainchain TestNet](https://explorer-testnet.unification.io/).

## Running Finchain Locally

Finchain can be run locally as a completely self-contained demo, to allow developers to play with different configurations, and see how the internals of a WRKChain work.

Docker and Docker Compose are required to run the localised, self-contained
Finchain demo.

Copy `example.env` to `.env` and make any required changes. API keys are required
for the demo to work - see `example.env` for details on where to obtain the
necessary API keys.

Run the demo using:

```bash
make
```

### Docker network issues

By default, the demo uses the `172.25.0.0/24` subnet. If this subnet overlaps with your own, run:

```bash
ifconfig
```

and look for a line for your connection similar to:

```bash
inet addr:192.168.1.2  Bcast:192.168.1.255  Mask:255.255.255.0
```

Look for the `network` (first 3 parts of the IP address)
value of `inet addr`. This is your current subnet. In the example above, this is `192.168.1`.

Then set the `SUBNET_IP` variable to any other subnet. For example, run:

```bash
SUBNET_IP=192.168.5 make
```

or edit `.env` changing the value for `SUBNET_IP`

to run the demo on the `192.168.5.0/24` subnet

### WRKChain: Finchain Info

Finchain is a `geth` based WRKChain.

Network ID: `2339117895`  

Block Explorer: [http://localhost:8081](http://localhost:8081)

JSON RPC: [http://localhost:8547](http://localhost:8547)

WRKChain Block Validation UI: [http://localhost:4040](http://localhost:4040)


### Local UND Mainchain DevNet

The Finchain demo contains a completely self-contained pre-configured local Mainchain DevNet:

Chain ID: `UND-Mainchain-DevNet  `

Block Explorer: [http://localhost:3000](http://localhost:3000)

RPC: [http://localhost:26661](http://localhost:26661)
REST: [https://localhost:1318](https://localhost:1318)
