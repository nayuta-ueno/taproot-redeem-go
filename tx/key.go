package tx

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
)

type Key struct {
	privKey *btcec.PrivateKey
	pubKey  *btcec.PublicKey
	net     *chaincfg.Params
}

func NewKey(privKeyBytes []byte, net *chaincfg.Params) *Key {
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)
	return &Key{
		privKey: privKey,
		pubKey:  pubKey,
		net:     net,
	}
}
