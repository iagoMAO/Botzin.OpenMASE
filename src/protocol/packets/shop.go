package packets

import (
	"bytes"

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

	buf.Write(data.SCR_PackInt(p.Item.Id))

	return protocol.EncryptPacket(protocol.ShopBuyAnswer, buf.Bytes(), protocol.MASE_OK)
}
