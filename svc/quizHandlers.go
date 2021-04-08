package svc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tchaudhry91/laqz/svc/models"
	"gorm.io/gorm"
)

func (s *QServer) CreateQuiz() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Name string   `json:"name,omitempty"`
			Tags []string `json:"tags,omitempty"`
		}
		type Response struct {
			Quiz *models.Quiz `json:"quiz,omitempty"`
		}
		r := Request{}

		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, err)
			return
		}
		qz, err := s.hub.CreateQuiz(req.Context(), r.Name, r.Tags)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		resp := Response{
			Quiz: qz,
		}
		s.respond(w, req, resp, http.StatusCreated, nil)
	}
}

func (s *QServer) GetQuiz() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			ID uint `json:"id,omitempty"`
		}
		type Response struct {
			Quiz *models.Quiz `json:"quiz,omitempty"`
		}

		r := Request{}
		params := mux.Vars(req)
		idStr := params["id"]
		var id int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad ID supplied"))
			return
		}
		r.ID = uint(id)

		qz, err := s.hub.GetQuiz(req.Context(), r.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.respond(w, req, nil, http.StatusNotFound, nil)
				return
			}
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		resp := &Response{
			Quiz: qz,
		}
		s.respond(w, req, resp, http.StatusOK, nil)

	}
}

func (s *QServer) GetMyQuizzes() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Response struct {
			Quizzes []*models.Quiz `json:"quizzes"`
		}
		resp := Response{
			Quizzes: []*models.Quiz{},
		}
		qqz, err := s.hub.GetMyQuizzes(req.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.respond(w, req, resp, http.StatusOK, nil)
				return
			}
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		resp.Quizzes = qqz
		s.respond(w, req, resp, http.StatusOK, nil)
	}
}

func (s *QServer) DeleteQuiz() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			ID uint `json:"id,omitempty"`
		}
		type Response struct {
			Err string `json:"err,omitempty"`
		}

		r := Request{}
		params := mux.Vars(req)
		idStr := params["id"]
		var id int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad ID supplied"))
			return
		}
		resp := Response{}
		r.ID = uint(id)
		err = s.hub.DeleteQuiz(req.Context(), r.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				resp.Err = err.Error()
				s.respond(w, req, resp, http.StatusNotFound, nil)
				return
			}
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *QServer) AddQuestion() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			QuizID       uint   `json:"quiz_id,omitempty"`
			Text         string `json:"text,omitempty"`
			ImageLink    string `json:"image_link,omitempty"`
			AudioLink    string `json:"audio_link,omitempty"`
			Answer       string `json:"answer,omitempty"`
			Points       uint   `json:"points,omitempty"`
			TimerSeconds uint   `json:"timer_seconds,omitempty"`
		}

		type Response struct {
			Err string `json:"err,omitempty"`
		}
		r := Request{}

		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, err)
			return
		}

		if r.Text == "" || r.Answer == "" {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("You must supply atleast some text and an answer"))
			return
		}

		q := models.NewQuestion(r.QuizID, r.Text, r.ImageLink, r.AudioLink, r.Answer, r.Points, r.TimerSeconds)

		err = s.hub.AddQuestion(req.Context(), q)

		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, req, nil, http.StatusCreated, nil)

	}
}
