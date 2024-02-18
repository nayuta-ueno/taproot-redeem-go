package keypath

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"taptx/config"
	"taptx/redeem"
	"taptx/tx"
)

const (
	privKeyStr = "112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00"

	// previous outpoint
	prevHashStr   = "994e2da234734d14ec61eb95d3076d82ef2b660c026fc0f6378e585cbd3a51bc"
	prevIndex     = uint32(1)
	prevAmountSat = int64(10_000_000)
	feeSat        = int64(200)

	// send address
	// P2WPKH: bitcoin-cli -regtest getnewaddress
	// P2TR: bitcoin-cli -regtest getnewaddress "" bech32m
	sendAddrStr = "bcrt1pypjucsfaqlfga7kxal0gfttpd95c8pe3vdexrgxjp5fh606mf09s7gvluq"
)

func KeyPath(redeemType redeem.RedeemType) {
	privKey, _ := hex.DecodeString(privKeyStr)
	key := tx.NewKey(privKey, config.Network)

	fmt.Printf("sendAddrStr: %s\n", sendAddrStr)

	var rawTx []byte
	var txid string
	var err error
	switch redeemType {
	case redeem.RedeemP2wpkh:
		rawTx, txid, err = prevP2pkh(key)
	case redeem.RedeemP2trKeyPath:
		rawTx, txid, err = prevP2tr(key)
	default:
		fmt.Printf("invalid redeemType\n")
		return
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return
	}
	fmt.Printf("raw tx: %x\n", rawTx)
	fmt.Printf("txid: %s\n", txid)
}

func prevP2pkh(key *tx.Key) ([]byte, string, error) {
	p2pkh, _ := key.CreateP2wpkh()
	fmt.Printf("prev address: %s\n", p2pkh)

	prevHash, _ := chainhash.NewHashFromStr(prevHashStr)
	rawTx, txid, err := key.RedeemP2wpkhTx(prevHash, prevIndex, prevAmountSat, sendAddrStr, feeSat)
	if err != nil {
		return nil, "", fmt.Errorf("fail CreateRawTxP2WPKH: %w", err)
	}
	return rawTx, txid, nil
}

func prevP2tr(key *tx.Key) ([]byte, string, error) {
	p2tr, _ := key.CreateP2TR()
	fmt.Printf("prev address: %s\n", p2tr)

	prevHash, _ := chainhash.NewHashFromStr(prevHashStr)
	rawTx, txid, err := key.CreateRawTxP2TR(prevHash, prevIndex, prevAmountSat, sendAddrStr, feeSat)
	if err != nil {
		return nil, "", fmt.Errorf("fail CreateRawTxP2TR: %w", err)
	}
	return rawTx, txid, nil
}
