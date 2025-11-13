package dto

import (
	"time"

	"problum/internal/api"
	"problum/internal/model"
)

type User struct {
	ID             int
	Login          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func ToModel(user *User) *model.User {
	return &model.User{
		ID:             user.ID,
		Login:          user.Login,
		HashedPassword: user.HashedPassword,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func ToDTO(user *model.User) *User {
	return &User{
		ID:             user.ID,
		Login:          user.Login,
		HashedPassword: user.HashedPassword,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func ToAPI(user *User) api.UserGetResponse {
	return api.UserGetResponse{
		ID:        user.ID,
		Login:     user.Login,
		CreatedAt: user.CreatedAt,
	}
}
