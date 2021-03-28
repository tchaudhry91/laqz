package svc

import (
	"fmt"
	"net/http"
)

func (s QServer) AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Header)
		idtoken, ok := req.Header["Token"]
		if !ok {
			s.respond(w, req, nil, http.StatusUnauthorized, fmt.Errorf("No Token Supplied"))
			return
		}
		token, err := s.authClient.VerifyIDToken(req.Context(), idtoken[0])
		if err != nil {
			s.respond(w, req, nil, http.StatusUnauthorized, fmt.Errorf("Could not verify token:%w", err))
		}
		fmt.Println(token)
		next.ServeHTTP(w, req)
	})
}
