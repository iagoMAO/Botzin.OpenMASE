package protocol

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/iagoMAO/Botzin.OpenMASE/security"
	"github.com/rs/zerolog/log"
)

type NetworkMessage struct {
	Type    PacketCode
	Payload []byte
}

type Packet interface {
	Compose() []byte
}

type PacketCode interface {
	Code() byte
}

func DecryptPacket(data []byte) NetworkMessage {
	// Decrypt the packet & rebuild it for request interpretation
	decrypt := security.DecryptXTEA(data[2:])

	if len(decrypt) < 2 {
		log.Error().Msgf("Received invalid packet data for decryption: %s", hex.EncodeToString(data))
		return NetworkMessage{}
	}

	packetLength := int(binary.BigEndian.Uint16(decrypt[0:2]))
	packetData := decrypt[2 : 2+packetLength]

	if packetLength < 2 {
		log.Error().Msgf("Received invalid packet length for decryption: %d", packetLength)
		return NetworkMessage{}
	}

	if len(packetData) < packetLength {
		log.Error().Msgf("Packet data length does not match sent length: (data) %d - (expected) %d", len(packetData), packetLength)
		return NetworkMessage{}
	}

	// Ok. We should now have the actual packet data.
	// Packet type
	rawType := packetData[0]
	packetType := PacketType(int(rawType))

	return NetworkMessage{packetType, packetData[1:]}
}

func EncryptPacket(packetType PacketCode, packetData []byte, statusCode StatusCode) []byte {
	// Prepend packetType + statusCode before packetData
	packetData = append([]byte{
		packetType.Code(),
		statusCode.Code(),
		0x09,
	}, packetData...)

	log.Debug().Msgf("before encryption: %s", hex.EncodeToString(packetData))

	length := len(packetData)
	lenBytes := make([]byte, 2)

	binary.BigEndian.PutUint16(lenBytes, uint16(length))

	md5 := security.EncryptMD5(packetData)

	input := append(append(lenBytes, packetData...), md5...)
	xtea := security.EncryptXTEA(input)

	length = len(xtea)
	binary.BigEndian.PutUint16(lenBytes, uint16(length))

	output := append(lenBytes, xtea...)

	log.Debug().Msgf("input: %s", hex.EncodeToString(input))
	log.Debug().Msgf("encrypted: %s", hex.EncodeToString(output))

	return output
}
