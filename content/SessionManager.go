package content

import (
	"log"
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
	log.Println("NEW SESSION", sm.CurrentSessionNum)
	sm.Sessions.Store(n, &Session{
		SessionId: n,
		Users:     []User{},
	})
}

func (sm *SessionManager) NewPlayer(conn net.Conn) int32 {
	n := sm.CurrentSessionNum
	if s, ok := sm.Sessions.Load(n); ok {
		if len(s.(*Session).Users) < MATCHINGNUM {
			s.(*Session).UserEnter(User{Conn: conn, SessionId: n})
		} else {
			sm.NewSession()
			sm.NewPlayer(conn)
		}
	}
	return 2
}

// func (sm *SessionManager) FindSession(conn net.Conn) Session {
// 	return sm.Sessions[conn]
// }
