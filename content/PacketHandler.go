package content

import (
	"FPSProject/pkt"
	"FPSProject/utils"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type PacketHandler struct {
	TCPHandlerFunc map[string]func(net.Conn, string)
	UDPHandlerFunc map[string]func(*net.UDPAddr, string)

	//Login -> Store
	IdMap sync.Map

	// ROOM
	Room    sync.Map
	RoomNum atomic.Int32
}

type RoomInfo struct {
	Title          string
	Id             string
	RoomNumber     int32
	NumberOfPeople int32
	Ids            []string
}

func (ph *PacketHandler) Init() {

	log.Println("INIT_PacketHandler")

	ph.TCPHandlerFunc = make(map[string]func(net.Conn, string))
	ph.UDPHandlerFunc = make(map[string]func(*net.UDPAddr, string))

	/* ------------------------------------------------------------
						TCP Packet Handler
	------------------------------------------------------------ */

	ph.TCPHandlerFunc["Login"] = ph.Handle_Login

	ph.TCPHandlerFunc["PlayerMove"] = ph.Handle_PlayerMove
	ph.TCPHandlerFunc["PlayerRotation"] = ph.Handle_PlayerRotation

	// ROOM
	ph.TCPHandlerFunc["RoomEnter"] = ph.Handle_RoomEnter
	ph.TCPHandlerFunc["GameStartButton"] = ph.Handle_GameStartButton
	ph.TCPHandlerFunc["LoadingComplete"] = ph.Handle_LoadingComplete
	ph.TCPHandlerFunc["RoomCreate"] = ph.Handle_RoomCreate
	ph.TCPHandlerFunc["RequestRoomList"] = ph.Handle_RequestRoomList

	// Health
	ph.TCPHandlerFunc["ChangeHealth"] = ph.Handle_ChangeHealth

	/* ------------------------------------------------------------
							CONTENTS
	------------------------------------------------------------ */
	ph.Room = sync.Map{}

	test()
}

/* ------------------------------------------------------------
					TCP Handler Function
------------------------------------------------------------ */

func (ph *PacketHandler) Handle_Login(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_Login](json)
	ph.IdMap.Store(recvpkt.Id, c)
	log.Println("[LOGIN]", recvpkt.Id)
}

func (ph *PacketHandler) Handle_LoadingComplete(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_LoadingComplete](json)

	sm.TempNewSessionEnter(recvpkt.RoomNumber, recvpkt.PlayerId, c)
}
func (ph *PacketHandler) Handle_PlayerMove(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerMove](json)

	// sm.Users[recvpkt.PlayerId].NeedSync = true
	// sm.Users[recvpkt.PlayerId].CurrentLocation = recvpkt.CurrentLocation
	pk := pkt.SR_PlayerMove{PlayerIndex: recvpkt.PlayerIndex,
		InputKey:        recvpkt.InputKey,
		IsPress:         recvpkt.IsPress,
		CurrentLocation: recvpkt.CurrentLocation}
	buffer := utils.MakeSendBuffer("PlayerMove", pk)

	if s, ok := sm.Sessions[recvpkt.RoomNumber]; ok {
		s.BroadCastExcpetMe(buffer, recvpkt.PlayerIndex)
	}
	// sm.Users[recvpkt.PlayerIndex].Session.BroadCast(buffer)
}
func (ph *PacketHandler) Handle_PlayerRotation(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerRotation](json)

	// sm.Users[recvpkt.PlayerId].NeedSync = true
	// sm.Users[recvpkt.PlayerId].RotationY = recvpkt.RotationY

	pk := pkt.SR_PlayerRotation{PlayerIndex: recvpkt.PlayerIndex,
		RotationY: recvpkt.RotationY,
	}
	buffer := utils.MakeSendBuffer("PlayerRotation", pk)

	if s, ok := sm.Sessions[recvpkt.RoomNumber]; ok {
		s.BroadCastExcpetMe(buffer, recvpkt.PlayerIndex)
	}
	// sm.Users[recvpkt.PlayerIndex].Session.BroadCast(buffer)
}

