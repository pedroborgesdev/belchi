package models

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"size:255;not null;unique"`
	Email    string `json:"email" gorm:"size:255;unique;not null"`
	Password string `json:"-" gorm:"size:255;not null"`

	CreatedAt       time.Time `json:"created_at"`
	PasswordVersion int       `json:"password_version"`

	Packages []Packages `gorm:"foreignKey:AuthorID"`
}
