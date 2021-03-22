package models

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Name          string `gorm:"uniqueIndex"`
	Private       bool
	Collaborators []User `gorm:"many2many:User"`
	Tags          []string
	Questions     []Question `gorm:"many2many:questions"`
}
