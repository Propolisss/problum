package service

import (
	"context"
	"encoding/json"
	"testing"

	"problum/internal/model"
	"problum/internal/test/service/mocks"

	"go.uber.org/mock/gomock"
)

func TestServiceGetByProblemID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().GetByProblemID(ctx, 6).Return(&model.Test{
		ID:        1,
		ProblemID: 6,
		Tests:     json.RawMessage(`[{"input":"1","output":"1"}]`),
	}, nil)

	got, err := New(repo).GetByProblemID(ctx, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 1 || got.ProblemID != 6 || len(got.Tests) != 1 {
		t.Fatalf("unexpected test: %+v", got)
	}
}
