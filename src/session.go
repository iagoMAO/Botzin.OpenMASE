package main

import (
	"net"
	"sync"
)

// Session for a connected Player
type Session struct {
	UserId int
	Status BuddyStatus
	Conn   net.Conn
}

func (s Session) SetStatus(status BuddyStatus) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	s.Status = status
	activeSessions[s.Conn] = &s
}

var (
	activeSessions = make(map[net.Conn]*Session)
	sessionMutex   sync.RWMutex
)

func CreateSession(conn net.Conn, userId int) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	activeSessions[conn] = &Session{
		UserId: userId,
		Conn:   conn,
	}
}

func GetSessionByUserId(userId int) *Session {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	for _, session := range activeSessions {
		if session.UserId == userId {
			return session
		}
		return nil
	}

	return nil
}

func GetSession(conn net.Conn) *Session {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	return activeSessions[conn]
}

func RemoveSession(conn net.Conn) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	delete(activeSessions, conn)
}
