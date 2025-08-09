package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"log"
	"my_blog_backend/internal/domain"
	"my_blog_backend/pkg/e"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u UserRepository) Create(ctx context.Context, user *domain.User) error {
	const op = "UserRepository.Create"

	userModel := toUserModel(user)
	result := u.DB.WithContext(ctx).Create(userModel)
	if err := result.Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return e.ErrUserDublicate
			}
		}

		return e.WrapDBError(op, err)
	}

	user.ID = userModel.ID
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt

	log.Printf("%s: user saved successfully", op)
	return nil
}

func (u UserRepository) GetById(ctx context.Context, id uint) (*domain.User, error) {
	const op = "UserRepository.GetById"

	var userModel UserModel
	result := u.DB.WithContext(ctx).First(&userModel, "id = ?", id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotFound
		}

		return nil, e.WrapDBError(op, err)
	}

	return toUserEntity(&userModel), nil
}

func (u UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "UserRepository.GetByEmail"

	var userModel UserModel
	result := u.DB.WithContext(ctx).First(&userModel, "email = ?", email)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotFound
		}
		wrappedErr := e.Wrap(op, err)
		log.Print(wrappedErr)
		return nil, wrappedErr
	}

	return toUserEntity(&userModel), nil
}

func toUserModel(u *domain.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Role:         u.Role,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
}

func toUserEntity(u *UserModel) *domain.User {
	return &domain.User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Role:         u.Role,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
}
