package main

import (
	"net"
	"sync"
)

type ConnType int

const (
	ConnTypeMASE ConnType = iota
	ConnTypeBUDDY
)

// Session for a connected Player
type Session struct {
	UserId    int
	Status    BuddyStatus
	Conn      net.Conn
	BuddyConn net.Conn
}

var (
	sessionsByUser = make(map[int]*Session)
	sessionsByConn = make(map[net.Conn]*Session)
	sessionMutex   sync.RWMutex
)

func (s *Session) SetStatus(status BuddyStatus) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	s.Status = status
}

func RegisterConnection(conn net.Conn, userId int, connType ConnType) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	session, exists := sessionsByUser[userId]
	if !exists {
		session = &Session{UserId: userId}
		sessionsByUser[userId] = session
	}

	if connType == ConnTypeMASE {
		session.Conn = conn
	} else if connType == ConnTypeBUDDY {
		session.BuddyConn = conn
	}

	sessionsByConn[conn] = session
}

func GetSessionByUserId(userId int) *Session {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	return sessionsByUser[userId]
}

func GetAllSessions() []*Session {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	sessions := make([]*Session, 0, len(sessionsByUser))
	for _, v := range sessionsByUser {
		sessions = append(sessions, v)
	}
	return sessions
}

func GetSession(conn net.Conn) *Session {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	return sessionsByConn[conn]
}

func RemoveSession(conn net.Conn) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	session, exists := sessionsByConn[conn]
	if !exists {
		return
	}

	delete(sessionsByConn, conn)

	if session.Conn == conn {
		session.Conn = nil
	} else if session.BuddyConn == conn {
		session.BuddyConn = nil
	}

	if session.Conn == nil && session.BuddyConn == nil {
		delete(sessionsByUser, session.UserId)
	}
}
