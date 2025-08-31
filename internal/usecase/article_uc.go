package usecase

import (
	"context"
	"errors"
	"my_blog_backend/internal/domain"
	"my_blog_backend/internal/repository"
	"my_blog_backend/pkg/e"
)

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

func (s *ArticleService) GetAllArticlesByUserId(ctx context.Context, userId uint) (*GetArticles, error) {
	const op = "ArticleService.GetAllArticlesByUserId"

	articles, err := s.articleRepo.ListByAuthor(ctx, userId)
	if err != nil {
		if errors.Is(err, e.ErrArticleNotFound) {
			return &GetArticles{[]*ArticleRes{}}, nil
		}

		return nil, e.Wrap(op, err)
	}

	res := make([]*ArticleRes, len(articles))
	for i, article := range articles {
		res[i] = toArticleRes(&article)
	}

	return toGetArticlesByUserRes(res), nil
}

func (s *ArticleService) Create(ctx context.Context, req *CreateArticleReq) (*CreateArticleRes, error) {
	const op = "ArticleService.Create"

	category, err := s.categoryRepo.GetBySlug(ctx, req.CategorySlug)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	newArticle := domain.NewArticle(req.Title, req.Content, req.UserId, category.ID)
	if err := newArticle.Validate(); err != nil {
		return nil, e.Wrap(op, e.ErrArticleDataIsInvalid)
	}

	if err := s.articleRepo.ExistsByTitleContentAuthor(ctx, newArticle); err != nil {
		return nil, e.Wrap(op, e.ErrArticleDataIsInvalid)
	}

	result, err := s.articleRepo.Create(ctx, newArticle)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toCreateArticleRes(result, category.Slug, category.Name), nil
}

func (s *ArticleService) GetById(ctx context.Context, id uint) (*ArticleRes, error) {
	const op = "ArticleService.GetById"

	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toArticleRes(article), nil
}

func (s *ArticleService) GetAllArticlesByCategory(ctx context.Context, slug string) (*GetArticles, error) {
	const op = "ArticleService.GetAllArticlesByCategoryId"

	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	articles, err := s.articleRepo.ListByCategory(ctx, category.ID)
	if err != nil {
		if errors.Is(err, e.ErrArticleNotFound) {
			return &GetArticles{[]*ArticleRes{}}, nil
		}

		return nil, e.Wrap(op, err)
	}

	res := make([]*ArticleRes, len(articles))
	for i, article := range articles {
		res[i] = toArticleRes(&article)
	}

	return toGetArticlesByUserRes(res), nil
}

func (s *ArticleService) Delete(ctx context.Context, req *DeleteArticleReq) error {
	const op = "ArticleService.Delete"

	article, err := s.articleRepo.GetByID(ctx, req.ArticleId)
	if err != nil {
		return e.Wrap(op, err)
	}

	if err := article.CheckAuthor(req.UserId); err != nil {
		return e.Wrap(op, e.ErrUserNotAuthor)
	}

	if err := s.articleRepo.Delete(ctx, req.ArticleId); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (s *ArticleService) Update(ctx context.Context, req *UpdateArticleReq) (*UpdateArticleRes, error) {
	const op = "ArticleService.Update"

	article, err := s.articleRepo.GetByID(ctx, req.ArticleId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	if err := article.CheckAuthor(req.UserId); err != nil {
		return nil, e.Wrap(op, e.ErrUserNotAuthor)
	}

	if req.Title == nil && req.Content == nil && req.CategorySlug == nil {
		return nil, e.Wrap(op, e.ErrNoDataToUpdate)
	}

	if req.Title != nil {
		if err := article.ChangeTitle(*req.Title); err != nil {
			return nil, e.Wrap(op, e.ErrArticleNameIsExists)
		}
	}

	if req.Content != nil {
		if err := article.ChangeContent(*req.Content); err != nil {
			return nil, e.Wrap(op, e.ErrArticleContentIsExists)
		}
	}

	if req.CategorySlug != nil {
		category, err := s.categoryRepo.GetBySlug(ctx, *req.CategorySlug)
		if err != nil {
			return nil, e.Wrap(op, err)
		}

		if err := article.ChangeCategory(category); err != nil {
			return nil, e.Wrap(op, e.ErrArticleCategoryIsExists)
		}
	}

	updArticle, err := s.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return toUpdateArticleRes(updArticle), nil
}

func (s *ArticleService) GetAll(ctx context.Context) (*GetArticles, error) {
	const op = "ArticleService.GetAll"

	articles, err := s.articleRepo.ListAll(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	res := make([]*ArticleRes, len(articles))
	for i, article := range articles {
		res[i] = toArticleRes(&article)
	}

	return toGetArticlesByUserRes(res), nil
}

func toCategoryRes(category *domain.Category) *CategoryRes {
	return &CategoryRes{
		CategoryId:   category.ID,
		CategoryName: category.Name,
		CategorySlug: category.Slug,
	}
}

func toUpdateArticleRes(a *domain.Article) *UpdateArticleRes {
	return &UpdateArticleRes{
		ArticleId: a.ID,
		Title:     a.Title,
		Content:   a.Content,
		Category:  *toCategoryRes(a.Category),
		AuthorID:  a.Author.ID,
		UpdatedAt: a.UpdatedAt,
	}
}

func toArticleRes(article *domain.Article) *ArticleRes {
	return &ArticleRes{
		ArticleId: article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Author:    *toUserResponse(article.Author),
		Category:  *toCategoryRes(article.Category),
	}
}

func toGetArticlesByUserRes(articles []*ArticleRes) *GetArticles {
	return &GetArticles{
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
