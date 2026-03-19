package protocol

type PacketType byte

func (p PacketType) Code() byte { return byte(p) }

const (
	AvatarAttribLoadRequest   PacketType = 108
	AvatarAttribSaveRequest   PacketType = 104
	AvatarCheckPayLoadRequest PacketType = 112
	AvatarSetupLoadRequest    PacketType = 120
	AvatarSetupSaveRequest    PacketType = 116
	BroadcastRequest          PacketType = 46
	DisconnectRequest         PacketType = 42
	HackBusterRequest         PacketType = 150
	LoginRequest              PacketType = 2
	PingRequest               PacketType = 38
	ScoreRequest              PacketType = 18
	ServerQueryAvatarRequest  PacketType = 14
	ShopBuyRequest            PacketType = 26
	ShopGiftRequest           PacketType = 34
	ShopInfoRequest           PacketType = 22
	ShopSellRequest           PacketType = 30
	UserBootRequest           PacketType = 10
	UserDataRequest           PacketType = 6

	BootBuddyAnswer   PacketType = 92
	BootBuddyRequest  PacketType = 90
	BootStatusAnswer  PacketType = 88
	BootStatusRequest PacketType = 86

	AvatarAttribLoadAnswer   PacketType = 110
	AvatarAttribSaveAnswer   PacketType = 106
	AvatarCheckPayLoadAnswer PacketType = 114
	AvatarSetupLoadAnswer    PacketType = 122
	AvatarSetupSaveAnswer    PacketType = 118
	BroadcastAnswer          PacketType = 48
	DisconnectAnswer         PacketType = 44
	HackBusterAnswer         PacketType = 151
	LoginAnswer              PacketType = 4
	MaseShowGUIAnswer        PacketType = 124
	PingAnswer               PacketType = 40
	ScoreAnswer              PacketType = 20
	ServerQueryAvatarAnswer  PacketType = 16
	ShopBuyAnswer            PacketType = 28
	ShopGiftAnswer           PacketType = 36
	ShopInfoAnswer           PacketType = 24
	ShopSellAnswer           PacketType = 32
	UserBootAnswer           PacketType = 12
	UserDataAnswer           PacketType = 8
)
