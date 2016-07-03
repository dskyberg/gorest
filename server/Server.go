// Package server defines the Server type that main launches
package server

import (
  "os"
  "log"
  "fmt"
  "net/http"

  "github.com/gorilla/handlers"
  "github.com/confyrm/gorest/router"
  "github.com/confyrm/gorest/config"
)

type Server struct {
  Config *config.Config
  Name string
  Port int
  RouteSet router.Routes
}

// Server.Run is called in gorest.main, and launches the http router
// for this Server.
func (s *Server) Run() {
  router := s.RouteSet.New(s.Config)
  loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, router)
  log.Printf("%s: Listening on %d...", s.Name, s.Port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), loggedRouter))
}
