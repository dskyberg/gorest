// Package router uses github.com/gorilla/mux to define
// REST paths.
package router

import (
  "net/http"

  "github.com/gorilla/mux"
  "github.com/confyrm/gorest/router/handler"
  "github.com/confyrm/gorest/config"
)

// The Route struct holds the info needed to register a path and the
// associated path handler with mux.  See admin/routes/Routes.go and
// app/routes/Routes.go for examples.
type Route struct {
  // Name of the route, for no apparent reason.
  Name        string
  // The HTTP method for this route
  Method      string
  // The (potentially) patterned route path
  Pattern     string
  // The handler function for the route.
  HandlerFunc handler.EnvHandlerFunc
}

// Helper type for a slice of Routes.
type Routes []Route

// New returns the fully configured router for the set of routes.
// The config is not use by the Router directly, but it is passed down to
// all handler functions.  This is the simplest way to handle global
// config data.
func (routes Routes) New(config *config.Config) *mux.Router {
  // Setting StrictSlash allow a path to be
  // accepted with, or without a trailing slash.
  // So, a path of '/thispath' == a path of '/thispath/'
  router := mux.NewRouter().StrictSlash(true)
  for _, route := range routes {
    var h http.Handler
    //handler = route.HandlerFunc
    h = handler.Handler{config, route.HandlerFunc}

    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(h)
  }

  return router
}
