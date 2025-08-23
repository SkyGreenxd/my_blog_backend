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

// TODO: вынести повторяющийся код в функции
// вынести логику создания сессий в LoginUser и RefreshSession в отдельный метод
// сделать одну функцию для ошибок, заменить handleUserError
// создание сессии повторяющийся код

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
		if errors.Is(err, e.ErrUsernameIsExists) {
			return nil, e.Wrap(op, e.ErrUsernameIsExists)
		}

		if errors.Is(err, e.ErrEmailIsExists) {
			return nil, e.Wrap(op, e.ErrEmailIsExists)
		}

		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	hash, err := s.hashManager.HashPassword(userDto.Password)
	if err != nil {
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	newUser := domain.NewUser(userDto.Username, userDto.Email, hash)

	if err := newUser.Validate(); err != nil {
		return nil, e.Wrap(op, e.ErrUsernameIsForbidden)
	}

	userEntity, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	return toUserResponse(userEntity), nil
}

func (s *UserService) LoginUser(ctx context.Context, userDto *LoginUserReq) (*LoginUserRes, error) {
	const op = "UserService.LoginUser"

	user, err := s.userRepo.GetByEmail(ctx, userDto.Email)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, e.Wrap(op, e.ErrInvalidCredentials)
		}
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	if err := s.hashManager.Compare(userDto.Password, user.PasswordHash); err != nil {
		if errors.Is(err, e.ErrMismatchedHashAndPassword) {
			return nil, e.Wrap(op, e.ErrInvalidCredentials)
		}
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	jwtStruct, refreshToken, refreshTokenHash, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	session, err := s.sessionRepo.Create(ctx, domain.NewSession(
		user.ID,
		refreshTokenHash,
		time.Now().UTC().Add(refreshTokenTTL)),
	)
	if err != nil {
		if errors.Is(err, e.ErrRefreshTokenHashDuplicate) {
			return nil, e.Wrap(op, e.ErrRefreshTokenHashDuplicate)
		}
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	return toLoginUserResponse(user, session, jwtStruct, refreshToken), nil
}

func (s *UserService) GetUserById(ctx context.Context, id uint) (*UserRes, error) {
	const op = "UserService.GetUser"

	user, err := s.getUser(ctx, UserFilter{Id: &id})
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, e.Wrap(op, e.ErrUserNotFound)
		}
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	return toUserResponse(user), nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*UserRes, error) {
	const op = "UserService.GetUserByUsername"

	user, err := s.getUser(ctx, UserFilter{Username: &username})
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, e.Wrap(op, e.ErrUserNotFound)
		}

		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	return toUserResponse(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, userId uint, req *UpdateUserReq) (*UserRes, error) {
	const op = "UserService.UpdateUser"

	user, err := s.getUser(ctx, UserFilter{Id: &userId})
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, e.Wrap(op, e.ErrUserNotFound)
		}
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	if req.Username == nil && req.Email == nil {
		return nil, e.Wrap(op, e.ErrNoDataToUpdate)
	}

	if req.Username != nil {
		if err := user.ChangeUsername(*req.Username); err != nil {
			return nil, e.Wrap(op, e.ErrUsernameIsSame)
		}
		user.Username = *req.Username
	}
	if req.Email != nil {
		if err := user.ChangeEmail(*req.Email); err != nil {
			return nil, e.Wrap(op, e.ErrEmailIsSame)
		}
		user.Email = *req.Email
	}

	if err := user.Validate(); err != nil {
		return nil, e.Wrap(op, err)
	}

	updateUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	return toUserResponse(updateUser), nil
}

