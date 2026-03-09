package shop

import (
	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
)

func BuyItem(userId int, itemId int) packets.ShopBuyAnswerPacket {
	var userCredits, userGold, itemCredits, itemGold int

	row := database.DB.QueryRow("SELECT credits, gold FROM users WHERE id = ?", userId)

	err := row.Scan(&userCredits, &userGold)

	if err != nil {
		return packets.ShopBuyAnswerPacket{}
	}

	row = database.DB.QueryRow("SELECT class, payload, st, dx, iq, ht, credits, gold FROM items WHERE id = ?", itemId)

	item := packets.AvatarItemData{
		Id:     itemId,
		UserId: userId,
	}

	err = row.Scan(&item.Class, &item.Payload, &item.ST, &item.DX, &item.IQ, &item.HT, &itemCredits, &itemGold)

	if err != nil {
		return packets.ShopBuyAnswerPacket{}
	}

	if userCredits < itemCredits {
		return packets.ShopBuyAnswerPacket{ShopBuyAnswerType: protocol.SHOP_NO_CREDITS}
	}

	if userGold < itemGold {
		return packets.ShopBuyAnswerPacket{ShopBuyAnswerType: protocol.SHOP_NO_GOLD}
	}

	query := "INSERT INTO user_items (item_id, user_id, class, st, dx, iq, ht, payload, the_gen, enabled) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	database.DB.Exec(query, itemId, userId, item.Class, item.ST, item.DX, item.IQ, item.HT, item.Payload, 1, 0)

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
