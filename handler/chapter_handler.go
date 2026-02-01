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

type ChapterHandler struct {
	chapterService *service.ChapterService
}

func NewChapterHandler(chapterService *service.ChapterService) *ChapterHandler {
	return &ChapterHandler{chapterService: chapterService}
}

// GetByStory godoc
// @Summary Get chapters list by story
// @Tags chapters
// @Produce json
// @Param slug path string true "Story slug"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/v1/stories/{slug}/chapters [get]
func (h *ChapterHandler) GetByStory(c *gin.Context) {
	storySlug := c.Param("slug")

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pagination.Normalize()

	chapters, total, err := h.chapterService.GetListByStory(c.Request.Context(), storySlug, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response := utils.NewPaginatedResponse(chapters, pagination.Page, pagination.PageSize, total)
	utils.SuccessResponse(c, http.StatusOK, "", response)
}

// GetChapter godoc
// @Summary Get a chapter content
// @Tags chapters
// @Produce json
// @Param slug path string true "Story slug"
// @Param chapter_num path int true "Chapter number"
// @Success 200 {object} dto.ChapterResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/v1/stories/{slug}/chapters/{chapter_num} [get]
func (h *ChapterHandler) GetChapter(c *gin.Context) {
	storySlug := c.Param("slug")
	chapterNumStr := c.Param("chapter_num")

	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid chapter number")
		return
	}

	chapter, err := h.chapterService.GetByStoryAndNumber(c.Request.Context(), storySlug, chapterNum)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", chapter)
}

// Create godoc
// @Summary Create a new chapter
// @Tags chapters
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param slug path string true "Story slug"
// @Param request body dto.CreateChapterRequest true "Chapter details"
// @Success 201 {object} models.Chapter
// @Failure 400 {object} utils.APIResponse
// @Router /api/v1/stories/{slug}/chapters [post]
func (h *ChapterHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userRole, _ := middleware.GetUserRole(c)

	storySlug := c.Param("slug")

	var req dto.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	chapter, err := h.chapterService.Create(c.Request.Context(), storySlug, userID, userRole, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Chapter created successfully", chapter)
}

// Update godoc
// @Summary Update a chapter
// @Tags chapters
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param slug path string true "Story slug"
// @Param chapter_num path int true "Chapter number"
// @Param request body dto.UpdateChapterRequest true "Chapter update details"
// @Success 200 {object} models.Chapter
// @Failure 400 {object} utils.APIResponse
// @Router /api/v1/stories/{slug}/chapters/{chapter_num} [put]
func (h *ChapterHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userRole, _ := middleware.GetUserRole(c)

	storySlug := c.Param("slug")
	chapterNumStr := c.Param("chapter_num")

	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid chapter number")
		return
	}

	var req dto.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	chapter, err := h.chapterService.Update(c.Request.Context(), storySlug, chapterNum, userID, userRole, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Chapter updated successfully", chapter)
}

// Delete godoc
// @Summary Delete a chapter
// @Tags chapters
// @Security BearerAuth
// @Param slug path string true "Story slug"
// @Param chapter_num path int true "Chapter number"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Router /api/v1/stories/{slug}/chapters/{chapter_num} [delete]
func (h *ChapterHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userRole, _ := middleware.GetUserRole(c)

	storySlug := c.Param("slug")
	chapterNumStr := c.Param("chapter_num")

	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid chapter number")
		return
	}

	err = h.chapterService.Delete(c.Request.Context(), storySlug, chapterNum, userID, userRole)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Chapter deleted successfully", nil)
}
