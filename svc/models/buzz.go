package models

import (
	"time"

	"gorm.io/gorm"
)

type Buzz struct {
	gorm.Model `json:"-"`
	TS         time.Time `json:"ts,omitempty"`
	UserID     uint      `json:"user_id,omitempty"`
	QuestionID uint      `json:"question_id,omitempty"`
	Question   Question  `json:"question,omitempty"`
}
