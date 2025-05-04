package services

import (
	"belchi/models"
	"belchi/repository"
	"belchi/security"
	"belchi/utils"
	"belchi/validation"
	"fmt"
	"mime/multipart"
	"os"
)

type PackagesService struct {
	packagesRepo *repository.PackagesRepository
	validator    *validation.PackagesValidator
	jwt          *security.TokenJWT
}

func NewPackagesService() *PackagesService {
	return &PackagesService{
		packagesRepo: repository.NewPackagesRepository(),
		validator:    validation.NewPackagesValidator(),
		jwt:          security.NewTokenJWT(),
	}
}

func (c *PackagesService) UploadPackages(email, name, version string, file *multipart.FileHeader) (string, error) {
	if err := c.validator.ValidatePackageUpload(name, version); err != nil {
		return "", err
	}

	user, err := c.packagesRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	packageExists, err := c.packagesRepo.GetPackageWithAuthorAndNameAndVersion(user.ID, name, version)
	if err != nil {
		return "", err
	}

	if packageExists != nil {
		return "", fmt.Errorf("package already exists with the same version")
	}

	packages := &models.Packages{
		Name:     name,
		Version:  version,
		AuthorID: user.ID,
	}

	err = c.packagesRepo.UploadPackages(packages, file)
	if err != nil {
		return "", err
	}

	packageRef := fmt.Sprintf("%s@%s", user.Username, name)

	return packageRef, nil
}

func (c *PackagesService) DownloadPackages(name, version string) (*os.File, error) {
	if err := c.validator.ValidatePackageDownload(name, version); err != nil {
		return nil, err
	}

	packAuthor, packName := utils.SplitPackageName(name)

	authorID, err := c.packagesRepo.GetUserByName(packAuthor)
	if err != nil {
		return nil, err
	}

	if authorID == nil {
		return nil, fmt.Errorf("package not found")
	}

	packages, err := c.packagesRepo.GetPackageWithAuthorAndNameAndVersion(authorID.ID, packName, version)
	if err != nil {
		return nil, err
	}

	if packages == nil {
		return nil, fmt.Errorf("package not found")
	}

	return c.packagesRepo.DownloadPackages(packages)
}
