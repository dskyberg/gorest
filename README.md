github.com/confyrm/gorest
========

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

Config variables can be placed in the environment, config the file, or both
