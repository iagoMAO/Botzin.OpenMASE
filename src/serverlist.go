package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
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

	// Lobby list
	lobbies := []packets.Lobby{}

	// Close the socket once we're done
	defer listener.Close()

	buffer := make([]byte, 1024)

PacketLoop:
	for {
		_, addr, err := listener.ReadFrom(buffer)

		if err != nil {
			log.Error().Msgf("Error thrown whilst reading packet: %s", err)
			continue PacketLoop
		}

		switch protocol.PacketType(buffer[0]) {
		default:
			log.Debug().Msgf("packet type %d", buffer[0])
		case protocol.GameMasterInfoResponse:
			log.Debug().Msg("server info response")

			var offset = 1

			response := packets.GameMasterInfoResponse{}

			response.QueryFlags = buffer[offset]
			offset++

			response.Session = binary.LittleEndian.Uint32(buffer[offset : offset+4])
			offset += 4

			response.Key = binary.LittleEndian.Uint32(buffer[offset : offset+4])
			offset += 4

			response.GameType = string(buffer[offset+1 : offset+int(buffer[offset])+1])
			offset += int(buffer[offset]) + 1

			response.MissionType = string(buffer[offset+1 : offset+int(buffer[offset])+1])
			offset += int(buffer[offset]) + 1

			response.MaxPlayers = uint8(buffer[offset])
			offset++

			response.RegionMask = binary.LittleEndian.Uint32(buffer[offset : offset+4])
			offset += 4

			response.Version = binary.LittleEndian.Uint32(buffer[offset : offset+4])
			offset += 4

			response.Status = uint8(buffer[offset])
			offset++

			response.BotCount = uint8(buffer[offset])
			offset++

			response.ProcessorMhz = binary.LittleEndian.Uint32(buffer[offset : offset+4])
			offset += 4

			response.PlayerCount = uint8(buffer[offset])
			offset++

			missionType, err := strconv.Atoi(response.MissionType)

			if err != nil {
				continue PacketLoop
			}

			for _, lobby := range lobbies {
				if lobby.Address == addr.String() {
					continue PacketLoop
				}
			}

			newLobby := packets.Lobby{
				Level:       packets.LobbyLevel(missionType),
				Address:     addr.String(),
				MaxPlayers:  response.MaxPlayers,
				Status:      response.Status,
				BotCount:    response.BotCount,
				PlayerCount: response.PlayerCount,
				LastSeen:    time.Now().Unix(),
			}

			lobbies = append(lobbies, newLobby)

			log.Debug().Msgf("%#v", response)
		case protocol.GameHeartbeat:
			log.Debug().Msg("server heartbeat")

			var activeLobby *packets.Lobby

			for _, lobby := range lobbies {
				if lobby.Address == addr.String() {
					activeLobby = &lobby
					break
				} else {
					continue
				}
			}

			if activeLobby != nil {
				activeLobby.LastSeen = time.Now().Unix()
			}

			response := packets.GameMasterInfoRequest{QueryFlags: uint8(10), Session: uint32(0), Key: uint32(0)}

			packet := response.Compose()

			_, err = listener.WriteTo(packet, addr)
			if err != nil {
				log.Error().Msgf("Failed to send response to %s: %s", addr.String(), err)
			} else {
				log.Debug().Msgf("Successfully sent %d bytes to %s - \n%s", len(packet), addr.String(), hex.Dump(packet))
			}
		case protocol.MasterServerListRequest:
			log.Debug().Msg("server query")

			var offset = 1

			req := packets.MasterServerListRequest{}

			req.QueryFlags = buffer[offset]
			offset++

			req.Session = binary.LittleEndian.Uint32(buffer[offset : offset+4])
			offset += 4

			req.Key = binary.LittleEndian.Uint32(buffer[offset : offset+4])
			offset += 4

			// random 0xFF byte here
			offset += 1

			req.GameType = string(buffer[offset+1 : offset+int(buffer[offset])+1])
			offset += int(buffer[offset]) + 1

			req.MissionType = string(buffer[offset+1 : offset+int(buffer[offset])+1])
			offset += int(buffer[offset]) + 1

			log.Debug().Msgf("\n%s\n", hex.Dump(buffer))
			log.Debug().Msgf("queryFlags: %d, session: %d, gameType: %s, missionType: %s", req.QueryFlags, req.Session, req.GameType, req.MissionType)

			lobbyLevel, err := strconv.Atoi(req.MissionType)
			if err != nil {
				log.Debug().Msgf("error: %s", err)
				break
			}

			var filteredLobbies []packets.Lobby
			var activeLobbies []packets.Lobby
			now := time.Now()

			for _, lobby := range lobbies {
				lastSeen := time.Unix(lobby.LastSeen, 0)
				diff := now.Sub(lastSeen)

				if diff.Minutes() <= 3 {
					activeLobbies = append(activeLobbies, lobby)

					if int(lobby.Level) == lobbyLevel {
						filteredLobbies = append(filteredLobbies, lobby)
					}
				}
			}

			lobbies = activeLobbies

			response := packets.MasterServerListResponse{
				PacketIndex: 0,
				PacketTotal: 1,
				Session:     req.Session,
				Key:         req.Key,
				Lobbies:     filteredLobbies,
			}

			packet := response.Compose()

			_, err = listener.WriteTo(packet, addr)
			if err != nil {
				log.Error().Msgf("Failed to send response to %s: %s", addr.String(), err)
			} else {
				log.Debug().Msgf("Successfully sent %d bytes to %s - \n%s", len(packet), addr.String(), hex.Dump(packet))
			}
		}
	}
}

type GameHeartbeat struct {
	queryFlags uint8
	session    uint32
	key        uint32
}
