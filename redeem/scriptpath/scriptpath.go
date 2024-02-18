package scriptpath

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"taptx/config"
	"taptx/redeem"
	"taptx/tx"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

const (
	// bitcoin-cli -regtest sendtoaddress <script_address> 0.00010000
	// previous outpoint
	prevHashStr   = "5c8be30096ff25db8a11958498a9953f38cd5a231c1ece676429b687397544c6"
	prevIndex     = uint32(1)
	prevAmountSat = int64(10_000)

	// send address
	// P2TR: bitcoin-cli -regtest getnewaddress "" bech32
	sendAddrStr = "bcrt1qtxftdnsphctle6jv0salhumdnm0rpdyuld445c"
	feeSat      = int64(330)
)

var (
	// HTLC test data
	//  <<signature>>
	//  <<preimage>>
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

func createScript(pubkeyA []byte, pubkeyB []byte) ([]byte, error) {

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

	return script, nil
}

func ScriptPath(redeemType redeem.RedeemType) {
	if redeemType != redeem.RedeemP2wsh {
		return
	}

	script, err := createScript(
		keyAlice.PubKey.SerializeCompressed(),
		keyBob.PubKey.SerializeCompressed())
	if err != nil {
		fmt.Printf("fail createScript(): %v\n", err)
		return
	}

	// このアドレスへの送金をredeemする。
	// 送金したTXIDなどは prevHashStr に反映させる。
	sc := tx.NewScript(script, config.Network)
	addr, err := sc.CreateP2wsh()
	if err != nil {
		fmt.Printf("fail CreateP2wsh(): %v\n", err)
		return
	}
	// bcrt1qun9uqd06puhr2pc4cuxtt5q577musjympqf962p9u743zratcnlseqwmtl
	fmt.Printf("send to this script address= %s\n", addr)

	// redeem
	prevHash, _ := chainhash.NewHashFromStr(prevHashStr)
	tx, txid, err := sc.RedeemP2wshTx(
		prevHash, prevIndex, prevAmountSat,
		sendAddrStr, feeSat, preimage, keyAlice)
	if err != nil {
		fmt.Printf("fail RedeemP2wshTx(): %v\n", err)
		return
	}
	fmt.Printf("txid=%s\n", txid)
	fmt.Printf("tx= %x\n", tx)
}