package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ServerHandler implements type http.Handler interface, with our logger
type ServerHandler struct {
	logger *log.Logger
	mux    *http.ServeMux
	once   sync.Once
}

func (s *ServerHandler) SetLogger(logger *log.Logger) {
	s.logger = logger
}

func (s *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// on the first request only, lazily initialize
	s.once.Do(func() {
		if s.logger == nil {
			fmt.Println("No logger")
		}
		s.mux = http.NewServeMux()
		s.mux.HandleFunc("/favicon.ico", s.RedirectToHome)
		s.mux.HandleFunc("/health", HealthCheck)
	})

	s.mux.ServeHTTP(w, r)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(200)
}

func (s *ServerHandler) RedirectToHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("location", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func main() {
	logger := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	err := runServer(logger)
	if err == nil {
		logger.Println("finished clean")
		os.Exit(0)
	} else {
		logger.Printf("Got error: %v", err)
		os.Exit(1)
	}
}

func runServer(logger *log.Logger) error {

	httpServer := NewHTTPServer(logger)
	// make a buffered channel for Signals
	quit := make(chan os.Signal, 1)

	// listen for all interrupt signals, send them to quit channel
	signal.Notify(quit,
		os.Interrupt,    // interrupt = SIGINT = Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)

	// receive signals on quit channel, tell server to shutdown
	go func() {
		//cleanup: on interrupt shutdown webserver
		<-quit
		err := httpServer.Shutdown(context.Background())

		if err != nil {
			logger.Printf("An error occurred on shutdown: %v", err)
		}
	}()

	// listen and serve until error or shutdown is called
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// NewHTTPServer is factory function to initialize a new server
func NewHTTPServer(logger *log.Logger) *http.Server {

	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8080"
	}

	s := &ServerHandler{}
	//pass logger
	s.SetLogger(logger)

	h := &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return h
}
