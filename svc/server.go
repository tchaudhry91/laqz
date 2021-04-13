package svc

import (
	"context"
	"encoding/json"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type QServer struct {
	hub        QuizHub
	router     *mux.Router
	server     *http.Server
	logger     log.Logger
	authClient *auth.Client
}

func NewQServer(hub QuizHub, listenAddr string, logger log.Logger, authClient *auth.Client) *QServer {
	router := mux.NewRouter()
	s := &QServer{
		authClient: authClient,
		hub:        hub,
		router:     router,
		logger:     logger,
	}
	s.server = &http.Server{Addr: listenAddr, Handler: s.CorsMW()}
	s.routes()
	return s
}

// Start begins listening for requests on the listenAddr. Blocks
func (s *QServer) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully terminates the server
func (s *QServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// respond is a internal utility to set proper HTTP responses
func (s *QServer) respond(w http.ResponseWriter, req *http.Request, data interface{}, statusCode int, err error) {
	w.WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.Log("path", req.URL.Path, "method", req.Method, "err", err, "code", statusCode)
		}
		return
	}
	if err != nil {
		err := json.NewEncoder(w).Encode(struct{ Err string }{Err: err.Error()})
		if err != nil {
			s.logger.Log("path", req.URL.Path, "method", req.Method, "err", err, "code", statusCode)
		}
		return
	}
}

// Health is the handler to check if the service is up
func (s *QServer) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s.respond(w, req, map[string]string{"name": "QuizHub"}, http.StatusOK, nil)
	}
}
