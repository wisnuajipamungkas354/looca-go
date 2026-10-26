package handlers

import (
	"fmt"
	"net/http"
	"time"

	"looca-go/models"
	"looca-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsHandler struct {
	DB *gorm.DB
}

func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	return &AnalyticsHandler{DB: db}
}

type VisitorStat struct {
	Label    string `json:"label"`
	Visitors int    `json:"visitors"`
}

type AnalyticsResponse struct {
	Daily   []VisitorStat `json:"daily"`
	Weekly  []VisitorStat `json:"weekly"`
	Monthly []VisitorStat `json:"monthly"`
}

type DashboardStatsResponse struct {
	TotalProducts  int64 `json:"total_products"`
	TotalWorkshops int64 `json:"total_workshops"`
	TotalEvents    int64 `json:"total_events"`
}

func (h *AnalyticsHandler) Stats(c *gin.Context) {
	var totalProducts int64
	var totalWorkshops int64
	var totalEvents int64

	if err := h.DB.Model(&models.Product{}).Count(&totalProducts).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to count products")
		return
	}
	if err := h.DB.Model(&models.Workshop{}).Count(&totalWorkshops).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to count workshops")
		return
	}
	if err := h.DB.Model(&models.Event{}).Count(&totalEvents).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to count events")
		return
	}

	c.JSON(http.StatusOK, DashboardStatsResponse{
		TotalProducts:  totalProducts,
		TotalWorkshops: totalWorkshops,
		TotalEvents:    totalEvents,
	})
}

func (h *AnalyticsHandler) Analytics(c *gin.Context) {
	now := time.Now()

	// Earliest date for 10 months history: first day of month (now - 9 months)
	earliestDate := time.Date(now.Year(), now.Month()-9, 1, 0, 0, 0, 0, now.Location())
	earliestDateStr := earliestDate.Format("2006-01-02")

	var visits []models.DailyVisit
	if err := h.DB.Where("visit_date >= ?", earliestDateStr).Find(&visits).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch analytics")
		return
	}

	// Map visit_date -> visitors count
	visitMap := make(map[string]int)
	for _, v := range visits {
		visitMap[v.VisitDate] = v.Visitors
	}

	// 1. Daily: last 14 days (inclusive of today), oldest to newest
	daily := make([]VisitorStat, 0, 14)
	for i := 13; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		dateStr := day.Format("2006-01-02")
		label := day.Format("2 Jan")
		daily = append(daily, VisitorStat{
			Label:    label,
			Visitors: visitMap[dateStr],
		})
	}

	// 2. Weekly: last 10 ISO weeks, oldest to newest
	// To get the last 10 weeks, we can step week by week back 9 weeks to now
	weekly := make([]VisitorStat, 0, 10)
	for i := 9; i >= 0; i-- {
		// Target week reference point
		targetRef := now.AddDate(0, 0, -i*7)
		targetYear, targetWeek := targetRef.ISOWeek()
		label := fmt.Sprintf("W%d", targetWeek)

		// Sum all visits in visitMap that have the same ISO year and week
		weekVisitors := 0
		for dateStr, count := range visitMap {
			t, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				y, w := t.ISOWeek()
				if y == targetYear && w == targetWeek {
					weekVisitors += count
				}
			}
		}

		weekly = append(weekly, VisitorStat{
			Label:    label,
			Visitors: weekVisitors,
		})
	}

	// 3. Monthly: last 10 months, oldest to newest
	monthly := make([]VisitorStat, 0, 10)
	for i := 9; i >= 0; i-- {
		monthRef := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, now.Location())
		targetYear, targetMonth := monthRef.Year(), monthRef.Month()
		label := monthRef.Format("Jan")

		monthVisitors := 0
		for dateStr, count := range visitMap {
			t, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				if t.Year() == targetYear && t.Month() == targetMonth {
					monthVisitors += count
				}
			}
		}

		monthly = append(monthly, VisitorStat{
			Label:    label,
			Visitors: monthVisitors,
		})
	}

	c.JSON(http.StatusOK, AnalyticsResponse{
		Daily:   daily,
		Weekly:  weekly,
		Monthly: monthly,
	})
}
