package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/utils"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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

	// TODO: Maybe make this configurable?
	buf := make([]byte, 1024)

	reader := bufio.NewReader(conn)

	for {
		n, err := reader.Read(buf)

		if err != nil {
			log.Error().Msgf("Read error: %s", err)
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
		case protocol.BootStatusRequest:
			log.Debug().Msgf("Received Boot Status: %s", hex.EncodeToString(message.Payload))
		default:
			// log.Debug().Msgf("Received Packet %s - %d", hex.Dump(message.Payload), message.Type.Code())
		}
	}
}
