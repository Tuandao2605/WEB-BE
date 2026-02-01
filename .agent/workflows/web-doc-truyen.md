---
description: Quy trình xây dựng Web Đọc Truyện với Gin, PostgreSQL, JWT, SQLX
---

# 📚 Quy Trình Xây Dựng Web Đọc Truyện

## 📁 Cấu Trúc Thư Mục

```
web-be/
├── main.go                     # Entry point
├── go.mod
├── go.sum
├── .env                        # Environment variables
├── .env.example                # Example environment file
│
├── config/                     # Configuration
│   └── config.go               # Load config từ .env
│
├── db/                         # Database
│   ├── migrations/             # SQL migrations
│   │   ├── 001_create_users_table.up.sql
│   │   ├── 001_create_users_table.down.sql
│   │   ├── 002_create_stories_table.up.sql
│   │   ├── 002_create_stories_table.down.sql
│   │   ├── 003_create_chapters_table.up.sql
│   │   ├── 003_create_chapters_table.down.sql
│   │   ├── 004_create_categories_table.up.sql
│   │   ├── 004_create_categories_table.down.sql
│   │   ├── 005_create_reading_history_table.up.sql
│   │   └── 005_create_reading_history_table.down.sql
│   └── postgres.go             # Database connection
│
├── models/                     # Data models
│   ├── user.go
│   ├── story.go
│   ├── chapter.go
│   ├── category.go
│   └── reading_history.go
│
├── repository/                 # Data access layer (SQLX queries)
│   ├── user_repository.go
│   ├── story_repository.go
│   ├── chapter_repository.go
│   ├── category_repository.go
│   └── reading_history_repository.go
│
├── service/                    # Business logic
│   ├── auth_service.go
│   ├── user_service.go
│   ├── story_service.go
│   ├── chapter_service.go
│   └── reading_history_service.go
│
├── handler/                    # HTTP handlers (Controllers)
│   ├── auth_handler.go
│   ├── user_handler.go
│   ├── story_handler.go
│   ├── chapter_handler.go
│   └── reading_history_handler.go
│
├── middleware/                 # Middlewares
│   ├── auth_middleware.go      # JWT authentication
│   ├── cors_middleware.go      # CORS handling
│   └── logger_middleware.go    # Request logging
│
├── router/                     # Route definitions
│   └── router.go
│
├── dto/                        # Data Transfer Objects
│   ├── auth_dto.go             # Login, Register requests/responses
│   ├── story_dto.go
│   ├── chapter_dto.go
│   └── pagination_dto.go
│
├── utils/                      # Utilities
│   ├── jwt.go                  # JWT helper functions
│   ├── password.go             # Password hashing
│   ├── response.go             # Standard API responses
│   └── validator.go            # Custom validators
│
└── docs/                       # API Documentation (Swagger)
    └── swagger.yaml
```

---

## 🗄️ Database Schema (PostgreSQL)

### 1. Bảng Users
```sql
-- 001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    avatar_url VARCHAR(500),
    role VARCHAR(20) DEFAULT 'user', -- 'user', 'admin', 'author'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

### 2. Bảng Categories
```sql
-- 004_create_categories_table.up.sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3. Bảng Stories
```sql
-- 002_create_stories_table.up.sql
CREATE TABLE stories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    cover_image_url VARCHAR(500),
    author_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    author_name VARCHAR(100), -- Tên tác giả gốc (nếu khác user)
    status VARCHAR(20) DEFAULT 'ongoing', -- 'ongoing', 'completed', 'dropped'
    total_chapters INTEGER DEFAULT 0,
    total_views BIGINT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_stories_slug ON stories(slug);
CREATE INDEX idx_stories_author ON stories(author_id);
CREATE INDEX idx_stories_status ON stories(status);

-- Bảng liên kết Story - Category (Many-to-Many)
CREATE TABLE story_categories (
    story_id INTEGER REFERENCES stories(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (story_id, category_id)
);
```

