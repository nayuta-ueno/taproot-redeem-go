package tx

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

type TapScript struct {
	Key     *Key
	Scripts []txscript.TapLeaf
	Net     *chaincfg.Params
}

func NewTapScript(internalKey *Key, scripts [][]byte, net *chaincfg.Params) *TapScript {
	taps := make([]txscript.TapLeaf, len(scripts))
	for i, v := range scripts {
		taps[i] = txscript.NewBaseTapLeaf(v)
	}
	return &TapScript{
		Key:     internalKey,
		Scripts: taps,
		Net:     net,
	}
}
