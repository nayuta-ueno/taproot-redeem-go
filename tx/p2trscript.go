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

// P2TR script
func (s *TapScript) CreateP2tr() (string, error) {
	addr, err := s.createP2tr()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (s *TapScript) createP2tr() (*btcutil.AddressTaproot, error) {
	tree := txscript.AssembleTaprootScriptTree(s.Scripts...)
	rootHash := tree.RootNode.TapHash()
	pubKey := txscript.ComputeTaprootOutputKey(s.Key.PubKey, rootHash[:])
	witnessProg := schnorr.SerializePubKey(pubKey)
	return btcutil.NewAddressTaproot(witnessProg, s.Net)
}

func (s *TapScript) CreateRawTxP2TR(
	prevHash *chainhash.Hash,
	prevIndex uint32,
	prevAmountSat int64,
	sendAddrStr string,
	feeSat int64,
	scriptIdx int,
	witnessStack [][]byte,
	key *Key,
) ([]byte, string, error) {
	originTx := wire.NewMsgTx(2)

	prevAddr, err := s.createP2tr()
	if err != nil {
		return nil, "", fmt.Errorf("fail createP2tr(prevAddr): %w", err)
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

	// control block
	indexedTree := txscript.AssembleTaprootScriptTree(s.Scripts...)
	settleMerkleProof := indexedTree.LeafMerkleProofs[scriptIdx]
	cb := settleMerkleProof.ToControlBlock(s.Key.PubKey)
	cbBytes, _ := cb.ToBytes()

	sigHashes := txscript.NewTxSigHashes(originTx, prevOutputFetcher)
	witSig, err := txscript.RawTxInTapscriptSignature(
		originTx,
		sigHashes,
		txinIndex,
		prevAmountSat,
		prevPkScript,
		s.Scripts[scriptIdx], // leaf
		txscript.SigHashDefault,
		key.privKey,
	)
	if err != nil {
		return nil, "", fmt.Errorf("fail RawTxInWitnessSignature: %w", err)
	}
	txIn.Witness = make([][]byte, len(witnessStack)+3) // signature + <witnessStack> + redeem_script + control_block
	txIn.Witness[0] = witSig
	for i := 0; i < len(witnessStack); i++ {
		txIn.Witness[1+i] = witnessStack[i]
	}
	txIn.Witness[len(txIn.Witness)-2] = s.Scripts[scriptIdx].Script
	txIn.Witness[len(txIn.Witness)-1] = cbBytes

	var buf bytes.Buffer
	originTx.Serialize(&buf)
	txid := originTx.TxHash()
	return buf.Bytes(), txid.String(), nil
}
