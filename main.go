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
	"taptx/redeem"
	"taptx/redeem/keypath"
	"taptx/redeem/scriptpath"
)

func main() {
	rt := redeem.RedeemP2wsh
	switch rt {
	case redeem.RedeemP2wpkh:
	case redeem.RedeemP2trKeyPath:
		keypath.KeyPath(rt)
	case redeem.RedeemP2wsh:
		scriptpath.ScriptPath(rt)
	}
}
