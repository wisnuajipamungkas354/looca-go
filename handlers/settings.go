package handlers

import (
	"errors"
	"net/http"
	"strings"

	"looca-go/models"
	"looca-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingsHandler struct {
	DB *gorm.DB
}

func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{DB: db}
}

type SettingsInput struct {
	WhatsappNumber string `json:"whatsapp_number" binding:"required"`
	InstagramURL   string `json:"instagram_url"`
	TiktokURL      string `json:"tiktok_url"`
	ShopeeURL      string `json:"shopee_url"`
}

func (h *SettingsHandler) Get(c *gin.Context) {
	var setting models.Setting
	if err := h.DB.First(&setting, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If not seeded yet, return empty object
			c.JSON(http.StatusOK, models.Setting{
				WhatsappNumber: "",
				InstagramURL:   "",
				TiktokURL:      "",
				ShopeeURL:      "",
			})
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	c.JSON(http.StatusOK, setting)
}

func (h *SettingsHandler) Update(c *gin.Context) {
	var input SettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	wa := strings.TrimSpace(input.WhatsappNumber)
	if wa == "" {
		utils.Error(c, http.StatusBadRequest, "whatsapp_number is required")
		return
	}

	setting := models.Setting{
		ID:             1,
		WhatsappNumber: wa,
		InstagramURL:   strings.TrimSpace(input.InstagramURL),
		TiktokURL:      strings.TrimSpace(input.TiktokURL),
		ShopeeURL:      strings.TrimSpace(input.ShopeeURL),
	}

	if err := h.DB.Save(&setting).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save settings")
		return
	}

	c.JSON(http.StatusOK, setting)
}
