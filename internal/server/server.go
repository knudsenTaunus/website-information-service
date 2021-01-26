package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/knudsenTaunus/website-information-service/internal/handler"
)

type WebsiteInformationHandler interface {
	GetWebsiteInformation(address string)
}

type Server struct {
	srv     *http.ServeMux
	handler http.HandlerFunc
}

func New(websiteInformationHandler http.HandlerFunc) *Server {
	mux := http.NewServeMux()
	return &Server{srv: mux, handler: websiteInformationHandler}
}

func (s *Server) Start() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.srv.Handle("/getWebsiteInformation", handler.InputValidationHandler(s.handler))
		http.ListenAndServe(":8080", s.srv)
	}()

	log.Printf("Server started")
	<-done
	log.Printf("Server stopped")
}
