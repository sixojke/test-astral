package config

import "time"

type Authorization struct {
	AdminToken string
	JWT        JWT
}

type JWT struct {
	SigningKey     string
	AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"`
}
