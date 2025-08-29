package usecase

import (
	"my_blog_backend/internal/domain"
	"time"
)

type Services struct {
	UserService     *UserService
	ArticleService  *ArticleService
	CategoryService *CategoryService
}

func NewServices(u *UserService, a *ArticleService, c *CategoryService) *Services {
	return &Services{
		UserService:     u,
		ArticleService:  a,
		CategoryService: c,
	}
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
	OldPassword string
	NewPassword string
}

type CreateCategoryReq struct {
	UserRole     domain.Role
	CategoryName string
	CategorySlug string
}

type UpdateCategoryReq struct {
	UserRole        domain.Role
	CategorySlug    string
	NewCategoryName *string
	NewCategorySlug *string
}

type UpdateCategoryRes struct {
	CategoryName string
	CategorySlug string
}

type DeleteCategoryReq struct {
	UserRole     domain.Role
	CategorySlug string
}

type GetAllCategoriesRes struct {
	CategoryId   uint
	CategoryName string
	Slug         string
}

type ArticleRes struct {
	ArticleId    uint
	UserId       uint
	Username     string
	Title        string
	Content      string
	CategoryName string
	CategorySlug string
}

type GetArticlesByUserRes struct {
	Articles []*ArticleRes
}

type CreateArticleReq struct {
	UserId       uint
	Title        string
	Content      string
	CategorySlug string
}

type CreateArticleRes struct {
	ArticleId    uint
	Title        string
	Content      string
	CategoryName string
	CategorySlug string
}

type GetArticleRes struct {
	Title        string
	Content      string
	CategoryName string
	CategorySlug string
	Username     string
}

type UpdateUserReq struct {
	Username *string
	Email    *string
}

type GetAllArticlesRes struct {
	Title        string
	Content      string
	AuthorID     uint
	CategorySlug string
	CategoryName string
}

type UpdateArticleReq struct {
	UserId       uint
	ArticleId    uint
	Title        *string
	Content      *string
	CategorySlug *string
}

type UpdateArticleRes struct {
	AuthorID  uint
	ArticleId uint
	Title     string
	Content   string
	Category  CategoryRes
	UpdatedAt time.Time
}

type CategoryRes struct {
	CategoryName string
	CategorySlug string
	CategoryId   uint
}
