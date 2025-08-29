package v1

import (
	"log"
	"my_blog_backend/internal/delivery"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		}
		return
	}

	user, err := h.services.UserService.GetUserById(c.Request.Context(), userId.(uint))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	var req delivery.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	newCategory, err := h.services.CategoryService.Create(c.Request.Context(), delivery.ToCreateCategoryReq(&req, user.Role))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Category": newCategory})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		}
		return
	}

	user, err := h.services.UserService.GetUserById(c.Request.Context(), userId.(uint))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	categorySlug := c.Param("slug")
	if err := h.services.CategoryService.Delete(c.Request.Context(), delivery.ToDeleteCategoryReq(categorySlug, user.Role)); err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": "true"})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		}
		return
	}

	user, err := h.services.UserService.GetUserById(c.Request.Context(), userId.(uint))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	var req delivery.UpdateCategoryReq
	categorySlug := c.Param("slug")
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	category, err := h.services.CategoryService.Update(c.Request.Context(), delivery.ToUpdateCategoryReq(req, user.Role, categorySlug))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Category": category})
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.services.CategoryService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
