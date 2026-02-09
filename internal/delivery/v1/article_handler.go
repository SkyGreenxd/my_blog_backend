package v1

import (
	"log"
	"my_blog_backend/internal/delivery"
	"my_blog_backend/pkg/e"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary		Create Article
// @Tags			articles
// @Description	Create a new blog article
// @Accept			json
// @Produce		json
// @Param			input	body		delivery.CreateArticleReq	true	"Article data"
// @Success		201		{object}	delivery.CreateArticleRes
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Router			/articles [post]
func (h *Handler) createArticle(c *gin.Context) {
	strUserId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") != "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
			return
		}
		ErrorToHttpRes(e.ErrInternalServerError, c)
		return
	}

	var req delivery.CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	res, err := h.services.ArticleService.Create(c.Request.Context(), delivery.ToCreateArticleReq(&req, strUserId.(uint)))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusCreated, delivery.ToCreateArticleRes(res))
}

// @Summary		Update Article
// @Tags			articles
// @Description	Update an existing blog article
// @Accept			json
// @Produce		json
// @Param			id		path		int							true	"Article ID"
// @Param			input	body		delivery.UpdateArticleReq	true	"New article data"
// @Success		200		{object}	delivery.UpdateArticleRes
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Failure		403		{object}	delivery.ErrResponse
// @Router			/articles/{id} [patch]
func (h *Handler) updateArticle(c *gin.Context) {
	strUserId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") != "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
			return
		}
		ErrorToHttpRes(e.ErrInternalServerError, c)
		return
	}

	strArticleId := c.Param("id")
	articleId, err := strconv.Atoi(strArticleId)
	if err != nil {
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	var req delivery.UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	res, err := h.services.ArticleService.Update(c.Request.Context(), delivery.ToUpdateArticleReq(&req, strUserId.(uint), uint(articleId)))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToUpdateArticleRes(res))
}

// @Summary		Delete Article
// @Tags			articles
// @Description	Delete a blog article
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Article ID"
// @Success		204	{object}	map[string]interface{}
// @Failure		400	{object}	delivery.ErrResponse
// @Failure		401	{object}	delivery.ErrResponse
// @Failure		403	{object}	delivery.ErrResponse
// @Security		ApiKeyAuth
// @Router			/articles/{id} [delete]
func (h *Handler) deleteArticle(c *gin.Context) {
	strUserId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") != "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
			return
		}
		ErrorToHttpRes(e.ErrInternalServerError, c)
		return
	}

	strArticleId := c.Param("id")
	articleId, err := strconv.Atoi(strArticleId)
	if err != nil {
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	if err := h.services.ArticleService.Delete(c.Request.Context(), delivery.ToDeleteArticleReq(strUserId.(uint), uint(articleId))); err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// @Summary		Get Articles By Username
// @Tags			articles
// @Description	Get all articles written by a specific user
// @Accept			json
// @Produce		json
// @Param			username	path		string	true	"Username"
// @Success		200			{object}	delivery.GetArticlesByUserRes
// @Failure		404			{object}	delivery.ErrResponse
// @Router			/users/{username}/articles [get]
func (h *Handler) getArticlesByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.services.UserService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	dto, err := h.services.ArticleService.GetAllArticlesByUserId(c.Request.Context(), user.Id)
	if err != nil {
		ErrorToHttpRes(err, c)
	}

	articles := make([]*delivery.ArticleRes, len(dto.Articles))
	for i, article := range dto.Articles {
		articles[i] = delivery.ToArticleRes(article)
	}

	res := delivery.ToGetArticlesByUserRes(articles)
	c.JSON(http.StatusOK, res)
}

// @Summary		Get My Articles
// @Tags			articles
// @Description	Get all articles written by the currently authenticated user
// @Accept			json
// @Produce		json
// @Success		200	{object}	delivery.GetArticlesByUserRes
// @Failure		401	{object}	delivery.ErrResponse
// @Router			/users/me/articles [get]
func (h *Handler) getArticlesByUserId(c *gin.Context) {
	strUserId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") != "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
			return
		}
		ErrorToHttpRes(e.ErrInternalServerError, c)
		return
	}

	dto, err := h.services.ArticleService.GetAllArticlesByUserId(c.Request.Context(), strUserId.(uint))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	articles := make([]*delivery.ArticleRes, len(dto.Articles))
	for i, article := range dto.Articles {
		articles[i] = delivery.ToArticleRes(article)
	}

	res := delivery.ToGetArticlesByUserRes(articles)
	c.JSON(http.StatusOK, res)
}

// @Summary		Get Article By ID
// @Tags			articles
// @Description	Get a single article by its ID
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Article ID"
// @Success		200	{object}	delivery.ArticleRes
// @Failure		404	{object}	delivery.ErrResponse
// @Router			/articles/{id} [get]
func (h *Handler) getArticleByID(c *gin.Context) {
	strArticleId := c.Param("id")
	articleId, err := strconv.Atoi(strArticleId)
	if err != nil {
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	article, err := h.services.ArticleService.GetById(c.Request.Context(), uint(articleId))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToArticleRes(article))
}

// @Summary		Get Articles By Category
// @Tags			articles
// @Description	Get all articles in a specific category
// @Accept			json
// @Produce		json
// @Param			slug	path		string	true	"Category Slug"
// @Success		200		{object}	delivery.GetArticlesByUserRes
// @Failure		404		{object}	delivery.ErrResponse
// @Router			/categories/{slug}/articles [get]
func (h *Handler) getArticlesByCategorySlug(c *gin.Context) {
	slug := c.Param("slug")
	dto, err := h.services.ArticleService.GetAllArticlesByCategory(c.Request.Context(), slug)
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	articles := make([]*delivery.ArticleRes, len(dto.Articles))
	for i, article := range dto.Articles {
		articles[i] = delivery.ToArticleRes(article)
	}

	res := delivery.ToGetArticlesByUserRes(articles)
	c.JSON(http.StatusOK, res)
}

// @Summary		Get All Articles
// @Tags			articles
// @Description	Get a list of all blog articles
// @Accept			json
// @Produce		json
// @Success		200	{object}	delivery.GetArticlesByUserRes
// @Router			/articles [get]
func (h *Handler) getAllArticles(c *gin.Context) {
	dto, err := h.services.ArticleService.GetAll(c.Request.Context())
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	articles := make([]*delivery.ArticleRes, len(dto.Articles))
	for i, article := range dto.Articles {
		articles[i] = delivery.ToArticleRes(article)
	}

	res := delivery.ToGetArticlesByUserRes(articles)
	c.JSON(http.StatusOK, res)
}
