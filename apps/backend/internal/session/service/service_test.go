package service

import (
	"context"
	"testing"

	"problum/internal/model"
	"problum/internal/session/service/mocks"

	"go.uber.org/mock/gomock"
)

func TestServiceCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	session := &model.UserSession{UserID: 9, RefreshHash: "hash"}
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().Create(ctx, session).Return(&model.UserSession{ID: 1, UserID: 9, RefreshHash: "hash"}, nil)

	got, err := New(repo).Create(ctx, session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 1 || got.UserID != 9 {
		t.Fatalf("unexpected session: %+v", got)
	}
}
