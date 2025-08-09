package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"log"
	"my_blog_backend/internal/entities"
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

func (c *CategoryRepository) Create(ctx context.Context, category *entities.Category) error {
	const op = "CategoryRepository.Create"

	categoryModel := toCategoryModel(category)
	result := c.DB.WithContext(ctx).Create(categoryModel)
	if err := result.Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return e.ErrCategoryDublicate
			}
		}

		return e.WrapDBError(op, err)
	}

	category.ID = categoryModel.ID
	category.CreatedAt = categoryModel.CreatedAt
	category.UpdatedAt = categoryModel.UpdatedAt

	log.Printf("%s: category saved successfully", op)
	return nil
}

func (c *CategoryRepository) GetByID(ctx context.Context, id uint) (*entities.Category, error) {
	const op = "CategoryRepository.GetByID"

	var categoryModel CategoryModel
	result := c.DB.WithContext(ctx).First(&categoryModel, "id = ?", id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrCategoryNotFound
		}

		return nil, e.WrapDBError(op, err)
	}

	log.Printf("%s: category found successfully", op)
	return toCategoryEntity(&categoryModel), nil
}

func (c *CategoryRepository) Update(ctx context.Context, category *entities.Category) error {
	const op = "CategoryRepository.Update"

	if err := category.Validate(); err != nil {
		return err
	}

	categoryModel := toCategoryModel(category)
	result := c.DB.Model(&CategoryModel{}).Where("id = ?", categoryModel.ID).Updates(categoryModel)
	if err := e.WrapDBError(op, result.Error); err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return e.ErrCategoryNotFound
	}

	log.Printf("%s: category updated successfully", op)
	return nil
}

func (c *CategoryRepository) Delete(ctx context.Context, id uint) error {
	const op = "CategoryRepository.Delete"

	result := c.DB.WithContext(ctx).Delete(&CategoryModel{}, id)
	if err := e.WrapDBError(op, result.Error); err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return e.ErrCategoryNotFound
	}

	log.Printf("%s: category deleted successfully", op)
	return nil
}

func (c *CategoryRepository) ListAll(ctx context.Context) ([]entities.Category, error) {
	const op = "CategoryRepository.ListAll"

	var categoryModels []CategoryModel
	result := c.DB.WithContext(ctx).Find(&categoryModels)
	if err := e.WrapDBError(op, result.Error); err != nil {
		return nil, err
	}

	var categories []entities.Category
	for _, categoryModel := range categoryModels {
		categories = append(categories, *toCategoryEntity(&categoryModel))
	}

	log.Printf("%s: categories find successfully", op)
	return categories, nil
}

func toCategoryModel(c *entities.Category) *CategoryModel {
	return &CategoryModel{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
	}
}

func toCategoryEntity(c *CategoryModel) *entities.Category {
	return &entities.Category{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
	}
}
