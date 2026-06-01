package service

import (
	"context"
	"testing"

	"problum/internal/model"
	"problum/internal/template/service/mocks"

	"go.uber.org/mock/gomock"
)

func TestServiceGetLanguagesByProblemID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().GetLanguagesByProblemID(ctx, 4).Return([]string{"go", "python"}, nil)

	got, err := New(repo).GetLanguagesByProblemID(ctx, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "go" {
		t.Fatalf("unexpected languages: %+v", got)
	}
}

func TestServiceGetByProblemIDAndLanguage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().GetByProblemIDAndLanguage(ctx, 4, "go").Return(&model.Template{ID: 8, ProblemID: 4, Language: "go"}, nil)

	got, err := New(repo).GetByProblemIDAndLanguage(ctx, 4, "go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 8 || got.Language != "go" {
		t.Fatalf("unexpected template: %+v", got)
	}
}
