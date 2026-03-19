package main

import (
	"fmt"
	"net"
	"os"

	"github.com/iagoMAO/Botzin.OpenMASE/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func StartHB() {
	// First and foremost, load our config.
	cfg := utils.GetConfig()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create the listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.HB_PORT))

	if err != nil {
		log.Error().Msgf("Listening error: %s", err)
		return
	}

	log.Info().Msgf("HB - Successfully started listening on port %s.", cfg.HB_PORT)

	// Close the socket once we're done
	defer listener.Close()

	for {
		_, err := listener.Accept()

		if err != nil {
			log.Error().Msgf("Error thrown whilst accepting connection: %s", err)
			continue
		}
	}
}
