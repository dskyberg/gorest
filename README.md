github.com/confyrm/gorest
========
The gorest package was developed to provide a simple custom Slack app.  It uses
the github.com/gorilla/mux package for http routing.  

There are a couple utility packages buried in gorest.  I make no apologies for not
making them top level libs. :)

- github.com/gorest/slack: Defines Slack custom app request and response structs,
and support for defining Slack slash commands, similar to http routing.  Ie,
a map of command to handler.

- github.com/gorest/slack: A very simple wrapper around Viper.

- github.com/gorest/router: A wrapper for gorilla/mux Router that lets you Define
a set of routes that use the following handler wrapper

- github.com/gorest/router/handler: A wrapper for http route handler that allows
the Config instance to be passed down.

- github.com/gorest/errors: A set of HTTP error responders

- github.com/gorest/server: Provides a simple Run comand that runs route handlers
wrapped with a logging handler, for consistent logging.

Required env variables:
-----

- SLACK_TOKEN: Requests from Slack will contain a token.  The server matches against this token for authorization.

- GITHUB_TOKEN: This is just for commands that interact with GitHub.  Currently, the server uses the user `devhub@confyrm.com`.

Consider using Sneaker to get secrets from S3
----
/<env>/console/github/token

/<env>/console/slack/token

Installation
-----

```
go get github.com/confyrm/gorest
go test -v
go build
go Installation
```

Configuration
---

Gorest uses `github.com/spf13/viper` for configuration management. The env variable GOREST_CONFIG can be set to tell the app where the config file is located.  If not set, the app will look for `./config.<ext>`.  

Viper supports JSON, TOML, YAML, HCL, and Java properties files.  So the file can have any appropriate extension, and Viper will find it.

Config variables can be placed in the environment, config the file, or both.  
Env takes precedence over config file.
