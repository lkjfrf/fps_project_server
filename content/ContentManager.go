package content

var (
	ph *PacketHandler
)

func ContentManagerInit() {
	ph = &PacketHandler{}

	ph.Init()
}

func GetPacketHandler() *PacketHandler {
	return ph
}
