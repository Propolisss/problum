package api

import "time"

type UserGetResponse struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}
