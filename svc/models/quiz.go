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

// NewQuiz is used to initialize an empty Quiz
func NewQuiz(name string, owner *User) *Quiz {
	return &Quiz{
		Name:          name,
		Collaborators: []*User{owner},
		Private:       true,
	}
}

// AddTags attaches tags to a given quiz
func (qz *Quiz) AddTags(tags []*Tag) {
	qz.Tags = append(qz.Tags, tags...)
}

func (qz *Quiz) AddQuestion(q *Question) {
	qz.Questions = append(qz.Questions, q)
}

// Publish makes a quiz public
func (qz *Quiz) Publish() {
	qz.Private = false
}

func (qz *Quiz) CanView(email string) bool {
	for i := range qz.Collaborators {
		if qz.Collaborators[i].Email == email {
			return true
		}
	}
	return false
}
