package tx

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg"
)

type Key struct {
	privKey *btcec.PrivateKey
	PubKey  *btcec.PublicKey
	Net     *chaincfg.Params
}

func NewKey(privKeyBytes []byte, net *chaincfg.Params) *Key {
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)
	return &Key{
		privKey: privKey,
		PubKey:  pubKey,
		Net:     net,
	}
}

func (k *Key) SerializeSchnorrPubKey() []byte {
	return schnorr.SerializePubKey(k.PubKey)
}
