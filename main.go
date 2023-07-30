package main

import (
	"FPSProject/content"
	"FPSProject/network"
)

func main() {
	server := &network.Server{}
	content.ContentManagerInit()
	server.RunTCP(":1998")
}
