package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"web-be/dto"
	"web-be/middleware"
	"web-be/service"
	"web-be/utils"

	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	bookmarkService *service.BookmarkService
}

func NewBookmarkHandler(bookmarkService *service.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{bookmarkService: bookmarkService}
}

// AddBookmark godoc
// @Summary Bookmark a story
// @Tags bookmarks
// @Security BearerAuth
// @Param story_id path int true "Story ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/v1/bookmarks/{story_id} [post]
func (h *BookmarkHandler) AddBookmark(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	storyIDStr := c.Param("story_id")
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid story ID")
		return
	}

	err = h.bookmarkService.AddBookmark(c.Request.Context(), userID, storyID)
	if err != nil {
		slog.Warn("add bookmark failed", "error", err, "user_id", userID, "story_id", storyID)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Story bookmarked successfully", nil)
}

// RemoveBookmark godoc
// @Summary Remove a bookmark
// @Tags bookmarks
// @Security BearerAuth
// @Param story_id path int true "Story ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/v1/bookmarks/{story_id} [delete]
func (h *BookmarkHandler) RemoveBookmark(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	storyIDStr := c.Param("story_id")
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid story ID")
		return
	}

	err = h.bookmarkService.RemoveBookmark(c.Request.Context(), userID, storyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bookmark removed successfully", nil)
}

// GetMyBookmarks godoc
// @Summary Get current user's bookmarks
// @Tags bookmarks
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/v1/bookmarks [get]
func (h *BookmarkHandler) GetMyBookmarks(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pagination.Normalize()

	bookmarks, total, err := h.bookmarkService.GetUserBookmarks(c.Request.Context(), userID, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get bookmarks")
		return
	}

	response := utils.NewPaginatedResponse(bookmarks, pagination.Page, pagination.PageSize, total)
	utils.SuccessResponse(c, http.StatusOK, "", response)
}

// GetBookmarkStatus godoc
// @Summary Check if user has bookmarked a story + total bookmarks
// @Tags bookmarks
// @Security BearerAuth
// @Param story_id path int true "Story ID"
// @Success 200 {object} dto.BookmarkStatusResponse
// @Router /api/v1/bookmarks/{story_id}/status [get]
func (h *BookmarkHandler) GetBookmarkStatus(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	storyIDStr := c.Param("story_id")
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid story ID")
		return
	}

	isBookmarked, totalBookmarks, err := h.bookmarkService.GetBookmarkStatus(c.Request.Context(), userID, storyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get bookmark status")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", dto.BookmarkStatusResponse{
		IsBookmarked:   isBookmarked,
		TotalBookmarks: totalBookmarks,
	})
}

// GetViewStats godoc
// @Summary Get view statistics for a story
// @Tags stories
// @Produce json
// @Param slug path string true "Story slug"
// @Success 200 {object} dto.ViewStatsResponse
// @Router /api/v1/stories/{slug}/stats [get]
func (h *BookmarkHandler) GetViewStats(c *gin.Context) {
	slug := c.Param("slug")

	story, err := h.bookmarkService.GetStoryViewStats(c.Request.Context(), slug)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", dto.ViewStatsResponse{
		StoryID:       story.ID,
		StoryTitle:    story.Title,
		StorySlug:     story.Slug,
		TotalViews:    story.TotalViews,
		TotalChapters: story.TotalChapters,
	})
}
