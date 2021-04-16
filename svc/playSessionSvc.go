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
	IncrementPSQuestion(ctx context.Context, code uint) (err error)
	DecrementPSQuestion(ctx context.Context, code uint) (err error)
	RevealPSCurrentAnswer(ctx context.Context, code uint) (err error)
	GetPS(ctx context.Context, code uint) (s *models.PlaySession, err error)
	UpdateTeamPoints(ctx context.Context, code uint, points int, teamName string) (err error)
	EndPlaySession(ctx context.Context, code uint) (err error)
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

	return ps.db.UpdatePlaySession(s)
}

func (ps *PlaySessionSvc) EndPlaySession(ctx context.Context, code uint) (err error) {
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
	s.SetFinished()

	return ps.db.UpdatePlaySession(s)
}

func (ps *PlaySessionSvc) IncrementPSQuestion(ctx context.Context, code uint) (err error) {
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
	// Increment Question
	qqs, err := ps.db.GetQuestionsByQuiz(s.Quiz.ID)
	if err != nil {
		return err
	}

	if s.CurrentQuestionIndex < len(qqs)-1 {
		s.CurrentQuestionIndex += 1
	}
	s.UpdateQuestion(qqs[s.CurrentQuestionIndex])
	s.ClearCurrentAnswer()
	return ps.db.UpdatePlaySession(s)
}

func (ps *PlaySessionSvc) DecrementPSQuestion(ctx context.Context, code uint) (err error) {
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
	// Increment Question
	qqs, err := ps.db.GetQuestionsByQuiz(s.Quiz.ID)
	if err != nil {
		return err
	}
	if s.CurrentQuestionIndex > 0 {
		s.CurrentQuestionIndex -= 1
	}
	s.UpdateQuestion(qqs[s.CurrentQuestionIndex])
	s.ClearCurrentAnswer()
	return ps.db.UpdatePlaySession(s)
}

func (ps *PlaySessionSvc) UpdateTeamPoints(ctx context.Context, code uint, points int, teamName string) (err error) {
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
	// Award Points
	t, err := s.GetTeam(teamName)
	if err != nil {
		return err
	}
	t.AddPoints(points)

	return ps.db.UpdateTeam(t)
}

func (ps *PlaySessionSvc) RevealPSCurrentAnswer(ctx context.Context, code uint) (err error) {
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
	// Return Answer
	qqs, err := ps.db.GetQuestionsByQuiz(s.Quiz.ID)
	if err != nil {
		return err
	}
	s.SetCurrentAnswer(qqs[s.CurrentQuestionIndex].Answer)
	return ps.db.UpdatePlaySession(s)
}

func (ps *PlaySessionSvc) GetPS(ctx context.Context, code uint) (s *models.PlaySession, err error) {
	s, err = ps.db.GetPlaySession(code)
	// SetQuestion if needed
	if s.State == models.StateInProgress || s.State == models.StateFinished {
		qqs, err := ps.db.GetQuestionsByQuiz(s.Quiz.ID)
		if err != nil {
			return s, err
		}
		s.UpdateQuestion(qqs[s.CurrentQuestionIndex])
	}
	return
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
