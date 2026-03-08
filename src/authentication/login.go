package authentication

import (
	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
)

// We receive a login request packet & return the answer packet
func Login(packet packets.LoginRequestPacket) (int, protocol.Packet) {
	var id int
	var username string
	var password string

	database.DB.QueryRow("SELECT id, username, password_md5 FROM users WHERE username_md5 = ?", packet.UsernameHash, packet.PasswordHash).Scan(&id, &username, &password)

	// Invalid password or inexistent user (the client doesn't really handle both seperately)
	if password != packet.PasswordHash {
		return 0, packets.LoginErrorPacket{
			StatusCode: protocol.MASE_INVALID_LOGINPASS,
		}
	} else {
		// We can login!
		return id, packets.LoginAnswerPacket{
			StatusCode:  protocol.MASE_OK,
			MagicNumber: uint32(id),
			ClientGUID:  uint32(id),
		}
	}

	return 0, packets.LoginAnswerPacket{}
}
