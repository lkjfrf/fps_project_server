package pkt

import "FPSProject/utils"

// 방입장할때
type S_RoomEnter struct {
	PlayerId   string
	RoomNumber int32
}

// 방에서 게임시작버튼
type S_GameStartButton struct {
	RoomNumber int32
}

// 게임시작되고 맵 로딩 끝날때
type S_LoadingComplete struct {
	PlayerId   string
	RoomNumber int32
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

// 방에 들어갈 사람이 받는 인원들 패킷
type R_RoomEnter struct { // 방입장시 유저리스트 서버가 보냄
	PlayerId   []string
	RoomNumber int32
}

// 방에 원래있던 인원들이 받는 패킷
type R_RoomInUser struct {
	PlayerId string
}

// LOGIN
type S_Login struct {
	Id string
}

type R_GameStartButton struct {
}
