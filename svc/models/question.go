package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	UserID    uint
	QuizID    uint
	Points    uint64
	Text      string
	ImageLink string
	AudioLink string
	Answer    string
	Tags      []*Tag `gorm:"many2many:question_tags"`
}
