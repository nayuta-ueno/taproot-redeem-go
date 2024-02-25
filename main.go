package main

import (
	"taptx/p2trkey"
	"taptx/p2wpkh"
	"taptx/p2wsh"
	"taptx/tx"
)

func main() {
	// rt := tx.RedeemP2wpkh
	// rt := tx.RedeemP2trKeyPath
	rt := tx.RedeemP2wsh

	switch rt {
	case tx.RedeemP2wpkh:
		p2wpkh.P2wpkh()
	case tx.RedeemP2trKeyPath:
		p2trkey.KeyPath()
	case tx.RedeemP2wsh:
		p2wsh.P2wsh()
	}
}
