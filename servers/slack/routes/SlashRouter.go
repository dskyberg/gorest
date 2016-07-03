package routes

import (
  "errors"
  "net/http"
  "fmt"
  "log"
  "strings"
  "encoding/json"
  "github.com/gorilla/schema"

  "github.com/confyrm/gorest/slack"
  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
  "github.com/confyrm/gorest/servers/slack/cmd"
)

var commands = cmd.SlashCommands {
  cmd.SlashCommand{
    "/devhub",
    false,
    cmd.DevHub,
  },
}

 var helpResponses = map[string]string {
   "/": `
  github accepts the following commands:
  new: create a new issue
  get: display an issue
  update: Update an issue
  close: Mark an issue as closed.
  help: this text
  `,
  "new": `
  github new accepts a set of <key>=<value> statements.  Do not use
  the '=' character within a value.  Unrecognized keys are ignored.
  Here are the supported keys:
  -
  title:      Required. Issue title
  body:       Optional. Text to add to the issue body
  labels:     Optionall A comma separated list of labels
  milestone:  Optional. The milestone must exist in the identified repository
  assignee:   Optional.  Github user to assign the issue to
  repository: Optional. If not provided, devhub will be used.
  -
  Example: /github new title = Here is my title labels = label1, label2
  `,
  "test.some.stuff": "Help for test some stuff",
}

// Make a static decoder for performance
var decoder = schema.NewDecoder()
var commandRouter = commands.New()

func SlashRouter(config *config.Config, rw http.ResponseWriter, req *http.Request) error {

    if err := req.ParseForm(); err != nil {
        return StatusError{http.StatusBadRequest, err}
    }

    sReq := &slack.Request{}

    if err := decoder.Decode(sReq, req.PostForm); err != nil {
      return StatusError{http.StatusBadRequest, err}
    }

    command, err := sReq.TextToCommand()
    if err != nil {
      return StatusError{http.StatusInternalServerError, err}
    }
    log.Printf("Received Slack slash command: %#v", command)
    // First see if this is a help request.  If so, return a help response
    if command.Commands[len(command.Commands)-1] == "help" {
      rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
      response := HelpResponse(command)
      if err := json.NewEncoder(rw).Encode(response); err != nil {
        return StatusError{http.StatusInternalServerError, err}
      }
      return nil
    }

    // At this point, we can either process the command and return a
    // slack.Response, or, we can  kick off a goroutine, and return a
    // quick response.  The goroutine can send the response as a POST
    // request to the response_url listed in the sReq.
    if sReq.Token != config.GetString("SLACK_TOKEN") {
      return StatusError{http.StatusUnauthorized, errors.New("Not authorized. Wrong Slack Token.")}
    }
    if handler, ok := commandRouter[sReq.Command]; ok {
      response, statusErr := handler.Handler(config, sReq)
      if statusErr != nil {
        return statusErr
      }
      rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
      if err := json.NewEncoder(rw).Encode(response); err != nil {
        return StatusError{http.StatusInternalServerError, err}
      }
      return nil
    } else {
      return StatusError{http.StatusBadRequest, errors.New(fmt.Sprintf("Command not found [%s]", sReq.Command))}
    }

}

/*
  Send a help text response, if the user requested it
*/

func HelpResponse(command *slack.DevHubCommand) *slack.Response {

  cLen := len(command.Commands)
  var text string
  // If the only command is help, then return the top level help
  if cLen == 1 && command.Commands[0] == "help" {
    text = helpResponses["/"]
  } else {
    // Create the help lookup key, by joining all the commands except the last
    // one, which should be "help"
    key := strings.Join(command.Commands[:cLen - 1], ".")
    t, ok := helpResponses[key]
    if !ok {
      // No help for this command.  Just return the top level help.
      log.Printf("Help request received for %s, but no help found\n", key)
      text = helpResponses["/"]
    } else {
      text = t
    }
  }

  response := slack.Response{slack.Ephemeral, text, nil}

  /*
  if err := json.NewEncoder(rw).Encode(response); err != nil {
    return StatusError{http.StatusInternalServerError, err}
  }

  return nil
  */
  return &response
}
