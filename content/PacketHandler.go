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
	ph.TCPHandlerFunc["RoomCreate"] = ph.Handle_RoomCreate
	ph.TCPHandlerFunc["RequestRoomList"] = ph.Handle_RequestRoomList

	/* ------------------------------------------------------------
							CONTENTS
	------------------------------------------------------------ */
	ph.Room = sync.Map{}

	//test()
}

/* ------------------------------------------------------------
					TCP Handler Function
------------------------------------------------------------ */

func (ph *PacketHandler) Handle_Login(c net.Conn, json string) {
	// recvpkt := utils.JsonStrToStruct[pkt.S_Login](json)

}

func (ph *PacketHandler) Handle_R_LodingComplete(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.R_LodingComplete](json)

	sm.TempNewSessionEnter(recvpkt.RoomNum, recvpkt.PlayerId, c)
}
func (ph *PacketHandler) Handle_PlayerMove(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerMove](json)

	// sm.Users[recvpkt.PlayerId].NeedSync = true
	// sm.Users[recvpkt.PlayerId].CurrentLocation = recvpkt.CurrentLocation
	pkt := pkt.SR_PlayerMove{PlayerId: recvpkt.PlayerId,
		InputKey:        recvpkt.InputKey,
		IsPress:         recvpkt.IsPress,
		CurrentLocation: recvpkt.CurrentLocation}
	buffer := utils.MakeSendBuffer("PlayerMove", pkt)

	sm.Users[recvpkt.PlayerId].Session.BroadCast(buffer)
}
func (ph *PacketHandler) Handle_PlayerRotation(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerRotation](json)

	// sm.Users[recvpkt.PlayerId].NeedSync = true
	// sm.Users[recvpkt.PlayerId].RotationY = recvpkt.RotationY

	pkt := pkt.SR_PlayerRotation{PlayerId: recvpkt.PlayerId,
		RotationY: recvpkt.RotationY,
	}
	buffer := utils.MakeSendBuffer("PlayerRotation", pkt)

	sm.Users[recvpkt.PlayerId].Session.BroadCast(buffer)
}

func (ph *PacketHandler) Handle_RoomCreate(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.S_RoomCreate](json)

	roomNum := ph.RoomNum.Add(1)
	ph.Room.Store(roomNum, pkt.FRoomInfo{Id: recvpkt.Id, Title: recvpkt.Title, RoomNumber: roomNum, NumberOfPeople: 0})
}

func (ph *PacketHandler) Handle_RequestRoomList(c net.Conn, json string) {
	pkt := pkt.R_RoomList{RoomList: ph.GetRoomList()}
	log.Println("SENDROOMLIST : ", pkt.RoomList)
	buffer := utils.MakeSendBuffer("RoomList", pkt)
	c.Write(buffer)
}

/* ------------------------------------------------------------
						CONTENTS
------------------------------------------------------------ */

func (ph *PacketHandler) GetRoomList() []pkt.FRoomInfo {
	roomList := []pkt.FRoomInfo{}
	ph.Room.Range(func(key, value any) bool {
		roomList = append(roomList, value.(pkt.FRoomInfo))
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
