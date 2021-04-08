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
	quizRoutes.Handle("/{id}/", s.AuthMW(s.GetQuiz())).Methods("GET")
	quizRoutes.Handle("/{id}/", s.AuthMW(s.DeleteQuiz())).Methods("DELETE")
	quizRoutes.Handle("/list/user/", s.AuthMW(s.GetMyQuizzes())).Methods("GET")
}

// CorsMW is a middleware to add CORS header to the response
func (s *QServer) CorsMW() http.Handler {
	headers := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Content-Type", "access-control-allow-origin", "access-control-allow-headers", "token", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"*"})
	return handlers.CORS(headers, methods, origins)(s.router)

}
