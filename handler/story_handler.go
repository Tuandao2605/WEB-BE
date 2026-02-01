package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"web-be/dto"
	"web-be/middleware"
	"web-be/service"
	"web-be/utils"
)

type StoryHandler struct {
	storyService *service.StoryService
}

func NewStoryHandler(storyService *service.StoryService) *StoryHandler {
	return &StoryHandler{storyService: storyService}
}

// GetAll godoc
// @Summary Get all stories
// @Tags stories
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/v1/stories [get]
func (h *StoryHandler) GetAll(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pagination.Normalize()

	stories, total, err := h.storyService.GetAll(c.Request.Context(), pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get stories")
		return
	}

	response := utils.NewPaginatedResponse(stories, pagination.Page, pagination.PageSize, total)
	utils.SuccessResponse(c, http.StatusOK, "", response)
}

// GetBySlug godoc
// @Summary Get a story by slug
// @Tags stories
// @Produce json
// @Param slug path string true "Story slug"
// @Success 200 {object} models.Story
// @Failure 404 {object} utils.APIResponse
// @Router /api/v1/stories/{slug} [get]
func (h *StoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	story, err := h.storyService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", story)
}

// Search godoc
// @Summary Search stories
// @Tags stories
// @Produce json
// @Param q query string true "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/v1/stories/search [get]
func (h *StoryHandler) Search(c *gin.Context) {
	var req dto.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req.Normalize()

	stories, total, err := h.storyService.Search(c.Request.Context(), req.Query, req.GetLimit(), req.GetOffset())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search stories")
		return
	}

	response := utils.NewPaginatedResponse(stories, req.Page, req.PageSize, total)
	utils.SuccessResponse(c, http.StatusOK, "", response)
}

// Create godoc
// @Summary Create a new story
// @Tags stories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateStoryRequest true "Story details"
// @Success 201 {object} models.Story
// @Failure 400 {object} utils.APIResponse
// @Router /api/v1/stories [post]
func (h *StoryHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	story, err := h.storyService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Story created successfully", story)
}

// Update godoc
// @Summary Update a story
// @Tags stories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param slug path string true "Story slug"
// @Param request body dto.UpdateStoryRequest true "Story update details"
// @Success 200 {object} models.Story
// @Failure 400 {object} utils.APIResponse
// @Router /api/v1/stories/{slug} [put]
func (h *StoryHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userRole, _ := middleware.GetUserRole(c)

	slug := c.Param("slug")

	var req dto.UpdateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	story, err := h.storyService.Update(c.Request.Context(), slug, userID, userRole, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Story updated successfully", story)
}

// Delete godoc
// @Summary Delete a story
// @Tags stories
// @Security BearerAuth
// @Param slug path string true "Story slug"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Router /api/v1/stories/{slug} [delete]
func (h *StoryHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userRole, _ := middleware.GetUserRole(c)

	slug := c.Param("slug")

	err := h.storyService.Delete(c.Request.Context(), slug, userID, userRole)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Story deleted successfully", nil)
}

// GetByCategory godoc
// @Summary Get stories by category
// @Tags categories
// @Produce json
// @Param slug path string true "Category slug"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/v1/categories/{slug}/stories [get]
func (h *StoryHandler) GetByCategory(c *gin.Context) {
	slug := c.Param("slug")

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pagination.Normalize()

	stories, total, err := h.storyService.GetByCategory(c.Request.Context(), slug, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response := utils.NewPaginatedResponse(stories, pagination.Page, pagination.PageSize, total)
	utils.SuccessResponse(c, http.StatusOK, "", response)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Router /api/v1/categories [get]
func (h *StoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.storyService.GetAllCategories(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get categories")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", categories)
}

// CreateCategory godoc
// @Summary Create a new category (admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category details"
// @Success 201 {object} models.Category
// @Router /api/v1/admin/categories [post]
func (h *StoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.storyService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created", category)
}

// Publish godoc
// @Summary Publish/Unpublish a story (admin only)
// @Tags admin
// @Security BearerAuth
// @Param id path int true "Story ID"
// @Param publish query bool true "Publish status"
// @Success 200 {object} utils.APIResponse
// @Router /api/v1/admin/stories/{id}/publish [put]
func (h *StoryHandler) Publish(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid story ID")
		return
	}

	publish := c.Query("publish") == "true"

	err = h.storyService.Publish(c.Request.Context(), id, publish)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	message := "Story unpublished"
	if publish {
		message = "Story published"
	}
	utils.SuccessResponse(c, http.StatusOK, message, nil)
}

// GetReadingHistory godoc
// @Summary Get user's reading history
// @Tags history
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/v1/history [get]
func (h *StoryHandler) GetReadingHistory(c *gin.Context) {
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

	history, total, err := h.storyService.GetReadingHistory(c.Request.Context(), userID, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get reading history")
		return
	}

	response := utils.NewPaginatedResponse(history, pagination.Page, pagination.PageSize, total)
	utils.SuccessResponse(c, http.StatusOK, "", response)
}

// UpdateReadingHistory godoc
// @Summary Update reading history
// @Tags history
// @Security BearerAuth
// @Param story_id path int true "Story ID"
// @Param chapter_id query int false "Last read chapter ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/v1/history/{story_id} [post]
func (h *StoryHandler) UpdateReadingHistory(c *gin.Context) {
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

	var chapterID *int
	if chapterIDStr := c.Query("chapter_id"); chapterIDStr != "" {
		id, err := strconv.Atoi(chapterIDStr)
		if err == nil {
			chapterID = &id
		}
	}

	err = h.storyService.UpdateReadingHistory(c.Request.Context(), userID, storyID, chapterID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reading history updated", nil)
}

// DeleteReadingHistory godoc
// @Summary Delete a reading history entry
// @Tags history
// @Security BearerAuth
// @Param story_id path int true "Story ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/v1/history/{story_id} [delete]
func (h *StoryHandler) DeleteReadingHistory(c *gin.Context) {
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

	err = h.storyService.DeleteReadingHistory(c.Request.Context(), userID, storyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reading history deleted", nil)
}

// GetMyStories godoc
// @Summary Get current user's stories
// @Tags stories
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/v1/my-stories [get]
func (h *StoryHandler) GetMyStories(c *gin.Context) {
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

	stories, total, err := h.storyService.GetByAuthor(c.Request.Context(), userID, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get stories")
		return
	}

	response := utils.NewPaginatedResponse(stories, pagination.Page, pagination.PageSize, total)
	utils.SuccessResponse(c, http.StatusOK, "", response)
}
