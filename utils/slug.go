package utils

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)
	multipleDashRegex    = regexp.MustCompile(`-+`)
)

func Slugify(text string) string {
	s := strings.ToLower(strings.TrimSpace(text))
	// Replace non-alphanumeric characters with dash
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")
	// Collapse multiple dashes
	s = multipleDashRegex.ReplaceAllString(s, "-")
	// Trim dashes from start and end
	s = strings.Trim(s, "-")
	if s == "" {
		return "item"
	}
	return s
}

func GenerateUniqueSlug(db *gorm.DB, table string, baseText string, excludeID uint) string {
	baseSlug := Slugify(baseText)
	slug := baseSlug
	suffix := 2

	for {
		var count int64
		query := db.Table(table).Where("slug = ?", slug)
		if excludeID > 0 {
			query = query.Where("id != ?", excludeID)
		}
		if err := query.Count(&count).Error; err != nil {
			return slug
		}
		if count == 0 {
			return slug
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, suffix)
		suffix++
	}
}
