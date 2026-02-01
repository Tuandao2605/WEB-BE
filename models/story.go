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

type StoryCategory struct {
	StoryID    int `db:"story_id"`
	CategoryID int `db:"category_id"`
}
