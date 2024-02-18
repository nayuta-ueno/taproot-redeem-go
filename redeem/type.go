package redeem

type RedeemType = int

const (
	RedeemP2wpkh      RedeemType = 1
	RedeemP2trKeyPath RedeemType = 2
	RedeemP2wsh       RedeemType = 3
)
