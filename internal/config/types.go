package config

type MySQLConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

type Config struct {
	Source MySQLConfig
	Target MySQLConfig
}
