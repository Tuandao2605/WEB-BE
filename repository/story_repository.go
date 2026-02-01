package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"web-be/models"
)

type StoryRepository struct {
	db *sqlx.DB
}

func NewStoryRepository(db *sqlx.DB) *StoryRepository {
	return &StoryRepository{db: db}
}

func (r *StoryRepository) Create(ctx context.Context, story *models.Story) error {
	query := `
		INSERT INTO stories (title, slug, description, cover_image_url, author_id, author_name, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query,
		story.Title, story.Slug, story.Description, story.CoverImageURL,
		story.AuthorID, story.AuthorName, story.Status,
	).Scan(&story.ID, &story.CreatedAt, &story.UpdatedAt)
}

func (r *StoryRepository) GetByID(ctx context.Context, id int) (*models.Story, error) {
	var story models.Story
	query := `SELECT * FROM stories WHERE id = $1`
	err := r.db.GetContext(ctx, &story, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &story, nil
}

func (r *StoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Story, error) {
	var story models.Story
	query := `SELECT * FROM stories WHERE slug = $1`
	err := r.db.GetContext(ctx, &story, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &story, nil
}

func (r *StoryRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64

	countQuery := `SELECT COUNT(*) FROM stories WHERE is_published = true`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT * FROM stories 
		WHERE is_published = true 
		ORDER BY updated_at DESC 
		LIMIT $1 OFFSET $2
	`
	err = r.db.SelectContext(ctx, &stories, query, limit, offset)
	return stories, total, err
}

func (r *StoryRepository) GetByCategory(ctx context.Context, categoryID, limit, offset int) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64

	countQuery := `
		SELECT COUNT(*) FROM stories s
		INNER JOIN story_categories sc ON s.id = sc.story_id
		WHERE sc.category_id = $1 AND s.is_published = true
	`
	err := r.db.GetContext(ctx, &total, countQuery, categoryID)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT s.* FROM stories s
		INNER JOIN story_categories sc ON s.id = sc.story_id
		WHERE sc.category_id = $1 AND s.is_published = true
		ORDER BY s.updated_at DESC
		LIMIT $2 OFFSET $3
	`
	err = r.db.SelectContext(ctx, &stories, query, categoryID, limit, offset)
	return stories, total, err
}

func (r *StoryRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64
	searchPattern := "%" + keyword + "%"

	countQuery := `
		SELECT COUNT(*) FROM stories 
		WHERE is_published = true 
		AND (title ILIKE $1 OR description ILIKE $1)
	`
	err := r.db.GetContext(ctx, &total, countQuery, searchPattern)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT * FROM stories 
		WHERE is_published = true 
		AND (title ILIKE $1 OR description ILIKE $1)
		ORDER BY total_views DESC 
		LIMIT $2 OFFSET $3
	`
	err = r.db.SelectContext(ctx, &stories, query, searchPattern, limit, offset)
	return stories, total, err
}

func (r *StoryRepository) Update(ctx context.Context, story *models.Story) error {
	query := `
		UPDATE stories 
		SET title = $1, slug = $2, description = $3, cover_image_url = $4, 
		    author_name = $5, status = $6, is_published = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		story.Title, story.Slug, story.Description, story.CoverImageURL,
		story.AuthorName, story.Status, story.IsPublished, story.ID,
	)
	return err
}

func (r *StoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM stories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *StoryRepository) IncrementViews(ctx context.Context, id int) error {
	query := `UPDATE stories SET total_views = total_views + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *StoryRepository) UpdateChapterCount(ctx context.Context, storyID int) error {
	query := `
		UPDATE stories 
		SET total_chapters = (SELECT COUNT(*) FROM chapters WHERE story_id = $1 AND is_published = true),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, storyID)
	return err
}

func (r *StoryRepository) Publish(ctx context.Context, id int, publish bool) error {
	query := `UPDATE stories SET is_published = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, publish, id)
	return err
}

// Category management
func (r *StoryRepository) AddCategory(ctx context.Context, storyID, categoryID int) error {
	query := `INSERT INTO story_categories (story_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, storyID, categoryID)
	return err
}

func (r *StoryRepository) RemoveCategory(ctx context.Context, storyID, categoryID int) error {
	query := `DELETE FROM story_categories WHERE story_id = $1 AND category_id = $2`
	_, err := r.db.ExecContext(ctx, query, storyID, categoryID)
	return err
}

func (r *StoryRepository) GetCategories(ctx context.Context, storyID int) ([]models.Category, error) {
	var categories []models.Category
	query := `
		SELECT c.* FROM categories c
		INNER JOIN story_categories sc ON c.id = sc.category_id
		WHERE sc.story_id = $1
		ORDER BY c.name
	`
	err := r.db.SelectContext(ctx, &categories, query, storyID)
	return categories, err
}

func (r *StoryRepository) SetCategories(ctx context.Context, storyID int, categoryIDs []int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing categories
	_, err = tx.ExecContext(ctx, `DELETE FROM story_categories WHERE story_id = $1`, storyID)
	if err != nil {
		return err
	}

	// Add new categories
	for _, catID := range categoryIDs {
		_, err = tx.ExecContext(ctx, `INSERT INTO story_categories (story_id, category_id) VALUES ($1, $2)`, storyID, catID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *StoryRepository) GetByAuthor(ctx context.Context, authorID, limit, offset int) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64

	countQuery := `SELECT COUNT(*) FROM stories WHERE author_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, authorID)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT * FROM stories WHERE author_id = $1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
	err = r.db.SelectContext(ctx, &stories, query, authorID, limit, offset)
	return stories, total, err
}
