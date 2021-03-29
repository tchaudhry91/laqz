package svc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tchaudhry91/laqz/svc/models"
)

func (s QServer) AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		idtoken, ok := req.Header["Token"]
		if !ok {
			s.respond(w, req, nil, http.StatusUnauthorized, fmt.Errorf("No Token Supplied"))
			return
		}
		token, err := s.authClient.VerifyIDToken(req.Context(), idtoken[0])
		if err != nil {
			s.respond(w, req, nil, http.StatusUnauthorized, fmt.Errorf("Could not verify token:%w", err))
		}
		user := models.User{Email: token.Claims["email"].(string), Name: token.Claims["name"].(string), AvatarURL: token.Claims["picture"].(string)}
		ctx := req.Context()
		ctx = context.WithValue(ctx, "user", user)
		req = req.WithContext(ctx)
		err = s.hub.LogIn(&user)
		if err != nil {
			s.respond(w, req, nil, http.StatusInternalServerError, err)
			return
		}
		next.ServeHTTP(w, req)
	})
}
