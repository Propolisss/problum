package service

import (
	"context"
	"testing"

	"problum/internal/model"
	"problum/internal/user/service/dto"
	"problum/internal/user/service/mocks"

	"go.uber.org/mock/gomock"
)

func TestServiceCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().Create(ctx, gomock.Any()).Return(&model.User{ID: 3, Login: "alice", HashedPassword: "hash"}, nil)

	got, err := New(repo).Create(ctx, &dto.User{Login: "alice", HashedPassword: "hash"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 3 || got.Login != "alice" {
		t.Fatalf("unexpected user: %+v", got)
	}
}
