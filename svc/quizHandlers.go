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
			Name string `json:"name,omitempty"`
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
		qz, err := s.hub.CreateQuiz(req.Context(), r.Name)
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
		}
		resp := &Response{
			Quiz: qz,
		}
		s.respond(w, req, resp, http.StatusOK, nil)

	}
}
