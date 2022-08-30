module singula.finance/netio

go 1.19

replace singula.finance/node => ../node

require singula.finance/node v0.0.0-00010101000000-000000000000


require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/pkg/errors v0.9.1
//TODO upgrade
//github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1 // indirect
//github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
// github.com/aws/aws-sdk-go v1.29.16
// github.com/davecgh/go-spew v1.1.1 // indirect
)
