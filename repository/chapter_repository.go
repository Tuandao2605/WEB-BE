package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"web-be/models"
)

type ChapterRepository struct {
	db *sqlx.DB
}

func NewChapterRepository(db *sqlx.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

func (r *ChapterRepository) Create(ctx context.Context, chapter *models.Chapter) error {
	query := `
		INSERT INTO chapters (story_id, chapter_number, title, slug, content, word_count, is_published)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query,
		chapter.StoryID, chapter.ChapterNumber, chapter.Title, chapter.Slug,
		chapter.Content, chapter.WordCount, chapter.IsPublished,
	).Scan(&chapter.ID, &chapter.CreatedAt, &chapter.UpdatedAt)
}

func (r *ChapterRepository) GetByID(ctx context.Context, id int) (*models.Chapter, error) {
	var chapter models.Chapter
	query := `SELECT * FROM chapters WHERE id = $1`
	err := r.db.GetContext(ctx, &chapter, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

func (r *ChapterRepository) GetByStoryAndNumber(ctx context.Context, storyID, chapterNumber int) (*models.Chapter, error) {
	var chapter models.Chapter
	query := `SELECT * FROM chapters WHERE story_id = $1 AND chapter_number = $2`
	err := r.db.GetContext(ctx, &chapter, query, storyID, chapterNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

func (r *ChapterRepository) GetByStoryAndSlug(ctx context.Context, storyID int, slug string) (*models.Chapter, error) {
	var chapter models.Chapter
	query := `SELECT * FROM chapters WHERE story_id = $1 AND slug = $2`
	err := r.db.GetContext(ctx, &chapter, query, storyID, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

func (r *ChapterRepository) GetListByStory(ctx context.Context, storyID int, limit, offset int) ([]models.ChapterListItem, int64, error) {
	var chapters []models.ChapterListItem
	var total int64

	countQuery := `SELECT COUNT(*) FROM chapters WHERE story_id = $1 AND is_published = true`
	err := r.db.GetContext(ctx, &total, countQuery, storyID)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, story_id, chapter_number, title, slug, word_count, views, is_published, published_at, created_at
		FROM chapters 
		WHERE story_id = $1 AND is_published = true
		ORDER BY chapter_number ASC
		LIMIT $2 OFFSET $3
	`
	err = r.db.SelectContext(ctx, &chapters, query, storyID, limit, offset)
	return chapters, total, err
}

func (r *ChapterRepository) GetAllByStory(ctx context.Context, storyID int) ([]models.ChapterListItem, error) {
	var chapters []models.ChapterListItem
	query := `
		SELECT id, story_id, chapter_number, title, slug, word_count, views, is_published, published_at, created_at
		FROM chapters 
		WHERE story_id = $1 AND is_published = true
		ORDER BY chapter_number ASC
	`
	err := r.db.SelectContext(ctx, &chapters, query, storyID)
	return chapters, err
}

func (r *ChapterRepository) Update(ctx context.Context, chapter *models.Chapter) error {
	query := `
		UPDATE chapters 
		SET title = $1, slug = $2, content = $3, word_count = $4, 
		    is_published = $5, published_at = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		chapter.Title, chapter.Slug, chapter.Content, chapter.WordCount,
		chapter.IsPublished, chapter.PublishedAt, chapter.ID,
	)
	return err
}

func (r *ChapterRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM chapters WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ChapterRepository) IncrementViews(ctx context.Context, id int) error {
	query := `UPDATE chapters SET views = views + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ChapterRepository) GetNextChapter(ctx context.Context, storyID, currentNumber int) (*models.ChapterListItem, error) {
	var chapter models.ChapterListItem
	query := `
		SELECT id, story_id, chapter_number, title, slug, word_count, views, is_published, published_at, created_at
		FROM chapters 
		WHERE story_id = $1 AND chapter_number > $2 AND is_published = true
		ORDER BY chapter_number ASC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &chapter, query, storyID, currentNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

func (r *ChapterRepository) GetPrevChapter(ctx context.Context, storyID, currentNumber int) (*models.ChapterListItem, error) {
	var chapter models.ChapterListItem
	query := `
		SELECT id, story_id, chapter_number, title, slug, word_count, views, is_published, published_at, created_at
		FROM chapters 
		WHERE story_id = $1 AND chapter_number < $2 AND is_published = true
		ORDER BY chapter_number DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &chapter, query, storyID, currentNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &chapter, nil
}

func (r *ChapterRepository) GetMaxChapterNumber(ctx context.Context, storyID int) (int, error) {
	var maxNum sql.NullInt64
	query := `SELECT MAX(chapter_number) FROM chapters WHERE story_id = $1`
	err := r.db.GetContext(ctx, &maxNum, query, storyID)
	if err != nil {
		return 0, err
	}
	if !maxNum.Valid {
		return 0, nil
	}
	return int(maxNum.Int64), nil
}
