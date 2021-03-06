package svc

import (
	"context"
	"errors"
	"fmt"

	"github.com/tchaudhry91/laqz/svc/models"
	"gorm.io/gorm"
)

var NotPermittedError = errors.New("User is not permitted for the following action")
var UserNotFound = errors.New("User not found")

type contextKey string

type QuizHub interface {
	UserContextKey() contextKey
	LogIn(ctx context.Context, user *models.User) (err error)
	CreateQuiz(ctx context.Context, name string, tags []string) (qz *models.Quiz, err error)
	DeleteQuiz(ctx context.Context, id uint) (err error)
	GetQuiz(ctx context.Context, id uint) (qz *models.Quiz, err error)
	GetMyQuizzes(ctx context.Context) (qqz []*models.Quiz, err error)
	GetPublicQuizzes(ctx context.Context) (qqz []*models.Quiz, err error)
	AddQuestion(ctx context.Context, q *models.Question) (err error)
	UpdateQuestion(ctx context.Context, id, quizID uint, q *models.Question) (err error)
	DeleteQuestion(ctx context.Context, id, quizID uint) (err error)
	GetQuestions(ctx context.Context, quizID uint) (qq []*models.Question, err error)
	GetQuestion(ctx context.Context, id, quizID uint) (q *models.Question, err error)
	ToggleQuizPrivacy(ctx context.Context, id uint) (err error)

	PlaySessionSVC
}

type QHub struct {
	*PlaySessionSvc
	db models.QuizStore
}

func NewQHub(db models.QuizStore) *QHub {
	psSvc := NewPlaySessionSvc(db)
	return &QHub{
		psSvc,
		db,
	}
}

func getUserFromContext(ctx context.Context, key contextKey) (u *models.User, err error) {
	user := ctx.Value(key)
	if user != nil {
		return user.(*models.User), nil
	}
	return nil, UserNotFound
}

func (hub *QHub) UserContextKey() contextKey {
	var userContextKey = contextKey("user")
	return userContextKey
}

func (hub *QHub) LogIn(ctx context.Context, user *models.User) (err error) {
	// Check if user exists
	u, err := hub.db.GetUserByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Error while fetching User from database: %w", err)
		} else {
			// User Not Found
			return hub.db.CreateUser(user)
		}
	}
	// User Found. Update as necessary
	return hub.db.UpdateUser(u)
}

func (hub *QHub) CreateQuiz(ctx context.Context, name string, tags []string) (qz *models.Quiz, err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return qz, err
	}
	u, err = hub.db.GetUserByEmail(u.Email)
	if err != nil {
		return qz, err
	}

	// Get Existing Tags
	tt := []*models.Tag{}
	for _, t := range tags {
		tag, err := hub.db.GetTagByName(t)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return qz, err
			}
			tt = append(tt, &models.Tag{Name: t})
			continue
		}
		tt = append(tt, tag)
	}

	qz = models.NewQuiz(name, u, tt)
	err = hub.db.CreateQuiz(qz)
	if err != nil {
		return qz, err
	}
	return qz, nil
}

func (hub *QHub) GetQuiz(ctx context.Context, id uint) (qz *models.Quiz, err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		if !errors.Is(err, UserNotFound) {
			return qz, err
		}
		// Use Guest User
		u = &models.User{}
	}
	qz, err = hub.db.GetQuiz(id)
	if !qz.CanView(u.Email) {
		return nil, fmt.Errorf("This quiz has been kept private by the collaborators")
	}
	return
}

func (hub *QHub) ToggleQuizPrivacy(ctx context.Context, id uint) (err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return err
	}
	qz, err := hub.db.GetQuiz(id)
	if !qz.IsCollaborator(u.Email) {
		return NotPermittedError
	}
	qz.TogglePrivacy()
	return hub.db.UpdateQuiz(qz)
}

func (hub *QHub) GetMyQuizzes(ctx context.Context) (qqz []*models.Quiz, err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return qqz, err
	}
	qqz, err = hub.db.GetQuizzesByUser(u.Email)
	return
}

func (hub *QHub) GetPublicQuizzes(ctx context.Context) (qqz []*models.Quiz, err error) {
	qqz, err = hub.db.GetAllPublicQuizzes()
	return
}

func (hub *QHub) DeleteQuiz(ctx context.Context, id uint) (err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return err
	}
	qz, err := hub.db.GetQuiz(id)
	if err != nil {
		return err
	}
	if !qz.IsCollaborator(u.Email) {
		return NotPermittedError
	}
	return hub.db.DeleteQuiz(id)
}

func (hub *QHub) AddQuestion(ctx context.Context, q *models.Question) (err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return err
	}
	u, err = hub.db.GetUserByEmail(u.Email)
	if err != nil {
		return err
	}

	q.UserID = u.ID

	qz, err := hub.db.GetQuiz(q.QuizID)
	if err != nil {
		return err
	}
	if !qz.IsCollaborator(u.Email) {
		return NotPermittedError
	}
	qz.AddQuestion(q)
	return hub.db.UpdateQuiz(qz)
}

func (hub *QHub) GetQuestions(ctx context.Context, quizID uint) (qq []*models.Question, err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return qq, err
	}
	qz, err := hub.db.GetQuiz(quizID)
	if err != nil {
		return qq, err
	}
	if !qz.IsCollaborator(u.Email) {
		return qq, NotPermittedError
	}

	return hub.db.GetQuestionsByQuiz(quizID)
}

func (hub *QHub) GetQuestion(ctx context.Context, id uint, quizID uint) (q *models.Question, err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return q, err
	}
	qz, err := hub.db.GetQuiz(quizID)
	if err != nil {
		return q, err
	}
	if !qz.IsCollaborator(u.Email) {
		return q, NotPermittedError
	}
	return hub.db.GetQuestion(id)
}

func (hub *QHub) UpdateQuestion(ctx context.Context, id, quizID uint, q *models.Question) (err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return err
	}
	qz, err := hub.db.GetQuiz(quizID)
	if err != nil {
		return err
	}
	if !qz.IsCollaborator(u.Email) {
		return NotPermittedError
	}
	return hub.db.UpdateQuestion(id, q)
}

func (hub *QHub) DeleteQuestion(ctx context.Context, id, quizID uint) (err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return err
	}
	qz, err := hub.db.GetQuiz(quizID)
	if err != nil {
		return err
	}
	if !qz.IsCollaborator(u.Email) {
		return NotPermittedError
	}
	return hub.db.DeleteQuestion(id, quizID)
}
