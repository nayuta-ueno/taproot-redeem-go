package tx

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// P2WSH
// sha256(witness script) => witness program
func (s *Script) CreateP2wsh() (string, error) {
	addr, err := s.createP2wsh()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (s *Script) createP2wsh() (*btcutil.AddressWitnessScriptHash, error) {
	witnessProg := sha256.Sum256(s.Script)
	return btcutil.NewAddressWitnessScriptHash(witnessProg[:], s.Net)
}

func (s *Script) RedeemP2wshTx(
	prevHash *chainhash.Hash,
	prevIndex uint32,
	prevAmountSat int64,
	sendAddrStr string,
	feeSat int64,
	preimage []byte,
	key *Key,
) ([]byte, string, error) {
	originTx := wire.NewMsgTx(2)
	// originTx.LockTime = 0xffffffff

	prevAddr, err := s.createP2wsh()
	if err != nil {
		return nil, "", fmt.Errorf("fail createP2wpkh(prevAddr): %w", err)
	}
	prevPkScript, err := txscript.PayToAddrScript(prevAddr)
	if err != nil {
		return nil, "", fmt.Errorf("fail PayToAddrScript(prevAddr): %w", err)
	}
	prevOutputFetcher := txscript.NewCannedPrevOutputFetcher(prevPkScript, prevAmountSat)
	txinIndex := int(0)

	sendAddr, err := btcutil.DecodeAddress(sendAddrStr, s.Net)
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
	witSig, err := txscript.RawTxInWitnessSignature(
		originTx,
		sigHashes,
		txinIndex,
		prevAmountSat,
		s.Script,
		txscript.SigHashAll,
		key.privKey,
	)
	if err != nil {
		return nil, "", fmt.Errorf("fail RawTxInWitnessSignature: %w", err)
	}
	//  <<signature>>
	//  <<preimage>>
	//  <<script>
	txIn.Witness = [][]byte{witSig, preimage, s.Script}

	var buf bytes.Buffer
	originTx.Serialize(&buf)
	txid := originTx.TxHash()
	return buf.Bytes(), txid.String(), nil
}
