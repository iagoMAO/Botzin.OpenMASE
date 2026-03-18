package packets

import (
	"bytes"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
)

type MaseShowGUIAnswerPacket struct {
	StatusCode  protocol.StatusCode
	MagicNumber uint32
	ClientGUID  uint32
}

type BroadcastAnswerPacket struct {
	StatusCode   protocol.StatusCode
	MessageColor int
	MessageText  string
}

func (p BroadcastAnswerPacket) Compose() []byte {
	var buf bytes.Buffer
	buf.WriteString(p.MessageText)
	buf.WriteByte(0x09)
	buf.WriteByte(byte(p.MessageColor))

	return protocol.EncryptPacket(protocol.BroadcastAnswer, buf.Bytes(), p.StatusCode)
}

func (p MaseShowGUIAnswerPacket) Compose() []byte {
	return protocol.EncryptPacket(protocol.MaseShowGUIAnswer, []byte{}, p.StatusCode)
}
