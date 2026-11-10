package handlers

import (
	"errors"
	"net/http"
	"net/url"
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

func ParseVisitorSource(referrer string) string {
	ref := strings.ToLower(strings.TrimSpace(referrer))
	if ref == "" {
		return "Direct"
	}

	if strings.Contains(ref, "threads.net") || ref == "threads" {
		return "Threads"
	}
	if strings.Contains(ref, "instagram.com") || ref == "instagram" || ref == "ig" {
		return "Instagram"
	}
	if strings.Contains(ref, "google.") || ref == "google" {
		return "Google"
	}
	if strings.Contains(ref, "tiktok.com") || ref == "tiktok" {
		return "TikTok"
	}
	if strings.Contains(ref, "whatsapp.com") || strings.Contains(ref, "wa.me") || ref == "whatsapp" || ref == "wa" {
		return "WhatsApp"
	}
	if strings.Contains(ref, "twitter.com") || strings.Contains(ref, "t.co") || strings.Contains(ref, "x.com") || ref == "twitter" || ref == "x" {
		return "Twitter / X"
	}
	if strings.Contains(ref, "facebook.com") || strings.Contains(ref, "fb.com") || ref == "facebook" || ref == "fb" {
		return "Facebook"
	}
	if strings.Contains(ref, "youtube.com") || strings.Contains(ref, "youtu.be") || ref == "youtube" {
		return "YouTube"
	}
	if strings.Contains(ref, "shopee.") || ref == "shopee" {
		return "Shopee"
	}
	if strings.Contains(ref, "tokopedia.com") || ref == "tokopedia" {
		return "Tokopedia"
	}
	if strings.Contains(ref, "pinterest.") || ref == "pinterest" {
		return "Pinterest"
	}
	if strings.Contains(ref, "linkedin.com") || ref == "linkedin" {
		return "LinkedIn"
	}
	if strings.Contains(ref, "looca-id.com") || strings.Contains(ref, "looca.id") || strings.Contains(ref, "localhost") || strings.Contains(ref, "127.0.0.1") || ref == "direct" {
		return "Direct"
	}

	u, err := url.Parse(referrer)
	if err == nil && u.Host != "" {
		host := strings.TrimPrefix(u.Host, "www.")
		if host != "" {
			return host
		}
	}

	return "Direct"
}

type VisitPayload struct {
	Referrer string `json:"referrer"`
}

func (h *PublicHandler) Visit(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	// Upsert daily visit: increment visitors count for today
	err := h.DB.Exec(
		"INSERT INTO `daily_visits` (`visit_date`, `visitors`) VALUES (?, 1) ON DUPLICATE KEY UPDATE `visitors` = `visitors` + 1",
		today,
	).Error

	if err != nil {
		// Fallback for different dialects or driver configurations
		visit := models.DailyVisit{
			VisitDate: today,
			Visitors:  1,
		}
		_ = h.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "visit_date"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"visitors": gorm.Expr("visitors + ?", 1),
			}),
		}).Create(&visit).Error
	}

	// Read referrer from request body or header or query params
	var payload VisitPayload
	_ = c.ShouldBindJSON(&payload)

	rawRef := payload.Referrer
	if rawRef == "" {
		rawRef = c.Query("ref")
	}
	if rawRef == "" {
		rawRef = c.Query("utm_source")
	}
	if rawRef == "" {
		rawRef = c.GetHeader("Referer")
	}

	source := ParseVisitorSource(rawRef)

	// Upsert daily visitor source
	_ = h.DB.Exec(
		"INSERT INTO `daily_visitor_sources` (`visit_date`, `source`, `visitors`) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE `visitors` = `visitors` + 1",
		today, source,
	).Error

	// Upsert total visitor source
	sourceErr := h.DB.Exec(
		"INSERT INTO `visitor_sources` (`source`, `visitors`) VALUES (?, 1) ON DUPLICATE KEY UPDATE `visitors` = `visitors` + 1",
		source,
	).Error

	if sourceErr != nil {
		vs := models.VisitorSource{
			Source:   source,
			Visitors: 1,
		}
		_ = h.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "source"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"visitors": gorm.Expr("visitors + ?", 1),
			}),
		}).Create(&vs).Error
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *PublicHandler) RedirectEvent(c *gin.Context) {
	rawSlug := strings.TrimSpace(c.Param("slug"))
	slug := strings.TrimSuffix(strings.TrimSuffix(rawSlug, "/"), ".html")
	if slug == "" {
		c.Redirect(http.StatusFound, "https://looca-id.com")
		return
	}

	var event models.Event
	// Match exact slug, full URL, ending with /slug, containing /slug, or numeric ID
	err := h.DB.Where(
		"short_link = ? OR short_link = ? OR short_link = ? OR short_link LIKE ? OR short_link LIKE ? OR id = ?",
		slug,
		"https://looca-id.com/event/"+slug,
		"https://looca.id/event/"+slug,
		"%/"+slug,
		"%/"+slug+"%",
		slug,
	).First(&event).Error

	if err != nil {
		c.Redirect(http.StatusFound, "https://looca-id.com")
		return
	}

	if event.OriginLink == nil || strings.TrimSpace(*event.OriginLink) == "" {
		c.Redirect(http.StatusFound, "https://looca-id.com")
		return
	}

	dest := strings.TrimSpace(*event.OriginLink)
	if !strings.HasPrefix(dest, "http://") && !strings.HasPrefix(dest, "https://") {
		dest = "https://" + dest
	}

	c.Redirect(http.StatusFound, dest)
}