### 4. Bảng Chapters
```sql
-- 003_create_chapters_table.up.sql
CREATE TABLE chapters (
    id SERIAL PRIMARY KEY,
    story_id INTEGER REFERENCES stories(id) ON DELETE CASCADE,
    chapter_number INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    word_count INTEGER DEFAULT 0,
    views BIGINT DEFAULT 0,
    is_published BOOLEAN DEFAULT false,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(story_id, chapter_number)
);

CREATE INDEX idx_chapters_story ON chapters(story_id);
CREATE INDEX idx_chapters_slug ON chapters(story_id, slug);
```

### 5. Bảng Reading History
```sql
-- 005_create_reading_history_table.up.sql
CREATE TABLE reading_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    story_id INTEGER REFERENCES stories(id) ON DELETE CASCADE,
    last_chapter_id INTEGER REFERENCES chapters(id) ON DELETE SET NULL,
    last_read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, story_id)
);

CREATE INDEX idx_reading_history_user ON reading_history(user_id);
```

---

## 📦 Dependencies Cần Cài Đặt

```bash
# Framework & Router
go get github.com/gin-gonic/gin

# Database - PostgreSQL & SQLX
go get github.com/jmoiron/sqlx
go get github.com/lib/pq

# JWT Authentication
go get github.com/golang-jwt/jwt/v5

# Password Hashing
go get golang.org/x/crypto/bcrypt

# Environment Variables
go get github.com/joho/godotenv

# Validation
go get github.com/go-playground/validator/v10

# UUID (optional)
go get github.com/google/uuid

# Database Migrations
go get -u github.com/golang-migrate/migrate/v4/cmd/migrate
```

---

## 🔧 Các Bước Triển Khai Chi Tiết

### Bước 1: Cấu hình Environment

```env
# .env
# Server
PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=story_reader
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRY_HOURS=24
```

### Bước 2: Config Loader (`config/config.go`)

```go
package config

import (
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    Port            string
    GinMode         string
    DBHost          string
    DBPort          string
    DBUser          string
    DBPassword      string
    DBName          string
    DBSSLMode       string
    JWTSecret       string
    JWTExpiryHours  int
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        // .env file not found, use environment variables
    }

    jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

    return &Config{
        Port:           getEnv("PORT", "8080"),
        GinMode:        getEnv("GIN_MODE", "debug"),
        DBHost:         getEnv("DB_HOST", "localhost"),
        DBPort:         getEnv("DB_PORT", "5432"),
        DBUser:         getEnv("DB_USER", "postgres"),
        DBPassword:     getEnv("DB_PASSWORD", ""),
        DBName:         getEnv("DB_NAME", "story_reader"),
        DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
        JWTSecret:      getEnv("JWT_SECRET", "default-secret"),
        JWTExpiryHours: jwtExpiry,
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
```

### Bước 3: Database Connection (`db/postgres.go`)

```go
package db

import (
    "fmt"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "web-be/config"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
    )

    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)

    return db, nil
}
```

### Bước 4: Models (`models/`)

```go
// models/user.go
package models

import "time"

type User struct {
    ID           int       `db:"id" json:"id"`
    Username     string    `db:"username" json:"username"`
    Email        string    `db:"email" json:"email"`
    PasswordHash string    `db:"password_hash" json:"-"` // Never expose
    FullName     *string   `db:"full_name" json:"full_name,omitempty"`
    AvatarURL    *string   `db:"avatar_url" json:"avatar_url,omitempty"`
    Role         string    `db:"role" json:"role"`
    IsActive     bool      `db:"is_active" json:"is_active"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
```

```go
// models/story.go
package models

import "time"

