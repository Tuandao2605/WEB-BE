package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"web-be/models"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (name, slug, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowxContext(ctx, query,
		category.Name, category.Slug, category.Description,
	).Scan(&category.ID, &category.CreatedAt)
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	var category models.Category
	query := `SELECT * FROM categories WHERE id = $1`
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category
	query := `SELECT * FROM categories WHERE slug = $1`
	err := r.db.GetContext(ctx, &category, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	query := `SELECT * FROM categories ORDER BY name`
	err := r.db.SelectContext(ctx, &categories, query)
	return categories, err
}

func (r *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories 
		SET name = $1, slug = $2, description = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.Slug, category.Description, category.ID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
