package services

import (
	"belchi/models"
	"belchi/repository"
	"belchi/security"
	"belchi/validation"

	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	validator *validation.UserValidator
	hasher    *security.StringHasher
	jwt       *security.TokenJWT
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:  repository.NewUserRepository(),
		validator: validation.NewUserValidator(),
		hasher:    security.NewStringHasher(),
		jwt:       security.NewTokenJWT(),
	}
}

func (s *AuthService) RegisterUser(username, email, password string) (*models.User, string, error) {
	if err := s.validator.ValidateRegistration(username, email, password); err != nil {
		return nil, "", fmt.Errorf("invalid credentials: %s", err)
	}

	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, "", fmt.Errorf("error checking email existence: %s", err)
	}
	if exists {
		return nil, "", errors.New("email already registered")
	}

	exists, err = s.userRepo.UsernameExists(username)
	if err != nil {
		return nil, "", fmt.Errorf("error checking username existence: %s", err)
	}
	if exists {
		return nil, "", errors.New("username already registered")
	}

	hashedPassword, err := s.hasher.MakeHash(password)
	if err != nil {
		return nil, "", fmt.Errorf("password hashing failed: %s", err)
	}

	user := &models.User{
		Username:        username,
		Email:           email,
		Password:        hashedPassword,
		PasswordVersion: 1,
	}

	if err := s.userRepo.RegisterUser(user); err != nil {
		return nil, "", fmt.Errorf("user registration failed: %s", err)
	}

	token, err := s.jwt.GenerateTokenJWT(email, password)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil

}

func (s *AuthService) LoginUser(email, password string) (*models.User, string, error) {
	if err := s.validator.ValidateLogin(email, password); err != nil {
		return nil, "", fmt.Errorf("invalid credentials: %s", err)
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("error getting user: %s", err)
	}

	if user == nil {
		return nil, "", fmt.Errorf("user does not exist")
	}

	match, err := s.hasher.CompareHash(user.Password, password)

	if err != nil {
		return nil, "", fmt.Errorf("error compare password hash: %s", err)
	}

	if !match {
		return nil, "", fmt.Errorf("incorrect password")
	}

	token, err := s.jwt.GenerateTokenJWT(email, password)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GetPasswordByEmail(email string) (string, error) {
	if err := s.validator.ValidateEmail(email); err != nil {
		return "", fmt.Errorf("invalid email: %s", err)
	}

	password, err := s.userRepo.GetPasswordByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", fmt.Errorf("error getting password: %s", err)
	}

	return password, nil
}

func (s *AuthService) CompareHashAndPassword(email, plain string) (bool, error) {
	currentPassword, err := s.GetPasswordByEmail(email)
	if err != nil {
		return false, err
	}

	match, err := s.hasher.CompareHash(currentPassword, plain)

	if err != nil {
		return false, err
	}

	if !match {
		return false, nil
	}

	return true, nil
}
