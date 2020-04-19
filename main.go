package main

import (
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const defaultPort = "8080"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	router.HandleFunc("/", helloHandler).Methods(http.MethodGet)

	serverPort := envString("HEALTHCHECKER_PORT", defaultPort)
	srv := startServer(serverPort, router)
	waitForShutdown(srv)
}

func startServer(serverPort string, router *mux.Router) *http.Server {
	serverAddr := ":" + serverPort
	srv := &http.Server{Addr: serverAddr, Handler: router}

	go func() {
		log.WithFields(log.Fields{"url": serverAddr}).Info("starting the server")
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("{\"status\": \"OK\"}")); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello from " + os.Getenv("HOSTNAME"))); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func waitForShutdown(srv *http.Server) {
	killSignalChan := make(chan os.Signal, 1)
	signal.Notify(killSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)

	<-killSignalChan
	srv.Shutdown(context.Background())
	log.Info("Shutting down")
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
