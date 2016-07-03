// Package routes defines the http routes for this server, and the
// route handlers. See github.com/confyrm/gorest/router.
// For more info on routes, see github.com/gorilla/mux
package routes

import (
  "github.com/confyrm/gorest/router"
)

// RouteSet is the static set of http routes.  To add a new route:
// 1. Create a new route handler.  See Index.go in this package for Example.
// 2. Add a router.Route to RouteSet.
var RouteSet = router.Routes{
  router.Route {
    "Index",
    "GET",
    "/",
    Index,
  },
  router.Route {
    "SlashCommand",
    "POST",
    "/cmd",
    SlashRouter,
  },
}
