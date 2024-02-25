package tx

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// P2TR
// privkey -> pubkey -> shnorr-pubkey => witness program
func (k *Key) CreateP2TR() (string, error) {
	addr, err := k.createP2tr()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (k *Key) createP2tr() (*btcutil.AddressTaproot, error) {
	pubKey := txscript.ComputeTaprootKeyNoScript(k.PubKey)
	witnessProg := schnorr.SerializePubKey(pubKey)
	return btcutil.NewAddressTaproot(witnessProg, k.Net)
}

func (k *Key) CreateRawTxP2TR(
	prevHash *chainhash.Hash,
	prevIndex uint32,
	prevAmountSat int64,
	sendAddrStr string,
	feeSat int64,
) ([]byte, string, error) {
	originTx := wire.NewMsgTx(2)

	prevAddr, err := k.createP2tr()
	if err != nil {
		return nil, "", fmt.Errorf("fail createP2tr(prevAddr): %w", err)
	}
	prevPkScript, err := txscript.PayToAddrScript(prevAddr)
	if err != nil {
		return nil, "", fmt.Errorf("fail PayToAddrScript(prevAddr): %w", err)
	}
	prevOutputFetcher := txscript.NewCannedPrevOutputFetcher(prevPkScript, prevAmountSat)

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

	prevInIndex := int(0)
	sigHashes := txscript.NewTxSigHashes(originTx, prevOutputFetcher)
	witness, err := txscript.TaprootWitnessSignature(
		originTx,
		sigHashes,
		prevInIndex,
		prevAmountSat,
		prevPkScript,
		txscript.SigHashDefault,
		k.privKey,
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
