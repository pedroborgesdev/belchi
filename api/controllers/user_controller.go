package controllers

import (
	"belchi/logger"
	"belchi/services"
	"belchi/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	var request utils.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})
		return
	}

	tokenEmail, exists := ctx.Get("userEmail")

	if !exists {
		utils.BadRequest(ctx, gin.H{
			"error": "no information for this user",
		})
		return
	}

	result, err := c.userService.ChangePassword(tokenEmail.(string), request.Password, request.NewPassword)
	if err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})

		logger.Log("ERROR", "Change password failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Email", Value: tokenEmail.(string)},
			{Key: "Password", Value: request.Password},
			{Key: "New Password", Value: request.NewPassword},
		})
		return
	}

	if !result {
		utils.BadRequest(ctx, gin.H{
			"error": "password cannot update",
		})

		logger.Log("ERROR", "Change password failed", []logger.LogDetail{
			{Key: "Status", Value: result},
			{Key: "Email", Value: tokenEmail.(string)},
			{Key: "Password", Value: request.Password},
			{Key: "New Password", Value: request.NewPassword},
		})
	}

	utils.Success(ctx, gin.H{
		"email":        tokenEmail.(string),
		"new_password": request.NewPassword,
	})

	logger.Log("DEBUG", "Change password successfuly", []logger.LogDetail{
		{Key: "Status", Value: result},
		{Key: "Email", Value: tokenEmail.(string)},
		{Key: "Password", Value: request.Password},
		{Key: "New Password", Value: request.NewPassword},
	})
}
