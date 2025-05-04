package routes

import (
	"belchi/controllers"
	"belchi/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	packagesController := controllers.NewPackagesController()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	auth := router.Group("/auth")
	auth.Use(
		middlewares.RateLimiter(middlewares.RateLimiterConfig{
			RequestLimit:     15,
			WindowMinutes:    60,
			WarningThreshold: 3,
			BlockTimeout:     20,
		}),
	)
	{
		auth.POST("/register", authController.RegisterUser)
		auth.POST("/login", authController.LoginUser)
	}

	user := router.Group("/user")
	user.Use(
		middlewares.AuthMiddleware(),
		middlewares.RateLimiter(middlewares.RateLimiterConfig{
			RequestLimit:     15,
			WindowMinutes:    60,
			WarningThreshold: 3,
			BlockTimeout:     20,
		}),
	)

	sensitiveUser := user.Group("/")
	{
		sensitiveUser.PATCH("/password", userController.ChangePassword)
	}

	packages := router.Group("/packages")
	packages.Use(
		middlewares.RateLimiter(middlewares.RateLimiterConfig{
			RequestLimit:     15,
			WindowMinutes:    60,
			WarningThreshold: 3,
			BlockTimeout:     20,
		}),
	)

	packages.GET("/:name/:version", packagesController.DownloadPackages)

	authPackages := packages.Group("/")
	authPackages.Use(middlewares.AuthMiddleware())
	authPackages.POST("/", packagesController.UploadPackages)
}
