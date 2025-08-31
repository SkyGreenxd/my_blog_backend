package v1

import (
	"my_blog_backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   *usecase.Services
	middleware *Middleware
}

func NewHandler(services *usecase.Services, middleware *Middleware) *Handler {
	return &Handler{
		services:   services,
		middleware: middleware,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
			auth.POST("/refresh", h.refreshSession)
			auth.POST("/logout", h.logout)

			auth.Use(h.middleware.AuthMiddleware())
			{
				auth.POST("/password/change", h.changePassword)
			}
		}

		users := v1.Group("/users")
		{
			// users.GET("/:id", h.getUserById)
			users.GET("/:username", h.getUserByUsername)
			users.GET("/:username/articles", h.getArticlesByUsername)

			users.Use(h.middleware.AuthMiddleware())
			{
				users.GET("/me", h.getCurrentUser)
				users.GET("/me/articles", h.getArticlesByUserId)
				users.PATCH("/me/update", h.updateUser)
				// users.PATCH("me/admin", h.setAdminRole)
			}
		}

		categories := v1.Group("/categories")
		{
			categories.GET("", h.GetAllCategories)
			categories.GET("/:slug/articles", h.getArticlesByCategorySlug)

			categories.Use(h.middleware.AuthMiddleware())
			{
				categories.POST("", h.CreateCategory)
				categories.PATCH("/:slug", h.UpdateCategory)
				categories.DELETE("/:slug", h.DeleteCategory)
			}
		}

		articles := v1.Group("/articles")
		{
			articles.GET("/:id", h.getArticleByID)
			articles.GET("", h.getAllArticles)

			articles.Use(h.middleware.AuthMiddleware())
			{
				articles.POST("", h.createArticle)
				articles.PATCH("/:id", h.updateArticle)
				articles.DELETE("/:id", h.deleteArticle)
			}
		}
	}
}
