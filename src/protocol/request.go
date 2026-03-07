package protocol

type RequestPacketType int

const (
	AvatarAttribLoadRequest   RequestPacketType = 108
	AvatarAttribSaveRequest   RequestPacketType = 104
	AvatarCheckPayLoadRequest RequestPacketType = 112
	AvatarSetupLoadRequest    RequestPacketType = 120
	AvatarSetupSaveRequest    RequestPacketType = 116
	BroadcastRequest          RequestPacketType = 46
	DisconnectRequest         RequestPacketType = 42
	HackBusterRequest         RequestPacketType = 150
	LoginRequest              RequestPacketType = 2
	PingRequest               RequestPacketType = 38
	ScoreRequest              RequestPacketType = 18
	ServerQueryAvatarRequest  RequestPacketType = 14
	ShopBuyRequest            RequestPacketType = 26
	ShopGiftRequest           RequestPacketType = 34
	ShopInfoRequest           RequestPacketType = 22
	ShopSellRequest           RequestPacketType = 30
	UserBootRequest           RequestPacketType = 10
	UserDataRequest           RequestPacketType = 6
)
