package buddy

import (
	"github.com/iagoMAO/Botzin.OpenMASE/core"
	"github.com/iagoMAO/Botzin.OpenMASE/database"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol"
	"github.com/iagoMAO/Botzin.OpenMASE/protocol/packets"
)

func GetCategorizedContacts(userId int) packets.BuddyStatusList {
	contacts := GetContacts(userId)
	capacity := len(contacts)

	onlineContacts := make([]packets.BuddyContactInfo, 0, capacity)
	ingameContacts := make([]packets.BuddyContactInfo, 0, capacity)
	offlineContacts := make([]packets.BuddyContactInfo, 0, capacity)

	for _, contact := range contacts {
		session := core.GetSessionByUserId(contact.GUID)
		if session == nil {
			offlineContacts = append(offlineContacts, contact)
			continue
		}

		switch session.Status {
		case protocol.BUDDY_STATUS_ONLINE:
			onlineContacts = append(onlineContacts, contact)
		case protocol.BUDDY_STATUS_INGAME:
			ingameContacts = append(ingameContacts, contact)
		default:
			offlineContacts = append(offlineContacts, contact)
		}
	}

	return packets.BuddyStatusList{Online: onlineContacts, Ingame: ingameContacts, Offline: offlineContacts}
}

func RespondContact(userId int, contactId int, status protocol.BuddyStatus) {
	query := `SELECT COUNT(*) FROM contacts WHERE user_id = ? AND contact_id = ?`

	var existing int
	err := database.DB.QueryRow(query, contactId, userId).Scan(&existing)

	if err != nil {
		return
	}

	if existing <= 0 {
		return
	}

	query = `UPDATE contacts SET status = ? WHERE user_id = ? AND contact_id = ?`

	_, err = database.DB.Exec(query, status, contactId, userId)

	if err != nil {
		return
	}
}

func AddContact(userId int, contactId int) error {
	err := database.DB.QueryRow("SELECT id FROM users WHERE id = ?", contactId).Scan(&contactId)

	if err != nil {
		return err
	}

	var existing int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM contacts WHERE (user_id = ? AND contact_id = ?) OR (contact_id = ? AND user_id = ?)", userId, contactId, userId, contactId).Scan(&existing)

	if err != nil {
		return err
	}

	if existing > 0 {
		return err
	}

	query := "INSERT INTO contacts (user_id, contact_id, status) VALUES(?, ?, ?)"

	_, err = database.DB.Exec(query, userId, contactId, protocol.BUDDY_ANSWER_REQUEST)

	if err != nil {
		return err
	}

	return nil
}

func GetIncomingRequests(userId int) []packets.BuddyContactInfo {
	query := `
		SELECT c.contact_id AS incoming_id, u.username AS incoming_name
		FROM contacts c
		JOIN users u ON c.contact_id = u.id
		WHERE c.user_id = ? AND c.status = ?
	`
	rows, err := database.DB.Query(query, userId, protocol.BUDDY_ANSWER_REQUEST)

	if err != nil {
		return []packets.BuddyContactInfo{}
	}

	defer rows.Close()

	var contacts []packets.BuddyContactInfo

	for rows.Next() {
		var contact packets.BuddyContactInfo
		if err := rows.Scan(&contact.GUID, &contact.Name); err != nil {
			return []packets.BuddyContactInfo{}
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return []packets.BuddyContactInfo{}
	}

	return contacts
}

func GetOutgoingRequests(userId int, status protocol.BuddyStatus) []packets.BuddyContactInfo {
	query := `
		SELECT c.user_id AS friend_id, u.username AS friend_name 
		FROM contacts c 
		JOIN users u ON c.user_id = u.id 
		WHERE c.contact_id = ? AND c.status = ?`

	rows, err := database.DB.Query(query, userId, status)

	if err != nil {
		return []packets.BuddyContactInfo{}
	}

	defer rows.Close()

	var contacts []packets.BuddyContactInfo

	for rows.Next() {
		var contact packets.BuddyContactInfo
		if err := rows.Scan(&contact.GUID, &contact.Name); err != nil {
			return []packets.BuddyContactInfo{}
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return []packets.BuddyContactInfo{}
	}

	return contacts
}

func GetContacts(userId int) []packets.BuddyContactInfo {
	query := `
		SELECT c.contact_id AS friend_id, u.username AS friend_name 
		FROM contacts c 
		JOIN users u ON c.contact_id = u.id 
		WHERE c.user_id = ? AND status = ?
		
		UNION
		
		SELECT c.user_id AS friend_id, u.username AS friend_name 
		FROM contacts c 
		JOIN users u ON c.user_id = u.id 
		WHERE c.contact_id = ? AND status = ?`

	rows, err := database.DB.Query(query, userId, protocol.BUDDY_ANSWER_ACCEPTED, userId, protocol.BUDDY_ANSWER_ACCEPTED)

	if err != nil {
		return []packets.BuddyContactInfo{}
	}

	defer rows.Close()

	var contacts []packets.BuddyContactInfo

	for rows.Next() {
		var contact packets.BuddyContactInfo
		if err := rows.Scan(&contact.GUID, &contact.Name); err != nil {
			return []packets.BuddyContactInfo{}
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return []packets.BuddyContactInfo{}
	}

	return contacts
}

func QueryContacts(query string) []packets.BuddyContactInfo {
	rows, err := database.DB.Query("SELECT id, username FROM users WHERE username LIKE ?", "%"+query+"%")

	if err != nil {
		return []packets.BuddyContactInfo{}
	}

	defer rows.Close()

	var contacts []packets.BuddyContactInfo

	for rows.Next() {
		var contact packets.BuddyContactInfo
		if err := rows.Scan(&contact.GUID, &contact.Name); err != nil {
			return []packets.BuddyContactInfo{}
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return []packets.BuddyContactInfo{}
	}

	return contacts
}
