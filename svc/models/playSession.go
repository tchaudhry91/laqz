package models

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

const StateInitialized = "INITIALIZED"
const StateInProgress = "INPROGRESS"
const StateFinished = "FINISHED"

type PlaySession struct {
	gorm.Model           `json:"-"`
	Code                 uint `gorm:"uniqueIndex" json:"code"`
	QuizID               uint
	Quiz                 *Quiz     `json:"quiz"`
	State                string    `json:"state"`
	CurrentQuestionIndex int       `json:"-"`
	CurrentQuestion      *Question `gorm:"-" json:"current_question"`
	QuizMaster           string    `json:"quiz_master"`
	Users                []*User   `gorm:"many2many:session_users" json:"users"`
	Teams                []*Team   `gorm:"many2many:session_teams" json:"teams"`
}

func NewPlaySession(qm string, q *Quiz) (s *PlaySession) {
	rand.Seed(time.Now().UTC().Unix())
	return &PlaySession{
		Code:                 10000 + uint(rand.Intn(89999)),
		CurrentQuestionIndex: 0,
		Quiz:                 q,
		QuizMaster:           qm,
		State:                StateInitialized,
		Users:                []*User{},
		Teams:                []*Team{},
	}
}

func (s *PlaySession) AddUser(u *User) {
	s.Users = append(s.Users, u)
}

func (s *PlaySession) CreateTeam(t *Team) {
	s.Teams = append(s.Teams, t)
}

func (s *PlaySession) UpdateQuestion(q *Question) {
	s.CurrentQuestion = q
}

func (s *PlaySession) SetFinished() {
	s.State = StateFinished
}

func (s *PlaySession) SetInProgress() {
	s.State = StateInProgress
}