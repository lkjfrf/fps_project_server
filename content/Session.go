package content

import "net"

type Session struct {
	SessionId int
	Users     []User
}

type User struct {
	Conn net.Conn
	Id   string
}

func (s *Session) Init() {

	go s.SyncMove()
}

func (s *Session) UserEnter(usr *User) {

}

func (s *Session) SyncMove() {

}
