package protocol

type StatusCode byte

func (p StatusCode) Code() byte { return byte(p) }

const (
	MASE_ACCOUNT_BANNED    StatusCode = 106
	MASE_ACCOUNT_BLOCKED   StatusCode = 104
	MASE_ACCOUNT_INACTIVE  StatusCode = 108
	MASE_ALREADY_LOGGED    StatusCode = 102
	MASE_ATTRIBS_LOADED    StatusCode = 0
	MASE_ERROR             StatusCode = 100
	MASE_HACK_DETECTED     StatusCode = 171
	MASE_INVALID_LOGINPASS StatusCode = 100
	MASE_OK                StatusCode = 102

	SHOP_ALREADY_HAVE StatusCode = 103
	SHOP_BUY_DONE     StatusCode = 104
	SHOP_CANT_SELL    StatusCode = 106
	SHOP_DONT_HAVE    StatusCode = 105
	SHOP_NO_CREDITS   StatusCode = 100
	SHOP_NO_GOLD      StatusCode = 101
	SHOP_NO_STOCK     StatusCode = 102
	SHOP_SELL_DONE    StatusCode = 107
)
