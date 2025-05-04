package controllers

import (
	"belchi/logger"
	"belchi/services"
	"belchi/utils"
	"fmt"
	"io"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type PackagesController struct {
	packagesService *services.PackagesService
}

func NewPackagesController() *PackagesController {
	return &PackagesController{
		packagesService: services.NewPackagesService(),
	}
}

func (c *PackagesController) UploadPackages(ctx *gin.Context) {
	var request utils.UploadPackagesRequest

	if err := ctx.ShouldBind(&request); err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": "file is required",
		})
		return
	}

	version := ctx.PostForm("version")
	if version == "" {
		utils.BadRequest(ctx, gin.H{
			"error": "version is required",
		})
		return
	}

	name := ctx.PostForm("name")
	if name == "" {
		utils.BadRequest(ctx, gin.H{
			"error": "name is required",
		})
		return
	}

	tokenEmail, exists := ctx.Get("userEmail")
	if !exists {
		utils.BadRequest(ctx, gin.H{
			"error": "email not found on token",
		})
		return
	}

	packageName, err := c.packagesService.UploadPackages(tokenEmail.(string), name, version, file)
	if err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})
		logger.Log("ERROR", "Upload failed", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Email", Value: tokenEmail.(string)},
		})
		return
	}

	utils.Success(ctx, gin.H{
		"message": "upload successful",
		"package": packageName,
	})
}

func (c *PackagesController) DownloadPackages(ctx *gin.Context) {
	packageName := ctx.Param("name")
	version := ctx.Param("version")

	if packageName == "" {
		utils.BadRequest(ctx, gin.H{
			"error": "package name is required",
		})
		return
	}

	if version == "" {
		utils.BadRequest(ctx, gin.H{
			"error": "version is required",
		})
		return
	}

	file, err := c.packagesService.DownloadPackages(packageName, version)
	if err != nil {
		utils.BadRequest(ctx, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(file.Name())))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")

	io.Copy(ctx.Writer, file)
}
