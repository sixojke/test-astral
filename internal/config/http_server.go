package config

import "time"

type HTTPServer struct {
	Port               string        `mapstructure:"port"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
	MaxHeaderMegabytes int           `mapstructure:"max_header_megabytes"`
}
