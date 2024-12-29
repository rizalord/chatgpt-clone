package entity

import "time"

type Chat struct {
	ID        	uint   
	UserID		uint
	Topic 	 	*string
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	Messages	[]Message `gorm:"foreignKey:ChatID"`
}