func (s *UserService) ChangePassword(ctx context.Context, userId uint, changePassword *ChangePasswordReq) error {
	const op = "UserService.ChangePassword"

	user, err := s.getUser(ctx, UserFilter{Id: &userId})
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return e.Wrap(op, e.ErrUserNotFound)
		}
		return e.Wrap(op, e.ErrInternalServer)
	}

	if err := s.hashManager.Compare(changePassword.OldPassword, user.PasswordHash); err != nil {
		if errors.Is(err, e.ErrMismatchedHashAndPassword) {
			return e.Wrap(op, e.ErrInvalidCredentials)
		}

		return e.Wrap(op, e.ErrInternalServer)
	}

	newPassHash, err := s.hashManager.HashPassword(changePassword.NewPassword)
	if err != nil {
		return e.Wrap(op, e.ErrInternalServer)
	}

	if err := user.ChangePassword(newPassHash); err != nil {
		if errors.Is(err, e.ErrPasswordIsSame) {
			return e.Wrap(op, e.ErrPasswordIsSame)
		}
		return e.Wrap(op, e.ErrInternalServer)
	}

	if _, err := s.userRepo.Update(ctx, user); err != nil {
		return e.Wrap(op, e.ErrInternalServer)
	}

	return nil
}

func (s *UserService) RefreshSession(ctx context.Context, userRefreshToken string) (*LoginUserRes, error) {
	const op = "UserService.RefreshSession"

	oldSession, err := s.verifyRefreshToken(ctx, userRefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrRefreshTokenInvalid),
			errors.Is(err, e.ErrSessionRevoked),
			errors.Is(err, e.ErrSessionExpired):
			return nil, e.Wrap(op, e.ErrUnauthorized)
		default:
			return nil, e.Wrap(op, e.ErrInternalServer)
		}
	}

	if err := s.sessionRepo.RevokeSession(ctx, oldSession.Id); err != nil {
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	user, err := s.userRepo.GetById(ctx, oldSession.UserId)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, e.Wrap(op, e.ErrUserNotFound)
		}
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	jwtStruct, refreshToken, refreshTokenHash, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	newSession, err := s.sessionRepo.Create(ctx, domain.NewSession(
		user.ID,
		refreshTokenHash,
		time.Now().UTC().Add(refreshTokenTTL)),
	)
	if err != nil {
		if errors.Is(err, e.ErrRefreshTokenHashDuplicate) {
			return nil, e.Wrap(op, e.ErrRefreshTokenHashDuplicate)
		}
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	return toLoginUserResponse(user, newSession, jwtStruct, refreshToken), nil
}

func (s *UserService) LogoutUser(ctx context.Context, userRefreshToken string) error {
	const op = "UserService.LogoutUser"

	session, err := s.verifyRefreshToken(ctx, userRefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrRefreshTokenInvalid),
			errors.Is(err, e.ErrSessionRevoked),
			errors.Is(err, e.ErrSessionExpired):
			return e.Wrap(op, e.ErrUnauthorized)
		default:
			return e.Wrap(op, e.ErrInternalServer)
		}
	}

	if err := s.sessionRepo.RevokeSession(ctx, session.Id); err != nil {
		return e.Wrap(op, e.ErrInternalServer)
	}

	return nil
}

func (s *UserService) SetAdminRole(ctx context.Context, userId uint) error {
	const op = "UserService.SetAdminRole"

	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return e.Wrap(op, e.ErrUserNotFound)
		}

		return e.Wrap(op, e.ErrInternalServer)
	}

	if err := user.SetAdminRole(); err != nil {
		return e.Wrap(op, e.ErrUserAlreadyAdmin)
	}

	if _, err := s.userRepo.Update(ctx, user); err != nil {
		return e.Wrap(op, e.ErrInternalServer)
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
	Id       *uint
	Username *string
}

func (s *UserService) getUser(ctx context.Context, filter UserFilter) (*domain.User, error) {
	if filter.Id != nil {
		user, err := s.userRepo.GetById(ctx, *filter.Id)
		return handleUserError(user, err)
	}

	if filter.Username != nil {
		user, err := s.userRepo.GetByEmail(ctx, *filter.Username)
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
