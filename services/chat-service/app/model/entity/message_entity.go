package entity

import "time"

type Message struct {
	ID        	uint   
	ChatID		uint
	Chat 	 	Chat	`gorm:"foreignKey:ChatID"`
	UserID		uint
	Role		string
	Content		string
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}
