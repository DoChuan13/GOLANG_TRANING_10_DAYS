package model

type KafkaDB struct {
	Port              int    `json:"TK_DB_PORT"`
	User              string `json:"TK_DB_USER"`
	Password          string `json:"TK_DB_PASSWORD"`
	MaxOpenConnection int    `json:"TK_DB_MAX_OPEN_CONNECTION"`
	MaxIdleConnection int    `json:"TK_DB_MAX_IDLE_CONNECTION"`
	DriverName        string `json:"TK_DB_DRIVER_NAME"`
	RetryTimes        int    `json:"TK_DB_RETRY_TIMES"`
	RetryWaitMs       int    `json:"TK_DB_RETRY_WAIT_MS"`
}

// EndPoint contains DB configuration host and DB name
type EndPoint struct {
	TKQKBN            string `json:"TKQKBN"`
	SNDC              string `json:"SNDC"`
	DBMasterEndpoint  string `json:"TKDB_MASTER_ENDPOINT"`
	DBReplicaEndpoint string `json:"TKDB_REPLICA_ENDPOINT"`
	DBName            string `json:"TKDBNAME"`
}

// EndPointCustom only contains DB configuration host and DB name
type EndPointCustom struct {
	Endpoint string
	DBName   string
}
