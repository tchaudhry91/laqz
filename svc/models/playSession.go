package models

import (
	"fmt"
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
	CurrentQuestionIndex int       `json:"current_question_index"`
	CurrentQuestion      *Question `gorm:"-" json:"current_question"`
	CurrentAnswer        string    `json:"current_answer,omitempty"`
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

func (s *PlaySession) AddTeam(t *Team) {
	s.Teams = append(s.Teams, t)
}

func (s *PlaySession) UpdateQuestion(q *Question) {
	s.CurrentQuestion = q
}

func (s *PlaySession) SetCurrentAnswer(answer string) {
	s.CurrentAnswer = answer
}

func (s *PlaySession) ClearCurrentAnswer() {
	s.CurrentAnswer = ""
}

func (s *PlaySession) SetFinished() {
	s.State = StateFinished
}

func (s *PlaySession) GetTeam(name string) (t *Team, err error) {
	for i := range s.Teams {
		if s.Teams[i].Name == name {
			return s.Teams[i], nil
		}
	}
	return nil, fmt.Errorf("Team not found")

}

func (s *PlaySession) AddTeamPoints(points int, teamName string) (err error) {
	targetTeamIndex := 0
	found := false
	for i := range s.Teams {
		if s.Teams[i].Name == teamName {
			targetTeamIndex = i
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Team not found")
	}
	s.Teams[targetTeamIndex].AddPoints(points)
	return nil
}

func (s *PlaySession) SetInProgress() {
	s.State = StateInProgress
}

func (s *PlaySession) AssignUserToTeam(teamName string, email string) (err error) {
	var targetTeam *Team
	var targetUser *User
	for i := range s.Teams {
		if s.Teams[i].Name == teamName {
			targetTeam = s.Teams[i]
			break
		}
	}
	if targetTeam == nil {
		return fmt.Errorf("Team not found")
	}
	for i := range s.Users {
		if s.Users[i].Email == email {
			targetUser = s.Users[i]
			break
		}
	}
	if targetTeam == nil {
		return fmt.Errorf("Team not found")
	}
	targetTeam.Users = append(targetTeam.Users, targetUser)
	return nil
}
