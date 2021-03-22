package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type QuizStore interface {
	CreateUser(u *User) error
	UpdateUser(u *User) error
	GetUserByEmail(email string) (u *User, err error)
	CreateQuiz(qz *Quiz) error
	GetQuiz(name string) (qz *Quiz, err error)
	GetAllQuizzes() (qzs []*Quiz, err error)
	GetQuizzesByUser(email string) (qzs []*Quiz, err error)
	UpdateQuiz(qz *Quiz) error
	CreateQuestion(q *Question) error
	GetQuestion(id uint) (q *Question, err error)
	GetQuestionsByQuiz(qzID uint) (qq []Question, err error)
}

type QuizPGStore struct {
	client *gorm.DB
}

func NewQuizPGStore(dsn string) (*QuizPGStore, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return &QuizPGStore{}, err
	}
	st := &QuizPGStore{client: db}
	return st, nil

}

func (db *QuizPGStore) Migrate() error {
	return db.client.AutoMigrate(
		&User{},
		&Question{},
		&Quiz{},
		&Buzz{},
	)
}

func (db *QuizPGStore) CreateUser(u *User) error {
	return db.client.Create(u).Error
}

func (db *QuizPGStore) UpdateUser(u *User) error {
	return db.client.Save(u).Error
}

func (db *QuizPGStore) GetUserByEmail(email string) (*User, error) {
	u := User{}
	err := db.client.Where("email = ?", email).First(&u).Error
	return &u, err
}

func (db *QuizPGStore) CreateQuiz(qz *Quiz) error {
	return db.client.Create(qz).Error
}

func (db *QuizPGStore) GetQuiz(name string) (qz *Quiz, err error) {
	qz = &Quiz{}
	err = db.client.Where("name = ?", name).First(qz).Error
	return
}

func (db *QuizPGStore) GetAllQuizzes() (qzs []*Quiz, err error) {
	qzs = make([]*Quiz, 0)
	err = db.client.Find(qzs).Error
	return
}

func (db *QuizPGStore) GetQuizzesByUser(email string) (qzs []*Quiz, err error) {
	qzs = make([]*Quiz, 0)
	err = db.client.Where("email = ?", email).Find(qzs).Error
	return
}

func (db *QuizPGStore) UpdateQuiz(qz *Quiz) error {
	return db.client.Save(qz).Error
}

func (db *QuizPGStore) CreateQuestion(q *Question) error {
	return db.client.Create(q).Error
}

func (db *QuizPGStore) GetQuestion(id uint) (q *Question, err error) {
	q = &Question{}
	err = db.client.First(q, id).Error
	return
}

func (db *QuizPGStore) GetQuestionsByQuiz(qzID uint) (qq []Question, err error) {
	qz := &Quiz{}
	err = db.client.Preload("Questions").First(qz, qzID).Error
	if err != nil {
		return qq, err
	}
	return qz.Questions, nil
}
