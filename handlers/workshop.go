package handlers

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"looca-go/models"
	"looca-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkshopHandler struct {
	DB *gorm.DB
}

func NewWorkshopHandler(db *gorm.DB) *WorkshopHandler {
	return &WorkshopHandler{DB: db}
}

type WorkshopInput struct {
	Title       string  `json:"title" binding:"required"`
	SubTitle    *string `json:"sub_title"`
	Description string  `json:"description" binding:"required"`
	Thumbnail   string  `json:"thumbnail" binding:"required"`
	Video       *string `json:"video"`
}

func (h *WorkshopHandler) validateInput(input *WorkshopInput) (string, bool) {
	title := strings.TrimSpace(input.Title)
	if title == "" || utf8.RuneCountInString(title) > 50 {
		return "title must be between 1 and 50 characters", false
	}

	if input.SubTitle != nil {
		st := strings.TrimSpace(*input.SubTitle)
		if st == "" {
			input.SubTitle = nil
		} else if utf8.RuneCountInString(st) > 50 {
			return "sub_title max 50 characters", false
		} else {
			input.SubTitle = &st
		}
	}

	desc := strings.TrimSpace(input.Description)
	if desc == "" || utf8.RuneCountInString(desc) > 150 {
		return "description must be between 1 and 150 characters", false
	}

	thumb := strings.TrimSpace(input.Thumbnail)
	if thumb == "" {
		return "thumbnail is required", false
	}

	if input.Video != nil {
		v := strings.TrimSpace(*input.Video)
		if v == "" {
			input.Video = nil
		} else {
			input.Video = &v
		}
	}

	return "", true
}

func (h *WorkshopHandler) List(c *gin.Context) {
	workshops := make([]models.Workshop, 0)
	if err := h.DB.Order("id DESC").Find(&workshops).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve workshops")
		return
	}
	c.JSON(http.StatusOK, workshops)
}

func (h *WorkshopHandler) Get(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var workshop models.Workshop
	if err := h.DB.First(&workshop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Workshop not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	c.JSON(http.StatusOK, workshop)
}

func (h *WorkshopHandler) Create(c *gin.Context) {
	var input WorkshopInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if msg, valid := h.validateInput(&input); !valid {
		utils.Error(c, http.StatusBadRequest, msg)
		return
	}

	workshop := models.Workshop{
		Title:       strings.TrimSpace(input.Title),
		SubTitle:    input.SubTitle,
		Description: strings.TrimSpace(input.Description),
		Thumbnail:   strings.TrimSpace(input.Thumbnail),
		Video:       input.Video,
		TotalPlay:   0,
	}

	if err := h.DB.Create(&workshop).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create workshop")
		return
	}

	c.JSON(http.StatusCreated, workshop)
}

func (h *WorkshopHandler) Update(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var workshop models.Workshop
	if err := h.DB.First(&workshop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Workshop not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	var input WorkshopInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if msg, valid := h.validateInput(&input); !valid {
		utils.Error(c, http.StatusBadRequest, msg)
		return
	}

	workshop.Title = strings.TrimSpace(input.Title)
	workshop.SubTitle = input.SubTitle
	workshop.Description = strings.TrimSpace(input.Description)
	workshop.Thumbnail = strings.TrimSpace(input.Thumbnail)
	workshop.Video = input.Video
	// TotalPlay is untouched

	if err := h.DB.Save(&workshop).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update workshop")
		return
	}

	c.JSON(http.StatusOK, workshop)
}

func (h *WorkshopHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var workshop models.Workshop
	if err := h.DB.First(&workshop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Workshop not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	if err := h.DB.Delete(&workshop).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete workshop")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
