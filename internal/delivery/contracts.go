package delivery

import (
	"my_blog_backend/internal/domain"
	"my_blog_backend/internal/usecase"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=5,max=32,nospaces"`
	Email    string `json:"email" binding:"required,email,min=3,max=320,nospaces"`
	Password string `json:"password" binding:"required,min=8,max=128,nospaces"`
}

type UserRes struct {
	Id       uint        `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     domain.Role `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,min=3,max=320,nospaces"`
	Password string `json:"password" binding:"required,min=8,max=128,nospaces"`
}

type ErrResponse struct {
	ErrCode string `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type LoginUserRes struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserRes   `json:"user"`
}

type UpdateUserReq struct {
	Username *string `json:"username" binding:"omitempty,min=5,max=32,nospaces"`
	Email    *string `json:"email" binding:"omitempty,email,min=3,max=32,nospaces"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=128,nospaces"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128,nospaces"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutUserReq struct {
	RefreshToken string `json:"refresh_token"`
}

type CreateCategoryReq struct {
	CategoryName string `json:"category_name" binding:"required,min=3,max=128,nospaces"`
	CategorySlug string `json:"category_slug" binding:"required,min=3,max=128,nospaces"`
}

type UpdateCategoryReq struct {
	NewCategoryName *string `json:"new_category_name" binding:"omitempty,min=3,max=128,nospaces"`
	NewCategorySlug *string `json:"new_category_slug" binding:"omitempty,min=3,max=128,nospaces"`
}

type CreateArticleReq struct {
	Title        string `json:"title" binding:"required,min=3,max=100"`
	Content      string `json:"content" binding:"required,min=10,max=16000"`
	CategorySlug string `json:"category_slug" binding:"required,min=3,max=128,nospaces"`
}

type CreateArticleRes struct {
	ArticleId    uint   `json:"article_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	CategoryName string `json:"category_name"`
	CategorySlug string `json:"category_slug"`
}

type UpdateArticleReq struct {
	Title        *string `json:"title" binding:"omitempty,min=3,max=100"`
	Content      *string `json:"content" binding:"omitempty,min=10,max=16000"`
	CategorySlug *string `json:"category_slug" binding:"omitempty,min=3,max=128,nospaces"`
}

func ToDeleteArticleReq(userId, articleId uint) *usecase.DeleteArticleReq {
	return &usecase.DeleteArticleReq{
		UserId:    userId,
		ArticleId: articleId,
	}
}

type ArticleRes struct {
	ArticleId uint
	Title     string
	Content   string
	Author    UserRes
	Category  CategoryRes
}

type GetArticlesByUserRes struct {
	Articles []*ArticleRes `json:"articles"`
}

func ToArticleRes(res *usecase.ArticleRes) *ArticleRes {
	return &ArticleRes{
		ArticleId: res.ArticleId,
		Title:     res.Title,
		Content:   res.Content,
		Author:    *ToUserRes(&res.Author),
		Category:  *ToCategoryRes(&res.Category),
	}
}

func ToGetArticlesByUserRes(articles []*ArticleRes) *GetArticlesByUserRes {
	return &GetArticlesByUserRes{
		Articles: articles,
	}
}

func ToUpdateArticleReq(req *UpdateArticleReq, userId, articleId uint) *usecase.UpdateArticleReq {
	return &usecase.UpdateArticleReq{
		UserId:       userId,
		ArticleId:    articleId,
		Title:        req.Title,
		Content:      req.Content,
		CategorySlug: req.CategorySlug,
	}
}

type UpdateArticleRes struct {
	AuthorID  uint        `json:"author_id"`
	ArticleId uint        `json:"article_id"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Category  CategoryRes `json:"category"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CategoryRes struct {
	CategoryName string `json:"category_name"`
	CategorySlug string `json:"category_slug"`
	CategoryId   uint   `json:"category_id"`
}

func ToUpdateArticleRes(res *usecase.UpdateArticleRes) *UpdateArticleRes {
	return &UpdateArticleRes{
		AuthorID:  res.AuthorID,
		ArticleId: res.ArticleId,
		Title:     res.Title,
		Content:   res.Content,
		Category:  *ToCategoryRes(&res.Category),
		UpdatedAt: res.UpdatedAt,
	}
}

func ToCategoryRes(res *usecase.CategoryRes) *CategoryRes {
	return &CategoryRes{
		CategoryName: res.CategoryName,
		CategorySlug: res.CategorySlug,
		CategoryId:   res.CategoryId,
	}
}

func ToChangePasswordReq(req *ChangePasswordReq) *usecase.ChangePasswordReq {
	return &usecase.ChangePasswordReq{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
}

func ToUpdateUserReq(req *UpdateUserReq) *usecase.UpdateUserReq {
	return &usecase.UpdateUserReq{
		Username: req.Username,
		Email:    req.Email,
	}
}

func ToLoginUserRes(res *usecase.LoginUserRes) *LoginUserRes {
	return &LoginUserRes{
		SessionID:             res.SessionID,
		AccessToken:           res.AccessToken,
		RefreshToken:          res.RefreshToken,
		AccessTokenExpiresAt:  res.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: res.RefreshTokenExpiresAt,
		User:                  *ToUserRes(&res.User),
	}
}

func ToServiceCreateUserReq(request *CreateUserRequest) *usecase.CreateUserReq {
	return &usecase.CreateUserReq{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
}

func ToUserRes(res *usecase.UserRes) *UserRes {
	return &UserRes{
		Id:       res.Id,
		Username: res.Username,
		Email:    res.Email,
		Role:     res.Role,
	}
}

func ToLoginUserReq(req *LoginRequest) *usecase.LoginUserReq {
	return &usecase.LoginUserReq{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToCreateCategoryReq(req *CreateCategoryReq, userRole domain.Role) *usecase.CreateCategoryReq {
	return &usecase.CreateCategoryReq{
		CategoryName: req.CategoryName,
		CategorySlug: req.CategorySlug,
		UserRole:     userRole,
	}
}

func ToDeleteCategoryReq(categorySlug string, userRole domain.Role) *usecase.DeleteCategoryReq {
	return &usecase.DeleteCategoryReq{
		UserRole:     userRole,
		CategorySlug: categorySlug,
	}
}

func ToUpdateCategoryReq(req UpdateCategoryReq, userRole domain.Role, categorySlug string) *usecase.UpdateCategoryReq {
	return &usecase.UpdateCategoryReq{
		UserRole:        userRole,
		CategorySlug:    categorySlug,
		NewCategoryName: req.NewCategoryName,
		NewCategorySlug: req.NewCategorySlug,
	}
}

func ToCreateArticleReq(req *CreateArticleReq, userId uint) *usecase.CreateArticleReq {
	return &usecase.CreateArticleReq{
		UserId:       userId,
		Title:        req.Title,
		Content:      req.Content,
		CategorySlug: req.CategorySlug,
	}
}
func ToCreateArticleRes(res *usecase.CreateArticleRes) *CreateArticleRes {
	return &CreateArticleRes{
		ArticleId:    res.ArticleId,
		Title:        res.Title,
		Content:      res.Content,
		CategoryName: res.CategoryName,
		CategorySlug: res.CategorySlug,
	}
}
