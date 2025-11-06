package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"problum/internal/auth/service/dto"
	"problum/internal/model"
	"problum/internal/redis"
	"problum/internal/utils"

	sessionRepo "problum/internal/session/repository"
	userRepo "problum/internal/user/repository"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var HMACRefreshTokenKey = []byte("refresh_token_key")

type UserService interface {
	FindByLogin(context.Context, string) (*model.User, error)
	Create(context.Context, *model.User) (*model.User, error)
}

type SessionService interface {
	Create(context.Context, *model.UserSession) (*model.UserSession, error)
	GetByRefreshHash(context.Context, string) (*model.UserSession, error)
	GetByPreviousRefreshHash(context.Context, string) (*model.UserSession, error)
	Update(context.Context, *model.UserSession) (*model.UserSession, error)
	LogoutAll(context.Context, int) error
}

type Service struct {
	rdb        *redis.Redis
	userSvc    UserService
	sessionSvc SessionService
}

func New(rdb *redis.Redis, userSvc UserService, sessionSvc SessionService) *Service {
	return &Service{
		rdb:        rdb,
		userSvc:    userSvc,
		sessionSvc: sessionSvc,
	}
}

func (s *Service) Register(ctx context.Context, login, password, repeatedPassword string) (*dto.RegisterDTO, error) {
	if password != repeatedPassword {
		log.Error().Msg("Passwords mismatch")
		return nil, fmt.Errorf("passwords mismatch")
	}

	if _, err := s.userSvc.FindByLogin(ctx, login); !errors.Is(err, userRepo.ErrNotFound) {
		log.Error().Err(err).Msg("user already exists")
		return nil, fmt.Errorf("user already exists")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password")
	}

	user := &model.User{
		Login:          login,
		HashedPassword: string(hashedPass),
	}

	if _, err := s.userSvc.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	loginResp, err := s.Login(ctx, login, password)

	return &dto.RegisterDTO{
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		ExpiresAt:    loginResp.ExpiresAt,
	}, err
}

func (s *Service) Login(ctx context.Context, login, password string) (*dto.LoginDTO, error) {
	user, err := s.userSvc.FindByLogin(ctx, login)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		return nil, fmt.Errorf("user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		log.Error().Err(err).Msg("Invalid credentials")
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken := utils.GenerateToken(32)
	refreshToken := utils.GenerateToken(32)

	hash := utils.GenerateHMAC(HMACRefreshTokenKey, []byte(refreshToken))

	us, err := s.sessionSvc.Create(ctx, &model.UserSession{
		UserID:      user.ID,
		RefreshHash: string(hash),
		ExpiresAt:   time.Now().AddDate(0, 0, 14),
		Revoked:     false,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	sessionJSON, err := sonic.Marshal(us)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal session")
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	if err = s.rdb.Set(ctx, fmt.Sprintf("user_sessions:%s", accessToken), sessionJSON, 15*time.Minute); err != nil {
		log.Error().Err(err).Msg("Failed to create session in Redis")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err = s.rdb.SAdd(ctx, fmt.Sprintf("user_access_tokens_%d", user.ID), accessToken); err != nil {
		log.Error().Err(err).Msg("Failed to add member in Redis")
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	return &dto.LoginDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    15 * time.Minute,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refresh string) (*dto.RefreshDTO, error) {
	hash := utils.GenerateHMAC(HMACRefreshTokenKey, []byte(refresh))

	session, err := s.sessionSvc.GetByRefreshHash(ctx, hash)

	if errors.Is(err, sessionRepo.ErrNotFound) {
		if us, _ := s.sessionSvc.GetByPreviousRefreshHash(ctx, hash); us != nil {
			log.Warn().Int("user_id", us.UserID).Msg("Data compromise detected, logout all")

			if e := s.sessionSvc.LogoutAll(ctx, us.UserID); e != nil {
				log.Error().Err(e).Int("user_id", us.UserID).Msg("Failed to logout all")
			}

			accessTokens, e := s.rdb.SMembers(ctx, fmt.Sprintf("user_access_tokens_%d", us.UserID))
			if e != nil {
				log.Error().Err(e).Int("user_id", us.UserID).Msg("Failed to get members")
			} else {
				for _, access := range accessTokens {
					if e = s.rdb.Delete(ctx, fmt.Sprintf("user_sessions:%s", access)); e != nil {
						log.Error().Err(e).Int("user_id", us.UserID).Msg("Failed to delete access token")
					}
				}
			}
			if e = s.rdb.Delete(ctx, fmt.Sprintf("user_access_tokens_%d", us.UserID)); e != nil {
				log.Error().Err(e).Int("user_id", us.UserID).Msg("Failed to delete members")
			}

			return nil, fmt.Errorf("data compromise: %w", err)
		}

		return nil, err
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to get user session by refresh hash")
		return nil, fmt.Errorf("failed to get user session by refresh hash")
	}

	if !(session.Revoked == false && time.Now().Before(session.ExpiresAt) &&
		time.Now().Before(session.LastActivityAt.Add(7*24*time.Hour))) {
		return nil, fmt.Errorf("invalid refresh")
	}

	newAccess := utils.GenerateToken(32)
	newRefresh := utils.GenerateToken(32)
	newHash := utils.GenerateHMAC(HMACRefreshTokenKey, []byte(newRefresh))

	session.PreviousRefreshHash = hash
	session.RefreshHash = newHash
	session.LastActivityAt = time.Now()

	if _, err := s.sessionSvc.Update(ctx, session); err != nil {
		log.Error().Err(err).Msg("Failed to refresh session")
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	sessionJSON, err := sonic.Marshal(session)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal session")
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	if err = s.rdb.Set(ctx, fmt.Sprintf("user_sessions:%s", newAccess), sessionJSON, 15*time.Minute); err != nil {
		log.Error().Err(err).Msg("Failed to create session in Redis")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err = s.rdb.SAdd(ctx, fmt.Sprintf("user_access_tokens_%d", session.UserID), newAccess); err != nil {
		log.Error().Err(err).Msg("Failed to add member in Redis")
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	return &dto.RefreshDTO{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		ExpiresAt:    15 * time.Minute,
	}, nil
}

func (s *Service) Logout(ctx context.Context, access, refresh string) error {
	hash := utils.GenerateHMAC(HMACRefreshTokenKey, []byte(refresh))

	session, err := s.sessionSvc.GetByRefreshHash(ctx, hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user session by refresh hash")
		return fmt.Errorf("failed to get user session by refresh hash")
	}

	if session.Revoked {
		log.Warn().Int("user_id", session.UserID).Int("session_id", session.ID).Msg("Already logout")
		return fmt.Errorf("already logout")
	}

	session.Revoked = true
	session.LastActivityAt = time.Now()

	if _, err := s.sessionSvc.Update(ctx, session); err != nil {
		log.Error().Err(err).Msg("Failed to logout session")
		return fmt.Errorf("failed to logout session: %w", err)
	}

	if err = s.rdb.Delete(ctx, fmt.Sprintf("user_sessions:%s", access)); err != nil {
		log.Error().Err(err).Msg("Failed to delete session in Redis")
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if err = s.rdb.SRem(ctx, fmt.Sprintf("user_access_tokens_%d", session.UserID), access); err != nil {
		log.Error().Err(err).Msg("Failed to delete member in Redis")
		return fmt.Errorf("failed to delete member: %w", err)
	}

	return nil
}
