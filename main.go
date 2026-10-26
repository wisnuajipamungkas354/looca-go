package main

import (
	"log"
	"os"

	"looca-go/config"
	"looca-go/database"
	"looca-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload dir %s: %v", cfg.UploadDir, err)
	}

	db := database.Connect(cfg)
	database.Migrate(db)

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		log.Println("Running manual database seeding...")
		database.SeedEssential(db, cfg)
		database.SeedData(db, cfg)
		log.Println("Seeding completed successfully.")
		return
	}

	r := gin.Default()
	routes.Setup(r, db, cfg)

	log.Printf("Server listening on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
