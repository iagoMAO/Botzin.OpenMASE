package protocol

type AnswerPacketType int

const (
	AvatarAttribLoadAnswer   AnswerPacketType = 110
	AvatarAttribSaveAnswer   AnswerPacketType = 106
	AvatarCheckPayLoadAnswer AnswerPacketType = 114
	AvatarSetupLoadAnswer    AnswerPacketType = 122
	AvatarSetupSaveAnswer    AnswerPacketType = 118
	BroadcastAnswer          AnswerPacketType = 48
	DisconnectAnswer         AnswerPacketType = 44
	HackBusterAnswer         AnswerPacketType = 151
	LoginAnswer              AnswerPacketType = 4
	MaseShowGUIAnswer        AnswerPacketType = 124
	PingAnswer               AnswerPacketType = 40
	ScoreAnswer              AnswerPacketType = 20
	ServerQueryAvatarAnswer  AnswerPacketType = 16
	ShopBuyAnswer            AnswerPacketType = 28
	ShopGiftAnswer           AnswerPacketType = 36
	ShopInfoAnswer           AnswerPacketType = 24
	ShopSellAnswer           AnswerPacketType = 32
	UserBootAnswer           AnswerPacketType = 12
	UserDataAnswer           AnswerPacketType = 8
)
