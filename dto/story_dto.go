package dto

type CreateStoryRequest struct {
	Title         string  `json:"title" binding:"required,min=1,max=255"`
	Description   *string `json:"description"`
	CoverImageURL *string `json:"cover_image_url"`
	AuthorName    *string `json:"author_name"`
	Status        string  `json:"status" binding:"omitempty,oneof=ongoing completed dropped"`
	CategoryIDs   []int   `json:"category_ids"`
}

type UpdateStoryRequest struct {
	Title         *string `json:"title" binding:"omitempty,min=1,max=255"`
	Description   *string `json:"description"`
	CoverImageURL *string `json:"cover_image_url"`
	AuthorName    *string `json:"author_name"`
	Status        *string `json:"status" binding:"omitempty,oneof=ongoing completed dropped"`
	CategoryIDs   []int   `json:"category_ids"`
}

type StoryResponse struct {
	ID            int                `json:"id"`
	Title         string             `json:"title"`
	Slug          string             `json:"slug"`
	Description   *string            `json:"description,omitempty"`
	CoverImageURL *string            `json:"cover_image_url,omitempty"`
	AuthorID      *int               `json:"author_id,omitempty"`
	AuthorName    *string            `json:"author_name,omitempty"`
	Status        string             `json:"status"`
	TotalChapters int                `json:"total_chapters"`
	TotalViews    int64              `json:"total_views"`
	Rating        float64            `json:"rating"`
	IsPublished   bool               `json:"is_published"`
	Categories    []CategoryResponse `json:"categories,omitempty"`
	CreatedAt     string             `json:"created_at"`
	UpdatedAt     string             `json:"updated_at"`
}

type CategoryResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description *string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description"`
}
