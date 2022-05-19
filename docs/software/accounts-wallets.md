# Accounts and Wallets

The `und` CMD can be used to create new accounts, or import previous accounts
and keys.

#### Contents

[[toc]]

## Prerequisites

Before continuing, ensure you have gone through the following docs:

1. [Installing the software](installation.md)

## Create a new account

::: danger IMPORTANT
When you create a new account, the CMD will output a mnemonic. KEEP THIS SAFE - if you lose it, you will not be able to 
recover your account!
:::

To create a new account, run:

```bash
und keys add [account_name]
```

`[account_name]` is whatever name you would like to use as an identifier when signing transactions. For example:

```bash
und keys add my_new_wallet
```

::: tip
stick to alphanumeric characters, hyphens, underscores and full stops.
:::

If your OS keyring for `und` is not already unlocked/created, you will be prompted for a password. Once entered, 
the application will output your account details, including your wallet address, public key and mnemonic for 
recovery/future importing.

Accounts and keys are stored in your OS keyring by default.

## Import an account

The same command can be used to import a previously saved mnemonic by passing
the `--recover` flag:

```bash
und keys add [account_name] --recover
```

You will be prompted to enter your mnemonic.

## List & show accounts

You can list locally stored account with the following command:

```bash
und keys list
```

or show the details of a particular account:

```bash
und keys show my_new_wallet
```

Run `und keys --help` for further details
