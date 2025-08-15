package usecase

import "my_blog_backend/internal/repository"

// TODO: Реализовать сервис, перенести логику валидации сюда
type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(c repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: c}
}
