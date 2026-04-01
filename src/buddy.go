package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/iagoMAO/Botzin.OpenMASE/authentication"
	"github.com/iagoMAO/Botzin.OpenMASE/buddy"
	"github.com/iagoMAO/Botzin.OpenMASE/core"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
	"github.com/iagoMAO/Botzin.OpenMASE/utils"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	roomMessagePacketCache = make(map[int][]packets.PublicMessagePacket)
)

func BootBuddyRequest(session *core.Session) {
	buddies := buddy.GetCategorizedContacts(session.UserId)

	onlinePacket := packets.BootBuddyAnswerPacket{
		Status:              protocol.StatusCode(protocol.BUDDY_STATUS_ONLINE),
		TotalContactsOnList: len(buddies.Online),
		Contacts:            buddies.Online,
	}

	ingamePacket := packets.BootBuddyAnswerPacket{
		Status:              protocol.StatusCode(protocol.BUDDY_STATUS_INGAME),
		TotalContactsOnList: len(buddies.Ingame),
		Contacts:            buddies.Ingame,
	}

	offlinePacket := packets.BootBuddyAnswerPacket{
		Status:              protocol.StatusCode(protocol.BUDDY_STATUS_OFFLINE),
		TotalContactsOnList: len(buddies.Offline),
		Contacts:            buddies.Offline,
	}

	endPacket := packets.BootBuddyAnswerPacket{
		Status:              protocol.StatusCode(protocol.BUDDY_ENDOF_LIST),
		TotalContactsOnList: 0,
	}

	guiPacket := packets.MaseShowGUIAnswerPacket{StatusCode: protocol.MASE_OK}

	session.BuddyConn.Write(onlinePacket.Compose())
	session.BuddyConn.Write(ingamePacket.Compose())
	session.BuddyConn.Write(offlinePacket.Compose())
	session.BuddyConn.Write(endPacket.Compose())
	session.BuddyConn.Write(guiPacket.Compose())
}

func StartBuddyList() {
	// First and foremost, load our config.
	cfg := utils.GetConfig()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create the listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.BUDDY_PORT))

	if err != nil {
		log.Error().Msgf("Listening error: %s", err)
		return
	}

	log.Info().Msgf("BUDDY - Successfully started listening on port %s.", cfg.BUDDY_PORT)

	// Close the socket once we're done
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Error().Msgf("Error thrown whilst accepting connection: %s", err)
			continue
		}

		go handleBuddyConnection(conn)
	}
}

