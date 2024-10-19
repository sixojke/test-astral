package config

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string
	Password string
	DBName   string
	SSLMode  string `mapstructure:"ssl_mode"`
}
