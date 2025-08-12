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

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		DB: db,
	}
}

func (c *CategoryRepository) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	const op = "CategoryRepository.Create"

	categoryModel := toCategoryModel(category)
	result := c.DB.WithContext(ctx).Create(categoryModel)
	if err := result.Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, e.ErrCategoryDuplicate
			}
		}

		return nil, e.WrapDBError(op, err)
	}

	log.Printf("%s: category saved successfully", op)
	return toCategoryEntity(categoryModel), nil
}

func (c *CategoryRepository) GetByID(ctx context.Context, id uint) (*domain.Category, error) {
	const (
		op      = "CategoryRepository.GetByID"
		message = "category found successfully"
	)

	var categoryModel CategoryModel
	result := c.DB.WithContext(ctx).First(&categoryModel, "id = ?", id)
	if err := checkGetQueryResult(result, op, message, e.ErrCategoryNotFound); err != nil {
		return nil, err
	}

	return toCategoryEntity(&categoryModel), nil
}

func (c *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	const (
		op      = "CategoryRepository.Update"
		message = "category updated successfully"
	)

	categoryModel := toCategoryModel(category)
	result := c.DB.WithContext(ctx).Model(&CategoryModel{}).Where("id = ?", categoryModel.ID).Updates(categoryModel)

	return checkChangeQueryResult(result, op, message, e.ErrCategoryNotFound)
}

func (c *CategoryRepository) Delete(ctx context.Context, id uint) error {
	const op = "CategoryRepository.Delete"

	result := c.DB.WithContext(ctx).Delete(&CategoryModel{}, id)
	if err := result.Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return e.ErrCategoryInUse
		}
		return e.WrapDBError(op, err)
	}

	if result.RowsAffected == 0 {
		return e.ErrCategoryNotFound
	}

	log.Printf("%s: category deleted successfully", op)
	return nil
}

func (c *CategoryRepository) ListAll(ctx context.Context) ([]domain.Category, error) {
	const op = "CategoryRepository.ListAll"

	var categoryModels []CategoryModel
	result := c.DB.WithContext(ctx).Find(&categoryModels)
	if err := e.WrapDBError(op, result.Error); err != nil {
		return nil, err
	}

	categories := make([]domain.Category, 0, len(categoryModels))
	for _, categoryModel := range categoryModels {
		categories = append(categories, *toCategoryEntity(&categoryModel))
	}

	log.Printf("%s: categories find successfully", op)
	return categories, nil
}

func toCategoryModel(c *domain.Category) *CategoryModel {
	return &CategoryModel{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
	}
}

func toCategoryEntity(c *CategoryModel) *domain.Category {
	return &domain.Category{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
	}
}
