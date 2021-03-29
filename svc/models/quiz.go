package models

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Name          string `gorm:"uniqueIndex"`
	Private       bool
	Collaborators []*User     `gorm:"many2many:quiz_collaborators"`
	Tags          []*Tag      `gorm:"many2many:quiz_tags"`
	Questions     []*Question `gorm:"many2many:quiz_questions"`
}
