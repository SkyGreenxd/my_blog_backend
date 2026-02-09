package v1

import (
	"log"
	"my_blog_backend/internal/delivery"
	"my_blog_backend/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Create Category
// @Tags			categories
// @Description	Create a new category (Admin only)
// @Accept			json
// @Produce		json
// @Param			input	body		delivery.CreateCategoryReq	true	"Category data"
// @Success		200		{object}	map[string]interface{}
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Failure		403		{object}	delivery.ErrResponse
// @Router			/categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
		} else {
			ErrorToHttpRes(e.ErrInternalServer, c)
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
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	newCategory, err := h.services.CategoryService.Create(c.Request.Context(), delivery.ToCreateCategoryReq(&req, user.Role))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Category": newCategory})
}

// @Summary		Delete Category
// @Tags			categories
// @Description	Delete a category by its slug (Admin only)
// @Accept			json
// @Produce		json
// @Param			slug	path		string	true	"Category Slug"
// @Success		200		{object}	map[string]interface{}
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Failure		403		{object}	delivery.ErrResponse
// @Router			/categories/{slug} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
		} else {
			ErrorToHttpRes(e.ErrInternalServer, c)
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

// @Summary		Update Category
// @Tags			categories
// @Description	Update an existing category (Admin only)
// @Accept			json
// @Produce		json
// @Param			slug	path		string						true	"Category Slug"
// @Param			input	body		delivery.UpdateCategoryReq	true	"New category data"
// @Success		200		{object}	map[string]interface{}
// @Failure		400		{object}	delivery.ErrResponse
// @Failure		401		{object}	delivery.ErrResponse
// @Failure		403		{object}	delivery.ErrResponse
// @Router			/categories/{slug} [patch]
func (h *Handler) UpdateCategory(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		if c.GetHeader("Authorization") == "" {
			ErrorToHttpRes(e.ErrUnauthorized, c)
		} else {
			ErrorToHttpRes(e.ErrInternalServer, c)
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
		ErrorToHttpRes(e.ErrBadRequest, c)
		return
	}

	category, err := h.services.CategoryService.Update(c.Request.Context(), delivery.ToUpdateCategoryReq(req, user.Role, categorySlug))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Category": category})
}

// @Summary		Get All Categories
// @Tags			categories
// @Description	Get a list of all categories
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]interface{}
// @Router			/categories [get]
func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.services.CategoryService.GetAll(c.Request.Context())
	if err != nil {
		ErrorToHttpRes(e.ErrInternalServer, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
