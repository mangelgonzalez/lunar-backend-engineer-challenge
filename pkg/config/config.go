package config

type Config struct {
	ApplicationName string `env:"APPLICATION_NAME"`
	AppEnv          string `env:"APP_ENV"`

	HttpServerConfig
	MySqlConfig
	RedisConfig
	DatabaseMigrationConfig
}

type HttpServerConfig struct {
	ServerHost         string `env:"SERVER_HOST"`
	ServerPort         string `env:"SERVER_PORT"`
	ServerWriteTimeout int    `env:"SERVER_WRITE_TIMEOUT"`
	ServerReadTimeout  int    `env:"SERVER_READ_TIMEOUT"`
}

type RedisConfig struct {
	RedisHost               string `env:"REDIS_HOST"`
	RedisPort               int    `env:"REDIS_PORT"`
	RedisMaxIdleConnections int    `env:"REDIS_MAX_IDLE_CONNECTIONS"`
	RedisIdleTimeout        int    `env:"REDIS_IDLE_TIMEOUT"`
}

type MySqlConfig struct {
	MySqlWriterHost string `env:"MYSQL_HOST"`
	MySqlReaderHost string `env:"MYSQL_READER_HOST"`
	MySqlUser       string `env:"MYSQL_USER"`
	MySqlPass       string `env:"MYSQL_PASS"`
	MySqlPort       uint16 `env:"MYSQL_PORT"`
	MySqlDB         string `env:"MYSQL_DATABASE"`
}

type DatabaseMigrationConfig struct {
	MigrationsLocation  string `env:"DATABASE_MIGRATIONS_LOCATION"`
	MigrationsTableName string `env:"DATABASE_MIGRATIONS_TABLE_NAME"`
}
