/*
Botzin.OpenMASE
@author: maldoliver
*/
package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/iagoMAO/Botzin.OpenMASE/authentication"
	"github.com/iagoMAO/Botzin.OpenMASE/avatar"
	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
	"github.com/iagoMAO/Botzin.OpenMASE/shop"
	"github.com/iagoMAO/Botzin.OpenMASE/utils"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// First and foremost, load our config.
	cfg := utils.GetConfig()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load our database
	database.Initialize()

	defer database.DB.Close()

	// Create the listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.MASE_PORT))

	if err != nil {
		log.Error().Msgf("Listening error: %s", err)
		return
	}

	log.Info().Msgf("Successfully started listening on port %s.", cfg.MASE_PORT)

	// Close the socket once we're done
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Error().Msgf("Error thrown whilst accepting connection: %s", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Close once we're done, again
	defer RemoveSession(conn)
	defer conn.Close()

	// TODO: Maybe make this configurable?
	buf := make([]byte, 1024)

	reader := bufio.NewReader(conn)

	for {
		n, err := reader.Read(buf)

		if err != nil {
			log.Error().Msgf("Read error: %s", err)
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
				CreateSession(conn, id)
			}
		case protocol.UserBootRequest:
			session := GetSession(conn)

			if session == nil {
				return
			}

			user := authentication.GetUserInfoPacket(session.UserId)
			attribs := avatar.GetAvatarAttrib(session.UserId)
			avatar := avatar.GetAvatarInfo(session.UserId)
			guiPacket := packets.MaseShowGUIAnswerPacket{StatusCode: protocol.MASE_OK}

			conn.Write(user.Compose())
			conn.Write(attribs.Compose())
			conn.Write(avatar.Compose())
			conn.Write(guiPacket.Compose())

			log.Debug().Msgf("Received User Boot Request: %s from User %d", hex.EncodeToString(message.Payload), session.UserId)
		case protocol.UserDataRequest:
			session := GetSession(conn)

			if session == nil {
				return
			}

			user := authentication.GetUserInfoPacket(session.UserId)
			conn.Write(user.Compose())

			log.Debug().Msgf("Received User Data Request: %s from User %d", hex.EncodeToString(message.Payload), session.UserId)
		case protocol.ShopBuyRequest:
			session := GetSession(conn)

			if session == nil {
				return
			}

			log.Debug().Msgf("Received Shop Buy Request from User Id %d", session.UserId)

			itemId, err := strconv.Atoi(data.SCR_UnpackInt(message.Payload[1:]))

			if err != nil {
				return
			}

			packet := shop.BuyItem(session.UserId, itemId)
			conn.Write(packet.Compose())
		case protocol.ServerQueryAvatarRequest:
			session := GetSession(conn)

			if session == nil {
				return
			}

			userId, err := strconv.Atoi(data.SCR_UnpackInt(message.Payload[1:]))

			if err != nil {
				return
			}

			avatarData := avatar.GetAvatarSetupData(userId)

			conn.Write(avatarData.Compose())
		default:
			session := GetSession(conn)

			if session != nil {
				log.Debug().Msgf("Received Packet Type %d from User Id %d", message.Type.Code(), session.UserId)
			}
		}
	}
}
