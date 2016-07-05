package routes

import (
  "io/ioutil"
	"path/filepath"
  "errors"
  "net/http"
  "fmt"
  "log"
  "strings"
  "encoding/json"
  "github.com/gorilla/schema"
  "github.com/hashicorp/hcl"

  "github.com/confyrm/gorest/slack"
  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
  . "github.com/confyrm/gorest/servers/slack/commands"
)



// Make a static decoder for performance
var decoder = schema.NewDecoder()
var commandRouter = SlashCommands.New()

// This is set the first time HelpResponse is called.  It is set by Loading
// The help text file located in ../../../help/help.hcl
var helpResponses map[string]interface{}

// SlashRouter is the top level slash command router.
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
  if (len(command.Commands) == 0 && len(command.Params) == 0) ||
    command.Commands[len(command.Commands)-1] == "help" {
    rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
    response := HelpResponse(config, command)
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

// HelpResponse looks up help text from the global helpResponses variable.
// The slice of commands are joined to create the lookup key.
func HelpResponse(config *config.Config, command *slack.DevHubCommand) *slack.Response {
  // Although called every time, this will only attempted to read the help
  // text the first time its called.
  ImportHelpText(config)

  cLen := len(command.Commands)
  var text string
  // If the only command is help, then return the top level help
  if cLen == 1 && command.Commands[0] == "help" {
    text = helpResponses["/"].(string)
  } else {
    // Create the help lookup key, by joining all the commands except the last
    // one, which should be "help"
    key := strings.Join(command.Commands[:cLen - 1], ".")
    t, ok := helpResponses[key].(string)
    if !ok {
      // No help for this command.  Just return the top level help.
      log.Printf("Help request received for %s, but no help found\n", key)
      text = helpResponses["/"].(string)
    } else {
      text = t
    }
  }

  response := slack.Response{slack.Ephemeral, text, nil}

  return &response
}

// ImportHelpText attempts to read an HCL formated file located in
// APP_ROOT/help/help.hcl.
func ImportHelpText(config *config.Config) {
  if helpResponses != nil {
    return
  }

  helpPath := filepath.Join(config.GetString("APP_ROOT"), "help", "help.hcl" )
  var err error
  if helpResponses, err = ParseHelpText(helpPath); err != nil {
    helpResponses = make(map[string]interface{})
  }
}

func ParseHelpText(helpPath string) (map[string]interface{}, error) {

  helpText, err := ioutil.ReadFile(helpPath)
  if err != nil {
    log.Printf("Could not load help text from %s: %#v", helpPath, err)
    return nil, err
  }
  var out interface{}
  err = hcl.Decode(&out, string(helpText))
  if err != nil {
    log.Printf("Could not parse help text in %s: %#v", helpPath, err)
    return nil, err
  }
  return out.(map[string]interface{}), nil
}
