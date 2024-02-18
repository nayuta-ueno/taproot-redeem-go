package tx

import "github.com/btcsuite/btcd/chaincfg"

type Script struct {
	Script []byte
	Net    *chaincfg.Params
}

func NewScript(script []byte, net *chaincfg.Params) *Script {
	return &Script{
		Script: script,
		Net:    net,
	}
}
