package svc

import (
	"context"

	"github.com/tchaudhry91/laqz/svc/models"
)

type PlaySessionSVC interface {
	InitNewPS(ctx context.Context, quizID uint) (s *models.PlaySession, err error)
	StartPS(ctx context.Context, code uint) (err error)
	AddUserToPS(ctx context.Context, code uint) (err error)
	AddUserToTeam(ctx context.Context, code uint, teamName string, email string) (err error)
	AddTeamToPS(ctx context.Context, code uint, team *models.Team) (err error)
	GetPS(ctx context.Context, code uint) (s *models.PlaySession, err error)
}

type PlaySessionSvc struct {
	db models.QuizStore
}

func NewPlaySessionSvc(db models.QuizStore) *PlaySessionSvc {
	return &PlaySessionSvc{
		db: db,
	}
}

func (ps *PlaySessionSvc) UserContextKey() contextKey {
	var userContextKey = contextKey("user")
	return userContextKey
}

func (ps *PlaySessionSvc) InitNewPS(ctx context.Context, quizID uint) (s *models.PlaySession, err error) {
	u, err := getUserFromContext(ctx, ps.UserContextKey())
	if err != nil {
		return s, err
	}
	u, err = ps.db.GetUserByEmail(u.Email)
	if err != nil {
		return s, err
	}
	// Get Quiz
	qz, err := ps.db.GetQuiz(quizID)
	if err != nil {
		return s, err
	}
	if !qz.CanView(u.Email) {
		return s, NotPermittedError
	}
	s = models.NewPlaySession(u.Email, qz)
	err = ps.db.CreatePlaySession(s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (ps *PlaySessionSvc) StartPS(ctx context.Context, code uint) (err error) {
	u, err := getUserFromContext(ctx, ps.UserContextKey())
	if err != nil {
		return err
	}
	s, err := ps.db.GetPlaySession(code)
	if err != nil {
		return err
	}
	if s.QuizMaster != u.Email {
		return NotPermittedError
	}
	s.SetInProgress()

	// SetQuestion
	qqs, err := ps.db.GetQuestionsByQuiz(s.Quiz.ID)
	if err != nil {
		return err
	}
	s.UpdateQuestion(qqs[s.CurrentQuestionIndex])
	return ps.db.UpdatePlaySession(s)
}

func (ps *PlaySessionSvc) GetPS(ctx context.Context, code uint) (s *models.PlaySession, err error) {
	return ps.db.GetPlaySession(code)
}

func (ps *PlaySessionSvc) AddUserToPS(ctx context.Context, code uint) (err error) {
	u, err := getUserFromContext(ctx, ps.UserContextKey())
	if err != nil {
		return err
	}
	u, err = ps.db.GetUserByEmail(u.Email)
	if err != nil {
		return err
	}
	s, err := ps.db.GetPlaySession(code)
	if err != nil {
		return err
	}
	s.AddUser(u)
	return ps.db.UpdatePlaySession(s)
}

func (ps *PlaySessionSvc) AddTeamToPS(ctx context.Context, code uint, t *models.Team) (err error) {
	u, err := getUserFromContext(ctx, ps.UserContextKey())
	if err != nil {
		return err
	}
	s, err := ps.db.GetPlaySession(code)
	if err != nil {
		return err
	}
	if s.QuizMaster != u.Email {
		return NotPermittedError
	}
	s.AddTeam(t)
	return ps.db.UpdatePlaySession(s)
}

func (ps *PlaySessionSvc) AddUserToTeam(ctx context.Context, code uint, teamName string, email string) (err error) {
	u, err := getUserFromContext(ctx, ps.UserContextKey())
	if err != nil {
		return err
	}
	s, err := ps.db.GetPlaySession(code)
	if err != nil {
		return err
	}
	if email != u.Email {
		return NotPermittedError
	}
	s.AssignUserToTeam(teamName, email)
	return ps.db.UpdatePlaySession(s)
}
