package usecase

import (
	"context"
	"errors"
	"my_blog_backend/internal/domain"
	"my_blog_backend/internal/repository"
	"my_blog_backend/pkg/e"
	"time"
)

const (
	refreshTokenTTL = 30 * 24 * time.Hour
)

type UserService struct {
	userRepo     repository.UserRepository
	articleRepo  repository.ArticleRepository
	sessionRepo  repository.SessionRepository
	tokenManager TokenManager
	hashManager  HashManager
}

func NewUserService(u repository.UserRepository, a repository.ArticleRepository, s repository.SessionRepository, tm TokenManager, hm HashManager) *UserService {
	return &UserService{
		userRepo:     u,
		articleRepo:  a,
		sessionRepo:  s,
		tokenManager: tm,
		hashManager:  hm,
	}
}

func (s *UserService) CreateUser(ctx context.Context, userDto *CreateUserReq) (*UserRes, error) {
	const op = "UserService.CreateUser"

	err := s.userRepo.ExistsByEmailOrUsername(ctx, userDto.Email, userDto.Username)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	hash, err := s.hashManager.HashPassword(userDto.Password)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	newUser, err := domain.NewUser(userDto.Username, userDto.Email, hash)

	userEntity, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toUserResponse(userEntity), nil
}

func (s *UserService) LoginUser(ctx context.Context, userDto *LoginUserReq) (*LoginUserRes, error) {
	const op = "UserService.LoginUser"

	user, err := s.userRepo.GetByEmail(ctx, userDto.Email)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, e.Wrap(op, e.ErrInvalidEmail)
		}
		return nil, e.Wrap(op, err)
	}

	if err := s.hashManager.Compare(userDto.Password, user.PasswordHash); err != nil {
		if errors.Is(err, e.ErrMismatchedHashAndPassword) {
			return nil, e.Wrap(op, e.ErrInvalidPassword)
		}
		return nil, e.Wrap(op, err)
	}

	jwtStruct, refreshToken, refreshTokenHash, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	session, err := s.sessionRepo.Create(ctx, domain.NewSession(
		user.ID,
		refreshTokenHash,
		time.Now().UTC().Add(refreshTokenTTL)),
	)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toLoginUserResponse(user, session, jwtStruct, refreshToken), nil
}

func (s *UserService) GetUserById(ctx context.Context, id uint) (*UserRes, error) {
	const op = "UserService.GetUser"

	user, err := s.getUser(ctx, UserFilter{Id: &id})
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toUserResponse(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID uint) (*UserRes, error) {
	const op = "UserService.UpdateUser"

	user, err := s.getUser(ctx, UserFilter{Id: &userID})
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	updateUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toUserResponse(updateUser), nil
}

func (s *UserService) ChangePassword(ctx context.Context, changePassword *ChangePasswordReq) error {
	const op = "UserService.ChangePassword"

	user, err := s.getUser(ctx, UserFilter{Id: &changePassword.Id})
	if err != nil {
		return e.Wrap(op, err)
	}

	newPassHash, err := s.hashManager.HashPassword(changePassword.NewPassword)
	if err != nil {
		return e.Wrap(op, err)
	}

	if err := user.ChangePassword(newPassHash); err != nil {
		return e.Wrap(op, err)
	}

	if _, err := s.userRepo.Update(ctx, user); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (s *UserService) RefreshSession(ctx context.Context, userRefreshToken string) (*LoginUserRes, error) {
	const op = "UserService.RefreshSession"

	oldSession, err := s.verifyRefreshToken(ctx, userRefreshToken)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	if err := s.sessionRepo.RevokeSession(ctx, oldSession.Id); err != nil {
		return nil, e.Wrap(op, err)
	}

	user, err := s.userRepo.GetById(ctx, oldSession.UserId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	jwtStruct, refreshToken, refreshTokenHash, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	newSession, err := s.sessionRepo.Create(ctx, domain.NewSession(
		user.ID,
		refreshTokenHash,
		time.Now().UTC().Add(refreshTokenTTL)),
	)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toLoginUserResponse(user, newSession, jwtStruct, refreshToken), nil
}

func (s *UserService) LogoutUser(ctx context.Context, userRefreshToken string) error {
	const op = "UserService.LogoutUser"

	session, err := s.verifyRefreshToken(ctx, userRefreshToken)
	if err != nil {
		return e.Wrap(op, err)
	}

	if err := s.sessionRepo.RevokeSession(ctx, session.Id); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (s *UserService) verifyRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	tokenHash := s.tokenManager.HashRefreshToken(refreshToken)
	session, err := s.sessionRepo.GetByRefreshTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, e.ErrSessionNotFound) {
			return nil, e.ErrRefreshTokenInvalid
		}

		return nil, err
	}

	if err := session.ValidateState(); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *UserService) generateTokens(userId uint, email string, role domain.Role) (*TokenResponse, string, string, error) {
	jwtStruct, err := s.tokenManager.NewJWT(userId, email, role)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, refreshTokenHash, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, "", "", err
	}

	return jwtStruct, refreshToken, refreshTokenHash, nil
}

type UserFilter struct {
	Id    *uint
	Email *string
}

func (s *UserService) getUser(ctx context.Context, filter UserFilter) (*domain.User, error) {
	if filter.Id != nil {
		user, err := s.userRepo.GetById(ctx, *filter.Id)
		return handleUserError(user, err)
	}

	if filter.Email != nil {
		user, err := s.userRepo.GetByEmail(ctx, *filter.Email)
		return handleUserError(user, err)
	}

	return nil, e.ErrUserNotFound
}

func handleUserError(user *domain.User, err error) (*domain.User, error) {
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, e.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func toUserResponse(user *domain.User) *UserRes {
	return &UserRes{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}

func toLoginUserResponse(user *domain.User, session *domain.Session, accessToken *TokenResponse, refreshToken string) *LoginUserRes {
	return &LoginUserRes{
		SessionID:             session.Id.String(),
		AccessToken:           accessToken.Token,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessToken.ExpiresAt,
		RefreshTokenExpiresAt: session.ExpiresAt,
		User:                  *toUserResponse(user),
	}
}
