package content

import (
	"sync"
)

type SessionManager struct {
	Sessions          sync.Map
	CurrentSessionNum int
}

func (sm *SessionManager) Init() {
}

// func (sm *SessionManager) NewPlayer(conn net.Conn) Session {

// 	return sm.Sessions[conn]
// }

// func (sm *SessionManager) FindSession(conn net.Conn) Session {
// 	return sm.Sessions[conn]
// }
