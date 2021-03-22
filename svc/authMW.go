package svc

import (
	"fmt"
	"net/http"

	"google.golang.org/api/idtoken"
)

func (s QServer) AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		headers := req.Header
		token := headers.Get("token")
		if token == "" {
			s.respond(w, req, nil, http.StatusUnauthorized, fmt.Errorf("No token supplied"))
			return
		}
		tinfo, err := idtoken.Validate(req.Context(), token, "")
		if err != nil {
			s.respond(w, req, nil, http.StatusUnauthorized, fmt.Errorf("No token supplied"))
			return
		}
		fmt.Println(tinfo)
		next.ServeHTTP(w, req)
	})
}
