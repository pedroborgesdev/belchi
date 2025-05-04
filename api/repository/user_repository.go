package repository

import (
	"belchi/database"
	"belchi/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: database.DB,
	}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
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

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("username =?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) RegisterUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) ChangePassword(user *models.User, hashedPassword string) error {
	result := r.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("password", hashedPassword).
		Update("password_version", user.PasswordVersion+1)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) GetPasswordByEmail(email string) (string, error) {
	var user models.User
	err := r.DB.Where("email =?", email).Select("password").First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Password, nil
}
