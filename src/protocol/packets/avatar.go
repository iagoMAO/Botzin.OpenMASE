package packets

import (
	"bytes"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"
)

type AvatarItemData struct {
	Id      int
	UserId  int
	Class   int
	ST      int
	DX      int
	IQ      int
	HT      int
	Payload int
	TheGen  int
	Enabled int
}

type AvatarSetupLoadAnswerPacket struct {
	TotalAvatarItems int
	Items            []AvatarItemData
}

type ServerQueryAvatarAnswerPacket struct {
	ClientGUID       int
	TotalAvatarItems int
	Items            []AvatarItemData
	Nick             string
	ST               int
	DX               int
	IQ               int
	HT               int
	XP               int
	patent           int
	ranking          int
}

func createItemsPayload(totalItems int, items []AvatarItemData) []byte {
	var buf bytes.Buffer

	buf.Write(data.SCR_PackInt(totalItems))

	if totalItems <= 0 {
		return buf.Bytes()
	}

	for _, item := range items {
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.Class))
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.Id))
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.ST))
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.DX))
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.IQ))
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.HT))
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.Payload))
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.TheGen))
		buf.WriteByte(0x09)
		buf.Write(data.SCR_PackInt(item.Enabled))
		buf.WriteByte(0x0A)
	}

	return buf.Bytes()
}

func (p ServerQueryAvatarAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.Write(data.SCR_PackInt(p.ClientGUID))
	buf.WriteByte(0x09)
	buf.Write([]byte(p.Nick))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.ST))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.DX))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.IQ))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.HT))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.XP))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.patent))
	buf.WriteByte(0x09)
	buf.Write(data.SCR_PackInt(p.ranking))

	if p.TotalAvatarItems > 0 {
		buf.WriteByte(0x09)
		itemsData := createItemsPayload(p.TotalAvatarItems, p.Items)
		buf.Write(itemsData)
	}

	return protocol.EncryptPacket(protocol.ServerQueryAvatarAnswer, buf.Bytes(), protocol.MASE_OK)
}

func (p AvatarSetupLoadAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.Write(data.SCR_PackInt(p.TotalAvatarItems))

	if p.TotalAvatarItems <= 0 {
		return protocol.EncryptPacket(protocol.AvatarSetupLoadAnswer, buf.Bytes(), protocol.MASE_OK)
	}

	itemsData := createItemsPayload(p.TotalAvatarItems, p.Items)

	buf.Write(itemsData)

	return protocol.EncryptPacket(protocol.AvatarSetupLoadAnswer, buf.Bytes(), protocol.MASE_OK)
}

type AvatarAttribLoadAnswerPacket struct {
	XP string
	ST string
	DX string
	IQ string
	HT string
}

func (p AvatarAttribLoadAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.WriteString(p.XP)
	buf.WriteByte(0x09)
	buf.WriteString(p.ST)
	buf.WriteByte(0x09)
	buf.WriteString(p.DX)
	buf.WriteByte(0x09)
	buf.WriteString(p.IQ)
	buf.WriteByte(0x09)
	buf.WriteString(p.HT)
	buf.WriteByte(0x09)

	return protocol.EncryptPacket(protocol.AvatarAttribLoadAnswer, buf.Bytes(), protocol.MASE_OK)
}
