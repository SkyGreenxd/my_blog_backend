package usecase

import (
	"context"
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

	newCategory := domain.NewCategory(req.CategoryName, req.CategorySlug)

	categoryEntity, err := s.categoryRepo.Create(ctx, newCategory)
	if err != nil {
		return "", e.Wrap(op, err)
	}

	return categoryEntity.Name, nil
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*GetAllCategoriesRes, error) {
	const op = "CategoryService.GetAll"

	categories, err := s.categoryRepo.ListAll(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	result := make([]*GetAllCategoriesRes, len(categories))
	for i, category := range categories {
		result[i] = ToGetAllCategoriesRes(&category)
	}

	return result, nil
}

func (s *CategoryService) Update(ctx context.Context, req *UpdateCategoryReq) (*UpdateCategoryRes, error) {
	const op = "CategoryService.Update"

	if req.UserRole != domain.RoleAdmin {
		return nil, e.Wrap(op, e.ErrPermissionDenied)
	}

	category, err := s.categoryRepo.GetBySlug(ctx, req.CategorySlug)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	if req.NewCategoryName == nil && req.NewCategorySlug == nil {
		return nil, e.Wrap(op, e.ErrNoDataToUpdate)
	}

	if req.NewCategoryName != nil {
		if err := category.ChangeName(*req.NewCategoryName); err != nil {
			return nil, e.Wrap(op, err)
		}
	}

	if req.NewCategorySlug != nil {
		if err := category.ChangeSlug(*req.NewCategorySlug); err != nil {
			return nil, e.Wrap(op, err)
		}
	}

	updCategory, err := s.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return ToUpdateCategoryRes(updCategory), nil
}

func (s *CategoryService) Delete(ctx context.Context, req *DeleteCategoryReq) error {
	const op = "CategoryService.Delete"

	if req.UserRole != domain.RoleAdmin {
		return e.Wrap(op, e.ErrPermissionDenied)
	}

	category, err := s.categoryRepo.GetBySlug(ctx, req.CategorySlug)
	if err != nil {
		return e.Wrap(op, err)
	}

	if err := s.categoryRepo.Delete(ctx, category.ID); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func ToGetAllCategoriesRes(category *domain.Category) *GetAllCategoriesRes {
	return &GetAllCategoriesRes{
		CategoryId:   category.ID,
		CategoryName: category.Name,
		Slug:         category.Slug,
	}
}

func ToUpdateCategoryRes(category *domain.Category) *UpdateCategoryRes {
	return &UpdateCategoryRes{
		CategoryName: category.Name,
		CategorySlug: category.Slug,
	}
}
