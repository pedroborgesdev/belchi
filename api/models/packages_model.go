package models

import "time"

type Packages struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	AuthorID uint   `json:"author_id" gorm:"not null"`
	Author   User   `json:"author" gorm:"foreignKey:AuthorID"`
	Name     string `json:"name" gorm:"size:255;not null"`
	Version  string `json:"version" gorm:"size:255;not null"`
	Path     string `json:"path" gorm:"size:255;not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
