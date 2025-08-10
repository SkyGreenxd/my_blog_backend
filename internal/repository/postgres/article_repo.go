package postgres

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"my_blog_backend/internal/domain"
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

func (a *ArticleRepository) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	const op = "ArticleRepository.Create"

	articleModel := toArticleModel(article)
	result := a.DB.WithContext(ctx).Create(articleModel)
	if err := result.Error; err != nil {
		return nil, e.WrapDBError(op, err)
	}

	log.Printf("%s: article created successfully", op)
	return toArticleEntity(articleModel), nil
}

func (a *ArticleRepository) GetByID(ctx context.Context, id uint) (*domain.Article, error) {
	const op = "ArticleRepository.GetByID"

	var articleModel ArticleModel
	result := a.DB.WithContext(ctx).First(&articleModel, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, e.ErrArticleNotFound
	}

	if err := result.Error; err != nil {
		return nil, e.WrapDBError(op, err)
	}

	log.Printf("%s: article found successfully", op)
	return toArticleEntity(&articleModel), nil
}

func (a *ArticleRepository) Update(ctx context.Context, article *domain.Article) error {
	const op = "ArticleRepository.Update"

	articleModel := toArticleModel(article)
	result := a.DB.WithContext(ctx).Model(&ArticleModel{}).Where("id = ?", articleModel.ID).Updates(articleModel)
	if err := result.Error; err != nil {
		return e.WrapDBError(op, err)
	}

	if result.RowsAffected == 0 {
		return e.ErrArticleNotFound
	}

	log.Printf("%s: article updated successfully", op)
	return nil
}

func (a *ArticleRepository) Delete(ctx context.Context, id uint) error {
	const op = "ArticleRepository.Delete"

	result := a.DB.WithContext(ctx).Delete(&ArticleModel{}, id)
	if err := result.Error; err != nil {
		return e.WrapDBError(op, err)
	}

	if result.RowsAffected == 0 {
		return e.ErrArticleNotFound
	}

	log.Printf("%s: article deleted successfully", op)
	return nil
}

func (a *ArticleRepository) ListAll(ctx context.Context) ([]domain.Article, error) {
	const op = "ArticleRepository.ListAll"
	query := a.DB.WithContext(ctx)
	return a.listArticles(ctx, op, query)
}

func (a *ArticleRepository) ListByAuthor(ctx context.Context, authorID uint) ([]domain.Article, error) {
	const op = "ArticleRepository.ListByAuthor"
	query := a.DB.WithContext(ctx).Where("author_id = ?", authorID)
	return a.listArticles(ctx, op, query)
}

func (a *ArticleRepository) ListByCategory(ctx context.Context, categoryID uint) ([]domain.Article, error) {
	const op = "ArticleRepository.ListByCategory"
	query := a.DB.WithContext(ctx).Where("category_id = ?", categoryID)
	return a.listArticles(ctx, op, query)
}

func toArticleModel(a *domain.Article) *ArticleModel {
	return &ArticleModel{
		ID:         a.ID,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		Title:      a.Title,
		Content:    a.Content,
		AuthorID:   a.AuthorID,
		CategoryID: a.CategoryID,
	}
}

func toArticleEntity(a *ArticleModel) *domain.Article {
	return &domain.Article{
		ID:         a.ID,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		Title:      a.Title,
		Content:    a.Content,
		AuthorID:   a.AuthorID,
		CategoryID: a.CategoryID,
	}
}

func (a *ArticleRepository) listArticles(ctx context.Context, op string, query *gorm.DB) ([]domain.Article, error) {
	var articleModels []ArticleModel
	if err := query.Find(&articleModels).Error; err != nil {
		return nil, e.WrapDBError(op, err)
	}

	articles := make([]domain.Article, 0, len(articleModels))
	for _, model := range articleModels {
		articles = append(articles, *toArticleEntity(&model))
	}

	log.Printf("%s: found %d articles", op, len(articles))
	return articles, nil
}
