package postgres

import (
	"context"
	"my_blog_backend/internal/domain"
	"my_blog_backend/pkg/e"
	"time"

	"gorm.io/gorm"
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
	if err := postgresDuplicate(result, e.ErrCategoryIsExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toCategoryEntity(categoryModel), nil
}

func (c *CategoryRepository) GetByID(ctx context.Context, id uint) (*domain.Category, error) {
	const op = "CategoryRepository.GetByID"
	var categoryModel CategoryModel
	result := c.DB.WithContext(ctx).First(&categoryModel, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrCategoryNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toCategoryEntity(&categoryModel), nil
}

func (c *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	const op = "CategoryRepository.GetBySlug"
	var categoryModel CategoryModel
	result := c.DB.WithContext(ctx).First(&categoryModel, "slug = ?", slug)
	if err := checkGetQueryResult(result, e.ErrCategoryNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toCategoryEntity(&categoryModel), nil
}

func (c *CategoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	const op = "CategoryRepository.GetByName"
	var categoryModel CategoryModel
	result := c.DB.WithContext(ctx).First(&categoryModel, "name = ?", name)
	if err := checkGetQueryResult(result, e.ErrCategoryNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toCategoryEntity(&categoryModel), nil
}

func (c *CategoryRepository) Update(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	const op = "CategoryRepository.Update"

	categoryModel := toCategoryModel(category)
	updates := map[string]interface{}{
		"name":       categoryModel.Name,
		"updated_at": time.Now().UTC(),
		"slug":       categoryModel.Slug,
	}
	result := c.DB.WithContext(ctx).Model(&CategoryModel{}).Where("id = ?", category.ID).Updates(updates)
	if err := checkChangeQueryResult(result, e.ErrCategoryNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	updCategory, err := c.GetByID(ctx, category.ID)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updCategory, nil
}

func (c *CategoryRepository) Delete(ctx context.Context, id uint) error {
	const op = "CategoryRepository.Delete"
	result := c.DB.WithContext(ctx).Delete(&CategoryModel{}, id)
	if err := postgresForeignKeyViolation(result, e.ErrCategoryInUse); err != nil {
		return e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return e.Wrap(op, e.ErrCategoryNotFound)
	}

	return nil
}

func (c *CategoryRepository) ListAll(ctx context.Context) ([]domain.Category, error) {
	const op = "CategoryRepository.ListAll"
	var categoryModels []CategoryModel
	result := c.DB.WithContext(ctx).Find(&categoryModels)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	categories := make([]domain.Category, 0, len(categoryModels))
	for _, categoryModel := range categoryModels {
		categories = append(categories, *toCategoryEntity(&categoryModel))
	}

	return categories, nil
}

func toCategoryModel(c *domain.Category) *CategoryModel {
	return &CategoryModel{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
		Slug:      c.Slug,
	}
}

func toCategoryEntity(c *CategoryModel) *domain.Category {
	return &domain.Category{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
		Slug:      c.Slug,
	}
}
