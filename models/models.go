package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (s *StringList) Scan(value any) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan StringList: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

func (StringList) GormDataType() string {
	return "json"
}

type Admin struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Name         string `gorm:"size:255;not null" json:"name"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
}

func (Admin) TableName() string {
	return "admins"
}

type Category struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug           string `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	CategoriesName string `gorm:"column:categories_name;size:255;not null" json:"categories_name"`
}

func (Category) TableName() string {
	return "categories"
}

type Product struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug        string     `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	CategoryID  uint       `gorm:"not null;index" json:"category_id"`
	Images      StringList `gorm:"type:json;not null" json:"images"`
	Description *string    `gorm:"type:text" json:"description"`
	ActionLink  string     `gorm:"type:text;not null" json:"action_link"`
	ActionLabel string     `gorm:"size:30;not null" json:"action_label"`
	TotalClick  int        `gorm:"not null;default:0" json:"total_click"`
}

func (Product) TableName() string {
	return "products"
}

type Workshop struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string  `gorm:"size:50;not null" json:"title"`
	SubTitle    *string `gorm:"size:50" json:"sub_title"`
	Description string  `gorm:"size:150;not null" json:"description"`
	Thumbnail   string  `gorm:"type:text;not null" json:"thumbnail"`
	Video       *string `gorm:"type:text" json:"video"`
	TotalPlay   int     `gorm:"not null;default:0" json:"total_play"`
}

func (Workshop) TableName() string {
	return "workshops"
}

type Event struct {
	ID           uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string   `gorm:"size:50;not null" json:"title"`
	SubTitle     *string  `gorm:"size:50" json:"sub_title"`
	Description  *string  `gorm:"type:text" json:"description"`
	Poster       string   `gorm:"type:text;not null" json:"poster"`
	StartDate    *string  `gorm:"type:date" json:"start_date"`
	StartTime    *string  `gorm:"type:varchar(10)" json:"start_time"`
	Place        string   `gorm:"size:255;not null" json:"place"`
	AdmissionFee *float64 `gorm:"type:double" json:"admission_fee"`
	OriginLink   *string  `gorm:"type:text" json:"origin_link"`
	ShortLink    *string  `gorm:"type:text" json:"short_link"`
	IsComingSoon bool     `gorm:"not null" json:"is_coming_soon"`
	IsActive     bool     `gorm:"not null" json:"is_active"`
}

func (Event) TableName() string {
	return "events"
}

func (e *Event) AfterFind(tx *gorm.DB) error {
	if e.StartTime != nil && len(*e.StartTime) > 5 {
		trimmed := (*e.StartTime)[:5]
		e.StartTime = &trimmed
	}
	return nil
}

type Setting struct {
	ID             uint   `gorm:"primaryKey" json:"-"`
	WhatsappNumber string `gorm:"size:30;not null" json:"whatsapp_number"`
	InstagramURL   string `gorm:"type:text;not null" json:"instagram_url"`
	TiktokURL      string `gorm:"type:text;not null" json:"tiktok_url"`
	ShopeeURL      string `gorm:"type:text;not null" json:"shopee_url"`
}

func (Setting) TableName() string {
	return "settings"
}

type DailyVisit struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	VisitDate string `gorm:"type:date;uniqueIndex;not null" json:"visit_date"`
	Visitors  int    `gorm:"not null;default:0" json:"visitors"`
}

func (DailyVisit) TableName() string {
	return "daily_visits"
}
