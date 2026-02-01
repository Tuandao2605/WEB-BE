package service

import (
	"context"
	"errors"
	"time"

	"web-be/dto"
	"web-be/models"
	"web-be/repository"
	"web-be/utils"
)

type ChapterService struct {
	chapterRepo *repository.ChapterRepository
	storyRepo   *repository.StoryRepository
}

func NewChapterService(chapterRepo *repository.ChapterRepository, storyRepo *repository.StoryRepository) *ChapterService {
	return &ChapterService{
		chapterRepo: chapterRepo,
		storyRepo:   storyRepo,
	}
}

func (s *ChapterService) Create(ctx context.Context, storySlug string, userID int, userRole string, req *dto.CreateChapterRequest) (*models.Chapter, error) {
	story, err := s.storyRepo.GetBySlug(ctx, storySlug)
	if err != nil {
		return nil, err
	}
	if story == nil {
		return nil, errors.New("story not found")
	}

	// Check permission
	if userRole != "admin" && (story.AuthorID == nil || *story.AuthorID != userID) {
		return nil, errors.New("you don't have permission to add chapters to this story")
	}

	// Get next chapter number
	maxNum, err := s.chapterRepo.GetMaxChapterNumber(ctx, story.ID)
	if err != nil {
		return nil, err
	}
	nextNum := maxNum + 1

	slug := utils.GenerateSlug(req.Title)
	wordCount := utils.CountWords(req.Content)

	var publishedAt *time.Time
	if req.IsPublished {
		now := time.Now()
		publishedAt = &now
	}

	chapter := &models.Chapter{
		StoryID:       story.ID,
		ChapterNumber: nextNum,
		Title:         req.Title,
		Slug:          slug,
		Content:       req.Content,
		WordCount:     wordCount,
		IsPublished:   req.IsPublished,
		PublishedAt:   publishedAt,
	}

	err = s.chapterRepo.Create(ctx, chapter)
	if err != nil {
		return nil, errors.New("failed to create chapter")
	}

	// Update story chapter count
	_ = s.storyRepo.UpdateChapterCount(ctx, story.ID)

	return chapter, nil
}

func (s *ChapterService) GetByStoryAndNumber(ctx context.Context, storySlug string, chapterNum int) (*dto.ChapterResponse, error) {
	story, err := s.storyRepo.GetBySlug(ctx, storySlug)
	if err != nil {
		return nil, err
	}
	if story == nil {
		return nil, errors.New("story not found")
	}

	chapter, err := s.chapterRepo.GetByStoryAndNumber(ctx, story.ID, chapterNum)
	if err != nil {
		return nil, err
	}
	if chapter == nil {
		return nil, errors.New("chapter not found")
	}

	// Increment views
	_ = s.chapterRepo.IncrementViews(ctx, chapter.ID)
	_ = s.storyRepo.IncrementViews(ctx, story.ID)

	// Get navigation
	prevChapter, _ := s.chapterRepo.GetPrevChapter(ctx, story.ID, chapterNum)
	nextChapter, _ := s.chapterRepo.GetNextChapter(ctx, story.ID, chapterNum)

	response := &dto.ChapterResponse{
		ID:            chapter.ID,
		StoryID:       chapter.StoryID,
		ChapterNumber: chapter.ChapterNumber,
		Title:         chapter.Title,
		Slug:          chapter.Slug,
		Content:       chapter.Content,
		WordCount:     chapter.WordCount,
		Views:         chapter.Views,
		IsPublished:   chapter.IsPublished,
		CreatedAt:     chapter.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     chapter.UpdatedAt.Format(time.RFC3339),
	}

	if chapter.PublishedAt != nil {
		t := chapter.PublishedAt.Format(time.RFC3339)
		response.PublishedAt = &t
	}

	if prevChapter != nil {
		response.PrevChapter = &dto.ChapterNavItem{
			ChapterNumber: prevChapter.ChapterNumber,
			Title:         prevChapter.Title,
			Slug:          prevChapter.Slug,
		}
	}

	if nextChapter != nil {
		response.NextChapter = &dto.ChapterNavItem{
			ChapterNumber: nextChapter.ChapterNumber,
			Title:         nextChapter.Title,
			Slug:          nextChapter.Slug,
		}
	}

	return response, nil
}

