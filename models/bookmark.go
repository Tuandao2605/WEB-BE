package models

import "time"

type Bookmark struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	StoryID   int       `db:"story_id" json:"story_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// BookmarkWithStory - for display with story info
type BookmarkWithStory struct {
	Bookmark
	StoryTitle    string  `db:"story_title" json:"story_title"`
	StorySlug     string  `db:"story_slug" json:"story_slug"`
	CoverImageURL *string `db:"cover_image_url" json:"cover_image_url,omitempty"`
	AuthorName    *string `db:"author_name" json:"author_name,omitempty"`
	TotalChapters int     `db:"total_chapters" json:"total_chapters"`
	TotalViews    int64   `db:"total_views" json:"total_views"`
}
