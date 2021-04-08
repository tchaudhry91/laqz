package svc

import (
	"net/http"
	"time"
)

// LoggingMW is a logging middleware
func (s *QServer) LoggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				s.logger.Log(
					"err", err,
				)
			}
		}()
		defer func(start time.Time) {
			s.logger.Log("method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
		}(time.Now())
		next.ServeHTTP(w, r)
	})
}