type Story struct {
    ID            int       `db:"id" json:"id"`
    Title         string    `db:"title" json:"title"`
    Slug          string    `db:"slug" json:"slug"`
    Description   *string   `db:"description" json:"description,omitempty"`
    CoverImageURL *string   `db:"cover_image_url" json:"cover_image_url,omitempty"`
    AuthorID      *int      `db:"author_id" json:"author_id,omitempty"`
    AuthorName    *string   `db:"author_name" json:"author_name,omitempty"`
    Status        string    `db:"status" json:"status"`
    TotalChapters int       `db:"total_chapters" json:"total_chapters"`
    TotalViews    int64     `db:"total_views" json:"total_views"`
    Rating        float64   `db:"rating" json:"rating"`
    IsPublished   bool      `db:"is_published" json:"is_published"`
    CreatedAt     time.Time `db:"created_at" json:"created_at"`
    UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`

    // Relations (populated separately)
    Categories []Category `json:"categories,omitempty"`
}
```

```go
// models/chapter.go
package models

import "time"

type Chapter struct {
    ID            int        `db:"id" json:"id"`
    StoryID       int        `db:"story_id" json:"story_id"`
    ChapterNumber int        `db:"chapter_number" json:"chapter_number"`
    Title         string     `db:"title" json:"title"`
    Slug          string     `db:"slug" json:"slug"`
    Content       string     `db:"content" json:"content"`
    WordCount     int        `db:"word_count" json:"word_count"`
    Views         int64      `db:"views" json:"views"`
    IsPublished   bool       `db:"is_published" json:"is_published"`
    PublishedAt   *time.Time `db:"published_at" json:"published_at,omitempty"`
    CreatedAt     time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}
```

### Bước 5: Repository Layer (`repository/`)

```go
// repository/user_repository.go
package repository

import (
    "context"

    "github.com/jmoiron/sqlx"
    "web-be/models"
)

type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (username, email, password_hash, full_name, role)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
    return r.db.QueryRowxContext(ctx, query,
        user.Username, user.Email, user.PasswordHash, user.FullName, user.Role,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    query := `SELECT * FROM users WHERE email = $1 AND is_active = true`
    err := r.db.GetContext(ctx, &user, query, email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
    var user models.User
    query := `SELECT * FROM users WHERE id = $1`
    err := r.db.GetContext(ctx, &user, query, id)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
    var user models.User
    query := `SELECT * FROM users WHERE username = $1 AND is_active = true`
    err := r.db.GetContext(ctx, &user, query, username)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

```go
// repository/story_repository.go
package repository

import (
    "context"

    "github.com/jmoiron/sqlx"
    "web-be/models"
)

type StoryRepository struct {
    db *sqlx.DB
}

func NewStoryRepository(db *sqlx.DB) *StoryRepository {
    return &StoryRepository{db: db}
}

func (r *StoryRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Story, error) {
    var stories []models.Story
    query := `
        SELECT * FROM stories 
        WHERE is_published = true 
        ORDER BY updated_at DESC 
        LIMIT $1 OFFSET $2
    `
    err := r.db.SelectContext(ctx, &stories, query, limit, offset)
    return stories, err
}

func (r *StoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Story, error) {
    var story models.Story
    query := `SELECT * FROM stories WHERE slug = $1`
    err := r.db.GetContext(ctx, &story, query, slug)
    if err != nil {
        return nil, err
    }
    return &story, nil
}

func (r *StoryRepository) Create(ctx context.Context, story *models.Story) error {
    query := `
        INSERT INTO stories (title, slug, description, cover_image_url, author_id, author_name, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `
    return r.db.QueryRowxContext(ctx, query,
        story.Title, story.Slug, story.Description, story.CoverImageURL,
        story.AuthorID, story.AuthorName, story.Status,
    ).Scan(&story.ID, &story.CreatedAt, &story.UpdatedAt)
}

func (r *StoryRepository) IncrementViews(ctx context.Context, id int) error {
    query := `UPDATE stories SET total_views = total_views + 1 WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

func (r *StoryRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]models.Story, error) {
    var stories []models.Story
    query := `
        SELECT * FROM stories 
        WHERE is_published = true 
        AND (title ILIKE $1 OR description ILIKE $1)
        ORDER BY total_views DESC 
        LIMIT $2 OFFSET $3
    `
    searchPattern := "%" + keyword + "%"
    err := r.db.SelectContext(ctx, &stories, query, searchPattern, limit, offset)
    return stories, err
}
```

### Bước 6: JWT Utilities (`utils/jwt.go`)

```go
package utils

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

type JWTManager struct {
    secretKey  []byte
    expiryTime time.Duration
}

func NewJWTManager(secret string, expiryHours int) *JWTManager {
    return &JWTManager{
        secretKey:  []byte(secret),
        expiryTime: time.Duration(expiryHours) * time.Hour,
    }
}

func (j *JWTManager) GenerateToken(userID int, username, role string) (string, error) {
    claims := JWTClaims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiryTime)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*JWTClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}
```

### Bước 7: Password Utilities (`utils/password.go`)

```go
package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Bước 8: Auth Middleware (`middleware/auth_middleware.go`)

```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "web-be/utils"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Expected format: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        claims, err := jwtManager.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Set user info in context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        c.Next()
    }
}

