package dto

type BookmarkRequest struct {
	StoryID int `json:"story_id" binding:"required"`
}

type BookmarkStatusResponse struct {
	IsBookmarked   bool  `json:"is_bookmarked"`
	TotalBookmarks int64 `json:"total_bookmarks"`
}

type ViewStatsResponse struct {
	StoryID       int    `json:"story_id"`
	StoryTitle    string `json:"story_title"`
	StorySlug     string `json:"story_slug"`
	TotalViews    int64  `json:"total_views"`
	TotalChapters int    `json:"total_chapters"`
}
