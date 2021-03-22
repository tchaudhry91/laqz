package svc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type QServer struct {
	hub    QuizHub
	router *mux.Router
	server *http.Server
	logger log.Logger
}

func NewQServer(hub QuizHub, listenAddr string, logger log.Logger) *QServer {
	router := mux.NewRouter()
	s := &QServer{
		hub:    hub,
		router: router,
		server: &http.Server{
			Addr:    listenAddr,
			Handler: router,
		},
		logger: logger,
	}
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
	s.logger.Log("path", req.URL.Path, "method", req.Method, "err", err, "code", statusCode)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.Log("path", req.URL.Path, "method", req.Method, "err", err, "code", statusCode)
		}
	}
}

// Login is the handler to complete the login process
// GET /auth/google
func (s *QServer) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Yes. Quiz here"))
	}
}
