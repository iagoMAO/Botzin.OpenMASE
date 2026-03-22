package packets

import (
	"bytes"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"
)

type BuddyContactInfo struct {
	GUID int
	Name string
}

type BootBuddyAnswerPacket struct {
	Status              protocol.StatusCode
	TotalContactsOnList int
	Contacts            []BuddyContactInfo
}

type FindContactRequestPacket struct {
	GUID int
	Name string
}

type FindContactAnswerPacket struct {
	Status              protocol.StatusCode
	TotalContactsOnList int
	Contacts            []BuddyContactInfo
}

func (p FindContactAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.Write(data.SCR_PackInt(p.TotalContactsOnList))

	if p.TotalContactsOnList > 0 {
		for _, contact := range p.Contacts {
			buf.WriteByte(0x09)
			buf.Write(data.SCR_PackInt(contact.GUID))
			buf.WriteByte(0x09)
			buf.WriteString(contact.Name)
		}
	}

	return protocol.EncryptPacket(protocol.FindContactAnswer, buf.Bytes(), p.Status)
}

type AddContactAnswerPacket struct {
	Status              protocol.StatusCode
	TotalContactsOnList int
	Contacts            []BuddyContactInfo
}

type BootStatusAnswerPacket struct {
	Status              protocol.StatusCode
	TotalContactsOnList int
	Contacts            []BuddyContactInfo
}

func (p BootStatusAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.Write(data.SCR_PackInt(p.TotalContactsOnList))

	if p.TotalContactsOnList > 0 {
		for _, contact := range p.Contacts {
			buf.WriteByte(0x09)
			buf.Write(data.SCR_PackInt(contact.GUID))
			buf.WriteByte(0x09)
			buf.WriteString(contact.Name)
		}
	}

	return protocol.EncryptPacket(protocol.BootStatusAnswer, buf.Bytes(), p.Status)
}

func (p AddContactAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.Write(data.SCR_PackInt(p.TotalContactsOnList))

	if p.TotalContactsOnList > 0 {
		for _, contact := range p.Contacts {
			buf.WriteByte(0x09)
			buf.Write(data.SCR_PackInt(contact.GUID))
			buf.WriteByte(0x09)
			buf.WriteString(contact.Name)
		}
	}

	return protocol.EncryptPacket(protocol.AddContactAnswer, buf.Bytes(), p.Status)
}

func (p BootBuddyAnswerPacket) Compose() []byte {
	var buf bytes.Buffer

	buf.Write(data.SCR_PackInt(p.TotalContactsOnList))

	if p.TotalContactsOnList > 0 {
		for _, contact := range p.Contacts {
			buf.WriteByte(0x09)
			buf.Write(data.SCR_PackInt(contact.GUID))
			buf.WriteByte(0x09)
			buf.WriteString(contact.Name)
		}
	}

	return protocol.EncryptPacket(protocol.BootBuddyAnswer, buf.Bytes(), p.Status)
}

type PrivateMessagePacket struct {
	Status  protocol.StatusCode
	Contact BuddyContactInfo
	Message string
}

func (p PrivateMessagePacket) Compose() []byte {
	var buf bytes.Buffer

	buf.Write(data.SCR_PackInt(p.Contact.GUID))
	buf.WriteByte(0x09)
	buf.WriteString(p.Contact.Name)
	buf.WriteByte(0x09)
	buf.WriteString(p.Message)

	return protocol.EncryptPacket(protocol.PrivateMessage, buf.Bytes(), p.Status)
}
