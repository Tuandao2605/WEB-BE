package dto

type CreateChapterRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Content     string `json:"content" binding:"required"`
	IsPublished bool   `json:"is_published"`
}

type UpdateChapterRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=255"`
	Content     *string `json:"content"`
	IsPublished *bool   `json:"is_published"`
}

type ChapterResponse struct {
	ID            int     `json:"id"`
	StoryID       int     `json:"story_id"`
	ChapterNumber int     `json:"chapter_number"`
	Title         string  `json:"title"`
	Slug          string  `json:"slug"`
	Content       string  `json:"content"`
	WordCount     int     `json:"word_count"`
	Views         int64   `json:"views"`
	IsPublished   bool    `json:"is_published"`
	PublishedAt   *string `json:"published_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`

	// Navigation
	PrevChapter *ChapterNavItem `json:"prev_chapter,omitempty"`
	NextChapter *ChapterNavItem `json:"next_chapter,omitempty"`
}

type ChapterListResponse struct {
	ID            int     `json:"id"`
	StoryID       int     `json:"story_id"`
	ChapterNumber int     `json:"chapter_number"`
	Title         string  `json:"title"`
	Slug          string  `json:"slug"`
	WordCount     int     `json:"word_count"`
	Views         int64   `json:"views"`
	IsPublished   bool    `json:"is_published"`
	PublishedAt   *string `json:"published_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

type ChapterNavItem struct {
	ChapterNumber int    `json:"chapter_number"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
}
