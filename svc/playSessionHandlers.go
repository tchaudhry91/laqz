package svc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/tchaudhry91/laqz/svc/models"
)

func (s *QServer) UserContextKey() contextKey {
	var userContextKey = contextKey("user")
	return userContextKey
}

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
		// Create a new wsHub for the playSession
		s.wsHubs[ps.Code] = newHub()
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
		err = s.hub.AddUserToPS(req.Context(), r.Code)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastReload()
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) AddTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Code     uint   `json:"code,omitempty"`
			TeamName string `json:"team_name,omitempty"`
		}
		r := Request{}
		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, err)
			return
		}
		params := mux.Vars(req)
		idStr := params["code"]
		var id int
		id, err = strconv.Atoi(idStr)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad Code supplied"))
			return
		}
		r.Code = uint(id)
		team := models.NewTeam(r.TeamName)
		err = s.hub.AddTeamToPS(req.Context(), r.Code, team)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastReload()
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) AddUserToTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Code     uint   `json:"code,omitempty"`
			TeamName string `json:"team_name,omitempty"`
			Email    string `json:"email,omitempty"`
		}
		r := Request{}
		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, err)
			return
		}
		params := mux.Vars(req)
		idStr := params["code"]
		var id int
		id, err = strconv.Atoi(idStr)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad Code supplied"))
			return
		}
		r.Code = uint(id)
		err = s.hub.AddUserToTeam(req.Context(), r.Code, r.TeamName, r.Email)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastReload()
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) StartPS() http.HandlerFunc {
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
		err = s.hub.StartPS(req.Context(), r.Code)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastReload()
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) EndPS() http.HandlerFunc {
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
		err = s.hub.EndPlaySession(req.Context(), r.Code)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.respond(w, req, nil, http.StatusNoContent, nil)
			return
		}
		s.wsHubs[r.Code].BroadcastReload()
		// Give it a few seconds
		time.Sleep(3 * time.Second)

		// CleanUp Connections and delete Hub
		for c := range s.wsHubs[r.Code].connections {
			s.logger.Log("msg", "Closing WS Connections", "ps-code", r.Code)
			s.wsHubs[r.Code].removeConnection(c)
		}
		delete(s.wsHubs, r.Code)
		s.logger.Log("msg", "Deleted Websocket Session", "ps-code", r.Code)

		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) IncrementPSQuestion() http.HandlerFunc {
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
		err = s.hub.IncrementPSQuestion(req.Context(), r.Code)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastReload()
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) DecrementPSQuestion() http.HandlerFunc {
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
		err = s.hub.DecrementPSQuestion(req.Context(), r.Code)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastReload()
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) AddPSTeamPoints() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Code     uint   `json:"code,omitempty"`
			Points   int    `json:"points,omitempty"`
			TeamName string `json:"team_name,omitempty"`
		}
		r := Request{}
		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, err)
			return
		}
		params := mux.Vars(req)
		idStr := params["code"]
		var id int
		id, err = strconv.Atoi(idStr)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad Code supplied"))
			return
		}
		r.Code = uint(id)
		err = s.hub.UpdateTeamPoints(req.Context(), r.Code, r.Points, r.TeamName)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastReload()
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) BroadcastChatMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Code    uint   `json:"code,omitempty"`
			Message string `json:"message,omitempty"`
			User    string `json:"user,omitempty"`
		}
		r := Request{}
		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, err)
			return
		}
		params := mux.Vars(req)
		idStr := params["code"]
		var id int
		id, err = strconv.Atoi(idStr)
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad Code supplied"))
			return
		}
		r.Code = uint(id)
		fmt.Println(r)
		u, err := getUserFromContext(req.Context(), s.UserContextKey())
		if err != nil {
			s.respond(w, req, nil, http.StatusBadRequest, fmt.Errorf("Bad User Supplied"))
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastChat(u.Name, r.Message)
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) RevealPSCurrentAnswer() http.HandlerFunc {
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
		err = s.hub.RevealPSCurrentAnswer(req.Context(), r.Code)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.wsHubs[r.Code] = newHub()
		}
		s.wsHubs[r.Code].BroadcastReload()
		s.respond(w, req, nil, http.StatusNoContent, nil)
	}
}

func (s *QServer) WebSocketPS() http.HandlerFunc {
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
			return
		}
		r.Code = uint(id)
		if _, ok := s.wsHubs[r.Code]; !ok {
			s.logger.Log("msg", "Attempted WS Connection to non-existant hub")
			return
		}
		wsConn, err := s.wsUpgrader.Upgrade(w, req, nil)
		if err != nil {
			s.logger.Log("msg", "Failed to upgrade WS", "err", err)
			return
		}
		c := &connection{send: make(chan []byte, 256), h: s.wsHubs[r.Code]}
		c.h.addConnection(c)
		defer c.h.removeConnection(c)
		var wg sync.WaitGroup
		wg.Add(2)
		go c.writer(&wg, wsConn)
		wg.Wait()
		wsConn.Close()
	}
}
