package utils

import (
	"encoding/binary"
	"encoding/json"
	"log"
)

type Vec3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func ParseHeader(header []byte) (int, int) {
	namesize := binary.LittleEndian.Uint16(header[:2])
	datasize := binary.LittleEndian.Uint16(header[2:4])

	return int(namesize), int(datasize)
}

func ExtractData(namesize int, datasize int, recvdata []byte) (string, string) {
	totalsize := namesize + datasize

	namedata := recvdata[:namesize]
	jsondata := recvdata[namesize:totalsize]

	return string(namedata), string(jsondata)
}

func MakeSendBuffer[T any](pktname string, data T) []byte {
	sendData, err := json.Marshal(&data)
	if err != nil {
		log.Println("Marshal Error !", err)
	}
	packetname := []byte(pktname)

	namesize := len(packetname)
	datasize := len(sendData)

	sendBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint16(sendBuffer, uint16(namesize))
	binary.LittleEndian.PutUint16(sendBuffer[2:], uint16(datasize))

	sendBuffer = append(sendBuffer, packetname...)
	sendBuffer = append(sendBuffer, sendData...)

	return sendBuffer
}

func JsonStrToStruct[T any](jsonstr string) T {
	var data T
	err := json.Unmarshal([]byte(jsonstr), &data)
	if err != nil {
		log.Println("Unmarshal Error !", err)
	}
	return data
}
