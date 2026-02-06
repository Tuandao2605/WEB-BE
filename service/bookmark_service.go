package service

import (
	"context"
	"errors"
	"log/slog"

	"web-be/models"
	"web-be/repository"
)

type BookmarkService struct {
	bookmarkRepo *repository.BookmarkRepository
	storyRepo    *repository.StoryRepository
}

func NewBookmarkService(
	bookmarkRepo *repository.BookmarkRepository,
	storyRepo *repository.StoryRepository,
) *BookmarkService {
	return &BookmarkService{
		bookmarkRepo: bookmarkRepo,
		storyRepo:    storyRepo,
	}
}

func (s *BookmarkService) AddBookmark(ctx context.Context, userID, storyID int) error {
	// Validate story exists
	story, err := s.storyRepo.GetByID(ctx, storyID)
	if err != nil {
		slog.Error("failed to check story existence", "error", err, "story_id", storyID)
		return errors.New("failed to add bookmark")
	}
	if story == nil {
		return errors.New("story not found")
	}

	err = s.bookmarkRepo.Create(ctx, userID, storyID)
	if err != nil {
		slog.Error("failed to create bookmark", "error", err, "user_id", userID, "story_id", storyID)
		return errors.New("failed to add bookmark")
	}

	slog.Info("bookmark added", "user_id", userID, "story_id", storyID)
	return nil
}

func (s *BookmarkService) RemoveBookmark(ctx context.Context, userID, storyID int) error {
	err := s.bookmarkRepo.Delete(ctx, userID, storyID)
	if err != nil {
		slog.Error("failed to remove bookmark", "error", err, "user_id", userID, "story_id", storyID)
		return errors.New("failed to remove bookmark")
	}

	slog.Info("bookmark removed", "user_id", userID, "story_id", storyID)
	return nil
}

func (s *BookmarkService) GetUserBookmarks(ctx context.Context, userID, limit, offset int) ([]models.BookmarkWithStory, int64, error) {
	bookmarks, total, err := s.bookmarkRepo.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		slog.Error("failed to get user bookmarks", "error", err, "user_id", userID)
		return nil, 0, errors.New("failed to get bookmarks")
	}

	slog.Debug("fetched user bookmarks", "user_id", userID, "count", len(bookmarks))
	return bookmarks, total, nil
}

func (s *BookmarkService) GetBookmarkStatus(ctx context.Context, userID, storyID int) (bool, int64, error) {
	isBookmarked, err := s.bookmarkRepo.Exists(ctx, userID, storyID)
	if err != nil {
		slog.Error("failed to check bookmark status", "error", err, "user_id", userID, "story_id", storyID)
		return false, 0, err
	}

	totalBookmarks, err := s.bookmarkRepo.CountByStory(ctx, storyID)
	if err != nil {
		slog.Error("failed to count bookmarks", "error", err, "story_id", storyID)
		return false, 0, err
	}

	return isBookmarked, totalBookmarks, nil
}

func (s *BookmarkService) GetStoryViewStats(ctx context.Context, slug string) (*models.Story, error) {
	story, err := s.storyRepo.GetBySlug(ctx, slug)
	if err != nil {
		slog.Error("failed to get story for view stats", "error", err, "slug", slug)
		return nil, err
	}
	if story == nil {
		return nil, errors.New("story not found")
	}

	slog.Debug("fetched view stats", "slug", slug, "total_views", story.TotalViews)
	return story, nil
}
