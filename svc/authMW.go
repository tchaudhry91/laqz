package svc

import (
	"context"
	"net/http"
)

func (s QServer) AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := s.validator.ValidateRequest(req)
		if err != nil {
			s.respond(w, req, nil, http.StatusUnauthorized, err)
			return
		}
		claims := map[string]interface{}{}
		err = token.Claims(token, &claims)
		if err != nil {
			s.respond(w, req, nil, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(req.Context(), "claims", claims)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
