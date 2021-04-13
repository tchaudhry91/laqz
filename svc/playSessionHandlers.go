package svc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tchaudhry91/laqz/svc/models"
)

func (s *QServer) CreatePS() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			QuizID uint `json:"quiz_id,omitempty"`
		}
		type Response struct {
			PlaySession *models.PlaySession `json:"play_session,omitempty"`
			Err         error               `json:"err,omitempty"`
		}

		r := Request{}
		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, err)
			return
		}
		ps, err := s.hub.InitNewPS(req.Context(), r.QuizID)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		resp := Response{}
		resp.PlaySession = ps
		s.respond(w, req, resp, http.StatusOK, nil)
	}
}

func (s *QServer) GetPS() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Code uint `json:"code,omitempty"`
		}
		type Response struct {
			PlaySession *models.PlaySession `json:"play_session,omitempty"`
		}
		r := Request{}
		params := mux.Vars(req)
		idStr := params["code"]
		var id int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad Code supplied"))
			return
		}
		r.Code = uint(id)
		resp := Response{}
		ps, err := s.hub.GetPS(req.Context(), r.Code)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		resp.PlaySession = ps
		s.respond(w, req, resp, http.StatusOK, nil)
	}
}

func (s *QServer) JoinPS() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Code uint `json:"code,omitempty"`
		}
		r := Request{}
		params := mux.Vars(req)
		idStr := params["code"]
		var id int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad Code supplied"))
			return
		}
		r.Code = uint(id)
		s.hub.AddUserToPS(req.Context(), r.Code)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}
