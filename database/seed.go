package database

import (
	"log"

	"looca-go/config"
	"looca-go/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedEssential(db *gorm.DB, cfg *config.Config) {
	// 1. Admin
	var adminCount int64
	if err := db.Model(&models.Admin{}).Count(&adminCount).Error; err != nil {
		log.Fatalf("Failed to check admin count: %v", err)
	}

	if adminCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash admin password: %v", err)
		}

		admin := models.Admin{
			Email:        cfg.AdminEmail,
			Name:         cfg.AdminName,
			PasswordHash: string(hash),
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to seed admin: %v", err)
		}
		log.Println("seed admin: inserted 1")
	} else {
		log.Println("seed admin: skipped (table not empty)")
	}

	// 2. Categories
	var categoryCount int64
	if err := db.Model(&models.Category{}).Count(&categoryCount).Error; err != nil {
		log.Fatalf("Failed to check category count: %v", err)
	}

	if categoryCount == 0 {
		categories := []models.Category{
			{ID: 1, Slug: "fashion", CategoriesName: "Fashion"},
			{ID: 2, Slug: "accessories", CategoriesName: "Accessories"},
			{ID: 3, Slug: "bagcharm", CategoriesName: "Bagcharm"},
			{ID: 4, Slug: "souvenir", CategoriesName: "Souvenir"},
			{ID: 5, Slug: "merchandise", CategoriesName: "Merchandise"},
		}
		if err := db.Create(&categories).Error; err != nil {
			log.Fatalf("Failed to seed categories: %v", err)
		}
		log.Println("seed categories: inserted 5")
	} else {
		log.Println("seed categories: skipped (table not empty)")
	}

	// 3. Settings
	var settingCount int64
	if err := db.Model(&models.Setting{}).Count(&settingCount).Error; err != nil {
		log.Fatalf("Failed to check setting count: %v", err)
	}

	if settingCount == 0 {
		setting := models.Setting{
			ID:             1,
			WhatsappNumber: "6281215312939",
			InstagramURL:   "https://instagram.com/looca_id",
			TiktokURL:      "https://tiktok.com/@looca_id",
			ShopeeURL:      "https://shopee.co.id/looca_id",
		}
		if err := db.Create(&setting).Error; err != nil {
			log.Fatalf("Failed to seed settings: %v", err)
		}
		log.Println("seed settings: inserted 1")
	} else {
		log.Println("seed settings: skipped (table not empty)")
	}
}
