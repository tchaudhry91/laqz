package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model `json:"-"`
	Name       string `json:"name,omitempty"`
}
