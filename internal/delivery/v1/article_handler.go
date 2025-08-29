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

// только автор
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

// только автор
func (h *Handler) deleteArticle(c *gin.Context) {

}

// список статей юзера по нику юзера, в бд поиск по айди юзера
func (h *Handler) getArticlesByUsername(c *gin.Context) {

}

// список статей по айди юзера
func (h *Handler) getArticlesByUserId(c *gin.Context) {

}

// одна статья по айди СТАТЬИ
func (h *Handler) getArticleByID(c *gin.Context) {

}

// статьи по категории
func (h *Handler) getArticlesByCategorySlug(c *gin.Context) {

}

// все статьи
func (h *Handler) getAllArticles(c *gin.Context) {

}
