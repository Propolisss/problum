package service

import (
	"context"
	"errors"
	"testing"

	"problum/internal/model"
	"problum/internal/problem/service/mocks"
	templateDTO "problum/internal/template/service/dto"

	"go.uber.org/mock/gomock"
)

func TestServiceGetWithOptions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(*mocks.MockRepository, *mocks.MockTemplateService)
		opts    []Option
		wantErr bool
	}{
		{
			name: "with template and languages",
			setup: func(repo *mocks.MockRepository, templateSvc *mocks.MockTemplateService) {
				repo.EXPECT().Get(ctx, 1).Return(&model.Problem{ID: 1, LessonID: 2}, nil)
				templateSvc.EXPECT().GetByProblemIDAndLanguage(ctx, 1, "go").Return(&templateDTO.Template{ProblemID: 1, Language: "go"}, nil)
				templateSvc.EXPECT().GetLanguagesByProblemID(ctx, 1).Return([]string{"go"}, nil)
			},
			opts: []Option{WithLanguage("go"), WithTemplate(), WithLanguages()},
		},
		{
			name: "empty language",
			setup: func(repo *mocks.MockRepository, templateSvc *mocks.MockTemplateService) {
				repo.EXPECT().Get(ctx, 1).Return(&model.Problem{ID: 1, LessonID: 2}, nil)
			},
			opts:    []Option{WithLanguage(""), WithTemplate()},
			wantErr: true,
		},
		{
			name: "repo error",
			setup: func(repo *mocks.MockRepository, templateSvc *mocks.MockTemplateService) {
				repo.EXPECT().Get(ctx, 1).Return(nil, errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			attemptSvc := mocks.NewMockAttemptService(ctrl)
			templateSvc := mocks.NewMockTemplateService(ctrl)
			tt.setup(repo, templateSvc)

			got, err := New(repo, nil, attemptSvc, templateSvc).GetWithOptions(ctx, 1, tt.opts...)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != 1 || got.Template == nil || len(got.Languages) != 1 {
				t.Fatalf("unexpected problem: %+v", got)
			}
		})
	}
}
