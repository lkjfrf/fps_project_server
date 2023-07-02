package main

import (
	"FPSProject/network"
)

func main() {
	server := &network.Server{}

	server.RunTCP(":1998")
}