func handleBuddyConnection(conn net.Conn) {
	// Close once we're done, again
	defer core.RemoveSession(conn)
	defer conn.Close()

	// TODO: Maybe make this configurable?func (p PacketType) Code() byte { return byte(p) }

	buf := make([]byte, 1024)

	reader := bufio.NewReader(conn)

	for {
		n, err := reader.Read(buf)

		if err != nil {
			conn.Close()
			return
		}

		if reader.Size() <= 2 {
			return
		}

		message := protocol.DecryptPacket(buf[:n])

		switch message.Type {
		case protocol.LoginRequest:
			log.Debug().Msgf("Received Login Request: %s", hex.EncodeToString(message.Payload))

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			if len(parts) < 3 {
				conn.Write(protocol.EncryptPacket(protocol.LoginAnswer, []byte{}, protocol.MASE_ERROR))
				return
			}

			var id int

			id = data.SCR_StrToInt(parts[1])

			conn.Write(protocol.EncryptPacket(protocol.LoginAnswer, []byte{}, protocol.MASE_OK))

			if id != 0 {
				core.RegisterConnection(conn, id, core.ConnTypeBUDDY)
			}

		case protocol.AddContactRequest:
			log.Debug().Msgf("Received Add Contact Request: %s", hex.EncodeToString(message.Payload))

			session := core.GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			if len(parts) < 1 {
				return
			}

			userId := data.SCR_StrToInt(parts[1])

			err := buddy.AddContact(session.UserId, userId)

			if err != nil {
				return
			}

			contacts := buddy.GetContacts(session.UserId)

			response := packets.AddContactAnswerPacket{
				Status:              protocol.MASE_OK,
				TotalContactsOnList: len(contacts) + 1,
				Contacts:            contacts,
			}

			buddySession := core.GetSessionByUserId(userId)

			if buddySession != nil {
				pending := buddy.GetOutgoingRequests(userId, protocol.BUDDY_ANSWER_REQUEST)

				pendingPacket := packets.BootStatusAnswerPacket{
					Status:              protocol.StatusCode(protocol.BUDDY_ANSWER_REQUEST),
					TotalContactsOnList: len(pending),
					Contacts:            pending,
				}

				buddySession.BuddyConn.Write(pendingPacket.Compose())
			}

			conn.Write(response.Compose())
		case protocol.FindContactRequest:
			log.Debug().Msgf("Received Find Contact Request: %s", hex.EncodeToString(message.Payload))

			session := core.GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			request := packets.FindContactRequestPacket{
				GUID: data.SCR_StrToInt(parts[1]),
				Name: string(parts[2]),
			}

			contacts := buddy.QueryContacts(request.Name)

			response := packets.FindContactAnswerPacket{
				Status:              protocol.MASE_OK,
				TotalContactsOnList: len(contacts),
				Contacts:            contacts,
			}

			log.Debug().Msgf("Received Find Contact Request: %#v\n", request)
			log.Debug().Msgf("Sent Find Contact Response: %#v\n", response)

			conn.Write(response.Compose())
		case protocol.BootBuddyRequest:
			log.Debug().Msgf("Received Boot Buddy Request: %s", hex.EncodeToString(message.Payload))

			session := core.GetSession(conn)

			if session == nil {
				return
			}

			BootBuddyRequest(session)

			buddies := buddy.GetCategorizedContacts(session.UserId)

			for _, contact := range buddies.Online {
				session := core.GetSessionByUserId(contact.GUID)
				if session == nil {
					continue
				}

				BootBuddyRequest(session)
			}

			for _, contact := range buddies.Ingame {
				session := core.GetSessionByUserId(contact.GUID)
				if session == nil {
					continue
				}

				BootBuddyRequest(session)
			}
		case protocol.PrivateMessage:
			log.Debug().Msgf("Received Private Message: %s", hex.EncodeToString(message.Payload))

			session := core.GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			buddyId := data.SCR_StrToInt(parts[1])
			message := string(parts[2])

			buddySession := core.GetSessionByUserId(buddyId)

			if buddySession == nil {
				return
			}

			user := authentication.GetUserInfo(session.UserId)

			contact := packets.BuddyContactInfo{
				GUID: session.UserId,
				Name: user.Nick,
			}

			packet := packets.PrivateMessagePacket{
				Status:  protocol.MASE_OK,
				Contact: contact,
				Message: message,
			}

			if buddySession.BuddyConn != nil {
				buddySession.BuddyConn.Write(packet.Compose())
			}
		case protocol.BuddyResponse:
			log.Debug().Msgf("Received Buddy Response: %s", hex.EncodeToString(message.Payload))

			session := core.GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			buddyId := parts[1]
			status := parts[2]

			buddy.RespondContact(session.UserId, data.SCR_StrToInt(buddyId), protocol.BuddyStatus(status[0]))

			buddySession := core.GetSessionByUserId(data.SCR_StrToInt(buddyId))

			if buddySession != nil {
				contacts := buddy.GetOutgoingRequests(data.SCR_StrToInt(buddyId), protocol.BuddyStatus(status[0]))
				contactsPacket := packets.BootStatusAnswerPacket{
					Status:              protocol.StatusCode(protocol.BuddyStatus(status[0])),
					TotalContactsOnList: len(contacts),
					Contacts:            contacts,
				}

				buddySession.BuddyConn.Write(contactsPacket.Compose())
				buddySession.BuddyConn.Write(protocol.EncryptPacket(protocol.BootStatusAnswer, []byte{}, protocol.StatusCode(protocol.BUDDY_ENDOF_LIST)))
			}
		case protocol.BootStatusRequest:
			log.Debug().Msgf("Received Boot Status: %s", hex.EncodeToString(message.Payload))

			session := core.GetSession(conn)

			if session == nil {
				return
			}

			pending := buddy.GetIncomingRequests(session.UserId)
			rejected := buddy.GetOutgoingRequests(session.UserId, protocol.BUDDY_ANSWER_REJECTED)
			accepted := buddy.GetOutgoingRequests(session.UserId, protocol.BUDDY_ANSWER_ACCEPTED)
			removed := buddy.GetOutgoingRequests(session.UserId, protocol.BUDDY_ANSWER_REMOVED)

			pendingPacket := packets.BootStatusAnswerPacket{
				Status:              protocol.StatusCode(protocol.BUDDY_ANSWER_REQUEST),
				TotalContactsOnList: len(pending),
				Contacts:            pending,
			}

			rejectedPacket := packets.BootStatusAnswerPacket{
				Status:              protocol.StatusCode(protocol.BUDDY_ANSWER_REJECTED),
				TotalContactsOnList: len(rejected),
				Contacts:            rejected,
			}

			acceptedPacket := packets.BootStatusAnswerPacket{
				Status:              protocol.StatusCode(protocol.BUDDY_ANSWER_ACCEPTED),
				TotalContactsOnList: len(accepted),
				Contacts:            accepted,
			}

			removedPacket := packets.BootStatusAnswerPacket{
				Status:              protocol.StatusCode(protocol.BUDDY_ANSWER_REMOVED),
				TotalContactsOnList: len(removed),
				Contacts:            removed,
			}

			conn.Write(pendingPacket.Compose())
			conn.Write(rejectedPacket.Compose())
			conn.Write(acceptedPacket.Compose())
			conn.Write(removedPacket.Compose())
			conn.Write(protocol.EncryptPacket(protocol.BootStatusAnswer, []byte{}, protocol.StatusCode(protocol.BUDDY_ENDOF_LIST)))
		case protocol.PublicMessage:
			log.Debug().Msgf("Received Public Message: %s", hex.EncodeToString(message.Payload))

			session := core.GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			roomId := data.SCR_StrToInt(parts[1])
			message := string(parts[2])

			user := authentication.GetUserInfo(session.UserId)

			if len(message) <= 0 {
				for _, packet := range roomMessagePacketCache[roomId] {
					conn.Write(packet.Compose())
				}
				continue
			}

			packet := packets.PublicMessagePacket{
				Status:  protocol.MASE_OK,
				RoomId:  roomId,
				Nick:    user.Nick,
				Message: message,
			}

			if len(roomMessagePacketCache[roomId]) > 32 {
				roomMessagePacketCache[roomId] = nil
			}

			roomMessagePacketCache[roomId] = append(roomMessagePacketCache[roomId], packet)

			sessions := core.GetAllSessions()

			for _, session := range sessions {
				session.BuddyConn.Write(packet.Compose())
			}
		default:
			log.Debug().Msgf("Received Packet %s", hex.Dump(buf))
		}
	}
}
