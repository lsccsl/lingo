package dbaux


type DBBase interface {
	DBGetServerClusterEntryPoint(clientID int64)(srvID int64)
	DBGetClientData(clientID int64)*DBClientData
}

