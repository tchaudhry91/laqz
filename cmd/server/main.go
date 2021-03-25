package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/peterbourgon/ff"
	"github.com/tchaudhry91/laqz/svc"
	"github.com/tchaudhry91/laqz/svc/models"
)

func main() {
	fs := flag.NewFlagSet("qhub", flag.ExitOnError)
	var (
		listenAddr    = fs.String("listen-addr", "localhost:8080", "listen address")
		dbDSN         = fs.String("db-dsn", "postgresql://postgres:password@127.0.0.1:42261/laqz?sslmode=disable", "Database Connection String")
		auth0Domain   = fs.String("auth0-domain", "https://tux-sudo.us.auth0.com/", "auth0domain")
		auth0ClientID = fs.String("auth0-clientID", "9YNFdapfaDrGNu4ktFacwKpnHFK2hw8c", "auth0-client")
		auth0PEM      = fs.String("auth0-secret-pem", "tux-sudo.pem", "auth0-secret certificate")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("QHUB"))

	s, err := models.NewQuizPGStore(*dbDSN)
	if err != nil {
		panic(err)
	}
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	hub := svc.NewQHub(s)

	shutdown := make(chan error, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	server := svc.NewQServer(hub, *listenAddr, logger, *auth0ClientID, *auth0PEM, *auth0Domain)
	go func() {
		logger.Log("msg", "Starting server..", "listenAddr", *listenAddr)
		err = server.Start()
		shutdown <- err
	}()

	select {
	case signalKill := <-interrupt:
		logger.Log("msg", fmt.Sprintf("Stopping Server: %s", signalKill.String()))
	case err := <-shutdown:
		logger.Log("error", err)
	}

	err = server.Shutdown(context.TODO())
	if err != nil {
		logger.Log("error", err)
	}
}
