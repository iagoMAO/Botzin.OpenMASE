package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/iagoMAO/Botzin.OpenMASE/authentication"
	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
	"github.com/iagoMAO/Botzin.OpenMASE/utils"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type BuddyStatus byte

func (p BuddyStatus) Code() byte { return byte(p) }

const (
	BUDDY_ANSWER_ACCEPTED BuddyStatus = 102
	BUDDY_ANSWER_REJECTED BuddyStatus = 101
	BUDDY_ANSWER_REMOVED  BuddyStatus = 103
	BUDDY_ANSWER_REQUEST  BuddyStatus = 100
	BUDDY_ENDOF_LIST      BuddyStatus = 200
	BUDDY_SHOW_WINDOW     BuddyStatus = 0
	BUDDY_STATUS_INGAME   BuddyStatus = 101
	BUDDY_STATUS_OFFLINE  BuddyStatus = 102
	BUDDY_STATUS_ONLINE   BuddyStatus = 100
)

func RespondContact(userId int, contactId int, status BuddyStatus) {
	query := `SELECT COUNT(*) FROM contacts WHERE user_id = ? AND contact_id = ?`

	var existing int
	err := database.DB.QueryRow(query, contactId, userId).Scan(&existing)

	if err != nil {
		fmt.Println(err)
		return
	}

	if existing <= 0 {
		fmt.Println("exist")
		fmt.Println(userId)
		fmt.Println(contactId)

		return
	}

	query = `UPDATE contacts SET status = ? WHERE user_id = ? AND contact_id = ?`

	_, err = database.DB.Exec(query, status, contactId, userId)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func AddContact(userId int, contactId int) error {
	err := database.DB.QueryRow("SELECT id FROM users WHERE id = ?", contactId).Scan(&contactId)

	if err != nil {
		return err
	}

	var existing int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM contacts WHERE (user_id = ? AND contact_id = ?) OR (contact_id = ? AND user_id = ?)", userId, contactId, userId, contactId).Scan(&existing)

	if err != nil {
		return err
	}

	if existing > 0 {
		return err
	}

	query := "INSERT INTO contacts (user_id, contact_id, status) VALUES(?, ?, ?)"

	_, err = database.DB.Exec(query, userId, contactId, BUDDY_ANSWER_REQUEST)

	if err != nil {
		return err
	}

	return nil
}

func GetUserContacts(userId int, status BuddyStatus) []packets.BuddyContactInfo {
	query := `
		SELECT c.user_id AS friend_id, u.username AS friend_name 
		FROM contacts c 
		JOIN users u ON c.user_id = u.id 
		WHERE c.contact_id = ? AND c.status = ?`

	rows, err := database.DB.Query(query, userId, status)

	if err != nil {
		fmt.Println(err)
		return []packets.BuddyContactInfo{}
	}

	defer rows.Close()

	var contacts []packets.BuddyContactInfo

	for rows.Next() {
		var contact packets.BuddyContactInfo
		if err := rows.Scan(&contact.GUID, &contact.Name); err != nil {
			return []packets.BuddyContactInfo{}
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return []packets.BuddyContactInfo{}
	}

	fmt.Println(contacts)

	return contacts
}

func GetContacts(userId int) []packets.BuddyContactInfo {
	query := `
		SELECT c.contact_id AS friend_id, u.username AS friend_name 
		FROM contacts c 
		JOIN users u ON c.contact_id = u.id 
		WHERE c.user_id = ? AND status = ?
		
		UNION
		
		SELECT c.user_id AS friend_id, u.username AS friend_name 
		FROM contacts c 
		JOIN users u ON c.user_id = u.id 
		WHERE c.contact_id = ? AND status = ?`

	rows, err := database.DB.Query(query, userId, BUDDY_ANSWER_ACCEPTED, userId, BUDDY_ANSWER_ACCEPTED)

	if err != nil {
		fmt.Println(err)
		return []packets.BuddyContactInfo{}
	}

	defer rows.Close()

	var contacts []packets.BuddyContactInfo

	for rows.Next() {
		var contact packets.BuddyContactInfo
		if err := rows.Scan(&contact.GUID, &contact.Name); err != nil {
			return []packets.BuddyContactInfo{}
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return []packets.BuddyContactInfo{}
	}

	return contacts
}

func QueryContacts(query string) []packets.BuddyContactInfo {
	rows, err := database.DB.Query("SELECT id, username FROM users WHERE username LIKE ?", "%"+query+"%")

	if err != nil {
		return []packets.BuddyContactInfo{}
	}

	defer rows.Close()

	var contacts []packets.BuddyContactInfo

	for rows.Next() {
		var contact packets.BuddyContactInfo
		if err := rows.Scan(&contact.GUID, &contact.Name); err != nil {
			return []packets.BuddyContactInfo{}
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return []packets.BuddyContactInfo{}
	}

	return contacts
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
	defer RemoveSession(conn)
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

		if reader.Size() <= 1 {
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
				CreateSession(conn, id)
			}

		case protocol.AddContactRequest:
			log.Debug().Msgf("Received Add Contact Request: %s", hex.EncodeToString(message.Payload))

			session := GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			if len(parts) < 1 {
				return
			}

			userId := data.SCR_StrToInt(parts[1])

			err := AddContact(session.UserId, userId)

			if err != nil {
				return
			}

			contacts := GetContacts(session.UserId)

			response := packets.AddContactAnswerPacket{
				Status:              protocol.MASE_OK,
				TotalContactsOnList: len(contacts),
				Contacts:            contacts,
			}

			conn.Write(response.Compose())
		case protocol.FindContactRequest:
			log.Debug().Msgf("Received Find Contact Request: %s", hex.EncodeToString(message.Payload))

			session := GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			request := packets.FindContactRequestPacket{
				GUID: data.SCR_StrToInt(parts[1]),
				Name: string(parts[2]),
			}

			contacts := QueryContacts(request.Name)

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

			session := GetSession(conn)

			if session == nil {
				return
			}

			contacts := GetContacts(session.UserId)

			var onlineContacts []packets.BuddyContactInfo
			for _, contact := range contacts {
				session := GetSessionByUserId(contact.GUID)
				if session == nil {
					fmt.Println("session not found")
					continue
				}

				fmt.Printf("session status for %d: %d\n", session.UserId, session.Status)
				if session.Status == BUDDY_STATUS_ONLINE {
					onlineContacts = append(onlineContacts, contact)
				} else {
					continue
				}
			}

			fmt.Printf("online contacts: %v\n contacts: %v", onlineContacts, contacts)

			onlinePacket := packets.BootBuddyAnswerPacket{
				Status:              protocol.StatusCode(BUDDY_STATUS_INGAME),
				TotalContactsOnList: len(onlineContacts),
				Contacts:            onlineContacts,
			}

			endPacket := packets.BootBuddyAnswerPacket{
				Status:              protocol.StatusCode(BUDDY_ENDOF_LIST),
				TotalContactsOnList: 0,
			}

			guiPacket := packets.MaseShowGUIAnswerPacket{StatusCode: protocol.MASE_OK}

			conn.Write(onlinePacket.Compose())
			conn.Write(endPacket.Compose())
			conn.Write(guiPacket.Compose())
		case protocol.PrivateMessage:
			log.Debug().Msgf("Received Private Message: %s", hex.EncodeToString(message.Payload))

			session := GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			buddyId := data.SCR_StrToInt(parts[1])
			message := string(parts[2])

			buddySession := GetSessionByUserId(buddyId)

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

			buddySession.Conn.Write(packet.Compose())
		case protocol.BuddyResponse:
			log.Debug().Msgf("Received Buddy Response: %s", hex.EncodeToString(message.Payload))

			session := GetSession(conn)

			if session == nil {
				return
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			buddyId := parts[1]
			status := parts[2]

			RespondContact(session.UserId, data.SCR_StrToInt(buddyId), BuddyStatus(status[0]))

		case protocol.BootStatusRequest:
			log.Debug().Msgf("Received Boot Status: %s", hex.EncodeToString(message.Payload))

			session := GetSession(conn)

			if session == nil {
				return
			}

			pending := GetUserContacts(session.UserId, BUDDY_ANSWER_REQUEST)
			rejected := GetUserContacts(session.UserId, BUDDY_ANSWER_REJECTED)
			accepted := GetUserContacts(session.UserId, BUDDY_ANSWER_ACCEPTED)
			removed := GetUserContacts(session.UserId, BUDDY_ANSWER_REMOVED)

			pendingPacket := packets.BootStatusAnswerPacket{
				Status:              protocol.StatusCode(BUDDY_ANSWER_REQUEST),
				TotalContactsOnList: len(pending),
				Contacts:            pending,
			}

			rejectedPacket := packets.BootStatusAnswerPacket{
				Status:              protocol.StatusCode(BUDDY_ANSWER_REJECTED),
				TotalContactsOnList: len(rejected),
				Contacts:            rejected,
			}

			acceptedPacket := packets.BootStatusAnswerPacket{
				Status:              protocol.StatusCode(BUDDY_ANSWER_ACCEPTED),
				TotalContactsOnList: len(accepted),
				Contacts:            accepted,
			}

			removedPacket := packets.BootStatusAnswerPacket{
				Status:              protocol.StatusCode(BUDDY_ANSWER_REMOVED),
				TotalContactsOnList: len(removed),
				Contacts:            removed,
			}

			conn.Write(pendingPacket.Compose())
			conn.Write(rejectedPacket.Compose())
			conn.Write(acceptedPacket.Compose())
			conn.Write(removedPacket.Compose())
			conn.Write(protocol.EncryptPacket(protocol.BootStatusAnswer, []byte{}, protocol.StatusCode(BUDDY_ENDOF_LIST)))
		default:
			log.Debug().Msgf("Received Packet %s - %d", hex.Dump(message.Payload), message.Type.Code())
		}
	}
}
