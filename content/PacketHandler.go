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
	Room      sync.Map
	RoomNum   atomic.Int32
	RoomUsers sync.Map
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

	ph.TCPHandlerFunc["R_LodingComplete"] = ph.Handle_R_LodingComplete

	ph.TCPHandlerFunc["PlayerMove"] = ph.Handle_PlayerMove
	ph.TCPHandlerFunc["PlayerRotation"] = ph.Handle_PlayerRotation

	// ROOM
	ph.TCPHandlerFunc["RoomEnter"] = ph.Handle_RoomEnter
	// ph.TCPHandlerFunc["GameStartButton"] = ph.Handle_GameStartButton
	// ph.TCPHandlerFunc["LodingComplete"] = ph.Handle_LodingComplete
	ph.TCPHandlerFunc["RoomCreate"] = ph.Handle_RoomCreate
	ph.TCPHandlerFunc["RequestRoomList"] = ph.Handle_RequestRoomList

	/* ------------------------------------------------------------
							CONTENTS
	------------------------------------------------------------ */
	ph.Room = sync.Map{}
	ph.RoomUsers = sync.Map{}

	//test()
}

/* ------------------------------------------------------------
					TCP Handler Function
------------------------------------------------------------ */

func (ph *PacketHandler) Handle_Login(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_Login](json)
	ph.IdMap.Store(recvpkt.Id, c)
	log.Println("[LOGIN]", recvpkt.Id)
}

func (ph *PacketHandler) Handle_R_LodingComplete(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_LodingComplete](json)

	sm.TempNewSessionEnter(recvpkt.RoomNum, recvpkt.PlayerId, c)
}
func (ph *PacketHandler) Handle_PlayerMove(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerMove](json)

	// sm.Users[recvpkt.PlayerId].NeedSync = true
	// sm.Users[recvpkt.PlayerId].CurrentLocation = recvpkt.CurrentLocation
	pk := pkt.SR_PlayerMove{PlayerId: recvpkt.PlayerId,
		InputKey:        recvpkt.InputKey,
		IsPress:         recvpkt.IsPress,
		CurrentLocation: recvpkt.CurrentLocation}
	buffer := utils.MakeSendBuffer("PlayerMove", pk)

	sm.Users[recvpkt.PlayerId].Session.BroadCast(buffer)
}
func (ph *PacketHandler) Handle_PlayerRotation(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerRotation](json)

	// sm.Users[recvpkt.PlayerId].NeedSync = true
	// sm.Users[recvpkt.PlayerId].RotationY = recvpkt.RotationY

	pk := pkt.SR_PlayerRotation{PlayerId: recvpkt.PlayerId,
		RotationY: recvpkt.RotationY,
	}
	buffer := utils.MakeSendBuffer("PlayerRotation", pk)

	sm.Users[recvpkt.PlayerId].Session.BroadCast(buffer)
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
	if r, ok := ph.Room.Load(recvpkt.RoomNum); ok {
		r.(*RoomInfo).NumberOfPeople++
		r.(*RoomInfo).Ids = append(r.(*RoomInfo).Ids, recvpkt.PlayerId)

		pk1 := pkt.R_RoomEnter{}
		pk2 := pkt.R_RoomEnter{Id: []string{recvpkt.PlayerId}}
		// 기존 인원들 list 주기
		for _, id := range r.(*RoomInfo).Ids {
			pk1.Id = append(pk1.Id, id)

			// 기존 인원들 한테 새로들어온 인원 알려주기
			if c2, ok := ph.IdMap.Load(id); ok {
				utils.SendPacket("RoomEnter", pk2, c2.(net.Conn))
			}
		}
		utils.SendPacket("RoomEnter", pk1, c)
	}
}

// func (ph *PacketHandler) Handle_GameStartButton(c net.Conn, json string) {
// 	recvpkt := utils.JsonStrToStruct[pkt.S_GameStartButton](json)

// }
// func (ph *PacketHandler) Handle_LodingComplete(c net.Conn, json string) {
// 	recvpkt := utils.JsonStrToStruct[pkt.S_LodingComplete](json)

// }

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
		sm.TempNewSessionEnter(1, "hi", nil)
		sm.TempNewSessionEnter(1, "wow", nil)
		sm.TempNewSessionEnter(1, "wow123", nil)

		a := "etstestset"
		sm.Users["hi"].Session.BroadCast([]byte(a))
	}()
}
