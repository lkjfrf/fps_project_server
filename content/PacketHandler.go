package content

import (
	"FPSProject/pkt"
	"FPSProject/utils"
	"log"
	"net"
	"sync"
	"sync/atomic"
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

	ph.TCPHandlerFunc["EnterGame"] = ph.Handle_EnterGame
	ph.TCPHandlerFunc["PlayerMove"] = ph.Handle_PlayerMove

	ph.UDPHandlerFunc["PlayerRotation"] = ph.Handle_PlayerRotation

	// ROOM
	ph.TCPHandlerFunc["RoomCreate"] = ph.Handle_RoomCreate
	ph.TCPHandlerFunc["RequestRoomList"] = ph.Handle_RequestRoomList

	/* ------------------------------------------------------------
							CONTENTS
	------------------------------------------------------------ */
	ph.Room = sync.Map{}
}

/* ------------------------------------------------------------
					TCP Handler Function
------------------------------------------------------------ */

func (ph *PacketHandler) Handle_Login(c net.Conn, json string) {
	// recvpkt := utils.JsonStrToStruct[pkt.S_Login](json)

}

func (ph *PacketHandler) Handle_EnterGame(c net.Conn, json string) {
	pkt := pkt.R_EnterGmae{
		PlayerId: sm.NewPlayer(c),
	}
	buffer := utils.MakeSendBuffer("EnterGame", pkt)
	c.Write(buffer)
	log.Println("ENTER SEND", string(buffer))
}
func (ph *PacketHandler) Handle_PlayerMove(c net.Conn, json string) {
	recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerMove](json)

	pkt := pkt.SR_PlayerMove{
		PlayerId:        recvpkt.PlayerId,
		InputKey:        recvpkt.InputKey,
		IsPress:         recvpkt.IsPress,
		CurrentLocation: recvpkt.CurrentLocation,
	}

	buffer := utils.MakeSendBuffer("PlayerMove", pkt)
	c.Write(buffer)
}
func (ph *PacketHandler) Handle_PlayerRotation(c *net.UDPAddr, json string) {
	// recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerRotation](json)
	// GetGlobalSession().GStarSelect(c, recvpkt)
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
