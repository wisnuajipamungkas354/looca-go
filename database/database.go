// Package database provides MySQL database connection and migrations.
// NOTE: Database 'looca' must already exist in MySQL before connecting.
package database

import (
	"fmt"
	"log"

	"looca-go/config"
	"looca-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Admin{},
		&models.Category{},
		&models.Product{},
		&models.Workshop{},
		&models.Event{},
		&models.Setting{},
		&models.DailyVisit{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Ensure start_time in events is VARCHAR(10) instead of DATETIME if previously created with type:time
	_ = db.Exec("ALTER TABLE `events` MODIFY COLUMN `start_time` VARCHAR(10) NULL").Error
}
