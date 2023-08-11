package pkt

import "FPSProject/utils"

// 7명 입장 순서
type R_EnterGmae struct {
	PlayerId int32
}

type SR_PlayerMove struct {
	PlayerId        int32
	InputKey        int32
	IsPress         bool
	CurrentLocation utils.Vec3
}

type SR_PlayerRotation struct {
	PlayerId  int32
	RotationY float32
}

type R_PlayerSpawn struct {
	PlayerIds   []int32
	SpawnPoints []int32
}

type R_RoomCreate struct {
	BCreate    bool
	RoomNumber int32
}

type S_RequestRoomList struct {
}

type RoomInfo struct {
	Title          string
	Id             string
	RoomNumber     int32
	NumberOfPeople int32
}

type R_RoomList struct {
	RoomList []RoomInfo
}

type S_RoomCreate struct {
	Title string
	Id    string
}

type S_Login struct {
	Id string
}
