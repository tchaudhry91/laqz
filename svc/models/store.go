package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuizStore interface {
	CreateUser(u *User) error
	UpdateUser(u *User) error
	GetUserByEmail(email string) (u *User, err error)
	CreateQuiz(qz *Quiz) error
	DeleteQuiz(id uint) error
	GetQuizByName(name string) (qz *Quiz, err error)
	GetQuiz(id uint) (qz *Quiz, err error)
	GetPreloadedQuiz(id uint) (qz *Quiz, err error)
	GetAllPublicQuizzes() (qzs []*Quiz, err error)
	GetQuizzesByUser(email string) (qzs []*Quiz, err error)
	UpdateQuiz(qz *Quiz) error
	CreateQuestion(q *Question) error
	GetQuestion(id uint) (q *Question, err error)
	GetQuestionsByQuiz(qzID uint) (qq []*Question, err error)
	GetTagByName(name string) (t *Tag, err error)
	CreatePlaySession(s *PlaySession) error
	UpdatePlaySession(s *PlaySession) error
	GetPlaySession(code uint) (s *PlaySession, err error)
	DeletePlaySession(code uint) (err error)
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
	err = st.Migrate()
	return st, err

}

func (db *QuizPGStore) Migrate() error {
	return db.client.AutoMigrate(
		&User{},
		&Question{},
		&Quiz{},
		&PlaySession{},
		&Team{},
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

func (db *QuizPGStore) GetQuizByName(name string) (qz *Quiz, err error) {
	qz = &Quiz{}
	err = db.client.Where("name = ?", name).First(qz).Error
	return
}

func (db *QuizPGStore) GetQuiz(id uint) (qz *Quiz, err error) {
	qz = &Quiz{}
	err = db.client.Preload("Collaborators").Preload("Tags").Preload("Questions").Where("id = ?", id).First(qz).Error
	return
}

func (db *QuizPGStore) GetPreloadedQuiz(id uint) (qz *Quiz, err error) {
	qz = &Quiz{}
	err = db.client.Preload(clause.Associations).Where("id = ?", id).First(qz).Error
	return
}

func (db *QuizPGStore) GetAllPublicQuizzes() (qzs []*Quiz, err error) {
	qzs = make([]*Quiz, 0)
	err = db.client.Preload("Collaborators").Preload("Tags").Where("private = false").Find(&qzs).Error
	return
}

func (db *QuizPGStore) GetQuizzesByUser(email string) (qzs []*Quiz, err error) {
	u := &User{}
	err = db.client.Preload("Quizzes.Tags").Preload("Quizzes.Collaborators").Preload("Quizzes").Where("email = ?", email).First(u).Error
	qzs = u.Quizzes
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

func (db *QuizPGStore) DeleteQuiz(id uint) (err error) {
	err = db.client.Delete(&Quiz{}, id).Error
	if err != nil {
		return
	}
	return
}

func (db *QuizPGStore) GetQuestionsByQuiz(qzID uint) (qq []*Question, err error) {
	qz := &Quiz{}
	err = db.client.Preload("Questions").First(qz, qzID).Error
	if err != nil {
		return qq, err
	}
	return qz.Questions, nil
}

func (db *QuizPGStore) GetTagByName(name string) (t *Tag, err error) {
	t = &Tag{}
	err = db.client.Where("name = ?", name).First(t).Error
	return
}

func (db *QuizPGStore) CreatePlaySession(s *PlaySession) error {
	return db.client.Create(s).Error
}

func (db *QuizPGStore) GetPlaySession(code uint) (s *PlaySession, err error) {
	s = &PlaySession{}
	err = db.client.Preload("Quiz").Preload("Users").Preload("Teams").Preload("Teams.Users").Where("code = ?", code).First(s).Error
	return
}

func (db *QuizPGStore) DeletePlaySession(code uint) (err error) {
	return db.client.Where("code = ?", code).Delete(&PlaySession{}).Error
}

func (db *QuizPGStore) UpdatePlaySession(s *PlaySession) error {
	return db.client.Save(s).Error
}
