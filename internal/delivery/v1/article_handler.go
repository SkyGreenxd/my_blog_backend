package v1

import (
	"log"
	"my_blog_backend/internal/delivery"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createArticle(c *gin.Context) {
	strUserId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var req delivery.CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	res, err := h.services.ArticleService.Create(c.Request.Context(), delivery.ToCreateArticleReq(&req, strUserId.(uint)))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusCreated, delivery.ToCreateArticleRes(res))
}

func (h *Handler) updateArticle(c *gin.Context) {
	strUserId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	strArticleId := c.Param("id")
	articleId, err := strconv.Atoi(strArticleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	var req delivery.UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	res, err := h.services.ArticleService.Update(c.Request.Context(), delivery.ToUpdateArticleReq(&req, strUserId.(uint), uint(articleId)))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToUpdateArticleRes(res))
}

func (h *Handler) deleteArticle(c *gin.Context) {
	strUserId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	strArticleId := c.Param("id")
	articleId, err := strconv.Atoi(strArticleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err := h.services.ArticleService.Delete(c.Request.Context(), delivery.ToDeleteArticleReq(strUserId.(uint), uint(articleId))); err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

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

func (h *Handler) getArticlesByUserId(c *gin.Context) {
	strUserId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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

func (h *Handler) getArticleByID(c *gin.Context) {
	strArticleId := c.Param("id")
	articleId, err := strconv.Atoi(strArticleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	article, err := h.services.ArticleService.GetById(c.Request.Context(), uint(articleId))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, delivery.ToArticleRes(article))
}

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
