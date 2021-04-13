package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Email      string  `gorm:"uniqueIndex" json:"email,omitempty"`
	Name       string  `json:"name,omitempty"`
	Quizzes    []*Quiz `gorm:"many2many:quiz_collaborators" json:"quizzes,omitempty"`
	Teams      []*Team `gorm:"many2many:user_teams" json:"teams,omitempty"`
	AvatarURL  string  `json:"avatar_url,omitempty"`
}
