package content

import (
	"FPSProject/pkt"
	"FPSProject/utils"
	"log"
	"net"
	"sync"
	"time"
)

type Session struct {
	RoomNum   int32
	Users     []User
	IsRunning bool
	UserLock  sync.Mutex
}

type User struct {
	Conn            net.Conn
	Id              string
	RoomNum         int32
	NeedSync        bool
	CurrentLocation utils.Vec3
	RotationY       float32
	Session         *Session
}

func (s *Session) Init() {
	s.IsRunning = true
}

func (s *Session) UserEnter(usr User) {
	s.UserLock.Lock()
	s.Users = append(s.Users, usr)
	sm.Users[usr.Id] = &usr
	s.UserLock.Unlock()
	if len(s.Users) == MATCHINGNUM {
		log.Println("Match Complete Start Sync Move")
		s.StartSyncMove()
	} else if len(s.Users) > MATCHINGNUM {
		log.Println("MatchNum OverFlow")
	}
}

// func (s *Session) UserEnter(usr User) int32 {
// 	if len(s.Users) < MATCHINGNUM {
// 		id := len(s.Users) + 1
// 		usr.Id = int32(id)
// 		usr.NickName = fmt.Sprintf("%d", id)
// 		s.Users = append(s.Users, usr)
// 		return usr.Id
// 	} else {
// 		log.Println("UserEnter Overflow")
// 		sm.NewSession()
// 		s.UserEnter(usr)
// 		return -1
// 	}
// }

// func (s *Session) TempSessionEnter(usr User) {
// 	if
// }

func (s *Session) StartSyncMove() {

	// Spawn
	spawnPkt := pkt.R_PlayerSpawn{}
	for i, u := range s.Users {
		spawnPkt.PlayerIds = append(spawnPkt.PlayerIds, u.Id)
		spawnPkt.SpawnPoints = append(spawnPkt.SpawnPoints, int32(i))
	}
	spawnBuffer := utils.MakeSendBuffer("PlayerSpawn", spawnPkt)
	s.BroadCast(spawnBuffer)

	// PlayerMove, Rotation Sync
	go func() {
		for {
			time.Sleep(time.Millisecond * 200)

			if s.IsRunning {
				for _, u := range s.Users {
					if u.NeedSync {
						pkt := pkt.SR_PlayerRotation{PlayerId: u.Id, RotationY: u.RotationY}
						buffer := utils.MakeSendBuffer("PlayerRotation", pkt)
						s.BroadCast(buffer)
						u.NeedSync = false
					}
				}
			} else {
				log.Println("Game End", s.RoomNum)
				return
			}
		}
	}()
}

func (s *Session) BroadCast(buffer []byte) {
	for _, u := range s.Users {
		if u.Conn != nil {
			u.Conn.Write(buffer)
		} else {
			log.Println("No Connection:", string(buffer))
		}
	}
}
