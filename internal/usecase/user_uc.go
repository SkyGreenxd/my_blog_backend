package usecase

import (
	"context"
	"errors"
	"my_blog_backend/internal/domain"
	"my_blog_backend/internal/repository"
	"my_blog_backend/pkg/e"
	"time"
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

	if err := userDto.Validate(); err != nil {
		return nil, err
	}

	err := s.userRepo.ExistsByEmailOrUsername(ctx, userDto.Email, userDto.Username)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	hash, err := s.hashManager.HashPassword(userDto.Password)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	newUser := &domain.User{
		Role:         domain.RoleUser,
		Username:     userDto.Username,
		Email:        userDto.Email,
		PasswordHash: hash,
	}

	userEntity, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toUserResponse(userEntity), nil
}

func (s *UserService) LoginUser(ctx context.Context, userDto LoginUserReq) (*LoginUserRes, error) {
	const (
		op              = "UserService.LoginUser"
		refreshTokenTTL = 30 * 24 * time.Hour
	)
	// Проверить DTO.
	if err := userDto.Validate(); err != nil {
		return nil, e.Wrap(op, err)
	}
	//	Найти пользователя в базе по email.Если не найден — вернуть ошибку.
	user, err := s.userRepo.GetByEmail(ctx, userDto.Email)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, e.ErrInvalidCredentials
		}
		return nil, e.Wrap(op, err)
	}
	//	Сравнить пароль из DTO с хэшем из базы данных с помощью hashManager. Если пароли не совпадают — вернуть ошибку.
	if err := s.hashManager.Compare(userDto.Password, user.PasswordHash); err != nil {
		if errors.Is(err, e.ErrMismatchedHashAndPassword) {
			return nil, e.ErrInvalidCredentials
		}
		return nil, e.Wrap(op, err)
	}
	//	Если все хорошо — сгенерировать токены с помощью tokenManager.
	jwtStruct, err := s.tokenManager.NewJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	refreshToken, refreshTokenHash, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	//Создать сессию и загрузить ее в бд
	session, err := s.sessionRepo.Create(ctx, domain.NewSession(
		user.ID,
		refreshTokenHash,
		time.Now().UTC().Add(refreshTokenTTL)),
	)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	//	Вернуть response, заполненный данными
	return toLoginUserResponse(user, session, *jwtStruct, refreshToken), nil
}

func (s *UserService) GetUser(ctx context.Context, username string) (*UserRes, error) {
	const op = "UserService.GetUser"
}
func (s *UserService) UpdateUser(ctx context.Context, userID uint, updateUserDto *UpdateUserReq) (*UserRes, error)

// TODO: DTO должен содержать OldPassword и NewPassword.
func (s *UserService) ChangePassword(ctx context.Context, userID uint, changePassDto *ChangePasswordReq) error
func (s *UserService) RefreshSession(ctx context.Context, refreshToken string) (*LoginUserRes, error)
func (s *UserService) LogoutUser(ctx context.Context, refreshToken string) error

func toUserResponse(user *domain.User) *UserRes {
	return &UserRes{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}

func toLoginUserResponse(user *domain.User, session *domain.Session, accessToken TokenResponse, refreshToken string) *LoginUserRes {
	return &LoginUserRes{
		SessionID:             session.Id.String(),
		AccessToken:           accessToken.Token,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessToken.ExpiresAt,
		RefreshTokenExpiresAt: session.ExpiresAt,
		User:                  *toUserResponse(user),
	}
}
