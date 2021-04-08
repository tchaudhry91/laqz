package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model `json:"-"`
	UserID     uint   `json:"user_id,omitempty"`
	QuizID     uint   `json:"quiz_id,omitempty"`
	Points     uint64 `json:"points,omitempty"`
	Text       string `json:"text,omitempty"`
	ImageLink  string `json:"image_link,omitempty"`
	AudioLink  string `json:"audio_link,omitempty"`
	Answer     string `json:"answer,omitempty"`
	Tags       []*Tag `gorm:"many2many:question_tags" json:"tags,omitempty"`
}
