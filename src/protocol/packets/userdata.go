package packets

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
)

type UserDataRequestPacket struct {
	EXECRC       string
	PasswordHash string
	UsernameHash string
	ClanTag      string
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

	log.Println(hex.EncodeToString(buf.Bytes()))

	return protocol.EncryptPacket(protocol.AvatarAttribLoadAnswer, buf.Bytes(), protocol.MASE_OK)
}

type PacketUserInfo struct {
	Nick        string
	XP          int
	ST          int
	DX          int
	IQ          int
	HT          int
	PromoButton string
	Points      string
	Credits     string
	Gold        string
	Ranking     string
	TotalRK     string
	Level       string
	PMX         string
}

type UserDataAnswerPacket struct {
	StatusCode protocol.StatusCode
	UserInfo   PacketUserInfo
}

func (p UserDataAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.WriteString(p.UserInfo.Nick)
	buf.WriteByte(0x09)

	// oops. those are the attributes
	// binary.Write(&buf, binary.BigEndian, p.UserInfo.XP)
	// buf.WriteByte(0x09)
	// binary.Write(&buf, binary.BigEndian, p.UserInfo.ST)
	// buf.WriteByte(0x09)
	// binary.Write(&buf, binary.BigEndian, p.UserInfo.DX)
	// buf.WriteByte(0x09)
	// binary.Write(&buf, binary.BigEndian, p.UserInfo.IQ)
	// buf.WriteByte(0x09)
	// binary.Write(&buf, binary.BigEndian, p.UserInfo.HT)
	// buf.WriteByte(0x09)

	buf.WriteString(p.UserInfo.PromoButton)
	buf.WriteByte(0x09)
	buf.WriteString(p.UserInfo.Points)
	buf.WriteByte(0x09)
	buf.WriteString(p.UserInfo.Credits)
	buf.WriteByte(0x09)
	buf.WriteString(p.UserInfo.Gold)
	buf.WriteByte(0x09)
	buf.WriteString(p.UserInfo.Ranking)
	buf.WriteByte(0x09)
	buf.WriteString(p.UserInfo.TotalRK)
	buf.WriteByte(0x09)
	buf.WriteString(p.UserInfo.Level)
	buf.WriteByte(0x09)
	buf.WriteString(p.UserInfo.PMX)

	log.Println(hex.EncodeToString(buf.Bytes()))

	return protocol.EncryptPacket(protocol.UserDataAnswer, buf.Bytes(), p.StatusCode)
}
