package usecase

import (
	"context"
	"errors"
	"my_blog_backend/internal/domain"
	"my_blog_backend/internal/repository"
	"my_blog_backend/pkg/e"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(c repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: c}
}

func (s *CategoryService) Create(ctx context.Context, req *CreateCategoryReq) (string, error) {
	const op = "CategoryService.Create"

	if req.UserRole != domain.RoleAdmin {
		return "", e.Wrap(op, e.ErrPermissionDenied)
	}

	newCategory := domain.NewCategory(req.CategoryName)

	categoryEntity, err := s.categoryRepo.Create(ctx, newCategory)
	if err != nil {
		if errors.Is(err, e.ErrCategoryInUse) {
			return "", e.Wrap(op, err)
		}

		return "", e.Wrap(op, e.ErrInternalServer)
	}

	return categoryEntity.Name, nil
}

func (s *CategoryService) GetAll(ctx context.Context) ([]string, error) {
	const op = "CategoryService.GetAll"

	categories, err := s.categoryRepo.ListAll(ctx)
	if err != nil {
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	result := make([]string, len(categories))
	for i, category := range categories {
		result[i] = category.Name
	}

	return result, nil
}

func (s *CategoryService) Update(ctx context.Context, req *UpdateCategoryReq) error {
	const op = "CategoryService.Update"

	if req.UserRole != domain.RoleAdmin {
		return e.Wrap(op, e.ErrPermissionDenied)
	}

	category, err := s.categoryRepo.GetByID(ctx, req.CategoryId)
	if err != nil {
		if errors.Is(err, e.ErrCategoryNotFound) {
			return e.Wrap(op, e.ErrCategoryNotFound)
		}

		return e.Wrap(op, e.ErrInternalServer)
	}

	if err := category.ChangeName(req.NewCategoryName); err != nil {
		return e.Wrap(op, err)
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return e.Wrap(op, e.ErrInternalServer)
	}

	return nil
}

func (s *CategoryService) Delete(ctx context.Context, req *DeleteCategoryReq) error {
	const op = "CategoryService.Delete"

	if req.UserRole != domain.RoleAdmin {
		return e.Wrap(op, e.ErrPermissionDenied)
	}

	if err := s.categoryRepo.Delete(ctx, req.CategoryId); err != nil {
		if errors.Is(err, e.ErrCategoryNotFound) {
			return e.Wrap(op, e.ErrCategoryNotFound)
		}

		if errors.Is(err, e.ErrCategoryInUse) {
			return e.Wrap(op, e.ErrCategoryInUse)
		}

		return e.Wrap(op, e.ErrInternalServer)
	}

	return nil
}
