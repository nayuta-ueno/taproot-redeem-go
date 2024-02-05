// 1. Get address from Bitcoin Core.
// 		P2WPKH: bitcoin-cli -regtest getnewaddress
// 		P2TR: bitcoin-cli -regtest getnewaddress "" bech32m
// 2. Update `sendAddrStr` by previous address.
// 3. Update `privKeyStr` as you like.
// 3. execute "go run .", and get "prev address".
// 4. Send bitcoin to "prev address".
//		bitcoin-cli -regtest -named sendtoaddress address=<prev address> amount=0.1 fee_rate=1
// 5. Get transaction information from Bitcoin Core.
//		bitcoin-cli -regtest gettransaction <previous txid>
// 6. Update `prevHashStr` and `prevIndex` from "gettransaction" result.
// 7. execute "go run .", and get "raw tx".
// 8. Sned raw transaction.
//		bitcoin-cli -regtest sendrawtransaction <raw tx>

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"taptx/tx"
)

const (
	privKeyStr = "112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00"

	// previous outpoint
	prevHashStr   = "994e2da234734d14ec61eb95d3076d82ef2b660c026fc0f6378e585cbd3a51bc"
	prevIndex     = uint32(1)
	prevAmountSat = int64(1000_0000)
	feeSat        = int64(200)

	// send address
	// P2WPKH: bitcoin-cli -regtest getnewaddress
	// P2TR: bitcoin-cli -regtest getnewaddress "" bech32m
	sendAddrStr = "bcrt1pypjucsfaqlfga7kxal0gfttpd95c8pe3vdexrgxjp5fh606mf09s7gvluq"
)

func main() {
	privKey, _ := hex.DecodeString(privKeyStr)
	key := tx.NewKey(privKey, &chaincfg.RegressionNetParams)

	fmt.Printf("sendAddrStr: %s\n", sendAddrStr)

	// ToDo
	// rawTx, txid, err := prevP2pkh(key)
	rawTx, txid, err := prevP2tr(key)

	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return
	}
	fmt.Printf("raw tx: %x\n", rawTx)
	fmt.Printf("txid: %s\n", txid)
}

func prevP2pkh(key *tx.Key) ([]byte, string, error) {
	p2pkh, _ := key.CreateP2WPKH()
	fmt.Printf("prev address: %s\n", p2pkh)

	prevHash, _ := chainhash.NewHashFromStr(prevHashStr)
	rawTx, txid, err := key.CreateRawTxP2PKH(prevHash, prevIndex, prevAmountSat, sendAddrStr, feeSat)
	if err != nil {
		return nil, "", fmt.Errorf("fail CreateRawTxP2PKH: %w", err)
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
