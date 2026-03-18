package shop

import (
	"github.com/rs/zerolog/log"

	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
)

func BuyItem(userId int, itemId int) packets.ShopBuyAnswerPacket {
	var userCredits, userGold, itemCredits, itemGold, count int

	row := database.DB.QueryRow("SELECT COUNT(*) FROM user_items WHERE item_id = ? AND user_id = ?", itemId, userId)

	err := row.Scan(&count)

	if err != nil || count != 0 {
		log.Debug().Msgf("Error: %d - %d", itemId, userId)
		return packets.ShopBuyAnswerPacket{ShopBuyAnswerType: protocol.SHOP_ALREADY_HAVE}
	}

	row = database.DB.QueryRow("SELECT credits, gold FROM users WHERE id = ?", userId)

	err = row.Scan(&userCredits, &userGold)

	if err != nil {
		return packets.ShopBuyAnswerPacket{}
	}

	row = database.DB.QueryRow("SELECT credits, gold FROM items WHERE id = ?", itemId)

	item := packets.AvatarItemData{
		Id:     itemId,
		UserId: userId,
	}

	err = row.Scan(&itemCredits, &itemGold)

	if err != nil {
		log.Debug().Msgf("Error: %s", err.Error())
		return packets.ShopBuyAnswerPacket{}
	}

	if userCredits < itemCredits {
		return packets.ShopBuyAnswerPacket{ShopBuyAnswerType: protocol.SHOP_NO_CREDITS}
	}

	if userGold < itemGold {
		return packets.ShopBuyAnswerPacket{ShopBuyAnswerType: protocol.SHOP_NO_GOLD}
	}

	query := "UPDATE users SET credits = ?, gold = ? WHERE id = ?"

	_, err = database.DB.Exec(query, (userCredits - itemCredits), (userGold - itemGold), userId)

	if err != nil {
		log.Debug().Msgf("Error: %s", err.Error())
	}

	query = "INSERT INTO user_items (item_id, user_id, enabled) VALUES(?, ?, ?)"

	_, err = database.DB.Exec(query, itemId, userId, 0)

	if err != nil {
		log.Debug().Msgf("Error: %s", err.Error())
	}

	log.Debug().Msgf("Item: %d", itemId)

	return packets.ShopBuyAnswerPacket{
		ShopBuyAnswerType: protocol.SHOP_BUY_DONE,
		Item: packets.AvatarItemData{
			Id:      itemId,
			UserId:  userId,
			Class:   item.Class,
			ST:      item.ST,
			DX:      item.DX,
			IQ:      item.IQ,
			HT:      item.HT,
			Payload: item.Payload,
			TheGen:  1,
			Enabled: 0,
		},
	}
}
