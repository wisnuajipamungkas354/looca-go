package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"looca-go/models"
	"looca-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

type ProductInput struct {
	Name        string   `json:"name" binding:"required"`
	CategoryID  uint     `json:"category_id" binding:"required"`
	Images      []string `json:"images" binding:"required"`
	Description *string  `json:"description"`
	ActionLink  string   `json:"action_link" binding:"required"`
	ActionLabel string   `json:"action_label" binding:"required"`
}

func (h *ProductHandler) validateInput(input *ProductInput) (string, bool) {
	name := strings.TrimSpace(input.Name)
	if name == "" || utf8.RuneCountInString(name) > 255 {
		return "name must be between 1 and 255 characters", false
	}

	var categoryCount int64
	if err := h.DB.Model(&models.Category{}).Where("id = ?", input.CategoryID).Count(&categoryCount).Error; err != nil || categoryCount == 0 {
		return "category not found", false
	}

	if len(input.Images) < 1 || len(input.Images) > 3 {
		return "images must contain between 1 and 3 items", false
	}

	for _, img := range input.Images {
		if strings.TrimSpace(img) == "" {
			return "image url cannot be empty", false
		}
	}

	actionLink := strings.TrimSpace(input.ActionLink)
	if actionLink == "" {
		return "action_link is required", false
	}

	actionLabel := strings.TrimSpace(input.ActionLabel)
	if actionLabel == "" || utf8.RuneCountInString(actionLabel) > 30 {
		return "action_label must be between 1 and 30 characters", false
	}

	return "", true
}

func (h *ProductHandler) List(c *gin.Context) {
	products := make([]models.Product, 0)
	query := h.DB.Model(&models.Product{}).Order("id DESC")

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if catID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			query = query.Where("category_id = ?", catID)
		}
	}

	if err := query.Find(&products).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var product models.Product
	if err := h.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Product not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if msg, valid := h.validateInput(&input); !valid {
		utils.Error(c, http.StatusBadRequest, msg)
		return
	}

	name := strings.TrimSpace(input.Name)
	slug := utils.GenerateUniqueSlug(h.DB, "products", name, 0)

	product := models.Product{
		Slug:        slug,
		Name:        name,
		CategoryID:  input.CategoryID,
		Images:      models.StringList(input.Images),
		Description: input.Description,
		ActionLink:  strings.TrimSpace(input.ActionLink),
		ActionLabel: strings.TrimSpace(input.ActionLabel),
		TotalClick:  0,
	}

	if err := h.DB.Create(&product).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create product")
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var product models.Product
	if err := h.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Product not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if msg, valid := h.validateInput(&input); !valid {
		utils.Error(c, http.StatusBadRequest, msg)
		return
	}

	newName := strings.TrimSpace(input.Name)
	if product.Name != newName {
		product.Name = newName
		product.Slug = utils.GenerateUniqueSlug(h.DB, "products", newName, product.ID)
	}
	product.CategoryID = input.CategoryID
	product.Images = models.StringList(input.Images)
	product.Description = input.Description
	product.ActionLink = strings.TrimSpace(input.ActionLink)
	product.ActionLabel = strings.TrimSpace(input.ActionLabel)
	// total_click is untouched

	if err := h.DB.Save(&product).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update product")
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var product models.Product
	if err := h.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Product not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	if err := h.DB.Delete(&product).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
