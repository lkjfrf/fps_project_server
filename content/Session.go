package content

import (
	"FPSProject/pkt"
	"FPSProject/utils"
	"log"
	"net"
	"sync"
)

type Session struct {
	RoomNum    int32
	Users      []User
	IsRunning  bool
	UserLock   sync.Mutex
	SpawnIndex []int32
	PlayerNum  int32
	DieLock    sync.Mutex
}

type User struct {
	Conn            net.Conn
	Id              string
	RoomNum         int32
	NeedSync        bool
	CurrentLocation utils.Vec3
	RotationY       float32
	Session         *Session
	Health          int32

	SpawnIndex int32
	Dead       bool
}

func (s *Session) Init() {
	s.IsRunning = true
}

func (s *Session) UserEnter(usr User) {
	s.UserLock.Lock()
	s.Users = append(s.Users, usr)
	sm.Users[usr.Id] = &usr
	s.UserLock.Unlock()
	log.Println("UserEntered", s.Users, "/", s.PlayerNum)
	if len(s.Users) == int(s.PlayerNum) {
		log.Println("Match Complete Start Sync Move")
		s.StartSyncMove()
		ph.Room.Delete(s.RoomNum)
	} else if len(s.Users) > int(s.PlayerNum) {
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
	spawnPkt := pkt.R_PlayerSpawn{SpawnIndex: s.SpawnIndex, PlayerNum: s.PlayerNum}
	for _, u := range s.Users {
		spawnPkt.PlayerIds = append(spawnPkt.PlayerIds, u.Id)
	}
	for i, u := range s.Users {
		// u.SpawnIndex = int32(i)
		spawnPkt.PlayerIndex = int32(i)
		utils.SendPacket("PlayerSpawn", spawnPkt, u.Conn)
	}

	// PlayerMove, Rotation Sync
	//go func() {
	//	for {
	//		time.Sleep(time.Millisecond * 200)
	//
	//		if s.IsRunning {
	//			for _, u := range s.Users {
	//				if u.NeedSync {
	//					pkt := pkt.SR_PlayerRotation{PlayerIndex: u.Id, RotationY: u.RotationY}
	//					buffer := utils.MakeSendBuffer("PlayerRotation", pkt)
	//					s.BroadCast(buffer)
	//					u.NeedSync = false
	//				}
	//			}
	//		} else {
	//			log.Println("Game End", s.RoomNum)
	//			return
	//		}
	//	}
	//}()
}

func (s *Session) UserDie(index int32) {
	s.DieLock.Lock()

	pk := pkt.R_Die{PlayerIndex: index, Rank: int32(s.PlayerNum)}
	buffer := utils.MakeSendBuffer("Die", pk)
	s.BroadCast(buffer)
	log.Println("User", index, " Die")
	s.Users[index].Dead = true
	s.PlayerNum--

	if s.PlayerNum <= 1 {
		pk := pkt.R_GameEnd{}
		for _, u := range s.Users {
			if !u.Dead {
				utils.SendPacket("GameEnd", pk, u.Conn)
			}
		}
		log.Println(s.RoomNum, "ROOM GAME END!!")
	}
	s.DieLock.Unlock()
}

func (s *Session) BroadCast(buffer []byte) {
	for _, u := range s.Users {
		if u.Conn != nil {
			u.Conn.Write(buffer)
			//log.Println("BROADCASTED to ", u.Id)
		}
	}
}

func (s *Session) BroadCastExcpetMe(buffer []byte, index int32) {
	for _, u := range s.Users {
		if u.Conn != nil && u.SpawnIndex != index {
			u.Conn.Write(buffer)
			//log.Println("BROADCASTED to ", u.Id)
		}
	}
}

func (s *Session) ChangeHealth(index int32, value int32) int32 {
	s.UserLock.Lock()
	s.Users[index].Health += value
	s.UserLock.Unlock()
	return s.Users[index].Health
}
