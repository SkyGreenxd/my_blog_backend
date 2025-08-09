package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"log"
	"my_blog_backend/internal/entities"
	"my_blog_backend/pkg/e"
)

type ArticleRepository struct {
	DB *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{
		DB: db,
	}
}

func (a *ArticleRepository) Create(ctx context.Context, article *entities.Article) error {
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
		wrappedErr := e.Wrap(op, err)
		log.Print(wrappedErr)
		return wrappedErr
	}

	user.ID = userModel.ID
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt

	log.Printf("%s: user saved successfully", op)
	return nil
}
func (a *ArticleRepository) GetByID(ctx context.Context, id uint) (*entities.Article, error)
func (a *ArticleRepository) Update(ctx context.Context, article *entities.Article) error
func (a *ArticleRepository) Delete(ctx context.Context, id uint) error
func (a *ArticleRepository) ListAll(ctx context.Context) ([]entities.Article, error)
func (a *ArticleRepository) ListByAuthor(ctx context.Context, authorID uint) ([]entities.Article, error)
func (a *ArticleRepository) ListByCategory(ctx context.Context, categoryID uint) ([]entities.Article, error)

func toArticleModel(a *entities.Article) *ArticleModel {
	return &ArticleModel{
		ID:    a.ID,
		Title: a.Title,
	}
}
