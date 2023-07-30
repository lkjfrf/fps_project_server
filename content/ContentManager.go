package content

var (
	ph *PacketHandler
	sm *SessionManager
)

func ContentManagerInit() {
	ph = &PacketHandler{}
	sm = &SessionManager{}
	ph.Init()
	sm.Init()
}

func GetPacketHandler() *PacketHandler {
	return ph
}
