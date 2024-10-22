package domain

import "time"

type Session struct {
	UserId      string
	AccessToken string
	ExpiresAt   time.Time
}
