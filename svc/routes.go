package svc

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// routes registers the handlers to the specified routes
func (s *QServer) routes() {
	s.router.Use(s.LoggingMW)
	s.router.Handle("/healthz", s.Health()).Methods("GET")

	// Quiz Routes
	quizRoutes := s.router.PathPrefix("/quiz").Subrouter()
	quizRoutes.Handle("/", s.AuthMW(s.CreateQuiz())).Methods("POST")
	quizRoutes.Handle("/{id}/", s.OptionalAuthMW(s.GetQuiz())).Methods("GET")
	quizRoutes.Handle("/{id}/", s.AuthMW(s.DeleteQuiz())).Methods("DELETE")
	quizRoutes.Handle("/{id}/toggleVisibility", s.AuthMW(s.ToggleQuizPrivacy())).Methods("PATCH")
	quizRoutes.Handle("/{id}/addQuestion", s.AuthMW(s.AddQuestion())).Methods("POST")
	quizRoutes.Handle("/list/user/", s.AuthMW(s.GetMyQuizzes())).Methods("GET")
	quizRoutes.Handle("/list", s.GetQuizzes()).Methods("GET")

	// PlaySessionRoutes
	psRoutes := s.router.PathPrefix("/ps").Subrouter()
	psRoutes.Handle("/create", s.AuthMW(s.CreatePS())).Methods("POST")
	psRoutes.Handle("/{code}/", s.AuthMW(s.GetPS())).Methods("GET")
	psRoutes.Handle("/join/{code}", s.AuthMW(s.JoinPS())).Methods("POST")
	psRoutes.Handle("/{code}/addTeam", s.AuthMW(s.AddTeam())).Methods("POST")
	psRoutes.Handle("/{code}/assignTeamToUser", s.AuthMW(s.AddUserToTeam())).Methods("POST")
	psRoutes.Handle("/ws/{code}", s.WebSocketPS())
}

// CorsMW is a middleware to add CORS header to the response
func (s *QServer) CorsMW() http.Handler {
	headers := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Content-Type", "access-control-allow-origin", "content-type", "access-control-allow-headers", "token", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"*"})
	return handlers.CORS(headers, methods, origins)(s.router)

}
