package services

import (
	"belchi/repository"
	"belchi/security"
	"belchi/validation"

	"fmt"
)

type UserService struct {
	userRepo  *repository.UserRepository
	validator *validation.UserValidator
	hasher    *security.StringHasher
	jwt       *security.TokenJWT
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:  repository.NewUserRepository(),
		validator: validation.NewUserValidator(),
		hasher:    security.NewStringHasher(),
		jwt:       security.NewTokenJWT(),
	}
}

func (c *UserService) UserExists(email string) (bool, error) {
	return true, nil
}

func (s *UserService) ChangePassword(email, currentPassword, newPassword string) (bool, error) {
	if currentPassword == newPassword {
		return false, fmt.Errorf("current password and new password is match")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return false, fmt.Errorf("error getting user: %s", err)
	}

	if user == nil {
		return false, fmt.Errorf("user does not exist")
	}

	match, err := s.hasher.CompareHash(user.Password, currentPassword)
	if err != nil {
		return false, fmt.Errorf("error compare password hash: %s", err)
	}

	if !match {
		return false, fmt.Errorf("incorrect password")
	}

	hashedPassword, err := s.hasher.MakeHash(newPassword)
	if err != nil {
		return false, fmt.Errorf("password hashing failed: %s", err)
	}

	err = s.userRepo.ChangePassword(user, hashedPassword)

	if err != nil {
		return false, fmt.Errorf("failed to update password user: %s", err)
	}

	return true, nil
}
