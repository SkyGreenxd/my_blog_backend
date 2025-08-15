package usecase

import (
	"my_blog_backend/internal/repository"
)

// TODO: Реализовать сервис, перенести логику валидации сюда
type ArticleService struct {
	articleRepo repository.ArticleRepository
}

func NewArticleService(a repository.ArticleRepository) *ArticleService {
	return &ArticleService{articleRepo: a}
}
