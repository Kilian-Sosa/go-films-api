package domain

import "time"

type Film struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null"`
	Title       string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Director    string
	ReleaseDate time.Time
	Cast        string
	Genre       string
	Synopsis    string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User User `gorm:"foreignKey:UserID"`
}
