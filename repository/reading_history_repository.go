package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"web-be/models"
)

type ReadingHistoryRepository struct {
	db *sqlx.DB
}

func NewReadingHistoryRepository(db *sqlx.DB) *ReadingHistoryRepository {
	return &ReadingHistoryRepository{db: db}
}

func (r *ReadingHistoryRepository) Upsert(ctx context.Context, history *models.ReadingHistory) error {
	query := `
		INSERT INTO reading_history (user_id, story_id, last_chapter_id, last_read_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, story_id) 
		DO UPDATE SET last_chapter_id = $3, last_read_at = CURRENT_TIMESTAMP
		RETURNING id, last_read_at
	`
	return r.db.QueryRowxContext(ctx, query,
		history.UserID, history.StoryID, history.LastChapterID,
	).Scan(&history.ID, &history.LastReadAt)
}

func (r *ReadingHistoryRepository) GetByUserAndStory(ctx context.Context, userID, storyID int) (*models.ReadingHistory, error) {
	var history models.ReadingHistory
	query := `SELECT * FROM reading_history WHERE user_id = $1 AND story_id = $2`
	err := r.db.GetContext(ctx, &history, query, userID, storyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &history, nil
}

func (r *ReadingHistoryRepository) GetByUser(ctx context.Context, userID, limit, offset int) ([]models.ReadingHistoryWithDetails, int64, error) {
	var histories []models.ReadingHistoryWithDetails
	var total int64

	countQuery := `SELECT COUNT(*) FROM reading_history WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			rh.id, rh.user_id, rh.story_id, rh.last_chapter_id, rh.last_read_at,
			s.title as story_title, s.slug as story_slug, s.cover_image_url,
			c.chapter_number, c.title as chapter_title
		FROM reading_history rh
		INNER JOIN stories s ON rh.story_id = s.id
		LEFT JOIN chapters c ON rh.last_chapter_id = c.id
		WHERE rh.user_id = $1
		ORDER BY rh.last_read_at DESC
		LIMIT $2 OFFSET $3
	`
	err = r.db.SelectContext(ctx, &histories, query, userID, limit, offset)
	return histories, total, err
}

func (r *ReadingHistoryRepository) Delete(ctx context.Context, userID, storyID int) error {
	query := `DELETE FROM reading_history WHERE user_id = $1 AND story_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, storyID)
	return err
}

func (r *ReadingHistoryRepository) DeleteAllByUser(ctx context.Context, userID int) error {
	query := `DELETE FROM reading_history WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
