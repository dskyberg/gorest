// Package main provides the entrypoint into the gorest server.
package main

import (
  "os"
  "log"

  "github.com/confyrm/gorest/admin"
  "github.com/confyrm/gorest/servers/slack"
  "github.com/confyrm/gorest/config"
  github "github.com/confyrm/gorest/githubclient"
)

func main() {

  c := SetupConfig()
  //  Run the admin server as a go function.  Currently, all it supports is
  // /exit.
  go admin.Server(c)

  // Run the main app
  app := slack.New(c)
  app.Run()
}

func SetupConfig() *config.Config {
  // Check the env for the location of a config file. The configFile can be
  // just a path, or a path/file.  The configFile does not need an extension.
  defaultConfigName := "config"
  var configFile string
  if configFile = os.Getenv("GOREST_CONFIG"); configFile == "" {
    log.Printf("GOREST_CONFIG not set.  Using %s", defaultConfigName)
    configFile = defaultConfigName
  } else {
    log.Printf("GOREST_CONFIG is set to %s", configFile)
  }

  // Default configuration settings that need to have valid
  // values, even if not set in a config file or in the env.
  configDefaults := map[string]interface{} {
    "APP_PORT": 8080,
    "APP_NAME": "devhub",
    "APP_ROOT": "/Users/david/golang/src/github.com/confyrm/gorest",
    "ADMIN_PORT": 8001,
    github.DefaultOwner: "confyrm",
    github.DefaultRepo: "devhub",
  }
  log.Printf("Using the following defaults: %v.  To override, set in either the ENV or the config file", configDefaults)
  // Tell Viper to look for a file called 'config'.  The fileName can
  // be pathed, as well.  In which case, Viper will look in the path and
  // in the current directory '.', in that order.
  return config.New(configFile, configDefaults)
}
