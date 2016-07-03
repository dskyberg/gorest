// Package main provides the entrypoint into the gorest server.
package main

import (
    "github.com/confyrm/gorest/admin"
    "github.com/confyrm/gorest/servers/slack"
    "github.com/confyrm/gorest/config"
)

func main() {
  // Tell Viper to look for a file called 'config'.  The fileName can
  // be pathed, as well.  In which case, Viper will look in the path and
  // in the current directory '.', in that order.
  c := config.New("config", "devhub")

  //  Run the admin server as a go function.  Currently, all it supports is
  // /exit.
  go admin.Server(c)

  // Run the main app
  app := slack.New(c)
  app.Run()
}
