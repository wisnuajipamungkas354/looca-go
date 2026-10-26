package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"looca-go/config"
	"looca-go/models"

	"gorm.io/gorm"
)

func sp(s string) *string {
	return &s
}

func fp(f float64) *float64 {
	return &f
}

func copyAsset(cfg *config.Config, name, kind string) string {
	destDir := filepath.Join(cfg.UploadDir, kind)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Printf("WARN failed to create dest dir %s: %v", destDir, err)
	}

	srcPath := filepath.Join(cfg.SeedAssetsDir, name)
	destPath := filepath.Join(destDir, name)

	if _, err := os.Stat(destPath); err == nil {
		// File already exists at destination, skip copying
		return fmt.Sprintf("/uploads/%s/%s", kind, name)
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		log.Printf("WARN seed asset missing: %s (%v)", srcPath, err)
		return fmt.Sprintf("/uploads/%s/%s", kind, name)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		log.Printf("WARN failed to create dest file %s: %v", destPath, err)
		return fmt.Sprintf("/uploads/%s/%s", kind, name)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		log.Printf("WARN failed to copy asset %s: %v", name, err)
	}

	return fmt.Sprintf("/uploads/%s/%s", kind, name)
}

func SeedData(db *gorm.DB, cfg *config.Config) {
	// 1. Products
	var productCount int64
	if err := db.Model(&models.Product{}).Count(&productCount).Error; err != nil {
		log.Fatalf("Failed to check product count: %v", err)
	}

	if productCount == 0 {
		products := []models.Product{
			{
				ID:          1,
				Slug:        "polca-looca-scraft-hat",
				Name:        "Polca-Looca Scraft Hat",
				CategoryID:  1,
				Images:      models.StringList{copyAsset(cfg, "fashion2.jpeg", "images"), copyAsset(cfg, "fashion3.jpeg", "images")},
				Description: sp("2 in 1 scraft and hat that avail for 2 size small and medium."),
				ActionLink:  "https://id.shp.ee/AUnMoub8",
				ActionLabel: "Shop Now",
				TotalClick:  0,
			},
			{
				ID:          2,
				Slug:        "custom-attire",
				Name:        "Custom Attire",
				CategoryID:  1,
				Images:      models.StringList{copyAsset(cfg, "fashion1.jpeg", "images"), copyAsset(cfg, "fashion.jpeg", "images")},
				Description: sp("Custom-made kebaya, family outfits, bridesmaid dresses, and more. Tailored to your design, measurements, and style."),
				ActionLink:  "https://www.instagram.com/looca_tailor",
				ActionLabel: "Request Custom",
				TotalClick:  0,
			},
			{
				ID:          3,
				Slug:        "initial-bracelet",
				Name:        "Initial Bracelet",
				CategoryID:  2,
				Images:      models.StringList{copyAsset(cfg, "about1.jpeg", "images"), copyAsset(cfg, "about2.jpeg", "images")},
				Description: sp("Stainless steel bracelet, locally handmade by skilled crafters. Durable, timeless, and made to last."),
				ActionLink:  "https://id.shp.ee/HYVGmweV",
				ActionLabel: "Shop Now",
				TotalClick:  0,
			},
			{
				ID:          4,
				Slug:        "orchid-bagcharm",
				Name:        "Orchid Bagcharm",
				CategoryID:  3,
				Images:      models.StringList{copyAsset(cfg, "about3.jpeg", "images")},
				Description: sp("Handmade orchid charm with delicate bead details"),
				ActionLink:  "https://id.shp.ee/Egy3nwWV",
				ActionLabel: "Shop Now",
				TotalClick:  0,
			},
			{
				ID:          5,
				Slug:        "2-in-1-tasbeeh-bracelet",
				Name:        "2 in 1 Tasbeeh bracelet",
				CategoryID:  2,
				Images:      models.StringList{copyAsset(cfg, "about4.jpeg", "images"), copyAsset(cfg, "about6.jpeg", "images")},
				Description: sp("A versatile bracelet & tasbih in one with pastel tone acrylic beads."),
				ActionLink:  "https://id.shp.ee/yndBxeAL",
				ActionLabel: "Shop Now",
				TotalClick:  0,
			},
			{
				ID:          6,
				Slug:        "k-pop-bracelet-haechan",
				Name:        "K-pop bracelet – Haechan",
				CategoryID:  2,
				Images:      models.StringList{copyAsset(cfg, "acc-kpop-bracelet.jpg", "images")},
				Description: sp("Haechan-inspired bracelet with 2 designs: Bear & Sunflower."),
				ActionLink:  "https://id.shp.ee/D5oPDkBV",
				ActionLabel: "Shop Now",
				TotalClick:  0,
			},
			{
				ID:          7,
				Slug:        "coorporate-brand-merch",
				Name:        "Coorporate/Brand Merch",
				CategoryID:  5,
				Images:      models.StringList{copyAsset(cfg, "souvenir-institution-brand.jpeg", "images"), copyAsset(cfg, "souvenir-institution-brand-2.png", "images")},
				Description: sp("Custom designs with adjustable budget & quantity. Available products: keychains, stickers, mirrors, pouches, and more."),
				ActionLink:  "https://wa.me/6281215312939",
				ActionLabel: "Inquire Bulk Order",
				TotalClick:  0,
			},
			{
				ID:          8,
				Slug:        "wedding-souvenir",
				Name:        "Wedding Souvenir",
				CategoryID:  4,
				Images:      models.StringList{copyAsset(cfg, "souvenir-wedding.jpeg", "images"), copyAsset(cfg, "souvenir-wedding-2.jpeg", "images")},
				Description: sp("Custom-made souvenirs with personalized initials and cute packaging."),
				ActionLink:  "https://wa.me/6281215312939",
				ActionLabel: "Inquire Bulk Order",
				TotalClick:  0,
			},
			{
				ID:          9,
				Slug:        "concert-fan-merch-edition",
				Name:        "Concert & Fan Merch Edition",
				CategoryID:  5,
				Images:      models.StringList{copyAsset(cfg, "merch-concert.jpeg", "images"), copyAsset(cfg, "merch-concert-2.jpeg", "images")},
				Description: sp("Custom concert freebies, lightstick straps, and fan gathering merchandise."),
				ActionLink:  "https://wa.me/6281215312939",
				ActionLabel: "Inquire Bulk Order",
				TotalClick:  0,
			},
		}

		if err := db.Create(&products).Error; err != nil {
			log.Fatalf("Failed to seed products: %v", err)
		}
		log.Println("seed products: inserted 9")
	} else {
		log.Println("seed products: skipped (table not empty)")
	}

	// 2. Workshops
	var workshopCount int64
	if err := db.Model(&models.Workshop{}).Count(&workshopCount).Error; err != nil {
		log.Fatalf("Failed to check workshop count: %v", err)
	}

	if workshopCount == 0 {
		workshops := []models.Workshop{
			{
				ID:          1,
				Title:       "Princess Birthday Crown",
				SubTitle:    sp("Princess Craft"),
				Description: "Create magical sparkling princess crowns adorned with crystals and custom charms.",
				Thumbnail:   copyAsset(cfg, "diy1.jpeg", "images"),
				Video:       sp(copyAsset(cfg, "crown.mp4", "videos")),
				TotalPlay:   0,
			},
			{
				ID:          2,
				Title:       "Orchid Flower Bagcharm",
				SubTitle:    sp("Botanical Craft"),
				Description: "Learn the delicate art of assembling realistic beaded orchid blooms with durable stainless wire.",
				Thumbnail:   copyAsset(cfg, "diy2.jpeg", "images"),
				Video:       sp(copyAsset(cfg, "orchid.mp4", "videos")),
				TotalPlay:   0,
			},
			{
				ID:          3,
				Title:       "Painting Canvas Totebag",
				SubTitle:    sp("Painting Class"),
				Description: "Express yourself by painting artistic custom patterns on organic cotton canvas totebags.",
				Thumbnail:   copyAsset(cfg, "diy3.jpeg", "images"),
				Video:       sp(copyAsset(cfg, "painting-totebag.mp4", "videos")),
				TotalPlay:   0,
			},
			{
				ID:          4,
				Title:       "Paracord Bagcharm Craft",
				SubTitle:    sp("Weaving & Cord"),
				Description: "Knotting and styling colorful outdoor-grade paracords into aesthetic daily accessories.",
				Thumbnail:   copyAsset(cfg, "diy4.jpeg", "images"),
				Video:       sp(copyAsset(cfg, "paracord-bagcharm.mp4", "videos")),
				TotalPlay:   0,
			},
			{
				ID:          5,
				Title:       "Fluffy Moru Doll DIY",
				SubTitle:    sp("Plush & Wire"),
				Description: "Bend, twist, and fluff wire plush chenille stems into super cute portable companion dolls.",
				Thumbnail:   copyAsset(cfg, "diy5.jpeg", "images"),
				Video:       nil,
				TotalPlay:   0,
			},
		}

		if err := db.Create(&workshops).Error; err != nil {
			log.Fatalf("Failed to seed workshops: %v", err)
		}
		log.Println("seed workshops: inserted 5")
	} else {
		log.Println("seed workshops: skipped (table not empty)")
	}

	// 3. Events
	var eventCount int64
	if err := db.Model(&models.Event{}).Count(&eventCount).Error; err != nil {
		log.Fatalf("Failed to check event count: %v", err)
	}

	if eventCount == 0 {
		events := []models.Event{
			{
				ID:           1,
				Title:        "Weekend Beads & Blooms Workshop",
				SubTitle:     sp("Creative Gathering"),
				Description:  sp("Hands-on jewelry making with tea, treats, and fresh flowers at One Ten Cafe Karawang."),
				Poster:       copyAsset(cfg, "diy2.jpeg", "images"),
				StartDate:    sp("2026-11-15"),
				StartTime:    sp("14:00"),
				Place:        "One Ten Cafe, Galuh Mas Karawang",
				AdmissionFee: fp(125000),
				OriginLink:   sp("https://instagram.com/looca_id"),
				ShortLink:    sp("https://looca.id/event/beads-blooms"),
				IsComingSoon: false,
				IsActive:     true,
			},
			{
				ID:           2,
				Title:        "Moru Doll Pop-Up Studio",
				SubTitle:     sp("Craft & Sip"),
				Description:  sp("Bring your favorite fluffy friends to life in this cozy 2-hour guided crafting session."),
				Poster:       copyAsset(cfg, "diy5.jpeg", "images"),
				StartDate:    sp("2026-11-28"),
				StartTime:    sp("10:30"),
				Place:        "Klinik Bisnis IPB, Bogor",
				AdmissionFee: fp(85000),
				OriginLink:   sp("https://instagram.com/looca_id"),
				ShortLink:    sp("https://looca.id/event/moru-studio"),
				IsComingSoon: false,
				IsActive:     true,
			},
			{
				ID:           3,
				Title:        "Looca Year-End Art Bazaar",
				SubTitle:     sp("Annual Showcase"),
				Description:  sp("Exclusive new handcrafted lines, live custom tailor orders, and interactive booth games."),
				Poster:       copyAsset(cfg, "fashion1.jpeg", "images"),
				StartDate:    sp("2026-12-20"),
				StartTime:    sp("09:00"),
				Place:        "Outlet Dekranasda Karawang",
				AdmissionFee: nil,
				OriginLink:   sp("https://instagram.com/looca_id"),
				ShortLink:    sp("https://looca.id/event/art-bazaar"),
				IsComingSoon: true,
				IsActive:     true,
			},
			{
				ID:           4,
				Title:        "Exclusive Private Corporate Crafting",
				SubTitle:     sp("Team Building"),
				Description:  sp("Past event for corporate creative recreation and custom pouch painting."),
				Poster:       copyAsset(cfg, "diy3.jpeg", "images"),
				StartDate:    sp("2026-10-01"),
				StartTime:    sp("13:00"),
				Place:        "Jakarta Creative Hub",
				AdmissionFee: fp(250000),
				OriginLink:   nil,
				ShortLink:    nil,
				IsComingSoon: false,
				IsActive:     false,
			},
		}

		if err := db.Create(&events).Error; err != nil {
			log.Fatalf("Failed to seed events: %v", err)
		}
		log.Println("seed events: inserted 4")
	} else {
		log.Println("seed events: skipped (table not empty)")
	}
}
