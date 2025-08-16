package usecase

import (
	"my_blog_backend/internal/domain"
	"time"
)

type Services struct {
	UserService     UserService
	ArticleService  ArticleService
	CategoryService CategoryService
}

type AuthenticatedUser struct {
	ID    uint
	Role  domain.Role
	Email string
}

type TokenResponse struct {
	Token     string
	ExpiresAt time.Time
}

type CreateUserReq struct {
	Username string
	Email    string
	Password string
}

type LoginUserReq struct {
	Email    string
	Password string
}

type UserRes struct {
	Username string
	Email    string
	Role     domain.Role
}

type LoginUserRes struct {
	SessionID             string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	User                  UserRes
}

type ChangePasswordReq struct {
	Id          uint
	NewPassword string
}

type CreateCategoryReq struct {
	UserRole     domain.Role
	CategoryName string
}

type UpdateCategoryReq struct {
	UserRole        domain.Role
	CategoryId      uint
	NewCategoryName string
}

type DeleteCategoryReq struct {
	UserRole   domain.Role
	CategoryId uint
}
