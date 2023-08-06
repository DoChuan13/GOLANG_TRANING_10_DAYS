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

// EndPoint only contains DB configuration host and DB name
type EndPoint struct {
	DockerPath string
	LocalPath  string
	Endpoint   string
	DBName     string
}

type Goroutine struct {
	Limited int
}
