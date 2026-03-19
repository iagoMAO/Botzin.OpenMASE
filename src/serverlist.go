package main

import (
	"fmt"
	"net"
	"os"

	"github.com/iagoMAO/Botzin.OpenMASE/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func StartServerList() {
	// First and foremost, load our config.
	cfg := utils.GetConfig()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create the listener using ListenPacket for UDP
	listener, err := net.ListenPacket("udp", fmt.Sprintf("0.0.0.0:%s", cfg.SERVERLIST_PORT))

	if err != nil {
		log.Error().Msgf("Listening error: %s", err)
		return
	}

	log.Info().Msgf("SERVERLIST - Successfully started listening on port %s.", cfg.SERVERLIST_PORT)

	// Close the socket once we're done
	defer listener.Close()

	buffer := make([]byte, 1024)

	for {
		_, _, err := listener.ReadFrom(buffer)

		if err != nil {
			log.Error().Msgf("Error thrown whilst reading packet: %s", err)
			continue
		}
	}
}
