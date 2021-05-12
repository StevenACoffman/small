package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// =========================================================================
	// App Starting
	logger := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("main : Started")
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
	// =========================================================================
	// Start API Service
	api := NewHTTPServer(logger)
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logger.Printf("main : API listening on %s", api.Addr)
		// listen and serve blocks until error or shutdown is called
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	// listen for all interrupt signals, send them to quit channel
	signal.Notify(shutdown,
		os.Interrupt,    // interrupt = SIGINT = Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		logger.Fatalf("error: listening and serving: %s", err)
		return err

	case <-shutdown:
		logger.Println("runServer : Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			logger.Printf("runServer : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
			return err
		}
		return err
	}
}

// NewHTTPServer is factory function to initialize a new server
func NewHTTPServer(logger *log.Logger) *http.Server {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":3000"
	}

	s := &ServerHandler{}
	// pass logger
	s.SetLogger(logger)

	h := &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return h
}

// ServerHandler implements type http.Handler interface, with our logger
type ServerHandler struct {
	logger *log.Logger
	mux    *http.ServeMux
	once   sync.Once
}

// SetLogger provides external injection of logger
func (s *ServerHandler) SetLogger(logger *log.Logger) {
	s.logger = logger
}

// ServeHTTP satisfies Handler interface, sets up the Path Routing
func (s *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// on the first request only, lazily initialize
	s.once.Do(func() {
		if s.logger == nil {
			s.logger = log.New(os.Stdout,
				"INFO: ",
				log.Ldate|log.Ltime|log.Lshortfile)
			s.logger.Printf("Default Logger used")
		}
		s.mux = http.NewServeMux()
		s.mux.HandleFunc("/redirect", s.RedirectToHome)
		s.mux.HandleFunc("/health", HealthCheck)
		s.mux.HandleFunc("/", s.HelloHome)
	})

	s.mux.ServeHTTP(w, r)
}

func (s *ServerHandler) HelloHome(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("Hello, World!"))
	if err != nil {
		s.logger.Println("error writing hello world:", err)
	}
}

// HealthCheck verifies externally that the program is still responding
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(200)
}

// RedirectToHome Will Log the Request, and respond with a HTTP 303 to redirect to /
func (s *ServerHandler) RedirectToHome(w http.ResponseWriter, r *http.Request) {
	s.logger.Printf("Redirected request %v to /", r.RequestURI)
	w.Header().Add("location", "/")
	w.WriteHeader(http.StatusSeeOther)
}
