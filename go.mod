module github.com/unification-com/mainchain

go 1.13

require (
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	github.com/ChainSafe/go-schnorrkel v0.0.0-20200115165343-aa45d48b5ed6 // indirect
	github.com/btcsuite/btcutil v1.0.1 // indirect
	github.com/cosmos/cosmos-sdk v0.38.3
	github.com/gorilla/mux v1.7.3
	github.com/gtank/merlin v0.1.1 // indirect
	github.com/rakyll/statik v0.1.6
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.3
	github.com/tendermint/tm-db v0.5.1
	golang.org/x/crypto v0.0.0-20200128174031-69ecbb4d6d5d // indirect
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	google.golang.org/genproto v0.0.0-20200128133413-58ce757ed39b // indirect
	gopkg.in/ini.v1 v1.51.1 // indirect
)

replace github.com/cosmos/ledger-cosmos-go => github.com/unification-com/ledger-unification-go v0.11.2
