package packets

import (
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
)

type MaseShowGUIAnswerPacket struct {
	StatusCode  protocol.StatusCode
	MagicNumber uint32
	ClientGUID  uint32
}

func (p MaseShowGUIAnswerPacket) Compose() []byte {
	return protocol.EncryptPacket(protocol.MaseShowGUIAnswer, []byte{}, p.StatusCode)
}