func (s *ChapterService) GetListByStory(ctx context.Context, storySlug string, limit, offset int) ([]dto.ChapterListResponse, int64, error) {
	story, err := s.storyRepo.GetBySlug(ctx, storySlug)
	if err != nil {
		return nil, 0, err
	}
	if story == nil {
		return nil, 0, errors.New("story not found")
	}

	chapters, total, err := s.chapterRepo.GetListByStory(ctx, story.ID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.ChapterListResponse
	for _, ch := range chapters {
		response := dto.ChapterListResponse{
			ID:            ch.ID,
			StoryID:       ch.StoryID,
			ChapterNumber: ch.ChapterNumber,
			Title:         ch.Title,
			Slug:          ch.Slug,
			WordCount:     ch.WordCount,
			Views:         ch.Views,
			IsPublished:   ch.IsPublished,
			CreatedAt:     ch.CreatedAt.Format(time.RFC3339),
		}
		if ch.PublishedAt != nil {
			t := ch.PublishedAt.Format(time.RFC3339)
			response.PublishedAt = &t
		}
		responses = append(responses, response)
	}

	return responses, total, nil
}

func (s *ChapterService) Update(ctx context.Context, storySlug string, chapterNum int, userID int, userRole string, req *dto.UpdateChapterRequest) (*models.Chapter, error) {
	story, err := s.storyRepo.GetBySlug(ctx, storySlug)
	if err != nil {
		return nil, err
	}
	if story == nil {
		return nil, errors.New("story not found")
	}

	if userRole != "admin" && (story.AuthorID == nil || *story.AuthorID != userID) {
		return nil, errors.New("you don't have permission to edit this chapter")
	}

	chapter, err := s.chapterRepo.GetByStoryAndNumber(ctx, story.ID, chapterNum)
	if err != nil {
		return nil, err
	}
	if chapter == nil {
		return nil, errors.New("chapter not found")
	}

	if req.Title != nil {
		chapter.Title = *req.Title
		chapter.Slug = utils.GenerateSlug(*req.Title)
	}
	if req.Content != nil {
		chapter.Content = *req.Content
		chapter.WordCount = utils.CountWords(*req.Content)
	}
	if req.IsPublished != nil {
		chapter.IsPublished = *req.IsPublished
		if *req.IsPublished && chapter.PublishedAt == nil {
			now := time.Now()
			chapter.PublishedAt = &now
		}
	}

	err = s.chapterRepo.Update(ctx, chapter)
	if err != nil {
		return nil, errors.New("failed to update chapter")
	}

	_ = s.storyRepo.UpdateChapterCount(ctx, story.ID)

	return chapter, nil
}

func (s *ChapterService) Delete(ctx context.Context, storySlug string, chapterNum int, userID int, userRole string) error {
	story, err := s.storyRepo.GetBySlug(ctx, storySlug)
	if err != nil {
		return err
	}
	if story == nil {
		return errors.New("story not found")
	}

	if userRole != "admin" && (story.AuthorID == nil || *story.AuthorID != userID) {
		return errors.New("you don't have permission to delete this chapter")
	}

	chapter, err := s.chapterRepo.GetByStoryAndNumber(ctx, story.ID, chapterNum)
	if err != nil {
		return err
	}
	if chapter == nil {
		return errors.New("chapter not found")
	}

	err = s.chapterRepo.Delete(ctx, chapter.ID)
	if err != nil {
		return err
	}

	_ = s.storyRepo.UpdateChapterCount(ctx, story.ID)

	return nil
}
