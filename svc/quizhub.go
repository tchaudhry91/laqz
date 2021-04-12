package svc

import (
	"context"
	"errors"
	"fmt"

	"github.com/tchaudhry91/laqz/svc/models"
	"gorm.io/gorm"
)

var NotPermittedError = errors.New("User is not permitted for the following action")

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
	ToggleQuizPrivacy(ctx context.Context, id uint) (err error)
}

type QHub struct {
	db models.QuizStore
}

func NewQHub(db models.QuizStore) *QHub {
	return &QHub{
		db: db,
	}
}

func getUserFromContext(ctx context.Context, key contextKey) (u *models.User, err error) {
	user := ctx.Value(key)
	if user != nil {
		return user.(*models.User), nil
	}
	return nil, fmt.Errorf("User not found")
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
		return qz, err
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
