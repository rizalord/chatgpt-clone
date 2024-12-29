package entity

import "time"

type User struct {
	ID        uint   
	Name      string 
	Email     string 
	ImageURL  *string `gorm:"column:image_url"`
	Password  *string 
	CreatedAt time.Time
	UpdatedAt time.Time
}
