package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `gorm:"uniqueIndex"`
	Name      string
	Quizzes   []*Quiz `gorm:"many2many:quiz_collaborators"`
	AvatarURL string
}
