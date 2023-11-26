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
	RoomNumber      int32
	PlayerIndex     int32
	InputKey        int32
	IsPress         bool
	CurrentLocation utils.Vec3
}

type SR_PlayerRotation struct {
	RoomNumber  int32
	PlayerIndex int32
	RotationY   float32
}

type R_PlayerSpawn struct {
	PlayerIndex int32    // 자기 자신의 스폰포인트값
	SpawnIndex  []int32  // 모든사람의 스폰포인트 값 (무작위)
	PlayerIds   []string // 모든 사람 ID 값
	PlayerNum   int32    // 인원 수
	// 들어온 순 : 우현 영민 민석
	//2 0 1
	//우현 영민 민석
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

//Health
type S_ChangeHealth struct { // 내가 쏜 총알이 otherplayer에 맞으면 보내는 패킷
	PlayerIndex int32 // 총알 맞은 대상
	RoomNumber  int32
	Value       int32 // 총맞으면 -10 , 물약먹으면 10  이런식으로
}

type R_ChangeHealth struct { // 맵에서 누군가 피 닳으면 맵에 있는 모든사람이 받는 패킷 (자신포함)
	PlayerIndex   int32
	CurrentHealth int32
}

// 맵에서 누군가 죽으면 받는 패킷 (자신포함)
// 죽으면 바로 세션끊김
type R_Die struct {
	PlayerIndex int32
	Rank        int32
}

// 최후의 1인 되면 그 1인한테 주는 패킷
// 세션끊고 방 지워버림
type R_GameEnd struct {
}
