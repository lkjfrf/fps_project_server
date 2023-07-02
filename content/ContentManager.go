package content

var (
	ph *PacketHandler
	//gs  *GlobalSession
	// dbm *DBManager
)

func ContentManagerInit() {
	ph = &PacketHandler{}
	// gs = &GlobalSession{}
	// dbm = &DBManager{}

	ph.Init()
	// gs.Init()
	// dbm.Init()
}

func GetPacketHandler() *PacketHandler {
	return ph
}

// func GetGlobalSession() *GlobalSession {
// 	return gs
// }

// func GetDBManager() *DBManager {
// 	return dbm
// }
