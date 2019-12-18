# Accounts and Wallets

The `undcli` CMD can be used to create new accounts, or import previous accounts
and keys.

## Create a New Account

**IMPORTANT**: When you create a new account, the CMD will output a mnemonic. KEEP THIS
SAFE - if you lose it, you will not be able to recover your account!

To create a new account, run:

```bash
undcli keys add [account_name]
```

`[account_name]` is whatever name you would like to use as an identifier

You will be prompted to enter a password to secure the account. Once entered,
the application will output your account details, including your address, public key
and mnemonic for recovery/future importing.

Accounts and keys are stored in `~/.und_cli/keys` by default.

## Import an Account

The same command can be used to import a previously saved mnemonic by passing
the `--recover` flag:

```bash
undcli keys add [account_name] --recover
```

As when creating a new account, you will be prompted to enter a password
to secure the account, followed by a prompt to enter your mnemonic.

## List accounts

You can list locally stored account with the following command:

```bash
undcli keys list
```

Run `undcli keys --help` for further details
