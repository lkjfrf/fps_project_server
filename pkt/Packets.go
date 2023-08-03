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
