package svc

import (
	"fmt"

	"github.com/tchaudhry91/laqz/svc/models"
	"gorm.io/gorm"
)

type QuizHub interface {
	LogIn(email, name, avatarURL string) (err error)
}

type QHub struct {
	db models.QuizStore
}

func NewQHub(db models.QuizStore) *QHub {
	return &QHub{
		db: db,
	}
}

func (hub *QHub) LogIn(email, name, avatarURL string) (err error) {
	u, err := hub.db.GetUserByEmail(email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("Error while fetching User from database: %w", err)
		}
	}
	u = &models.User{Email: email}
	err = hub.db.CreateUser(u)
	return
}
