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

func Run() {
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

	if err := pgDatabase.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	articleRepo := postgres.NewArticleRepository(pgDatabase.Db)
	categoryRepo := postgres.NewCategoryRepository(pgDatabase.Db)
	sessionRepo := postgres.NewSessionRepository(pgDatabase.Db)
	userRepo := postgres.NewUserRepository(pgDatabase.Db)

	jwtTTLStr := os.Getenv("TOKEN_TTL")
	if jwtTTLStr == "" {
		log.Fatal("TOKEN_TTL must be set")
	}

	jwtTTL, err := time.ParseDuration(jwtTTLStr)
	if err != nil {
		log.Fatalf("failed to parse TOKEN_TTL: %v", err)
	}

	tokenManager := token.NewTokenManager(secret, jwtTTL)
	hashManager, err := hash.NewBcryptHashManager(bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to create hash manager: %v", err)
	}

	articleService := usecase.NewArticleService(articleRepo, userRepo, categoryRepo)
	categoryService := usecase.NewCategoryService(categoryRepo)
	userService := usecase.NewUserService(userRepo, articleRepo, sessionRepo, tokenManager, hashManager)
	services := usecase.NewServices(userService, articleService, categoryService)

	middleware := v1.NewMiddleware(tokenManager)
	handler := v1.NewHandler(services, middleware)

	r := gin.Default()
	handler.Init(r)

	serverCfg := config.LoadHttpServerConfig()
	srv := server.NewServer(r, serverCfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Printf("starting server on port %s", serverCfg.Port)
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped gracefully")
}
