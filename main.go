package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"web-be/config"
	"web-be/db"
	"web-be/handler"
	"web-be/repository"
	"web-be/router"
	"web-be/service"
	"web-be/utils"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Connect to database
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Connected to database successfully")

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	categoryRepo := repository.NewCategoryRepository(database)
	storyRepo := repository.NewStoryRepository(database)
	chapterRepo := repository.NewChapterRepository(database)
	historyRepo := repository.NewReadingHistoryRepository(database)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	storyService := service.NewStoryService(storyRepo, categoryRepo, historyRepo)
	chapterService := service.NewChapterService(chapterRepo, storyRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	storyHandler := handler.NewStoryHandler(storyService)
	chapterHandler := handler.NewChapterHandler(chapterService)

	// Setup router
	r := router.NewRouter(jwtManager, authHandler, storyHandler, chapterHandler)
	engine := r.Setup()

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Printf("📚 Story Reader API ready at http://localhost:%s/api/v1", cfg.Port)

	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
