package routes

import (
  "errors"
  "net/http"
  "fmt"

  "github.com/gorilla/schema"

  "github.com/confyrm/gorest/slack"
  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
  "github.com/confyrm/gorest/app/cmd"
)

var decoder = schema.NewDecoder()

type Command struct {
  Name string
  Handler func(config *config.Config, sReq *slack.Request, rw http.ResponseWriter, req *http.Request) error
}

var commands = map[string]Command {
  "/devhub": {"/devhub", cmd.DevHub },
}

func SlashCommand(config *config.Config, rw http.ResponseWriter, req *http.Request) error {

    if err := req.ParseForm(); err != nil {
        return StatusError{http.StatusBadRequest, err}
    }

    sReq := &slack.Request{}

    if err := decoder.Decode(sReq, req.PostForm); err != nil {
      return StatusError{http.StatusBadRequest, err}
    }

    if sReq.Token != config.GetString("SLACK_TOKEN") {
      return StatusError{http.StatusUnauthorized, errors.New("Not authorized")}
    }
    if command, ok := commands[sReq.Command]; ok {
      return command.Handler(config, sReq, rw, req)
    } else {
      return StatusError{http.StatusBadRequest, errors.New(fmt.Sprintf("Command not found [%s]", sReq.Command))}
    }

}
