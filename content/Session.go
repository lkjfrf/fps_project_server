package content

import (
	"log"
	"net"
)

type Session struct {
	SessionId int32
	Users     []User
	IsRunning bool
}

type User struct {
	Conn      net.Conn
	Id        int32
	NickName  string
	SessionId int32
}

func (s *Session) Init() {
	s.Users = []User{}
	s.IsRunning = true
	go s.SyncMove()
}

func (s *Session) UserEnter(usr User) int32 {
	if len(s.Users) <= 7 {
		id := len(s.Users) + 1
		usr.Id = int32(id)
		s.Users = append(s.Users, usr)
		return usr.Id
	} else {
		return -1
	}
}

func (s *Session) SyncMove() {
	for {
		if s.IsRunning {

		} else {
			log.Println("Game End", s.SessionId)
			return
		}
	}
}
