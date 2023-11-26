package content

import (
	"FPSProject/utils"
	"net"
	"sync"
)

type SessionManager struct {
	Sessions          map[int32]*Session
	CurrentSessionNum int32
	SessionLock       *sync.Mutex
	// PlayerSession     map[int32]*Session
	Users map[string]*User
}

func (sm *SessionManager) Init() {
	sm.SessionLock = &sync.Mutex{}
	sm.Sessions = map[int32]*Session{}
	sm.Users = map[string]*User{}
	sm.CurrentSessionNum = 0
	// sm.NewSession()
}

func (sm *SessionManager) TempNewSessionEnter(RoomNum int32, Id string, Conn net.Conn) {
	sm.SessionLock.Lock()

	if s, ok := sm.Sessions[RoomNum]; ok {
		s.UserEnter(User{Conn: Conn, Id: Id, RoomNum: RoomNum, Session: s, SpawnIndex: int32(len(s.Users)), Health: 100})
	} else {
		if r2, ok := ph.Room.Load(RoomNum); ok {
			matchingNum := r2.(*RoomInfo).NumberOfPeople

			s := Session{
				RoomNum:    RoomNum,
				Users:      []User{},
				SpawnIndex: utils.RandomInt32(int(matchingNum), 0, int(matchingNum)-1),
				PlayerNum:  matchingNum,
			}
			s.Init()
			s.UserEnter(User{Conn: Conn, Id: Id, RoomNum: RoomNum, Session: &s, SpawnIndex: 0, Health: 100})
			sm.Sessions[RoomNum] = &s
		}
	}
	sm.SessionLock.Unlock()
}

// func (sm *SessionManager) NewSession() {
// 	sm.SessionLock.Lock()
// 	n := sm.CurrentSessionNum
// 	sm.SessionLock.Unlock()
// 	log.Println("NEW SESSION", sm.CurrentSessionNum)
// 	sm.Sessions.Store(n, &Session{
// 		RoomNum: n,
// 		Users:   []User{},
// 	})
// }

// func (sm *SessionManager) NewPlayer(conn net.Conn) int32 {
// 	n := sm.CurrentSessionNum
// 	if s, ok := sm.Sessions.Load(n); ok {
// 		if len(s.(*Session).Users) < MATCHINGNUM {
// 			// sm.PlayerSession
// 			// return s.(*Session).UserEnter(User{Conn: conn, RoomNum: n})
// 		} else {
// 			sm.NewSession()
// 			sm.NewPlayer(conn)
// 		}
// 	}
// 	return -1
// }

// func (sm *SessionManager) FindSession(conn net.Conn) Session {
// 	return sm.Sessions[conn]
// }
