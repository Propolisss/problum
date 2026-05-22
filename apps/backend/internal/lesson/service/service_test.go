package service

import (
	"context"
	"errors"
	"testing"

	"problum/internal/lesson/service/mocks"
	"problum/internal/model"
	problemDTO "problum/internal/problem/service/dto"

	"go.uber.org/mock/gomock"
)

func TestServiceGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(*mocks.MockRepository, *mocks.MockProblemService)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockRepository, problemSvc *mocks.MockProblemService) {
				repo.EXPECT().Get(ctx, 3).Return(&model.Lesson{ID: 3, CourseID: 1}, nil)
				problemSvc.EXPECT().ListByLessonID(ctx, 3).Return([]*problemDTO.Problem{{ID: 9, LessonID: 3}}, nil)
			},
		},
		{
			name: "repo error",
			setup: func(repo *mocks.MockRepository, problemSvc *mocks.MockProblemService) {
				repo.EXPECT().Get(ctx, 3).Return(nil, errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			problemSvc := mocks.NewMockProblemService(ctrl)
			tt.setup(repo, problemSvc)

			got, err := New(repo, problemSvc).Get(ctx, 3)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != 3 || len(got.Problems) != 1 {
				t.Fatalf("unexpected lesson: %+v", got)
			}
		})
	}
}
