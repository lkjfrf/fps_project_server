package content

import (
	"net"
	"sync"
)

type SessionManager struct {
	Sessions          *sync.Map
	CurrentSessionNum int32
	SessionLock       *sync.Mutex
}

func (sm *SessionManager) Init() {
	sm.SessionLock = &sync.Mutex{}
	sm.Sessions = &sync.Map{}
	sm.CurrentSessionNum = 0
}

func (sm *SessionManager) NewSession() {
	sm.SessionLock.Lock()
	n := sm.CurrentSessionNum
	sm.SessionLock.Unlock()
	sm.Sessions.Store(n, &Session{
		SessionId: n,
		Users:     []User{},
	})
}

func (sm *SessionManager) NewPlayer(conn net.Conn) int32 {
	sm.Sessions.Load(sm.CurrentSessionNum)
	return 2
}

// func (sm *SessionManager) FindSession(conn net.Conn) Session {
// 	return sm.Sessions[conn]
// }
