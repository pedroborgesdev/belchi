package main

import (
	"belchi/config"
	"belchi/database"
	"belchi/debug"
	"belchi/logger"
	"belchi/middlewares"
	"belchi/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	debug.LoadDebugConfig()

	logger.Log("INFO", "Application has been started", []logger.LogDetail{})

	config.LoadAppConfig()

	database.InitDB()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(
		middlewares.CORSMiddleware(),
	)

	routes.SetupRoutes(router)

	router.Run(":" + config.AppConfig.HTTPPort)
}
