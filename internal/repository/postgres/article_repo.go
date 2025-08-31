package postgres

import (
	"context"
	"errors"
	"my_blog_backend/internal/domain"
	"my_blog_backend/pkg/e"

	"gorm.io/gorm"
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
		return nil, e.Wrap(op, err)
	}

	return toArticleEntity(articleModel), nil
}

func (a *ArticleRepository) GetByID(ctx context.Context, id uint) (*domain.Article, error) {
	const op = "ArticleRepository.GetByID"
	var articleModel ArticleModel

	result := a.DB.WithContext(ctx).
		Preload("Author").
		Preload("Category").
		First(&articleModel, "id = ?", id)

	if err := checkGetQueryResult(result, e.ErrArticleNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toArticleEntity(&articleModel), nil
}

func (a *ArticleRepository) Update(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	const op = "ArticleRepository.Update"
	articleModel := toArticleModel(article)
	updates := map[string]interface{}{
		"category_id": articleModel.Category.ID,
		"title":       articleModel.Title,
		"content":     articleModel.Content,
	}
	result := a.DB.WithContext(ctx).Model(&ArticleModel{}).Where("id = ?", articleModel.ID).Updates(updates)
	if err := checkChangeQueryResult(result, e.ErrArticleNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	updArticle, err := a.GetByID(ctx, articleModel.ID)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updArticle, nil
}

func (a *ArticleRepository) Delete(ctx context.Context, id uint) error {
	const op = "ArticleRepository.Delete"
	result := a.DB.WithContext(ctx).Delete(&ArticleModel{}, id)
	if err := checkChangeQueryResult(result, e.ErrArticleNotFound); err != nil {
		return e.Wrap(op, err)
	}

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

func (a *ArticleRepository) ExistsByTitleContentAuthor(ctx context.Context, article *domain.Article) error {
	const op = "ArticleRepository.ExistsByTitleContentAuthor"

	articleModel := toArticleModel(article)
	result := a.DB.WithContext(ctx).Where(map[string]interface{}{
		"title":     articleModel.Title,
		"content":   articleModel.Content,
		"author_id": articleModel.AuthorID,
	}).First(&articleModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		return e.Wrap(op, result.Error)
	}

	return e.Wrap(op, e.ErrArticleDuplicate)
}

func (a *ArticleRepository) listArticles(ctx context.Context, op string, query *gorm.DB) ([]domain.Article, error) {
	var articleModels []ArticleModel
	result := query.Preload("Author").Preload("Category").Find(&articleModels)
	if err := checkGetQueryResult(result, e.ErrArticleNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	articles := make([]domain.Article, 0, len(articleModels))
	for _, model := range articleModels {
		articles = append(articles, *toArticleEntity(&model))
	}

	return articles, nil
}

func toArticleModel(a *domain.Article) *ArticleModel {
	model := &ArticleModel{
		ID:         a.ID,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		Title:      a.Title,
		Content:    a.Content,
		AuthorID:   a.AuthorID,
		CategoryID: a.CategoryID,
	}

	if a.Author != nil {
		model.Author = toUserModel(a.Author)
	}
	if a.Category != nil {
		model.Category = toCategoryModel(a.Category)
	}

	return model
}

func toArticleEntity(a *ArticleModel) *domain.Article {
	entity := &domain.Article{
		ID:         a.ID,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		Title:      a.Title,
		Content:    a.Content,
		AuthorID:   a.AuthorID,
		CategoryID: a.CategoryID,
	}

	if a.Author != nil {
		entity.Author = toUserEntity(a.Author)
	}

	if a.Category != nil {
		entity.Category = toCategoryEntity(a.Category)
	}

	return entity
}
