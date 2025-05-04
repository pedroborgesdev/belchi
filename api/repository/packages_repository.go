package repository

import (
	"belchi/config"
	"belchi/database"
	"belchi/models"
	"errors"
	"io"
	"mime/multipart"
	"os"

	"gorm.io/gorm"
)

type PackagesRepository struct {
	DB *gorm.DB
}

func NewPackagesRepository() *PackagesRepository {
	return &PackagesRepository{
		DB: database.DB,
	}
}

func (r *PackagesRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := r.DB.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *PackagesRepository) GetUserByName(name string) (*models.User, error) {
	var user models.User

	result := r.DB.Where("username = ?", name).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *PackagesRepository) UploadPackages(packages *models.Packages, file *multipart.FileHeader) error {
	filePath := config.AppConfig.UploadPath + file.Filename

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	packages.Path = filePath

	if err := r.DB.Create(packages).Error; err != nil {
		return err
	}

	return nil
}

func (r *PackagesRepository) GetPackageWithAuthorAndNameAndVersion(author uint, name, version string) (*models.Packages, error) {
	var packages models.Packages

	result := r.DB.Where("author_id =? AND name =? AND version =?", author, name, version).First(&packages)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &packages, nil
}

func (r *PackagesRepository) GetPackageWithName(name string) (*models.Packages, error) {
	var packages models.Packages

	result := r.DB.Where("name =?", name).Find(&packages)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &packages, nil
}

func (r *PackagesRepository) DownloadPackages(packages *models.Packages) (*os.File, error) {
	file, err := os.Open(packages.Path)
	if err != nil {
		return nil, err
	}

	return file, nil
}
