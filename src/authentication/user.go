package authentication

import (
	"github.com/rs/zerolog/log"

	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"
)

type UserInfo struct {
	Nick        string
	XP          int
	ST          int
	DX          int
	IQ          int
	HT          int
	PromoButton string
	Points      int
	Credits     int
	Gold        int
	Ranking     int
	TotalRK     int
	Level       int
	PMX         int
}

func GetUserInfo(id int) packets.PacketUserInfo {
	row := database.DB.QueryRow("SELECT username, xp, st, dx, iq, ht, level, ranking FROM users WHERE id = ?", id)
	var u packets.PacketUserInfo
	err := row.Scan(
		&u.Nick,
		&u.XP,
		&u.ST,
		&u.DX,
		&u.IQ,
		&u.HT,
		&u.Level,
		&u.Ranking,
	)

	if err != nil {
		log.Error().Msgf("Error: %s", err)
		return packets.PacketUserInfo{}
	}

	return u
}

func GetUserInfoPacket(id int) packets.UserDataAnswerPacket {
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
			Points:      string(data.SCR_PackInt(u.Points)),
			Credits:     string(data.SCR_PackInt(u.Credits)),
			Gold:        string(data.SCR_PackInt(u.Gold)),
			Ranking:     string(data.SCR_PackInt(u.Ranking)),
			TotalRK:     string(data.SCR_PackInt(u.TotalRK)),
			Level:       string(data.SCR_PackInt(u.Level)),
			PMX:         string(data.SCR_PackInt(u.PMX)),
		},
	}
}