// Optional: Role-based authorization
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }

        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
        c.Abort()
    }
}
```

### Bước 9: DTOs (`dto/`)

```go
// dto/auth_dto.go
package dto

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  struct {
        ID       int    `json:"id"`
        Username string `json:"username"`
        Email    string `json:"email"`
        Role     string `json:"role"`
    } `json:"user"`
}
```

```go
// dto/pagination_dto.go
package dto

type PaginationRequest struct {
    Page     int `form:"page" binding:"min=1"`
    PageSize int `form:"page_size" binding:"min=1,max=100"`
}

func (p *PaginationRequest) GetOffset() int {
    if p.Page < 1 {
        p.Page = 1
    }
    if p.PageSize < 1 {
        p.PageSize = 20
    }
    return (p.Page - 1) * p.PageSize
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalItems int64       `json:"total_items"`
    TotalPages int         `json:"total_pages"`
}
```

### Bước 10: Auth Handler (`handler/auth_handler.go`)

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "web-be/dto"
    "web-be/service"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.authService.Register(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.authService.Login(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
    userID, _ := c.Get("user_id")
    profile, err := h.authService.GetProfile(c.Request.Context(), userID.(int))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, profile)
}
```

### Bước 11: Router Setup (`router/router.go`)

```go
package router

import (
    "github.com/gin-gonic/gin"
    "web-be/handler"
    "web-be/middleware"
    "web-be/utils"
)

type Router struct {
    engine        *gin.Engine
    jwtManager    *utils.JWTManager
    authHandler   *handler.AuthHandler
    storyHandler  *handler.StoryHandler
    chapterHandler *handler.ChapterHandler
}

func NewRouter(
    jwtManager *utils.JWTManager,
    authHandler *handler.AuthHandler,
    storyHandler *handler.StoryHandler,
    chapterHandler *handler.ChapterHandler,
) *Router {
    return &Router{
        engine:        gin.Default(),
        jwtManager:    jwtManager,
        authHandler:   authHandler,
        storyHandler:  storyHandler,
        chapterHandler: chapterHandler,
    }
}

func (r *Router) Setup() *gin.Engine {
    // CORS Middleware
    r.engine.Use(middleware.CORSMiddleware())

    // API v1
    api := r.engine.Group("/api/v1")
    {
        // Health check
        api.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })

        // Auth routes (public)
        auth := api.Group("/auth")
        {
            auth.POST("/register", r.authHandler.Register)
            auth.POST("/login", r.authHandler.Login)
        }

        // Protected routes
        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(r.jwtManager))
        {
            // User profile
            protected.GET("/me", r.authHandler.GetProfile)

            // Reading history
            protected.GET("/history", r.storyHandler.GetReadingHistory)
            protected.POST("/history/:story_id", r.storyHandler.UpdateReadingHistory)
        }

        // Story routes (public read, protected write)
        stories := api.Group("/stories")
        {
            stories.GET("", r.storyHandler.GetAll)
            stories.GET("/search", r.storyHandler.Search)
            stories.GET("/:slug", r.storyHandler.GetBySlug)
            stories.GET("/:slug/chapters", r.chapterHandler.GetByStory)
            stories.GET("/:slug/chapters/:chapter_num", r.chapterHandler.GetChapter)

            // Protected story management
            storiesAuth := stories.Group("")
            storiesAuth.Use(middleware.AuthMiddleware(r.jwtManager))
            {
                storiesAuth.POST("", r.storyHandler.Create)
                storiesAuth.PUT("/:slug", r.storyHandler.Update)
                storiesAuth.DELETE("/:slug", r.storyHandler.Delete)
            }
        }

        // Category routes
        categories := api.Group("/categories")
        {
            categories.GET("", r.storyHandler.GetAllCategories)
            categories.GET("/:slug/stories", r.storyHandler.GetByCategory)
        }

        // Admin routes
        admin := api.Group("/admin")
        admin.Use(middleware.AuthMiddleware(r.jwtManager))
        admin.Use(middleware.RequireRole("admin"))
        {
            admin.GET("/users", r.authHandler.GetAllUsers)
            admin.PUT("/stories/:id/publish", r.storyHandler.Publish)
        }
    }

    return r.engine
}
```

