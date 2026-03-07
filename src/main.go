/*
Botzin.OpenMASE
@author: maldoliver
*/
package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/iagoMAO/Botzin.OpenMASE/security"
	"github.com/iagoMAO/Botzin.OpenMASE/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// First and foremost, load our config.
	cfg := utils.GetConfig()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create the listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.MASE_PORT))

	if err != nil {
		log.Error().Msgf("Listening error: %s", err)
	}

	log.Info().Msgf("Successfully started listening on port %s.", cfg.MASE_PORT)

	packet := []byte("a")

	length := len(packet)
	lenBytes := make([]byte, 2)

	binary.BigEndian.PutUint16(lenBytes, uint16(length))

	md5 := security.EncryptMD5(packet)

	input := append(append(lenBytes, packet...), md5...)
	xtea := security.EncryptXTEA(input)

	length = len(xtea)
	binary.BigEndian.PutUint16(lenBytes, uint16(length))

	output := append(lenBytes, xtea...)

	log.Debug().Msgf("input: %s", hex.EncodeToString(input))
	log.Debug().Msgf("test: %s", hex.EncodeToString(output))

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

		log.Debug().Msgf("RCV: %s\n", hex.EncodeToString(buf[:n]))
	}
}