func (ph *PacketHandler) Handle_RoomCreate(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_RoomCreate](json)

	roomNum := ph.RoomNum.Add(1)
	roomInfo := RoomInfo{Title: recvpkt.Title, Id: recvpkt.Id, RoomNumber: roomNum, NumberOfPeople: 0, Ids: []string{}}
	ph.Room.Store(roomNum, &roomInfo)

	pk := pkt.R_RoomCreate{BCreate: true, RoomNumber: roomNum}
	utils.SendPacket("RoomCreate", pk, c)
}

func (ph *PacketHandler) Handle_RequestRoomList(c net.Conn, json string) {
	pkt := pkt.R_RoomList{RoomList: ph.GetRoomList()}
	log.Println("SENDROOMLIST : ", pkt.RoomList)

	utils.SendPacket("RoomList", pkt, c)
}

func (ph *PacketHandler) Handle_RoomEnter(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_RoomEnter](json)
	if r, ok := ph.Room.Load(recvpkt.RoomNumber); ok {

		r.(*RoomInfo).Ids = append(r.(*RoomInfo).Ids, recvpkt.PlayerId)

		pk1 := pkt.R_RoomEnter{RoomNumber: recvpkt.RoomNumber}
		pk2 := pkt.R_RoomInUser{PlayerId: recvpkt.PlayerId}
		// 기존 인원들 list 주기
		for _, id := range r.(*RoomInfo).Ids {
			pk1.PlayerId = append(pk1.PlayerId, id)

			// 기존 인원들 한테 새로들어온 인원 알려주기
			if c2, ok := ph.IdMap.Load(id); ok && id != recvpkt.PlayerId {
				utils.SendPacket("RoomInUser", pk2, c2.(net.Conn))
			}
		}
		if r.(*RoomInfo).NumberOfPeople != 0 { // 방금 방만든사람은 아무것도 안 받도록
			utils.SendPacket("RoomEnter", pk1, c)
		}

		r.(*RoomInfo).NumberOfPeople++
	}
}

func (ph *PacketHandler) Handle_GameStartButton(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_GameStartButton](json)

	if r, ok := ph.Room.Load(recvpkt.RoomNumber); ok {
		pk := pkt.R_GameStartButton{}

		for _, id := range r.(*RoomInfo).Ids {
			if c2, ok := ph.IdMap.Load(id); ok {
				utils.SendPacket("GameStartButton", pk, c2.(net.Conn))
			}
		}
		ph.Room.Delete(recvpkt.RoomNumber)
	}
}

func (ph *PacketHandler) Handle_ChangeHealth(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_ChangeHealth](json)

	if s, ok := sm.Sessions[recvpkt.RoomNumber]; ok {
		currentHealth := s.ChangeHealth(recvpkt.PlayerIndex, recvpkt.Value)
		log.Println("User", recvpkt.PlayerIndex, " Hit CurrentHealth :", currentHealth)
		if currentHealth == 0 {
			pk := pkt.R_Die{PlayerIndex: recvpkt.PlayerIndex}
			buffer := utils.MakeSendBuffer("Die", pk)
			s.BroadCast(buffer)
			log.Println("User", recvpkt.PlayerIndex, " Die")
		} else {
			pk := pkt.R_ChangeHealth{PlayerIndex: recvpkt.PlayerIndex, CurrentHealth: currentHealth}
			buffer := utils.MakeSendBuffer("ChangeHealth", pk)
			s.BroadCast(buffer)
		}
	}
}

/* ------------------------------------------------------------
						CONTENTS
------------------------------------------------------------ */

func (ph *PacketHandler) GetRoomList() []pkt.FRoomInfo {
	roomList := []pkt.FRoomInfo{}
	ph.Room.Range(func(key, value any) bool {
		r := pkt.FRoomInfo{
			Title:          value.(*RoomInfo).Title,
			Id:             value.(*RoomInfo).Id,
			RoomNumber:     value.(*RoomInfo).RoomNumber,
			NumberOfPeople: value.(*RoomInfo).NumberOfPeople}
		roomList = append(roomList, r)
		return true
	})
	return roomList
}

func test() {
	go func() {
		time.Sleep(time.Second * 2)
		//log.Print(utils.RandomInt32(7, 0, 6))
	}()
}
