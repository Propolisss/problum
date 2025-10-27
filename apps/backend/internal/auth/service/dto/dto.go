package dto

import "time"

type LoginDTO struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Duration
}

type RegisterDTO struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Duration
}

type RefreshDTO struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Duration
}
