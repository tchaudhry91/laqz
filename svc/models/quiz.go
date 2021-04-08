package models

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Name          string      `gorm:"uniqueIndex" json:"name,omitempty"`
	Private       bool        `json:"private,omitempty"`
	Collaborators []*User     `gorm:"many2many:quiz_collaborators" json:"collaborators,omitempty"`
	Tags          []*Tag      `gorm:"many2many:quiz_tags" json:"tags,omitempty"`
	Questions     []*Question `gorm:"many2many:quiz_questions" json:"questions,omitempty"`
}

// NewQuiz is used to initialize an empty Quiz
func NewQuiz(name string, owner *User, tt []*Tag) *Quiz {
	return &Quiz{
		Name:          name,
		Collaborators: []*User{owner},
		Tags:          tt,
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
	if !qz.Private {
		return true
	}
	return qz.IsCollaborator(email)
}

func (qz *Quiz) IsCollaborator(email string) bool {
	for i := range qz.Collaborators {
		if qz.Collaborators[i].Email == email {
			return true
		}
	}
	return false
}
