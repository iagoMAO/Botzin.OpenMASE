package packets

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"
)

type ShopBuyAnswerPacket struct {
	ShopBuyAnswerType protocol.StatusCode
	Item              AvatarItemData
}

func (p ShopBuyAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.WriteByte(p.ShopBuyAnswerType.Code())

	buf.WriteByte(0x09)

	buf.Write(data.SCR_PackInt(p.Item.Class))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.Item.Id))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.Item.ST))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.Item.DX))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.Item.IQ))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.Item.HT))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.Item.Payload))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.Item.TheGen))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.Item.Enabled))

	log.Println(hex.EncodeToString(buf.Bytes()))

	return protocol.EncryptPacket(protocol.ShopBuyAnswer, buf.Bytes(), protocol.MASE_OK)
}
