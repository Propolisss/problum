package service

import (
	"context"
	"errors"
	"testing"

	"problum/internal/auth/service/mocks"
	userRepo "problum/internal/user/repository"
	userDTO "problum/internal/user/service/dto"

	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestServiceRegister(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(*mocks.MockUserService)
		login   string
		pass    string
		repeat  string
		wantErr bool
	}{
		{
			name:    "passwords mismatch",
			login:   "user",
			pass:    "pass",
			repeat:  "other",
			wantErr: true,
		},
		{
			name: "already exists",
			setup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().FindByLogin(ctx, "user").Return(&userDTO.User{ID: 1, Login: "user"}, nil)
			},
			login:   "user",
			pass:    "pass",
			repeat:  "pass",
			wantErr: true,
		},
		{
			name: "find error",
			setup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().FindByLogin(ctx, "user").Return(nil, errors.New("db error"))
			},
			login:   "user",
			pass:    "pass",
			repeat:  "pass",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			userSvc := mocks.NewMockUserService(ctrl)
			sessionSvc := mocks.NewMockSessionService(ctrl)
			if tt.setup != nil {
				tt.setup(userSvc)
			}

			_, err := New(nil, userSvc, sessionSvc).Register(ctx, tt.login, tt.pass, tt.repeat)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestServiceLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	hash, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	tests := []struct {
		name    string
		setup   func(*mocks.MockUserService)
		pass    string
		wantErr bool
	}{
		{
			name: "not found",
			setup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().FindByLogin(ctx, "user").Return(nil, userRepo.ErrNotFound)
			},
			pass:    "pass",
			wantErr: true,
		},
		{
			name: "invalid credentials",
			setup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().FindByLogin(ctx, "user").Return(&userDTO.User{ID: 1, Login: "user", HashedPassword: string(hash)}, nil)
			},
			pass:    "wrong",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			userSvc := mocks.NewMockUserService(ctrl)
			sessionSvc := mocks.NewMockSessionService(ctrl)
			tt.setup(userSvc)

			_, err := New(nil, userSvc, sessionSvc).Login(ctx, "user", tt.pass)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
