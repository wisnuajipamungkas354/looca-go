package handlers

import (
	"errors"
	"net/http"
	"strings"

	"looca-go/models"
	"looca-go/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	JWTSecret string
}

func NewAuthHandler(db *gorm.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{DB: db, JWTSecret: jwtSecret}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	var admin models.Admin
	if err := h.DB.Where("LOWER(email) = ?", email).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(h.JWTSecret, admin.ID, admin.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserResponse{
			Email: admin.Email,
			Name:  admin.Name,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	adminIDVal, exists := c.Get("admin_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	adminID, ok := adminIDVal.(uint)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var admin models.Admin
	if err := h.DB.First(&admin, adminID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "Admin not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Database error")
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Email: admin.Email,
		Name:  admin.Name,
	})
}
