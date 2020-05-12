# Genesis Settings & Parameters

This document gives a brief explanation of each of the main parameters found in a network's `genesis.json`. Some parameters have been omitted for brevity.

With the exception of Chain Params, the following parameters are changeable on-chain via governance.

[[toc]]

## Chain Params

- `.genesis_time`  
Genesis block timestamp. E.g. 2020-01-06T12:00:00Z
- `.chain_id`  
ID of the current chain, E.g. `FUND-Mainchain-MainNet-v1`

## Auth Params

Auth parameters are changeable via governance

- `.app_state.auth.params.max_memo_characters`  
The maximum number of characters allowed in a Tx memo, e.g. 256
- `.app_state.auth.params.tx_sig_limit`  
The maximum number of permitted signatures on a `multisig` Tx. E.g. 7
- `.app_state.auth.params.tx_size_cost_per_byte`  
The gas cost per byte for a Tx. E.g. 10

## Beacon Params

Beacon parameters are changeable via governance

- `.app_state.beacon.params.denom`  
Denomination used for BEACON fees. E.g. `nund`
- `.app_state.beacon.params.fee_record`  
Fee for recording BEACON timestamps, in `params.denom`, E.g. 1000000000
- `.app_state.beacon.params.fee_register`  
One time fee for registering a BEACON with the network, in `params.denom`. E.g. 1000000000000

## Distribution Params

Distribution parameters are changeable via governance

- `.app_state.distribution.params.base_proposer_reward`  
Additional % of fees etc. given to the block proposer as a bonus. E.g. 0.010000000000000000
- `.app_state.distribution.params.bonus_proposer_reward`  
Additional % of fees etc. given to the block proposer, based on voting metrics calculated with each block. E.g. 0.040000000000000000
- `.app_state.distribution.params.community_tax`  
A fixed % of fees sent to the community pool each block. The spending of community pool coins can be decided via governance. E.g. 0.020000000000000000

## Enterprise Params

Enterprise parameters are changeable via governance

- `.app_state.enterprise.params.decision_time_limit`  
The time (in seconds) after which a raised Enterprise Purchase Order will automatically be rejected if no decision has been made. E.g. 259200
- `.app_state.enterprise.params.denom`  
Denomination in which Enterprise FUND will be issued. E.g. nund
- `.app_state.enterprise.params.ent_signers`  
Comma separated list of wallet addresses authorised to accept/reject Enterprise FUND Purchase Orders.
- `.app_state.enterprise.params.min_Accepts`  
Minimum number of `params.ent_signers` required to make a decision on a Purchase Order.

## Governance Params

Governance parameters are changeable via governance

- `.app_state.gov.deposit_params.max_deposit_period`  
Time in ns in which deposits are required to be made for a proposal. E.g. 1209600000000000
`.app_state.gov.deposit_params.min_deposit`  
Minumum deposit required in order for a proposal to be valid, and enter the voting period. E.g. 1000000000000nund
- `.app_state.gov.tally_params.quorum`  
Minimum percentage of total stake needed to vote for a result to be considered valid. E.g. 0.400000000000000000
- `.app_state.gov.tally_params.threshold`  
Minimum proportion of Yes votes for proposal to pass. E.g. 0.500000000000000000
- `.app_state.gov.tally_params.veto`  
Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. E.g. 0.334000000000000000
- `.app_state.gov.voting_params.voting_period `
Length of the voting period in ns. E.g. 1209600000000000


## Slashing Params

Slashing parameters are changeable via governance

- `.app_state.slashing.params.downtime_jail_duration`  
Cooldown period in ns, after being jailed for downtime that a node cannot unjail. E.g. 600000000000
- `.app_state.slashing.params.min_signed_per_window`  
min % of SignedBlocksWindow that must be signed in order to keep EV in pool of active EVs. Below this = jailed. E.g. 0.050000000000000000
- `.app_state.slashing.params.signed_blocks_window`  
number of blocks to monitor for missed/double sign etc. E.g. 10000
- `.app_state.slashing.params.slash_fraction_double_sign`  
% of stake to slash for double signing. E.g. 0.050000000000000000
- `.app_state.slashing.params.slash_fraction_downtime`  
% of stakes to slash for downtime
. E.g. 0.000100000000000000

## Staking Params

Staking parameters are changeable via governance

- `.app_state.staking.params.bond_denom`  
Staking denomination. E.g. nund
- `.app_state.staking.params.historical_entries`  
The n most recent historical entries are persisted. E.g. 3
- `.app_state.staking.params.max_entries`  
max simultaneous entries for either unbonding delegation or redelegation (per pair/trio). E.g. 7
- `.app_state.staking.params.max_validators`  
Maximum number of active validators in the current pool. E.g. 96
- `.app_state.staking.params.unbonding_time`  
duration in ns unbonding takes to complete. E.g. 1814400000000000

## WRKChain Params

WRKChain parameters are changeable via governance

- `.app_state.wrkchain.params.denom`  
Denomination used for WRKChain fees. E.g. `nund`
- `.app_state.wrkchain.params.fee_record`  
Fee for recording WRKChain hashes, in `params.denom`, E.g. 1000000000
- `.app_state.wrkchain.params.fee_register`  
One time fee for registering a WRKChain with the network, in `params.denom`. E.g. 1000000000000
