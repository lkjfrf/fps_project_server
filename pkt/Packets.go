package pkt

import "FPSProject/utils"

// 방입장할때
type R_RoomEnter struct {
	PlayerId string
	RoomNum  int32
}

// 방에서 게임시작버튼
type R_GameStartButton struct {
	PlayerId string
	RoomNum  int32
}

// 게임시작되고 맵 로딩 끝날때
type R_LodingComplete struct {
	PlayerId string
	RoomNum  int32
}

type SR_PlayerMove struct {
	PlayerId        string
	InputKey        int32
	IsPress         bool
	CurrentLocation utils.Vec3
}

type SR_PlayerRotation struct {
	PlayerId  string
	RotationY float32
}

type R_PlayerSpawn struct {
	PlayerIds   []string
	SpawnPoints []int32
}

// ROOM

type R_RoomCreate struct {
	BCreate    bool
	RoomNumber int32
}

type S_RequestRoomList struct {
}

type FRoomInfo struct {
	Title          string
	Id             string
	RoomNumber     int32
	NumberOfPeople int32
}

type R_RoomList struct {
	RoomList []FRoomInfo
}

type S_RoomCreate struct {
	Title string
	Id    string
}

// LOGIN
type S_Login struct {
	Id string
}
