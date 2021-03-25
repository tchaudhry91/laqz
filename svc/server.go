package svc

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2"
)

type QServer struct {
	hub       QuizHub
	router    *mux.Router
	server    *http.Server
	logger    log.Logger
	validator *auth0.JWTValidator
}

func NewQServer(hub QuizHub, listenAddr string, logger log.Logger, auth0ClientID, auth0PEM, auth0Domain string) *QServer {
	router := mux.NewRouter()
	// Create a configuration with the Auth0 information
	pem, err := ioutil.ReadFile(auth0PEM)
	if err != nil {
		panic(err)
	}
	secret, err := LoadPublicKey(pem)
	if err != nil {
		panic(err)
	}
	secretProvider := auth0.NewKeyProvider(secret)
	configuration := auth0.NewConfiguration(secretProvider, []string{}, auth0Domain, jose.RS256)

	validator := auth0.NewValidator(configuration, nil)
	s := &QServer{
		hub:       hub,
		router:    router,
		logger:    logger,
		validator: validator,
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
	s.logger.Log("path", req.URL.Path, "method", req.Method, "err", err, "code", statusCode)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.Log("path", req.URL.Path, "method", req.Method, "err", err, "code", statusCode)
		}
	}
}

// Extracted from https://github.com/square/go-jose/blob/master/utils.go
// LoadPublicKey loads a public key from PEM/DER-encoded data.
// You can download the Auth0 pem file from `applications -> your_app -> scroll down -> Advanced Settings -> certificates -> download`
func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s' and '%s'", err0, err1)
}

// Login is the handler to complete the login process
// GET /auth/google
func (s *QServer) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s.respond(w, req, map[string]string{"name": "QuizHub"}, http.StatusOK, nil)
	}
}
