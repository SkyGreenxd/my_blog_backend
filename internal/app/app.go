package app

import (
	"context"
	"log"
	"my_blog_backend/internal/config"
	"my_blog_backend/internal/delivery"
	v1 "my_blog_backend/internal/delivery/v1"
	"my_blog_backend/internal/repository/postgres"
	"my_blog_backend/internal/server"
	"my_blog_backend/internal/usecase"
	"my_blog_backend/pkg/auth/hash"
	"my_blog_backend/pkg/auth/token"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtTTL = 15 * time.Minute
)

func Run() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET must be set")
	}

	delivery.RegisterCustomValidators()

	pgDatabase, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := pgDatabase.Close(); err != nil {
			log.Printf("failed to close DB: %v", err)
		}
	}()

	articleRepo := postgres.NewArticleRepository(pgDatabase.Db)
	categoryRepo := postgres.NewCategoryRepository(pgDatabase.Db)
	sessionRepo := postgres.NewSessionRepository(pgDatabase.Db)
	userRepo := postgres.NewUserRepository(pgDatabase.Db)

	tokenManager := token.NewTokenManager(secret, jwtTTL)
	hashManager, err := hash.NewBcryptHashManager(bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	articleService := usecase.NewArticleService(articleRepo, userRepo, categoryRepo)
	categoryService := usecase.NewCategoryService(categoryRepo)
	userService := usecase.NewUserService(userRepo, articleRepo, sessionRepo, tokenManager, hashManager)
	services := usecase.NewServices(userService, articleService, categoryService)

	middleware := v1.NewMiddleware(tokenManager)
	handler := v1.NewHandler(services, middleware)

	r := gin.Default()
	api := r.Group("")
	handler.Init(api)

	serverCfg := config.LoadHttpServerConfig()
	srv := server.NewServer(r, serverCfg)

	// 9. Контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 10. Запуск сервера в горутине
	go func() {
		log.Printf("starting server on port %s", serverCfg.Port)
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// 11. Ожидание сигнала завершения
	<-ctx.Done()
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped gracefully")
}
