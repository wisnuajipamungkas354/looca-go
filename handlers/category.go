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

type CategoryHandler struct {
	DB *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

type CategoryInput struct {
	CategoriesName string `json:"categories_name" binding:"required"`
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories := make([]models.Category, 0)
	if err := h.DB.Order("id DESC").Find(&categories).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "categories_name is required")
		return
	}

	name := strings.TrimSpace(input.CategoriesName)
	if name == "" || utf8.RuneCountInString(name) > 255 {
		utils.Error(c, http.StatusBadRequest, "categories_name must be between 1 and 255 characters")
		return
	}

	slug := utils.GenerateUniqueSlug(h.DB, "categories", name, 0)
	category := models.Category{
		Slug:           slug,
		CategoriesName: name,
	}

	if err := h.DB.Create(&category).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create category")
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Category not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "categories_name is required")
		return
	}

	name := strings.TrimSpace(input.CategoriesName)
	if name == "" || utf8.RuneCountInString(name) > 255 {
		utils.Error(c, http.StatusBadRequest, "categories_name must be between 1 and 255 characters")
		return
	}

	if category.CategoriesName != name {
		category.CategoriesName = name
		category.Slug = utils.GenerateUniqueSlug(h.DB, "categories", name, category.ID)
	}
	if err := h.DB.Save(&category).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update category")
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Category not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	var productCount int64
	if err := h.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	if productCount > 0 {
		utils.Error(c, http.StatusConflict, "Category is still used by products")
		return
	}

	if err := h.DB.Delete(&category).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
