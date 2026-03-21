package packets

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"
	"github.com/rs/zerolog/log"
)

type LobbyLevel int

const (
	Novice   LobbyLevel = 1
	Advanced LobbyLevel = 2
	Pro      LobbyLevel = 3
	G2       LobbyLevel = 888
)

type Lobby struct {
	Level       LobbyLevel
	Address     string
	MaxPlayers  uint8
	Status      uint8
	BotCount    uint8
	PlayerCount uint8
	LastSeen    int64
}

type GameMasterInfoRequest struct {
	QueryFlags uint8
	Session    uint32
	Key        uint32
}

func (p GameMasterInfoRequest) Compose() []byte {
	var buf bytes.Buffer

	buf.WriteByte(byte(protocol.GameMasterInfoRequest))
	buf.WriteByte(p.QueryFlags) // flags
	buf.Write(data.U32ToBytes(p.Session))
	buf.Write(data.U32ToBytes(p.Key))

	return buf.Bytes()
}

type GameMasterInfoResponse struct {
	QueryFlags   uint8
	Session      uint32
	Key          uint32
	GameType     string
	MissionType  string
	MaxPlayers   uint8
	RegionMask   uint32
	Version      uint32
	Status       uint8
	BotCount     uint8
	ProcessorMhz uint32
	PlayerCount  uint8
}

type MasterServerListRequest struct {
	QueryFlags  uint8
	Session     uint32
	Key         uint32
	GameType    string
	MissionType string
	MinPlayers  uint8
	MaxPlayers  uint8
	RegionMask  uint8
	Version     uint32
	FilterFlags uint8
	MaxBots     uint8
	MinCPU      uint16
	BuddyCount  uint8
}

type MasterServerListResponse struct {
	PacketIndex uint8
	PacketTotal uint8
	Session     uint32
	Key         uint32
	Lobbies     []Lobby
}

func (p MasterServerListResponse) Compose() []byte {
	var buf bytes.Buffer

	buf.WriteByte(byte(protocol.MasterServerListResponse))
	buf.WriteByte(0) // flags
	buf.Write(data.U32ToBytes(p.Session))
	buf.Write(data.U32ToBytes(p.Key))
	buf.WriteByte(p.PacketIndex)
	buf.WriteByte(p.PacketTotal)
	buf.Write(data.U16ToBytes(uint16(len(p.Lobbies))))

	log.Info().Msgf("LOBBIES LENGTH %d", uint8(len(p.Lobbies)))

	for _, lobby := range p.Lobbies {
		parts := strings.Split(lobby.Address, ":")
		strAddr := strings.Split(parts[0], ".")
		strPort := parts[1]

		intAddr := []uint8{}

		for i := 0; i < len(strAddr); i++ {
			s := strAddr[i]
			i, err := strconv.Atoi(s)
			if err != nil {
				log.Debug().Msgf("error: %s", err)
				break
			}
			intAddr = append(intAddr, uint8(i))
		}

		intPort, err := strconv.Atoi(strPort)

		if err != nil {
			log.Debug().Msgf("error: %s", err)
			continue
		}

		buf.WriteByte(intAddr[0])
		buf.WriteByte(intAddr[1])
		buf.WriteByte(intAddr[2])
		buf.WriteByte(intAddr[3])
		buf.Write(data.U16ToBytes(uint16(intPort)))
	}

	if len(p.Lobbies) <= 0 {
		buf.Write([]byte{0, 0, 0, 0})
		buf.Write(data.U16ToBytes(0))
	}

	return buf.Bytes()
}
