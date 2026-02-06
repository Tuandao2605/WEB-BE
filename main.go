package main

import (
	"log"
	"log/slog"

	"web-be/config"
	"web-be/db"
	"web-be/handler"
	"web-be/repository"
	"web-be/router"
	"web-be/service"
	"web-be/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize structured logger
	utils.InitLogger(cfg.GinMode)
	slog.Info("config loaded successfully")

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Connect to database
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	slog.Info("connected to database successfully")

	// 🔥 RUN MIGRATIONS
	if err := db.RunMigrations(database, "db/migrations"); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	slog.Info("database migrated successfully")

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	categoryRepo := repository.NewCategoryRepository(database)
	storyRepo := repository.NewStoryRepository(database)
	chapterRepo := repository.NewChapterRepository(database)
	historyRepo := repository.NewReadingHistoryRepository(database)
	bookmarkRepo := repository.NewBookmarkRepository(database)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	storyService := service.NewStoryService(storyRepo, categoryRepo, historyRepo)
	chapterService := service.NewChapterService(chapterRepo, storyRepo)
	bookmarkService := service.NewBookmarkService(bookmarkRepo, storyRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	storyHandler := handler.NewStoryHandler(storyService)
	chapterHandler := handler.NewChapterHandler(chapterService)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkService)

	// Setup router
	r := router.NewRouter(jwtManager, authHandler, storyHandler, chapterHandler, bookmarkHandler)
	engine := r.Setup()

	// Start server
	slog.Info("server starting", "port", cfg.Port)
	slog.Info("Story Reader API ready", "url", "http://localhost:"+cfg.Port+"/api/v1")

	if err := engine.Run(":" + cfg.Port); err != nil {
		slog.Error("failed to start server", "error", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}
