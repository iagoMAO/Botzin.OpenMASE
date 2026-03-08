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

	"github.com/iagoMAO/Botzin.OpenMASE/authentication"
	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
	"github.com/iagoMAO/Botzin.OpenMASE/utils"

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
				user := authentication.GetUserInfo(id)
				conn.Write(user.Compose())

				guiPacket := packets.MaseShowGUIAnswerPacket{StatusCode: protocol.MASE_OK}
				conn.Write(guiPacket.Compose())
			}

		case protocol.UserBootRequest:
			log.Debug().Msgf("Received User Boot Request: %s", hex.EncodeToString(message.Payload))
		case protocol.UserDataRequest:
			log.Debug().Msgf("Received User Data Request: %s", hex.EncodeToString(message.Payload))
		}
	}
}
