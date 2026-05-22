package service

import (
	"context"
	"errors"
	"testing"

	"problum/internal/enrollment/service/mocks"
	"problum/internal/model"

	"go.uber.org/mock/gomock"
)

func TestServiceEnroll(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(*mocks.MockRepository)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Get(ctx, 1, 2).Return(nil, errors.New("not found"))
				repo.EXPECT().Enroll(ctx, 1, 2).Return(nil)
			},
		},
		{
			name: "already enrolled",
			setup: func(repo *mocks.MockRepository) {
				repo.EXPECT().Get(ctx, 1, 2).Return(&model.Enrollment{CourseID: 1, UserID: 2}, nil)
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

			err := New(repo).Enroll(ctx, 1, 2)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
