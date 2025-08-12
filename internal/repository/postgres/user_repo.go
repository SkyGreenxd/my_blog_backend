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

func (u *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	const op = "UserRepository.Create"

	userModel := toUserModel(user)
	if err := u.DB.WithContext(ctx).Create(userModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "idx_username":
				return nil, e.ErrUsernameDuplicate
			case "idx_email":
				return nil, e.ErrEmailDuplicate
			default:
				return nil, e.ErrUserDuplicate
			}
		}

		return nil, e.WrapDBError(op, err)
	}

	log.Printf("%s: user saved successfully", op)
	return toUserEntity(userModel), nil
}

func (u *UserRepository) GetById(ctx context.Context, id uint) (*domain.User, error) {
	const op = "UserRepository.GetById"
	query := u.DB.WithContext(ctx).Where("id = ?", id)
	return u.getUser(ctx, op, query)
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "UserRepository.GetByEmail"
	query := u.DB.WithContext(ctx).Where("email = ?", email)
	return u.getUser(ctx, op, query)
}

func (u *UserRepository) Update(ctx context.Context, user *domain.User) error {
	const (
		op      = "UserRepository.Update"
		message = "user updated successfully"
	)

	userModel := toUserModel(user)
	result := u.DB.WithContext(ctx).Model(&UserModel{}).Where("id = ?", userModel.ID).Updates(userModel)

	return checkChangeQueryResult(result, op, message, e.ErrUserNotFound)
}

func (u *UserRepository) Delete(ctx context.Context, id uint) error {
	const (
		op      = "UserRepository.Delete"
		message = "user deleted successfully"
	)

	result := u.DB.WithContext(ctx).Delete(&UserModel{}, id)
	return checkChangeQueryResult(result, op, message, e.ErrUserNotFound)
}

func (u *UserRepository) ExistsByEmailOrUsername(ctx context.Context, email, username string) error {
	const op = "UserRepository.ExistsByEmailOrUsername"

	var foundUser struct {
		Username string
		Email    string
	}

	err := u.DB.WithContext(ctx).Model(&UserModel{}).
		Select("username", "email").Where("email = ?", email).
		Or("username = ?", username).First(&foundUser).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if err != nil {
		return e.WrapDBError(op, err)
	}

	if foundUser.Username == username {
		return e.ErrUsernameDuplicate
	}

	if foundUser.Email == email {
		return e.ErrEmailDuplicate
	}

	return nil
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

func (u *UserRepository) getUser(ctx context.Context, op string, query *gorm.DB) (*domain.User, error) {
	const message = "user get successfully"
	var userModel UserModel
	result := query.First(&userModel)
	if err := checkGetQueryResult(result, op, message, e.ErrUserNotFound); err != nil {
		return nil, err
	}
	return toUserEntity(&userModel), nil
}
