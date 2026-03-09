package avatar

import (
	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
	"github.com/iagoMAO/Botzin.OpenMASE/utils/data"
	"github.com/rs/zerolog/log"
)

type AvatarAttrib struct {
	XP int
	ST int
	DX int
	IQ int
	HT int
}

func GetAvatarInfo(id int) packets.AvatarSetupLoadAnswerPacket {
	rows, err := database.DB.Query("SELECT item_id, user_id, class, st, dx, iq, ht, payload, the_gen, enabled FROM user_items WHERE user_id = ?", id)

	if err != nil {
		log.Error().Msgf("error: %s", err)
		return packets.AvatarSetupLoadAnswerPacket{}
	}

	defer rows.Close()

	var items []packets.AvatarItemData

	for rows.Next() {
		var item packets.AvatarItemData
		err := rows.Scan(
			&item.Id,
			&item.UserId,
			&item.Class,
			&item.ST,
			&item.DX,
			&item.IQ,
			&item.HT,
			&item.Payload,
			&item.TheGen,
			&item.Enabled,
		)

		if err != nil {
			log.Error().Msgf("error: %s", err)
			return packets.AvatarSetupLoadAnswerPacket{}
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Error().Msgf("error: %s", err)
		return packets.AvatarSetupLoadAnswerPacket{}
	}

	packet := packets.AvatarSetupLoadAnswerPacket{
		TotalAvatarItems: len(items),
		Items:            items,
	}

	return packet
}

func GetAvatarAttrib(id int) packets.AvatarAttribLoadAnswerPacket {
	row := database.DB.QueryRow("SELECT xp, st, dx, iq, ht FROM users WHERE id = ?", id)
	var u AvatarAttrib
	err := row.Scan(
		&u.XP,
		&u.ST,
		&u.DX,
		&u.IQ,
		&u.HT,
	)

	if err != nil {
		log.Error().Msgf("Error: %s", err)
		return packets.AvatarAttribLoadAnswerPacket{}
	}

	return packets.AvatarAttribLoadAnswerPacket{
		XP: string(data.SCR_PackInt(u.XP)),
		ST: string(data.SCR_PackInt(u.ST)),
		DX: string(data.SCR_PackInt(u.DX)),
		IQ: string(data.SCR_PackInt(u.IQ)),
		HT: string(data.SCR_PackInt(u.HT)),
	}
}
