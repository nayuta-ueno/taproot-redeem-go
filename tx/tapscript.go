package tx

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

type TapScript struct {
	Scripts []txscript.TapLeaf
	Net     *chaincfg.Params
}

func NewTapScript(scripts [][]byte, net *chaincfg.Params) *TapScript {
	taps := make([]txscript.TapLeaf, len(scripts))
	for i, v := range scripts {
		taps[i] = txscript.NewBaseTapLeaf(v)
	}
	return &TapScript{
		Scripts: taps,
		Net:     net,
	}
}
