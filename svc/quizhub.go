package svc

import (
	"context"
	"errors"
	"fmt"

	"github.com/tchaudhry91/laqz/svc/models"
	"gorm.io/gorm"
)

type contextKey string

type QuizHub interface {
	UserContextKey() contextKey
	LogIn(ctx context.Context, user *models.User) (err error)
	CreateQuiz(ctx context.Context, name string) (qz *models.Quiz, err error)
	GetQuiz(ctx context.Context, id uint) (qz *models.Quiz, err error)
	GetMyQuizzes(ctx context.Context) (qqz []*models.Quiz, err error)
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

func (hub *QHub) CreateQuiz(ctx context.Context, name string) (qz *models.Quiz, err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return qz, err
	}
	u, err = hub.db.GetUserByEmail(u.Email)
	if err != nil {
		return qz, err
	}
	qz = models.NewQuiz(name, u)
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
	if qz.CanView(u.Email) {
		return nil, fmt.Errorf("This quiz has been kept private by the collaborators")
	}
	return
}

func (hub *QHub) GetMyQuizzes(ctx context.Context) (qqz []*models.Quiz, err error) {
	u, err := getUserFromContext(ctx, hub.UserContextKey())
	if err != nil {
		return qqz, err
	}
	qqz, err = hub.db.GetQuizzesByUser(u.Email)
	return
}
