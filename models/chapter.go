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

// ChapterListItem - lighter version for listing
type ChapterListItem struct {
	ID            int        `db:"id" json:"id"`
	StoryID       int        `db:"story_id" json:"story_id"`
	ChapterNumber int        `db:"chapter_number" json:"chapter_number"`
	Title         string     `db:"title" json:"title"`
	Slug          string     `db:"slug" json:"slug"`
	WordCount     int        `db:"word_count" json:"word_count"`
	Views         int64      `db:"views" json:"views"`
	IsPublished   bool       `db:"is_published" json:"is_published"`
	PublishedAt   *time.Time `db:"published_at" json:"published_at,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}
