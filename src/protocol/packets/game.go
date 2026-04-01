package packets

import (
	"bytes"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
)

type ScoreAnswer struct {
	StatusCode   protocol.StatusCode
	MessageColor int
	MessageText  string
}

func (p ScoreAnswer) Compose() []byte {
	var buf bytes.Buffer
	buf.WriteString(p.MessageText)
	buf.WriteByte(0x09)
	buf.WriteByte(byte(p.MessageColor))

	return protocol.EncryptPacket(protocol.BroadcastAnswer, buf.Bytes(), p.StatusCode)
}
