package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"location-tracking-shortlink/config"
	"location-tracking-shortlink/db"
	"location-tracking-shortlink/handlers"
	"location-tracking-shortlink/middleware"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := db.Init(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(gin.Logger())

	r.SetFuncMap(template.FuncMap{
		"formatDate": formatDate,
	})

	r.LoadHTMLGlob("templates/*.html")

	r.Static("/static", "./static")

	apiHandler := handlers.NewAPIHandler()
	shortlinkHandler := handlers.NewShortlinkHandler()
	visitHandler := handlers.NewVisitHandler()
	adminHandler := handlers.NewAdminHandler()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/admin")
	})

	r.GET("/health", apiHandler.HealthCheck)

	r.GET("/:code", shortlinkHandler.Redirect)

	api := r.Group("/api")
	{
		api.POST("/shortlinks/create", shortlinkHandler.Create)
		api.GET("/shortlinks/templates", shortlinkHandler.GetTemplates)
		api.POST("/visits/submit", visitHandler.Submit)
		api.POST("/visits/duration", visitHandler.UpdateDuration)
		api.POST("/visits/heartbeat", func(c *gin.Context) {
			c.JSON(200, gin.H{"code": 0, "message": "ok"})
		})
		api.GET("/visits", visitHandler.List)
	}

	admin := r.Group("/api/admin")
	{
		admin.GET("/shortlinks", adminHandler.ListShortlinks)
		admin.DELETE("/shortlinks/:id", adminHandler.DeleteShortlink)
		admin.PUT("/shortlinks/:id/toggle", adminHandler.ToggleShortlink)
		admin.GET("/stats", adminHandler.GetStats)
		admin.GET("/export", adminHandler.ExportCSV)
	}

	r.GET("/admin", adminHandler.AdminPage)
	r.GET("/admin/visits", adminHandler.AdminPage)
	r.GET("/stats", adminHandler.StatsPage)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}

func formatDate(t interface{}) string {
	return "formatted"
}
