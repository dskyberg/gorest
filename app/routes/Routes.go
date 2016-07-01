// Package routes defines the routes for this server, and the route handlers.
package routes

import (
  "github.com/confyrm/gorest/router"
)

var RouteSet = router.Routes{
  router.Route{
    "Index",
    "GET",
    "/",
    Index,
  },
  router.Route{
    "SlashCommand",
    "POST",
    "/cmd",
    SlashCommand,
  },
}
