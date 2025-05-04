package controllers

import (
	"belchi/logger"
	"belchi/services"
	"belchi/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (c *AuthController) RegisterUser(ctx *gin.Context) {
	var request utils.RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, token, err := c.authService.RegisterUser(request.Username, request.Email, request.Password)

	if err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})

		logger.Log("ERROR", "Registration failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Email", Value: request.Email},
		})
		return
	}

	utils.Success(ctx, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"token":    token,
	})

	logger.Log("INFO", "User registered successfully", []logger.LogDetail{
		{Key: "Username", Value: user.Username},
		{Key: "Email", Value: user.Email},
		{Key: "Token", Value: token},
	})
}

func (c *AuthController) LoginUser(ctx *gin.Context) {
	var request utils.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, token, err := c.authService.LoginUser(request.Email, request.Password)

	if err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})

		logger.Log("ERROR", "Login failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Email", Value: request.Email},
		})
		return
	}

	utils.Success(ctx, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"token":    token,
	})

	logger.Log("INFO", "User logged in successfully", []logger.LogDetail{
		{Key: "Username", Value: user.Username},
		{Key: "Email", Value: user.Email},
		{Key: "Token", Value: token},
	})
}
