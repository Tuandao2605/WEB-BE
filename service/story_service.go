package service

import (
	"context"
	"errors"
	"log/slog"

	"web-be/dto"
	"web-be/models"
	"web-be/repository"
	"web-be/utils"
)

type StoryService struct {
	storyRepo    *repository.StoryRepository
	categoryRepo *repository.CategoryRepository
	historyRepo  *repository.ReadingHistoryRepository
}

func NewStoryService(
	storyRepo *repository.StoryRepository,
	categoryRepo *repository.CategoryRepository,
	historyRepo *repository.ReadingHistoryRepository,
) *StoryService {
	return &StoryService{
		storyRepo:    storyRepo,
		categoryRepo: categoryRepo,
		historyRepo:  historyRepo,
	}
}

func (s *StoryService) Create(ctx context.Context, authorID int, req *dto.CreateStoryRequest) (*models.Story, error) {
	slug := utils.GenerateSlug(req.Title)

	// Check if slug exists
	existing, _ := s.storyRepo.GetBySlug(ctx, slug)
	if existing != nil {
		slog.Warn("story creation failed: duplicate slug", "slug", slug)
		return nil, errors.New("a story with this title already exists")
	}

	status := "ongoing"
	if req.Status != "" {
		status = req.Status
	}

	story := &models.Story{
		Title:         req.Title,
		Slug:          slug,
		Description:   req.Description,
		CoverImageURL: req.CoverImageURL,
		AuthorID:      &authorID,
		AuthorName:    req.AuthorName,
		Status:        status,
	}

	err := s.storyRepo.Create(ctx, story)
	if err != nil {
		slog.Error("failed to create story", "error", err, "title", req.Title)
		return nil, errors.New("failed to create story")
	}

	// Add categories if provided
	if len(req.CategoryIDs) > 0 {
		err = s.storyRepo.SetCategories(ctx, story.ID, req.CategoryIDs)
		if err != nil {
			return nil, errors.New("failed to set categories")
		}
	}

	// Load categories
	categories, _ := s.storyRepo.GetCategories(ctx, story.ID)
	story.Categories = categories

	slog.Info("story created", "id", story.ID, "slug", story.Slug, "author_id", authorID)
	return story, nil
}

func (s *StoryService) GetBySlug(ctx context.Context, slug string) (*models.Story, error) {
	story, err := s.storyRepo.GetBySlug(ctx, slug)
	if err != nil {
		slog.Error("failed to get story by slug", "error", err, "slug", slug)
		return nil, err
	}
	if story == nil {
		slog.Debug("story not found", "slug", slug)
		return nil, errors.New("story not found")
	}

	// Increment views
	_ = s.storyRepo.IncrementViews(ctx, story.ID)
	slog.Debug("story views incremented", "story_id", story.ID, "slug", slug)

	// Load categories
	categories, _ := s.storyRepo.GetCategories(ctx, story.ID)
	story.Categories = categories

	return story, nil
}

func (s *StoryService) GetAll(ctx context.Context, limit, offset int) ([]models.Story, int64, error) {
	stories, total, err := s.storyRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Load categories for each story
	for i := range stories {
		categories, _ := s.storyRepo.GetCategories(ctx, stories[i].ID)
		stories[i].Categories = categories
	}

	return stories, total, nil
}

func (s *StoryService) Search(ctx context.Context, keyword string, limit, offset int) ([]models.Story, int64, error) {
	stories, total, err := s.storyRepo.Search(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for i := range stories {
		categories, _ := s.storyRepo.GetCategories(ctx, stories[i].ID)
		stories[i].Categories = categories
	}

	return stories, total, nil
}

func (s *StoryService) GetByCategory(ctx context.Context, categorySlug string, limit, offset int) ([]models.Story, int64, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, categorySlug)
	if err != nil {
		return nil, 0, err
	}
	if category == nil {
		return nil, 0, errors.New("category not found")
	}

	stories, total, err := s.storyRepo.GetByCategory(ctx, category.ID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for i := range stories {
		categories, _ := s.storyRepo.GetCategories(ctx, stories[i].ID)
		stories[i].Categories = categories
	}

	return stories, total, nil
}

func (s *StoryService) Update(ctx context.Context, slug string, userID int, userRole string, req *dto.UpdateStoryRequest) (*models.Story, error) {
	story, err := s.storyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if story == nil {
		return nil, errors.New("story not found")
	}

	// Check permission
	if userRole != "admin" && (story.AuthorID == nil || *story.AuthorID != userID) {
		return nil, errors.New("you don't have permission to edit this story")
	}

	if req.Title != nil {
		story.Title = *req.Title
		story.Slug = utils.GenerateSlug(*req.Title)
	}
	if req.Description != nil {
		story.Description = req.Description
	}
	if req.CoverImageURL != nil {
		story.CoverImageURL = req.CoverImageURL
	}
	if req.AuthorName != nil {
		story.AuthorName = req.AuthorName
	}
	if req.Status != nil {
		story.Status = *req.Status
	}

	err = s.storyRepo.Update(ctx, story)
	if err != nil {
		return nil, errors.New("failed to update story")
	}

	if len(req.CategoryIDs) > 0 {
		err = s.storyRepo.SetCategories(ctx, story.ID, req.CategoryIDs)
		if err != nil {
			return nil, errors.New("failed to update categories")
		}
	}

	categories, _ := s.storyRepo.GetCategories(ctx, story.ID)
	story.Categories = categories

	return story, nil
}

func (s *StoryService) Delete(ctx context.Context, slug string, userID int, userRole string) error {
	story, err := s.storyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if story == nil {
		return errors.New("story not found")
	}

	if userRole != "admin" && (story.AuthorID == nil || *story.AuthorID != userID) {
		slog.Warn("unauthorized story deletion attempt", "user_id", userID, "story_id", story.ID)
		return errors.New("you don't have permission to delete this story")
	}

	err = s.storyRepo.Delete(ctx, story.ID)
	if err == nil {
		slog.Info("story deleted", "story_id", story.ID, "slug", slug, "by_user", userID)
	}
	return err
}

func (s *StoryService) Publish(ctx context.Context, id int, publish bool) error {
	story, err := s.storyRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if story == nil {
		return errors.New("story not found")
	}

	return s.storyRepo.Publish(ctx, id, publish)
}

func (s *StoryService) GetByAuthor(ctx context.Context, authorID, limit, offset int) ([]models.Story, int64, error) {
	return s.storyRepo.GetByAuthor(ctx, authorID, limit, offset)
}

// Category methods
func (s *StoryService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.categoryRepo.GetAll(ctx)
}

func (s *StoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*models.Category, error) {
	slug := utils.GenerateSlug(req.Name)

	existing, _ := s.categoryRepo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("category already exists")
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
	}

	err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, errors.New("failed to create category")
	}

	return category, nil
}

// Reading history methods
func (s *StoryService) GetReadingHistory(ctx context.Context, userID, limit, offset int) ([]models.ReadingHistoryWithDetails, int64, error) {
	return s.historyRepo.GetByUser(ctx, userID, limit, offset)
}

func (s *StoryService) UpdateReadingHistory(ctx context.Context, userID, storyID int, chapterID *int) error {
	history := &models.ReadingHistory{
		UserID:        userID,
		StoryID:       storyID,
		LastChapterID: chapterID,
	}
	return s.historyRepo.Upsert(ctx, history)
}

func (s *StoryService) DeleteReadingHistory(ctx context.Context, userID, storyID int) error {
	return s.historyRepo.Delete(ctx, userID, storyID)
}
