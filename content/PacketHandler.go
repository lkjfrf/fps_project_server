package content

import (
	"FPSProject/pkt"
	"FPSProject/utils"
	"log"
	"net"
)

type PacketHandler struct {
	TCPHandlerFunc map[string]func(net.Conn, string)
	UDPHandlerFunc map[string]func(*net.UDPAddr, string)
}

func (ph *PacketHandler) Init() {

	log.Println("INIT_PacketHandler")

	ph.TCPHandlerFunc = make(map[string]func(net.Conn, string))
	ph.UDPHandlerFunc = make(map[string]func(*net.UDPAddr, string))

	/* ------------------------------------------------------------
						TCP Packet Handler
	------------------------------------------------------------ */

	ph.TCPHandlerFunc["EnterGame"] = ph.Handle_EnterGame
	ph.TCPHandlerFunc["PlayerMove"] = ph.Handle_PlayerMove
	ph.TCPHandlerFunc["PlayerRotation"] = ph.Handle_PlayerRotation
	//ph.UDPHandlerFunc["PlayerRotation"] = ph.Handle_PlayerRotation

}

/* ------------------------------------------------------------
					TCP Handler Function
------------------------------------------------------------ */

func (ph *PacketHandler) Handle_EnterGame(c net.Conn, json string) {
	pkt := pkt.R_EnterGmae{
		PlayerId: sm.NewPlayer(c),
	}
	buffer := utils.MakeSendBuffer("EnterGame", pkt)
	c.Write(buffer)
	log.Println("ENTER SEND", string(buffer))
}
func (ph *PacketHandler) Handle_PlayerMove(c net.Conn, json string) {
	//recvpkt := utils.JsonStrToStruct[pkt.SR_PlayerMove](json)

	// pkt := pkt.SR_PlayerMove{
	// 	PlayerId:        recvpkt.PlayerId,
	// 	InputKey:        recvpkt.InputKey,
	// 	IsPress:         recvpkt.IsPress,
	// 	CurrentLocation: recvpkt.CurrentLocation,
	// }

	//buffer := utils.MakeSendBuffer("PlayerMove", pkt)
	//c.Write(buffer)
}
func (ph *PacketHandler) Handle_PlayerRotation(c net.Conn, json string) {

	//recvpkt := utils.JsonStrToStruct[pkt.C_GStarSelect](json)
	//GetGlobalSession().GStarSelect(c, recvpkt)
}
