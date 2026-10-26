package routes

import (
	"net/http"

	"looca-go/config"
	"looca-go/handlers"
	"looca-go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	r.Use(middleware.CORS(cfg.CORSOrigin))
	r.MaxMultipartMemory = 8 << 20
	r.Static("/uploads", cfg.UploadDir)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	categoryHandler := handlers.NewCategoryHandler(db)
	productHandler := handlers.NewProductHandler(db)
	workshopHandler := handlers.NewWorkshopHandler(db)
	eventHandler := handlers.NewEventHandler(db)
	settingsHandler := handlers.NewSettingsHandler(db)
	analyticsHandler := handlers.NewAnalyticsHandler(db)
	uploadHandler := handlers.NewUploadHandler(cfg.UploadDir)
	publicHandler := handlers.NewPublicHandler(db)

	// Short link redirects for events
	r.GET("/event/:slug", publicHandler.RedirectEvent)
	r.GET("/r/:slug", publicHandler.RedirectEvent)

	api := r.Group("/api")
	{
		// Auth
		api.POST("/auth/login", authHandler.Login)

		// Public endpoints
		pub := api.Group("/public")
		{
			pub.GET("/categories", publicHandler.Categories)
			pub.GET("/products", publicHandler.Products)
			pub.GET("/workshops", publicHandler.Workshops)
			pub.GET("/events", publicHandler.Events)
			pub.GET("/settings", publicHandler.Settings)
			pub.POST("/products/:id/click", publicHandler.ClickProduct)
			pub.POST("/workshops/:id/play", publicHandler.PlayWorkshop)
			pub.POST("/visit", publicHandler.Visit)
		}

		// Admin endpoints (protected)
		adm := api.Group("/admin", middleware.Auth(cfg.JWTSecret))
		{
			adm.GET("/me", authHandler.Me)
			adm.GET("/stats", analyticsHandler.Stats)
			adm.GET("/analytics", analyticsHandler.Analytics)
			adm.GET("/settings", settingsHandler.Get)
			adm.PUT("/settings", settingsHandler.Update)

			adm.POST("/uploads/image", uploadHandler.UploadImage)
			adm.POST("/uploads/video", uploadHandler.UploadVideo)

			// Categories
			adm.GET("/categories", categoryHandler.List)
			adm.POST("/categories", categoryHandler.Create)
			adm.PUT("/categories/:id", categoryHandler.Update)
			adm.DELETE("/categories/:id", categoryHandler.Delete)

			// Products
			adm.GET("/products", productHandler.List)
			adm.POST("/products", productHandler.Create)
			adm.GET("/products/:id", productHandler.Get)
			adm.PUT("/products/:id", productHandler.Update)
			adm.DELETE("/products/:id", productHandler.Delete)

			// Workshops
			adm.GET("/workshops", workshopHandler.List)
			adm.POST("/workshops", workshopHandler.Create)
			adm.GET("/workshops/:id", workshopHandler.Get)
			adm.PUT("/workshops/:id", workshopHandler.Update)
			adm.DELETE("/workshops/:id", workshopHandler.Delete)

			// Events
			adm.GET("/events", eventHandler.List)
			adm.POST("/events", eventHandler.Create)
			adm.GET("/events/:id", eventHandler.Get)
			adm.PUT("/events/:id", eventHandler.Update)
			adm.DELETE("/events/:id", eventHandler.Delete)
		}
	}
}