### Bước 12: Main Entry Point (`main.go`)

```go
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

    // Initialize JWT Manager
    jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)

    // Initialize repositories
    userRepo := repository.NewUserRepository(database)
    storyRepo := repository.NewStoryRepository(database)
    chapterRepo := repository.NewChapterRepository(database)
    // ... other repositories

    // Initialize services
    authService := service.NewAuthService(userRepo, jwtManager)
    storyService := service.NewStoryService(storyRepo)
    chapterService := service.NewChapterService(chapterRepo)

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)
    storyHandler := handler.NewStoryHandler(storyService)
    chapterHandler := handler.NewChapterHandler(chapterService)

    // Setup router
    r := router.NewRouter(jwtManager, authHandler, storyHandler, chapterHandler)
    engine := r.Setup()

    // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    if err := engine.Run(":" + cfg.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

---

## 🔌 API Endpoints Summary

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | Đăng ký user mới | ❌ |
| POST | `/api/v1/auth/login` | Đăng nhập | ❌ |
| GET | `/api/v1/me` | Lấy profile user | ✅ |
| GET | `/api/v1/stories` | Danh sách truyện | ❌ |
| GET | `/api/v1/stories/search?q=keyword` | Tìm kiếm truyện | ❌ |
| GET | `/api/v1/stories/:slug` | Chi tiết truyện | ❌ |
| GET | `/api/v1/stories/:slug/chapters` | Danh sách chương | ❌ |
| GET | `/api/v1/stories/:slug/chapters/:num` | Đọc chương | ❌ |
| POST | `/api/v1/stories` | Tạo truyện mới | ✅ |
| PUT | `/api/v1/stories/:slug` | Cập nhật truyện | ✅ |
| DELETE | `/api/v1/stories/:slug` | Xóa truyện | ✅ |
| GET | `/api/v1/categories` | Danh sách thể loại | ❌ |
| GET | `/api/v1/categories/:slug/stories` | Truyện theo thể loại | ❌ |
| GET | `/api/v1/history` | Lịch sử đọc | ✅ |
| POST | `/api/v1/history/:story_id` | Cập nhật lịch sử đọc | ✅ |

---

## 🚀 Lệnh Chạy Dự Án

```bash
# 1. Tạo database PostgreSQL
createdb story_reader

# 2. Chạy migrations
migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/story_reader?sslmode=disable" up

# 3. Cài đặt dependencies
go mod tidy

# 4. Chạy server
go run main.go

# 5. Chạy với hot reload (optional - cần cài air)
air
```

---

## ✅ Checklist Triển Khai

- [ ] Tạo cấu trúc thư mục
- [ ] Cấu hình .env
- [ ] Kết nối PostgreSQL với SQLX
- [ ] Tạo và chạy migrations
- [ ] Implement models
- [ ] Implement repositories
- [ ] Implement JWT utilities
- [ ] Implement password utilities
- [ ] Implement auth middleware
- [ ] Implement services
- [ ] Implement handlers
- [ ] Setup router
- [ ] Test API endpoints
- [ ] Thêm CORS middleware
- [ ] Thêm logging middleware
- [ ] Viết API documentation (Swagger)
