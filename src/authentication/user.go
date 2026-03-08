package authentication

import (
	"github.com/rs/zerolog/log"

	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
)

type UserInfo struct {
	Nick        string
	XP          int
	ST          int
	DX          int
	IQ          int
	HT          int
	PromoButton string
	Points      string
	Credits     string
	Gold        string
	Ranking     string
	TotalRK     string
	Level       string
	PMX         string
}

func GetUserInfo(id int) packets.UserDataAnswerPacket {
	row := database.DB.QueryRow("SELECT username, xp, st, dx, iq, ht, points, credits, gold, ranking, totalRK, level, pmx FROM users WHERE id = ?", id)
	var u UserInfo
	err := row.Scan(
		&u.Nick,
		&u.XP,
		&u.ST,
		&u.DX,
		&u.IQ,
		&u.HT,
		&u.Points,
		&u.Credits,
		&u.Gold,
		&u.Ranking,
		&u.TotalRK,
		&u.Level,
		&u.PMX,
	)

	if err != nil {
		log.Error().Msgf("Error: %s", err)
		return packets.UserDataAnswerPacket{}
	}

	return packets.UserDataAnswerPacket{
		StatusCode: protocol.MASE_OK,
		UserInfo: packets.PacketUserInfo{
			Nick:        u.Nick,
			XP:          u.XP,
			ST:          u.ST,
			DX:          u.DX,
			IQ:          u.IQ,
			HT:          u.HT,
			PromoButton: "",
			Points:      u.Points,
			Credits:     u.Credits,
			Gold:        u.Gold,
			Ranking:     u.Ranking,
			TotalRK:     u.TotalRK,
			Level:       u.Level,
			PMX:         u.PMX,
		},
	}
}
