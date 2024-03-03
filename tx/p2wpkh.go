package tx

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// P2WPKH
// privkey -> pubkey
// hash160(pubkey) => witness program
func (k *Key) CreateP2wpkh() (string, error) {
	addr, err := k.createP2wpkh()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (k *Key) createP2wpkh() (*btcutil.AddressWitnessPubKeyHash, error) {
	witnessProg := btcutil.Hash160(k.PubKey.SerializeCompressed())
	return btcutil.NewAddressWitnessPubKeyHash(witnessProg, k.Net)
}

func (k *Key) RedeemP2wpkhTx(
	prevHash *chainhash.Hash,
	prevIndex uint32,
	prevAmountSat int64,
	sendAddrStr string,
	feeSat int64,
) ([]byte, string, error) {
	originTx := wire.NewMsgTx(2)
	// originTx.LockTime = 0xffffffff

	prevAddr, err := k.createP2wpkh()
	if err != nil {
		return nil, "", fmt.Errorf("fail createP2wpkh(prevAddr): %w", err)
	}
	prevPkScript, err := txscript.PayToAddrScript(prevAddr)
	if err != nil {
		return nil, "", fmt.Errorf("fail PayToAddrScript(prevAddr): %w", err)
	}
	prevOutputFetcher := txscript.NewCannedPrevOutputFetcher(prevPkScript, prevAmountSat)
	txinIndex := int(0)

	sendAddr, err := btcutil.DecodeAddress(sendAddrStr, k.Net)
	if err != nil {
		return nil, "", fmt.Errorf("fail DecodeAddress(sendAddr): %w", err)
	}
	sendPkScript, err := txscript.PayToAddrScript(sendAddr)
	if err != nil {
		return nil, "", fmt.Errorf("fail PayToAddrScript(sendAddr): %w", err)
	}
	sendAmountSat := prevAmountSat - feeSat
	txOut := wire.NewTxOut(sendAmountSat, sendPkScript)
	originTx.AddTxOut(txOut)

	prevOut := wire.NewOutPoint(prevHash, prevIndex)
	txIn := wire.NewTxIn(prevOut, nil, nil)
	originTx.AddTxIn(txIn)

	sigHashes := txscript.NewTxSigHashes(originTx, prevOutputFetcher)
	witness, err := txscript.WitnessSignature(
		originTx,
		sigHashes,
		txinIndex,
		prevAmountSat,
		prevPkScript,
		txscript.SigHashAll,
		k.privKey,
		true, // compress pubkey
	)
	if err != nil {
		return nil, "", fmt.Errorf("fail RawTxInWitnessSignature: %w", err)
	}
	txIn.Witness = witness

	var buf bytes.Buffer
	originTx.Serialize(&buf)
	txid := originTx.TxHash()
	return buf.Bytes(), txid.String(), nil
}
