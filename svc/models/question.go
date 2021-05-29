package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	UserID       uint   `json:"-"`
	QuizID       uint   `json:"-"`
	Points       uint   `json:"points,omitempty"`
	Text         string `json:"text,omitempty"`
	ImageLink    string `json:"image_link,omitempty"`
	AudioLink    string `json:"audio_link,omitempty"`
	Answer       string `json:"answer,omitempty"`
	TimerSeconds uint   `json:"timer_seconds,omitempty"`
	Tags         []*Tag `gorm:"many2many:question_tags" json:"tags,omitempty"`
}

func NewQuestion(quizID uint, text, imageLink, audioLink, answer string, points, timerSeconds uint) *Question {
	return &Question{
		QuizID:       quizID,
		Text:         text,
		ImageLink:    imageLink,
		AudioLink:    audioLink,
		Answer:       answer,
		Points:       points,
		TimerSeconds: timerSeconds,
	}
}
