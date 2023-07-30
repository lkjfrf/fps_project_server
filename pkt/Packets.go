package pkt

import "FPSProject/utils"

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
