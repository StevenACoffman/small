package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

)

type Server struct {
	server *http.Server
	log *log.Logger
}

func main() {

	err := runServer()
	if err == nil {
		log.Println("finished clean")
		os.Exit(0)
	}	else {
		log.Printf("Got error: %v", err)
		os.Exit(1)
	}
}

func runServer() error {
	httpServer := newHTTPServer()

	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit,
		os.Interrupt,    // interrupt = SIGINT = Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)

	//defer func() {
	//	signal.Stop(quit)
	//	cancel()
	//}()

	go httpServer.shutdowner(ctx, cancel, quit)

	return httpServer.ListenAndServe()
}

// shutdowner listens for interrupt signals, and will trigger the server to gracefully terminate.
 func (httpServer *Server) shutdowner(ctx context.Context, cancel context.CancelFunc, quit chan os.Signal) {
	for {
		select {
		case sig := <-quit:
			log.Printf("Got %s signal. Aborting...\n", sig)
			cancel()
		case <-ctx.Done():
			//cleanup: on interrupt shutdown webserver
			httpServer.server.Shutdown(ctx)
		}
	}
}
func (httpServer *Server) ListenAndServe() error {
	if err := httpServer.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func newHTTPServer() *Server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":8080"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.HandleFunc("/", HealthCheck)
	logger := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	logger.Printf("HTTP Metrics server serving at %s", ":8080")
	return  &Server{httpServer, logger}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(200)
}
