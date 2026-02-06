package router

import (
	"web-be/handler"
	"web-be/middleware"
	"web-be/utils"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine          *gin.Engine
	jwtManager      *utils.JWTManager
	authHandler     *handler.AuthHandler
	storyHandler    *handler.StoryHandler
	chapterHandler  *handler.ChapterHandler
	bookmarkHandler *handler.BookmarkHandler
}

func NewRouter(
	jwtManager *utils.JWTManager,
	authHandler *handler.AuthHandler,
	storyHandler *handler.StoryHandler,
	chapterHandler *handler.ChapterHandler,
	bookmarkHandler *handler.BookmarkHandler,

) *Router {
	return &Router{
		engine:          gin.Default(),
		jwtManager:      jwtManager,
		authHandler:     authHandler,
		storyHandler:    storyHandler,
		chapterHandler:  chapterHandler,
		bookmarkHandler: bookmarkHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	// CORS Middleware
	//r.engine.Use(middleware.CORSMiddleware())

	// API v1
	api := r.engine.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "message": "Story Reader API is running"})
		})

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
		}

		// User profile (protected)
		api.GET("/me", middleware.AuthMiddleware(r.jwtManager), r.authHandler.GetProfile)
		api.PUT("/me", middleware.AuthMiddleware(r.jwtManager), r.authHandler.UpdateProfile)

		// My stories (protected)
		api.GET("/my-stories", middleware.AuthMiddleware(r.jwtManager), r.storyHandler.GetMyStories)

		// Reading history (protected)
		history := api.Group("/history")
		history.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			history.GET("", r.storyHandler.GetReadingHistory)
			history.POST("/:story_id", r.storyHandler.UpdateReadingHistory)
			history.DELETE("/:story_id", r.storyHandler.DeleteReadingHistory)
		}

		// Category routes (public)
		categories := api.Group("/categories")
		{
			categories.GET("", r.storyHandler.GetAllCategories)
			categories.GET("/:slug/stories", r.storyHandler.GetByCategory)
		}

		// Story routes
		stories := api.Group("/stories")
		{
			// Public routes
			stories.GET("", r.storyHandler.GetAll)
			stories.GET("/search", r.storyHandler.Search)
			stories.GET("/:slug", r.storyHandler.GetBySlug)
			stories.GET("/:slug/chapters", r.chapterHandler.GetByStory)
			stories.GET("/:slug/chapters/:chapter_num", r.chapterHandler.GetChapter)
			stories.GET("/:slug/stats", r.bookmark_handler.GetViewStats)

			// Protected routes
			storiesAuth := stories.Group("")
			storiesAuth.Use(middleware.AuthMiddleware(r.jwtManager))
			{
				storiesAuth.POST("", r.storyHandler.Create)
				storiesAuth.PUT("/:slug", r.storyHandler.Update)
				storiesAuth.DELETE("/:slug", r.storyHandler.Delete)

				// Chapter management
				storiesAuth.POST("/:slug/chapters", r.chapterHandler.Create)
				storiesAuth.PUT("/:slug/chapters/:chapter_num", r.chapterHandler.Update)
				storiesAuth.DELETE("/:slug/chapters/:chapter_num", r.chapterHandler.Delete)
			}
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(r.jwtManager))
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/users", r.authHandler.GetAllUsers)
			admin.POST("/categories", r.storyHandler.CreateCategory)
			admin.PUT("/stories/:id/publish", r.storyHandler.Publish)
		}

	}

	return r.engine
}
