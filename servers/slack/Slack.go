// Package slack is a gorest server configuration that is runnable by main.
package slack

import (
  "log"
  "github.com/confyrm/gorest/config"
  "github.com/confyrm/gorest/server"
  . "github.com/confyrm/gorest/servers/slack/routes"
)

// Prefix gets added to the config lookups.
var Prefix = "APP_"

// New returns a configured server.Server that can be Run by main.  It would
// be great if I knew how to reflect a package, so that main can just find
// this New.  But for now, the package has to be loaded in main manually, and
// this New specifically called.
func New(c *config.Config) *server.Server {
  // Do some checks to make sure all required configs are present, etc.
  if !c.IsSet("SLACK_TOKEN") {
    log.Fatal("No Slack Token found. Check your config.")
  }
  if !c.IsSet("GITHUB_TOKEN") {
    log.Fatal("No GitHub Token found. Check your config.")
  }

  s := server.Server {
    c,
    c.GetString(config.Key(Prefix, "NAME")),
    c.GetInt(config.Key(Prefix, "PORT")),
    RouteSet,
  }
  return &s
}
