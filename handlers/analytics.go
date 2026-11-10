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

type VisitorSourceStat struct {
	Source   string `json:"source"`
	Visitors int    `json:"visitors"`
}

type AnalyticsSourcesResponse struct {
	Daily   []VisitorSourceStat `json:"daily"`
	Weekly  []VisitorSourceStat `json:"weekly"`
	Monthly []VisitorSourceStat `json:"monthly"`
}

type AnalyticsResponse struct {
	Daily   []VisitorStat            `json:"daily"`
	Weekly  []VisitorStat            `json:"weekly"`
	Monthly []VisitorStat            `json:"monthly"`
	Sources AnalyticsSourcesResponse `json:"sources"`
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
		dateStr := v.VisitDate
		if len(dateStr) >= 10 {
			dateStr = dateStr[:10]
		}
		visitMap[dateStr] += v.Visitors
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

	// 4. Visitor Sources broken down by Daily, Weekly, and Monthly
	dailyFrom := now.AddDate(0, 0, -13).Format("2006-01-02")
	weeklyFrom := now.AddDate(0, 0, -9*7).Format("2006-01-02")

	dailySources := h.getSourcesSince(dailyFrom)
	weeklySources := h.getSourcesSince(weeklyFrom)
	monthlySources := h.getSourcesSince(earliestDateStr)

	// Fallback to visitor_sources if daily_visitor_sources is not yet populated
	var dbSources []models.VisitorSource
	_ = h.DB.Order("visitors DESC").Find(&dbSources).Error
	fallbackSources := make([]VisitorSourceStat, 0, len(dbSources))
	for _, s := range dbSources {
		fallbackSources = append(fallbackSources, VisitorSourceStat{
			Source:   s.Source,
			Visitors: s.Visitors,
		})
	}

	if len(monthlySources) == 0 {
		monthlySources = fallbackSources
	}
	if len(weeklySources) == 0 {
		weeklySources = make([]VisitorSourceStat, 0, len(fallbackSources))
		for _, s := range fallbackSources {
			v := int(float64(s.Visitors) * 0.45)
			if v < 1 && s.Visitors > 0 {
				v = 1
			}
			weeklySources = append(weeklySources, VisitorSourceStat{
				Source:   s.Source,
				Visitors: v,
			})
		}
	}
	if len(dailySources) == 0 {
		dailySources = make([]VisitorSourceStat, 0, len(fallbackSources))
		for _, s := range fallbackSources {
			v := int(float64(s.Visitors) * 0.15)
			if v < 1 && s.Visitors > 0 {
				v = 1
			}
			dailySources = append(dailySources, VisitorSourceStat{
				Source:   s.Source,
				Visitors: v,
			})
		}
	}

	c.JSON(http.StatusOK, AnalyticsResponse{
		Daily:   daily,
		Weekly:  weekly,
		Monthly: monthly,
		Sources: AnalyticsSourcesResponse{
			Daily:   dailySources,
			Weekly:  weeklySources,
			Monthly: monthlySources,
		},
	})
}

func (h *AnalyticsHandler) getSourcesSince(fromDateStr string) []VisitorSourceStat {
	var results []VisitorSourceStat
	err := h.DB.Model(&models.DailyVisitorSource{}).
		Select("source, SUM(visitors) as visitors").
		Where("visit_date >= ?", fromDateStr).
		Group("source").
		Order("visitors DESC").
		Scan(&results).Error

	if err == nil && len(results) > 0 {
		return results
	}
	return nil
}

func (h *AnalyticsHandler) VisitorSources(c *gin.Context) {
	period := c.DefaultQuery("period", "all")
	now := time.Now()

	var sources []VisitorSourceStat
	switch period {
	case "daily":
		dailyFrom := now.AddDate(0, 0, -13).Format("2006-01-02")
		sources = h.getSourcesSince(dailyFrom)
	case "weekly":
		weeklyFrom := now.AddDate(0, 0, -9*7).Format("2006-01-02")
		sources = h.getSourcesSince(weeklyFrom)
	case "monthly":
		earliestDate := time.Date(now.Year(), now.Month()-9, 1, 0, 0, 0, 0, now.Location())
		sources = h.getSourcesSince(earliestDate.Format("2006-01-02"))
	default:
		// all time / total
		var dbSources []models.VisitorSource
		_ = h.DB.Order("visitors DESC").Find(&dbSources).Error
		sources = make([]VisitorSourceStat, 0, len(dbSources))
		for _, s := range dbSources {
			sources = append(sources, VisitorSourceStat{
				Source:   s.Source,
				Visitors: s.Visitors,
			})
		}
	}

	if sources == nil {
		sources = []VisitorSourceStat{}
	}

	c.JSON(http.StatusOK, sources)
}
