package service

import (
	"context"
	"errors"
	"testing"

	"problum/internal/course/service/mocks"
	enrollmentDTO "problum/internal/enrollment/service/dto"
	lessonDTO "problum/internal/lesson/service/dto"
	"problum/internal/model"

	"go.uber.org/mock/gomock"
)

func TestServiceList(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "user_id", 7)
	errRepo := errors.New("repo error")

	tests := []struct {
		name    string
		setup   func(*mocks.MockRepository, *mocks.MockEnrollmentService)
		wantErr bool
		want    int
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockRepository, enrollmentSvc *mocks.MockEnrollmentService) {
				repo.EXPECT().List(ctx).Return([]*model.Course{{ID: 1}, {ID: 2}}, nil)
				enrollmentSvc.EXPECT().GetListByUserID(ctx, 7).Return([]*enrollmentDTO.Enrollment{{CourseID: 2, UserID: 7}}, nil)
			},
			want: 2,
		},
		{
			name: "repo error",
			setup: func(repo *mocks.MockRepository, enrollmentSvc *mocks.MockEnrollmentService) {
				repo.EXPECT().List(ctx).Return(nil, errRepo)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			lessonSvc := mocks.NewMockLessonService(ctrl)
			enrollmentSvc := mocks.NewMockEnrollmentService(ctrl)
			tt.setup(repo, enrollmentSvc)

			got, err := New(repo, lessonSvc, enrollmentSvc).List(ctx)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.want {
				t.Fatalf("got %d courses, want %d", len(got), tt.want)
			}
			if !got[1].Enrolled {
				t.Fatal("expected second course to be enrolled")
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "user_id", 7)

	tests := []struct {
		name    string
		setup   func(*mocks.MockRepository, *mocks.MockLessonService, *mocks.MockEnrollmentService)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockRepository, lessonSvc *mocks.MockLessonService, enrollmentSvc *mocks.MockEnrollmentService) {
				repo.EXPECT().Get(ctx, 1).Return(&model.Course{ID: 1}, nil)
				lessonSvc.EXPECT().ListByCourseID(ctx, 1).Return([]*lessonDTO.Lesson{{ID: 3, CourseID: 1}}, nil)
				enrollmentSvc.EXPECT().Get(ctx, 1, 7).Return(&enrollmentDTO.Enrollment{CourseID: 1, UserID: 7}, nil)
			},
		},
		{
			name: "lesson error",
			setup: func(repo *mocks.MockRepository, lessonSvc *mocks.MockLessonService, enrollmentSvc *mocks.MockEnrollmentService) {
				repo.EXPECT().Get(ctx, 1).Return(&model.Course{ID: 1}, nil)
				lessonSvc.EXPECT().ListByCourseID(ctx, 1).Return(nil, errors.New("lesson error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			lessonSvc := mocks.NewMockLessonService(ctrl)
			enrollmentSvc := mocks.NewMockEnrollmentService(ctrl)
			tt.setup(repo, lessonSvc, enrollmentSvc)

			got, err := New(repo, lessonSvc, enrollmentSvc).Get(ctx, 1)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != 1 || len(got.Lessons) != 1 || !got.Enrolled {
				t.Fatalf("unexpected course: %+v", got)
			}
		})
	}
}
