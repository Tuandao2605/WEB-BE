package models

import "time"

type ReadingHistory struct {
	ID            int       `db:"id" json:"id"`
	UserID        int       `db:"user_id" json:"user_id"`
	StoryID       int       `db:"story_id" json:"story_id"`
	LastChapterID *int      `db:"last_chapter_id" json:"last_chapter_id,omitempty"`
	LastReadAt    time.Time `db:"last_read_at" json:"last_read_at"`
}

// ReadingHistoryWithDetails - for display with story info
type ReadingHistoryWithDetails struct {
	ReadingHistory
	StoryTitle    string  `db:"story_title" json:"story_title"`
	StorySlug     string  `db:"story_slug" json:"story_slug"`
	CoverImageURL *string `db:"cover_image_url" json:"cover_image_url,omitempty"`
	ChapterNumber *int    `db:"chapter_number" json:"chapter_number,omitempty"`
	ChapterTitle  *string `db:"chapter_title" json:"chapter_title,omitempty"`
}
