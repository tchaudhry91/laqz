package svc

import (
	"errors"
	"fmt"

	"github.com/tchaudhry91/laqz/svc/models"
	"gorm.io/gorm"
)

type QuizHub interface {
	LogIn(user *models.User) (err error)
}

type QHub struct {
	db models.QuizStore
}

func NewQHub(db models.QuizStore) *QHub {
	return &QHub{
		db: db,
	}
}

func (hub *QHub) LogIn(user *models.User) (err error) {
	// Check if user exists
	u, err := hub.db.GetUserByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Error while fetching User from database: %w", err)
		} else {
			// User Not Found
			return hub.db.CreateUser(user)
		}
	}
	// User Found. Update as necessary
	return hub.db.UpdateUser(u)
}
