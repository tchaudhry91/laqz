package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name   string  `json:"name"`
	Users  []*User `gorm:"many2many:user_teams" json:"users"`
	Points uint    `json:"points"`
}

func NewTeam(name string) *Team {
	return &Team{
		Name: name,
	}
}
