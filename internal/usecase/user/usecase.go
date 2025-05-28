package user

import (
	"context"
	"diprec_api/internal/domain"
	"diprec_api/internal/repository/user"
	"diprec_api/internal/service"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userUseCase struct {
	repo   user.IUserRepository
	auth   *service.AuthService
	logger *zap.Logger
}

type IUserUseCase interface {
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	Authenticate(ctx context.Context, username, password string) (*domain.User, error)
	GetMe(ctx context.Context, userID uint) (*domain.User, error)
	GenerateTokens(user *domain.User) (*domain.TokenPair, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
}

func NewUserUseCase(repo user.IUserRepository, auth *service.AuthService, logger *zap.Logger) IUserUseCase {
	return &userUseCase{
		repo:   repo,
		auth:   auth,
		logger: logger.Named("UserUseCase"),
	}
}

func (uc *userUseCase) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	existing, err := uc.repo.GetByUsername(ctx, user.Username)
	if existing != nil {
		uc.logger.Error("user with username "+user.Username+" already exists", zap.Any("user", user))
		return nil, domain.ErrUserExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		uc.logger.Error("error while seeking for user", zap.Error(err))
		return nil, err
	}

	if err := user.SetPassword(user.Password); err != nil {
		uc.logger.Error("failed to set password", zap.Error(err))
		return nil, err
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		uc.logger.Error("user creation failed", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) Authenticate(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	if err := user.CheckPassword(password); err != nil {
		uc.logger.Error("invalid password attempt", zap.String("username", username))
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (uc *userUseCase) GenerateTokens(user *domain.User) (*domain.TokenPair, error) {
	return uc.auth.GenerateTokens(user)
}

func (uc *userUseCase) GetMe(ctx context.Context, userID uint) (*domain.User, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (uc *userUseCase) RefreshTokens(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	userID, role, err := uc.auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	pair, err := uc.auth.GenerateTokens(&domain.User{ID: userID, Role: domain.Role(role)})
	if err != nil {
		return nil, err
	}

	return pair, nil
}
