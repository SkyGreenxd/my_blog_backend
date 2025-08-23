package usecase

import (
	"context"
	"errors"
	"my_blog_backend/internal/domain"
	"my_blog_backend/internal/repository"
	"my_blog_backend/pkg/e"
)

// TODO: Реализовать сервис, перенести логику валидации сюда
type ArticleService struct {
	articleRepo  repository.ArticleRepository
	userRepo     repository.UserRepository
	categoryRepo repository.CategoryRepository
}

func NewArticleService(a repository.ArticleRepository, u repository.UserRepository, c repository.CategoryRepository) *ArticleService {
	return &ArticleService{
		articleRepo:  a,
		userRepo:     u,
		categoryRepo: c,
	}
}

func (s *ArticleService) GetAllArticlesByUserId(ctx context.Context, userId uint) (*GetArticlesByUserRes, error) {
	const op = "ArticleService.GetAllArticlesByUserId"

	articles, err := s.articleRepo.ListByAuthor(ctx, userId)
	if err != nil {
		if errors.Is(err, e.ErrArticleNotFound) {
			return &GetArticlesByUserRes{[]*ArticleRes{}}, nil
		}

		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	res := make([]*ArticleRes, len(articles))
	for i, article := range articles {
		res[i] = toArticleRes(&article)
	}

	return toGetArticlesByUserRes(res), nil
}

func (s *ArticleService) Create(ctx context.Context, req *CreateArticleReq) (*CreateArticleRes, error) {
	const op = "ArticleService.Create"

	category, err := s.categoryRepo.GetByName(ctx, req.CategoryName)
	if err != nil {
		if errors.Is(err, e.ErrCategoryNotFound) {
			return nil, e.Wrap(op, e.ErrCategoryNotFound)
		}

		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	newArticle := domain.NewArticle(req.Title, req.Content, req.UserId, category.ID)

	result, err := s.articleRepo.Create(ctx, newArticle)
	if err != nil {
		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	return toCreateArticleRes(result, category.Slug, category.Name), nil
}

func (s *ArticleService) GetById(ctx context.Context, id uint) (*GetArticleRes, error) {
	const op = "ArticleService.GetById"

	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, e.ErrArticleNotFound) {
			return nil, e.Wrap(op, e.ErrArticleNotFound)
		}

		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	return toGetArticleRes(article), nil
}

func (s *ArticleService) GetAllArticlesByCategory(ctx context.Context, slug string) ([]GetAllArticlesRes, error) {
	const op = "ArticleService.GetAllArticlesByCategoryId"

	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, e.ErrCategoryNotFound) {
			return nil, e.Wrap(op, e.ErrCategoryNotFound)
		}

		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	articles, err := s.articleRepo.ListByCategory(ctx, category.ID)
	if err != nil {
		if errors.Is(err, e.ErrArticleNotFound) {
			return nil, e.Wrap(op, e.ErrArticleNotFound)
		}

		return nil, e.Wrap(op, e.ErrInternalServer)
	}

	res := make([]GetAllArticlesRes, len(articles))
	for i, article := range articles {
		res[i] = toGetAllArticlesRes(article)
	}

	return res, nil
}

func toGetAllArticlesRes(a domain.Article) GetAllArticlesRes {
	return GetAllArticlesRes{
		Title:        a.Title,
		Content:      a.Content,
		AuthorID:     a.AuthorID,
		CategorySlug: a.Category.Slug,
		CategoryName: a.Category.Name,
	}
}

func toArticleRes(article *domain.Article) *ArticleRes {
	return &ArticleRes{
		ArticleId:    article.ID,
		UserId:       article.AuthorID,
		Username:     article.Author.Username,
		Title:        article.Title,
		Content:      article.Content,
		CategoryName: article.Category.Name,
		CategorySlug: article.Category.Slug,
	}
}

func toGetArticlesByUserRes(articles []*ArticleRes) *GetArticlesByUserRes {
	return &GetArticlesByUserRes{
		Articles: articles,
	}
}

func toCreateArticleRes(article *domain.Article, categorySlug, categoryName string) *CreateArticleRes {
	return &CreateArticleRes{
		ArticleId:    article.ID,
		Title:        article.Title,
		Content:      article.Content,
		CategorySlug: categorySlug,
		CategoryName: categoryName,
	}
}

func toGetArticleRes(article *domain.Article) *GetArticleRes {
	return &GetArticleRes{
		Title:        article.Title,
		Content:      article.Content,
		CategorySlug: article.Category.Slug,
		Username:     article.Author.Username,
	}
}
