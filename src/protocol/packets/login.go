package packets

import (
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
)

type LoginRequestPacket struct {
	EXECRC       string
	PasswordHash string
	UsernameHash string
	ClanTag      string
}

type LoginAnswerPacket struct {
	StatusCode  protocol.StatusCode
	MagicNumber uint32
	ClientGUID  uint32
}

type LoginErrorPacket struct {
	StatusCode protocol.StatusCode
}

func (p LoginErrorPacket) Compose() []byte {
	data := []byte{
		uint8(p.StatusCode.Code()),
	}

	return protocol.EncryptPacket(protocol.LoginAnswer, data, p.StatusCode)
}

func (p LoginAnswerPacket) Compose() []byte {
	data := []byte{
		uint8(p.MagicNumber),
		0x09,
		uint8(p.ClientGUID),
	}

	return protocol.EncryptPacket(protocol.LoginAnswer, data, protocol.MASE_OK)
}
