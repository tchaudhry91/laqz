package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	firebase "firebase.google.com/go"
	"github.com/go-kit/kit/log"
	"github.com/peterbourgon/ff"
	"github.com/tchaudhry91/laqz/svc"
	"github.com/tchaudhry91/laqz/svc/models"
)

func main() {
	fs := flag.NewFlagSet("qhub", flag.ExitOnError)
	var (
		listenAddr          = fs.String("listen-addr", "0.0.0.0:8080", "listen address")
		dbDSN               = fs.String("db-dsn", "postgresql://postgres:password@127.0.0.1:42261/laqz?sslmode=disable", "Database Connection String")
		firebaseKeyFile     = fs.String("firebase-admin-key", "", "Firebase Admin Private Key")
		fileUploadDirectory = fs.String("file-upload-dir", "/app/uploads", "Place to put uploaded assets")
		externalURL         = fs.String("external-url", "https://laqz-fs.tux-sudo.com", "External URL for uploaded assets")
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

	if *firebaseKeyFile != "" {
		initFirebase(*firebaseKeyFile)
	}

	firebaseApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	server := svc.NewQServer(hub, *listenAddr, logger, authClient, *fileUploadDirectory, *externalURL)
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

func initFirebase(firebaseKeyFile string) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", firebaseKeyFile)
}
