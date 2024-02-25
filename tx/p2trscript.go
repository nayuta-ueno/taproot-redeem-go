package tx

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

// P2TR script
func (s *TapScript) CreateP2tr(internalKey *btcec.PublicKey) (string, error) {
	addr, err := s.createP2tr(internalKey)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (s *TapScript) createP2tr(internalKey *btcec.PublicKey) (*btcutil.AddressTaproot, error) {
	tree := txscript.AssembleTaprootScriptTree(s.Scripts...)
	rootHash := tree.RootNode.TapHash()
	pubKey := txscript.ComputeTaprootOutputKey(internalKey, rootHash[:])
	witnessProg := schnorr.SerializePubKey(pubKey)
	return btcutil.NewAddressTaproot(witnessProg, s.Net)
}
