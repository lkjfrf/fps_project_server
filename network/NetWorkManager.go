package network

import (
	"FPSProject/pkt"
	"FPSProject/utils"
	"errors"
	"io"
	"log"
	"net"
	"syscall"
)

type Server struct {
}

func (s *Server) RunTCP(addr string) {

	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server Running at : " + listener.Addr().String())

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("New Connection : " + conn.RemoteAddr().String())
		//s.testSend(conn)

		go s.EventLoop(conn)
	}
}

func (s *Server) testSend(conn net.Conn) {
	recvpkt := pkt.S_Test{Msg: "testrrr"}
	sendbuffer := utils.MakeSendBuffer("TestPacket", recvpkt)
	conn.Write(sendbuffer)
}

func (s *Server) EventLoop(conn net.Conn) {

	header := make([]byte, 4)

	for {
		n, err := conn.Read(header)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, syscall.Errno(10054)) {
				//content.GetGlobalSession().Disconnect(conn)
				conn.Close()
			}
			return
		}

		if n == 4 {
			namesize, datasize := utils.ParseHeader(header)
			recvdata := make([]byte, namesize+datasize)
			recvn, err := conn.Read(recvdata)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, syscall.Errno(10054)) {
					log.Println("Connection is closed from client : " + conn.RemoteAddr().String())
				}
				return
			}

			if recvn == namesize+datasize {
				pktname, jsondata := utils.ExtractData(namesize, datasize, recvdata)
				log.Println(pktname, jsondata)
				s.testSend(conn)
				// if _, ok := content.GetPacketHandler().TCPHandlerFunc[pktname]; ok {
				// 	log.Println(pktname, jsondata)
				// 	content.GetPacketHandler().TCPHandlerFunc[pktname](conn, jsondata)
				// } else {
				// 	log.Println("Unknown Packet : " + pktname)
				// }

			} else {
				log.Println("Recv Wrong Datasize !")
			}
		} else {
			log.Println("Recv Wrong Header !")
		}
	}
}
