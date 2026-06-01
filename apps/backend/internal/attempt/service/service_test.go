package service

import (
	"context"
	"errors"
	"testing"

	"problum/internal/attempt/service/dto"
	"problum/internal/attempt/service/mocks"
	"problum/internal/model"

	"go.uber.org/mock/gomock"
)

func TestServiceSubmit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	errRepo := errors.New("repo error")

	tests := []struct {
		name    string
		setup   func(*mocks.MockRepository)
		want    int
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Submit(ctx, gomock.Any()).Return(11, nil)
			},
			want: 11,
		},
		{
			name: "repo error",
			setup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Submit(ctx, gomock.Any()).Return(0, errRepo)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			tt.setup(repo)

			got, err := New(repo).Submit(ctx, &dto.Attempt{UserID: 1, ProblemID: 2, Language: "go", Code: "package main"})
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().Get(ctx, 5).Return(&model.Attempt{ID: 5, UserID: 1, ProblemID: 2, Status: "accepted"}, nil)

	got, err := New(repo).Get(ctx, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 5 || got.Status != "accepted" {
		t.Fatalf("unexpected attempt: %+v", got)
	}
}
