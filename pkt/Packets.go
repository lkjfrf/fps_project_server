package pkt

import "FPSProject/utils"

// 방입장할때
type S_RoomEnter struct {
	PlayerId string
	RoomNum  int32
}

// 방에서 게임시작버튼
type S_GameStartButton struct {
	PlayerId string
	RoomNum  int32
}

// 게임시작되고 맵 로딩 끝날때
type S_LodingComplete struct {
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
type S_RoomCreate struct { // 방 생성시 서버 받음
	Title string
	Id    string
}

type R_RoomCreate struct { // 방 생성시 서버가 보냄
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

type R_RoomEnter struct { // 방입장시 유저리스트 서버가 보냄
	Id []string // 방에 있던사람은 방금들어온사람 1명, 방에 없던 사람은 방에 있는 사람들
}

// LOGIN
type S_Login struct {
	Id string
}
