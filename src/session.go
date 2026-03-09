package main

import (
	"net"
	"sync"
)

// Session for a connected Player
type Session struct {
	UserId int
	Conn   net.Conn
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
