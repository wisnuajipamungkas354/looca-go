package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"looca-go/models"
	"looca-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PublicHandler struct {
	DB *gorm.DB
}

func NewPublicHandler(db *gorm.DB) *PublicHandler {
	return &PublicHandler{DB: db}
}

func (h *PublicHandler) Categories(c *gin.Context) {
	categories := make([]models.Category, 0)
	if err := h.DB.Order("id DESC").Find(&categories).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *PublicHandler) Products(c *gin.Context) {
	products := make([]models.Product, 0)
	query := h.DB.Model(&models.Product{}).Order("id DESC")

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if catID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			query = query.Where("category_id = ?", catID)
		}
	} else if categorySlug := c.Query("category"); categorySlug != "" {
		var cat models.Category
		if err := h.DB.Where("slug = ?", categorySlug).First(&cat).Error; err == nil {
			query = query.Where("category_id = ?", cat.ID)
		} else {
			// If category slug not found, return empty list directly
			c.JSON(http.StatusOK, products)
			return
		}
	}

	if err := query.Find(&products).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *PublicHandler) Workshops(c *gin.Context) {
	workshops := make([]models.Workshop, 0)
	if err := h.DB.Order("id DESC").Find(&workshops).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve workshops")
		return
	}
	c.JSON(http.StatusOK, workshops)
}

func (h *PublicHandler) Events(c *gin.Context) {
	events := make([]models.Event, 0)
	if err := h.DB.Where("is_active = ?", true).Order("id DESC").Find(&events).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}
	c.JSON(http.StatusOK, events)
}

func (h *PublicHandler) Settings(c *gin.Context) {
	var setting models.Setting
	if err := h.DB.First(&setting, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (h *PublicHandler) ClickProduct(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	res := h.DB.Model(&models.Product{}).Where("id = ?", id).Update("total_click", gorm.Expr("total_click + ?", 1))
	if res.Error != nil {
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}
	if res.RowsAffected == 0 {
		utils.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *PublicHandler) PlayWorkshop(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	res := h.DB.Model(&models.Workshop{}).Where("id = ?", id).Update("total_play", gorm.Expr("total_play + ?", 1))
	if res.Error != nil {
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}
	if res.RowsAffected == 0 {
		utils.Error(c, http.StatusNotFound, "Workshop not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *PublicHandler) Visit(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	visit := models.DailyVisit{
		VisitDate: today,
		Visitors:  1,
	}

	err := h.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "visit_date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"visitors": gorm.Expr("visitors + ?", 1),
		}),
	}).Create(&visit).Error

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to record visit")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *PublicHandler) RedirectEvent(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if slug == "" {
		utils.Error(c, http.StatusBadRequest, "slug is required")
		return
	}

	var event models.Event
	// Match exact slug, or short_link ending with /slug, or short_link equal to slug
	err := h.DB.Where(
		"short_link = ? OR short_link LIKE ? OR short_link LIKE ?",
		slug,
		"%/"+slug,
		"%/"+slug+"%",
	).First(&event).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Event not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	if event.OriginLink == nil || strings.TrimSpace(*event.OriginLink) == "" {
		utils.Error(c, http.StatusNotFound, "Origin link not available for this event")
		return
	}

	c.Redirect(http.StatusFound, *event.OriginLink)
}
