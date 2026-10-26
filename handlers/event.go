package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"looca-go/models"
	"looca-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventHandler struct {
	DB *gorm.DB
}

func NewEventHandler(db *gorm.DB) *EventHandler {
	return &EventHandler{DB: db}
}

type EventInput struct {
	Title        string   `json:"title" binding:"required"`
	SubTitle     *string  `json:"sub_title"`
	Description  *string  `json:"description"`
	Poster       string   `json:"poster" binding:"required"`
	StartDate    *string  `json:"start_date"`
	StartTime    *string  `json:"start_time"`
	Place        string   `json:"place" binding:"required"`
	AdmissionFee *float64 `json:"admission_fee"`
	OriginLink   *string  `json:"origin_link"`
	ShortLink    *string  `json:"short_link"`
	IsComingSoon *bool    `json:"is_coming_soon"`
	IsActive     *bool    `json:"is_active"`
}

func emptyToNil(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (h *EventHandler) validateInput(input *EventInput) (string, bool) {
	title := strings.TrimSpace(input.Title)
	if title == "" || utf8.RuneCountInString(title) > 50 {
		return "title must be between 1 and 50 characters", false
	}

	input.SubTitle = emptyToNil(input.SubTitle)
	if input.SubTitle != nil && utf8.RuneCountInString(*input.SubTitle) > 50 {
		return "sub_title max 50 characters", false
	}

	input.Description = emptyToNil(input.Description)

	poster := strings.TrimSpace(input.Poster)
	if poster == "" {
		return "poster is required", false
	}

	place := strings.TrimSpace(input.Place)
	if place == "" || utf8.RuneCountInString(place) > 255 {
		return "place must be between 1 and 255 characters", false
	}

	input.StartDate = emptyToNil(input.StartDate)
	if input.StartDate != nil {
		if _, err := time.Parse("2006-01-02", *input.StartDate); err != nil {
			return "start_date must be in YYYY-MM-DD format", false
		}
	}

	input.StartTime = emptyToNil(input.StartTime)
	if input.StartTime != nil {
		if _, err := time.Parse("15:04", *input.StartTime); err != nil {
			return "start_time must be in HH:MM format", false
		}
	}

	if input.AdmissionFee != nil && *input.AdmissionFee < 0 {
		return "admission_fee cannot be negative", false
	}

	input.OriginLink = emptyToNil(input.OriginLink)
	input.ShortLink = emptyToNil(input.ShortLink)

	return "", true
}

func (h *EventHandler) List(c *gin.Context) {
	events := make([]models.Event, 0)
	if err := h.DB.Order("id DESC").Find(&events).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}
	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) Get(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var event models.Event
	if err := h.DB.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Event not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Create(c *gin.Context) {
	var input EventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if msg, valid := h.validateInput(&input); !valid {
		utils.Error(c, http.StatusBadRequest, msg)
		return
	}

	isComingSoon := false
	if input.IsComingSoon != nil {
		isComingSoon = *input.IsComingSoon
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	event := models.Event{
		Title:        strings.TrimSpace(input.Title),
		SubTitle:     input.SubTitle,
		Description:  input.Description,
		Poster:       strings.TrimSpace(input.Poster),
		StartDate:    input.StartDate,
		StartTime:    input.StartTime,
		Place:        strings.TrimSpace(input.Place),
		AdmissionFee: input.AdmissionFee,
		OriginLink:   input.OriginLink,
		ShortLink:    input.ShortLink,
		IsComingSoon: isComingSoon,
		IsActive:     isActive,
	}

	if err := h.DB.Create(&event).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create event")
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) Update(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var event models.Event
	if err := h.DB.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Event not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	var input EventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if msg, valid := h.validateInput(&input); !valid {
		utils.Error(c, http.StatusBadRequest, msg)
		return
	}

	event.Title = strings.TrimSpace(input.Title)
	event.SubTitle = input.SubTitle
	event.Description = input.Description
	event.Poster = strings.TrimSpace(input.Poster)
	event.StartDate = input.StartDate
	event.StartTime = input.StartTime
	event.Place = strings.TrimSpace(input.Place)
	event.AdmissionFee = input.AdmissionFee
	event.OriginLink = input.OriginLink
	event.ShortLink = input.ShortLink

	if input.IsComingSoon != nil {
		event.IsComingSoon = *input.IsComingSoon
	}
	if input.IsActive != nil {
		event.IsActive = *input.IsActive
	}

	if err := h.DB.Save(&event).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update event")
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var event models.Event
	if err := h.DB.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Event not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	if err := h.DB.Delete(&event).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
