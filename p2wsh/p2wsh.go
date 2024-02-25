// 1. Get address from Bitcoin Core and update `sendAddrStr`.
// 	    bitcoin-cli -regtest getnewaddress "" bech32m
// 2. Update `privkeyHexAlice` and `privkeyHexBob` as you like.
// 3. execute "go run .", and get "script address".
// 4. Send bitcoin to "script address".
//		bitcoin-cli -regtest -named sendtoaddress address=<script address> amount=0.0001 fee_rate=1
// 5. Get transaction information from Bitcoin Core.
//		bitcoin-cli -regtest gettransaction <previous txid>
// 6. Update `prevHashStr` and `prevIndex` from "gettransaction" result.
// 7. execute "go run .", and get "txid" and "raw tx".
// 8. Send raw transaction.
//		bitcoin-cli -regtest sendrawtransaction <raw tx>

package p2wsh

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"taptx/config"
	"taptx/tx"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

const (
	// bitcoin-cli -regtest sendtoaddress <script_address> 0.00010000
	// previous outpoint
	prevHashStr   = "5c8be30096ff25db8a11958498a9953f38cd5a231c1ece676429b687397544c6"
	prevIndex     = uint32(1)
	prevAmountSat = int64(10_000)

	// send address: bitcoin-cli -regtest getnewaddress "" bech32
	sendAddrStr = "bcrt1qtxftdnsphctle6jv0salhumdnm0rpdyuld445c"
	feeSat      = int64(330)
)

var (
	//  <<signature>>
	//  <<preimage>>
	//
	//  OP_SHA256 <payment_hash> OP_EQUAL
	//  OP_IF
	//     <alicePubkey>
	//  OP_ELSE
	//     <bobPubkey>
	//  OP_ENDIF
	//  OP_CHKSIG
	preimage, _ = hex.DecodeString("00112233445566778899aabbccddeeff")
	paymentHash = sha256.Sum256(preimage)

	privkeyHexAlice, _ = hex.DecodeString("00112233445566778899aabbccddee00")
	privkeyHexBob, _   = hex.DecodeString("00112233445566778899aabbccddee01")
	keyAlice           = tx.NewKey(privkeyHexAlice, config.Network)
	keyBob             = tx.NewKey(privkeyHexBob, config.Network)
)

func createScript(pubkeyA []byte, pubkeyB []byte) []byte {
	const (
		OP_IF     = 0x63
		OP_ELSE   = 0x67
		OP_ENDIF  = 0x68
		OP_DROP   = 0x75
		OP_EQUAL  = 0x87
		OP_SHA256 = 0xa8
		OP_CHKSIG = 0xac
		OP_CSV    = 0xb2
	)

	part1 := []byte{OP_SHA256, byte(len(paymentHash))}
	// paymentHash[:]
	part2 := []byte{OP_EQUAL, OP_IF, byte(len(pubkeyA))}
	// pubkeyA
	part3 := []byte{OP_ELSE, byte(len(pubkeyB))}
	// pubkeyB
	part4 := []byte{OP_ENDIF, OP_CHKSIG}
	script := make(
		[]byte,
		0,
		len(part1)+
			len(paymentHash)+
			len(part2)+
			len(pubkeyA)+
			len(part3)+
			len(pubkeyB)+
			len(part4))
	script = append(script, part1...)
	script = append(script, paymentHash[:]...)
	script = append(script, part2...)
	script = append(script, pubkeyA...)
	script = append(script, part3...)
	script = append(script, pubkeyB...)
	script = append(script, part4...)

	return script
}

func P2wsh() {
	script := createScript(
		keyAlice.PubKey.SerializeCompressed(),
		keyBob.PubKey.SerializeCompressed())
	sc := tx.NewScript(script, config.Network)

	addr, err := sc.CreateP2wsh()
	if err != nil {
		fmt.Printf("fail CreateP2wsh(): %v\n", err)
		return
	}
	fmt.Printf("send to this script address= %s\n\n", addr)

	// redeem
	prevHash, _ := chainhash.NewHashFromStr(prevHashStr)
	tx, txid, err := sc.RedeemP2wshTx(
		// previous output
		prevHash, prevIndex, prevAmountSat,
		// current output
		sendAddrStr, feeSat,
		// unlock
		preimage, keyAlice,
	)
	if err != nil {
		fmt.Printf("fail RedeemP2wshTx(): %v\n", err)
		return
	}
	fmt.Printf("txid=%s\n", txid)
	fmt.Printf("tx= %x\n", tx)
}
