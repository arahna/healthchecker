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

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	srv := startServer(":8000")
	waitForShutdown(srv)
}

func startServer(serverUrl string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/health/", healthHandler).Methods(http.MethodGet)

	srv := &http.Server{Addr: serverUrl, Handler: router}

	go func() {
		log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
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

func waitForShutdown(srv *http.Server) {
	killSignalChan := make(chan os.Signal, 1)
	signal.Notify(killSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)

	<-killSignalChan
	srv.Shutdown(context.Background())
	log.Info("Shutting down")
}
