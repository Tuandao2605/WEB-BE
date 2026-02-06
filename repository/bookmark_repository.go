package repository

import (
	"context"
	"database/sql"

	"web-be/models"

	"github.com/jmoiron/sqlx"
)

type BookmarkRepository struct {
	db *sqlx.DB
}

func NewBookmarkRepository(db *sqlx.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) Create(ctx context.Context, userID, storyID int) error {
	query := `INSERT INTO bookmarks (user_id, story_id) VALUES ($1, $2) ON CONFLICT (user_id, story_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, storyID)
	return err
}

func (r *BookmarkRepository) Delete(ctx context.Context, userID, storyID int) error {
	query := `DELETE FROM bookmarks WHERE user_id = $1 AND story_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, storyID)
	return err
}

func (r *BookmarkRepository) Exists(ctx context.Context, userID, storyID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND story_id = $2)`
	err := r.db.GetContext(ctx, &exists, query, userID, storyID)
	return exists, err
}

func (r *BookmarkRepository) GetByUser(ctx context.Context, userID, limit, offset int) ([]models.BookmarkWithStory, int64, error) {
	var bookmarks []models.BookmarkWithStory
	var total int64

	countQuery := `SELECT COUNT(*) FROM bookmarks WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT b.id, b.user_id, b.story_id, b.created_at,
		       s.title AS story_title, s.slug AS story_slug, 
		       s.cover_image_url, s.author_name,
		       s.total_chapters, s.total_views
		FROM bookmarks b
		INNER JOIN stories s ON b.story_id = s.id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`
	err = r.db.SelectContext(ctx, &bookmarks, query, userID, limit, offset)
	return bookmarks, total, err
}

func (r *BookmarkRepository) CountByStory(ctx context.Context, storyID int) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM bookmarks WHERE story_id = $1`
	err := r.db.GetContext(ctx, &count, query, storyID)
	return count, err
}

// GetStoryByID - helper to validate story exists
func (r *BookmarkRepository) GetStoryByID(ctx context.Context, storyID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM stories WHERE id = $1)`
	err := r.db.GetContext(ctx, &exists, query, storyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
