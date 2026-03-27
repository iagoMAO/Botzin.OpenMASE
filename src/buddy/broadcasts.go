package buddy

import (
	"fmt"

	"github.com/iagoMAO/Botzin.OpenMASE/core"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
)

func BroadcastMessage(message string, color string) {
	if color != "" {
		message = fmt.Sprintf("<color:%s>%s", color, message)
	}

	packet := packets.BroadcastAnswerPacket{StatusCode: protocol.MASE_OK, MessageColor: 102, MessageText: message}
	sessions := core.GetAllSessions()

	for _, session := range sessions {
		session.Conn.Write(packet.Compose())
	}
}
