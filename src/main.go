/*
Botzin.OpenMASE
@author: maldoliver
*/
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/iagoMAO/Botzin.OpenMASE/authentication"
	"github.com/iagoMAO/Botzin.OpenMASE/avatar"
	"github.com/iagoMAO/Botzin.OpenMASE/buddy"
	"github.com/iagoMAO/Botzin.OpenMASE/core"
	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
	"github.com/iagoMAO/Botzin.OpenMASE/shop"
	"github.com/iagoMAO/Botzin.OpenMASE/utils"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfg utils.Config

func main() {
	// Initialize readline for input
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/mase.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Error().Msgf("Error initializing readline: %s", err)
	}
	defer rl.Close()

	// First and foremost, load our config.
	cfg = utils.GetConfig()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: rl.Stdout()})

	// Load our database
	database.Initialize()

	defer database.DB.Close()

	// Create the listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.MASE_PORT))

	if err != nil {
		log.Error().Msgf("Listening error: %s", err)
		return
	}

	log.Info().Msgf("MASE - Successfully started listening on port %s.", cfg.MASE_PORT)

	// Close the socket once we're done
	defer listener.Close()

	// HackBuster - has no functionality (for this). Sole purpose is so the client connects and allows rounds to complete.
	go StartHB()
	go StartBuddyList()
	go StartServerList()

	go commandInput(rl)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Error().Msgf("Error thrown whilst accepting connection: %s", err)
			continue
		}

		go handleConnection(conn)
	}
}

func commandInput(rl *readline.Instance) {
	var regex = regexp.MustCompile(`^([a-zA-Z0-9_]+)\((.*)\)$`)

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := regex.FindStringSubmatch(line)
		if len(matches) < 3 {

		}

		command := matches[1]
		argsStr := matches[2]

		var args []string
		if strings.TrimSpace(argsStr) != "" {
			reader := csv.NewReader(strings.NewReader(argsStr))
			reader.TrimLeadingSpace = true

			parsed, err := reader.Read()
			if err != nil {
				continue
			}
			args = parsed
		}
		switch command {
		case "broadcast":
			if len(args) != 2 {
				continue
			}

			message := args[0]
			color := args[1]

			buddy.BroadcastMessage(message, color)
		}
	}
}

func handleConnection(conn net.Conn) {
	// Close once we're done, again
	defer core.RemoveSession(conn)
	defer conn.Close()

	// TODO: Maybe make this configurable?
	buf := make([]byte, 1024)

	reader := bufio.NewReader(conn)

PacketLoop:
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
			payloadStr := string(message.Payload)

			request := packets.LoginRequestPacket{
				EXECRC:       payloadStr[0:64],
				PasswordHash: payloadStr[64:96],
				UsernameHash: payloadStr[96:128],
				ClanTag:      payloadStr[128:],
			}

			// Handle Login request
			id, login := authentication.Login(request)
			log.Debug().Msgf("Received Login Answer: %s", hex.EncodeToString(login.Compose()))

			conn.Write(login.Compose())

			if id != 0 {
				core.RegisterConnection(conn, id, core.ConnTypeMASE)
			} else {
				conn.Close()
			}
		case protocol.UserBootRequest:
			session := core.GetSession(conn)

			if session == nil {
				continue PacketLoop
			}

			session.SetStatus(protocol.BUDDY_STATUS_ONLINE)

			user := authentication.GetUserInfoPacket(session.UserId)
			attribs := avatar.GetAvatarAttrib(session.UserId)
			avatar := avatar.GetAvatarInfo(session.UserId)
			guiPacket := packets.MaseShowGUIAnswerPacket{StatusCode: protocol.MASE_OK}
			broadcastPacket := packets.BroadcastAnswerPacket{StatusCode: protocol.MASE_OK, MessageColor: 102, MessageText: cfg.WELCOME_BROADCAST_MESSAGE}

			conn.Write(user.Compose())
			conn.Write(attribs.Compose())
			conn.Write(avatar.Compose())
			conn.Write(guiPacket.Compose())
			conn.Write(broadcastPacket.Compose())

			log.Debug().Msgf("Received User Boot Request: %s from User %d", hex.EncodeToString(message.Payload), session.UserId)
		case protocol.UserDataRequest:
			session := core.GetSession(conn)

			if session == nil {
				continue PacketLoop
			}

			user := authentication.GetUserInfoPacket(session.UserId)
			conn.Write(user.Compose())

			log.Debug().Msgf("Received User Data Request: %s from User %d", hex.EncodeToString(message.Payload), session.UserId)
		case protocol.ShopBuyRequest:
			session := core.GetSession(conn)

			if session == nil {
				continue PacketLoop
			}

			log.Debug().Msgf("Received Shop Buy Request from User Id %d", session.UserId)

			itemId, err := strconv.Atoi(data.SCR_UnpackInt(message.Payload[1:]))

			if err != nil {
				continue PacketLoop
			}

			packet := shop.BuyItem(session.UserId, itemId)
			conn.Write(packet.Compose())
		case protocol.ServerQueryAvatarRequest:
			session := core.GetSession(conn)

			if session == nil {
				continue PacketLoop
			}

			userId, err := strconv.Atoi(data.SCR_UnpackInt(message.Payload[1:]))

			if err != nil {
				continue PacketLoop
			}

			session.SetStatus(protocol.BUDDY_STATUS_INGAME)
			contacts := buddy.GetContacts(userId)
			for _, contact := range contacts {
				userSession := core.GetSessionByUserId(contact.GUID)
				if userSession != nil {
					BootBuddyRequest(userSession)
				}
			}

			avatarData := avatar.GetAvatarSetupData(userId)

			conn.Write(avatarData.Compose())
		case protocol.AvatarSetupSaveRequest:
			session := core.GetSession(conn)

			if session == nil {
				continue PacketLoop
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			var itemIds []int
			for _, raw := range parts {
				itemIds = append(itemIds, data.SCR_StrToInt(raw))
			}

			log.Debug().Msgf("Received Avatar Setup Save Request from User Id %d %s", session.UserId, itemIds)

			request := packets.AvatarSetupSaveRequestPacket{
				ItemIds: itemIds,
			}

			avatarSaveResponse := avatar.SaveAvatarSetup(session.UserId, request)

			conn.Write(avatarSaveResponse.Compose())
		case protocol.AvatarAttribSaveRequest:
			session := core.GetSession(conn)

			if session == nil {
				continue PacketLoop
			}

			if err != nil {
				continue PacketLoop
			}

			parts := bytes.Split(message.Payload[1:], []byte{'\t'})

			if len(parts) < 5 {
				continue PacketLoop
			}

			request := packets.AvatarAttribSaveRequestPacket{
				BotId: data.SCR_StrToInt(parts[0]),
				ST:    data.SCR_StrToInt(parts[1]),
				DX:    data.SCR_StrToInt(parts[2]),
				IQ:    data.SCR_StrToInt(parts[3]),
				HT:    data.SCR_StrToInt(parts[4]),
			}

			attribSaveResponse := avatar.SaveAvatarAttrib(session.UserId, request)

			conn.Write(attribSaveResponse.Compose())

		default:
			session := core.GetSession(conn)

			if session != nil {
				log.Debug().Msgf("Received Packet %s", hex.Dump(message.Payload))
			}
		}
	}
}
